package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"golang.org/x/crypto/bcrypt"
)

type SubscriptionRepo struct {
	db *database.DB
}

func NewSubscriptionRepo(db *database.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) createTracking(tx *sql.Tx, campaignHash string, userIP, userAgent string, referrerURL string) (int, error) {
	var platformLimitID int
	// Get platform limit ID from hash
	err := tx.QueryRow(`
		SELECT id FROM campaign_platform_limits 
		WHERE hashed = $1`, campaignHash).Scan(&platformLimitID)
	if err != nil {
		return 0, err
	}

	// Increment platform usage counter atomically
	_, err = tx.Exec(`
		UPDATE campaign_platform_limits 
		SET used_count = used_count + 1 
		WHERE id = $1`, platformLimitID)
	if err != nil {
		return 0, err
	}

	// Create tracking record
	var trackingID int
	err = tx.QueryRow(`
		INSERT INTO campaign_tracking (
			platform_limit_id, ip_address, user_agent, 
			referrer_url, converted
		) VALUES ($1, $2, $3, $4, false)
		RETURNING id`,
		platformLimitID, userIP, userAgent, referrerURL).Scan(&trackingID)
	if err != nil {
		return 0, err
	}

	return trackingID, nil
}

func (r *SubscriptionRepo) processVoucher(tx *sql.Tx, campaignID int, trackingID int, pkg struct {
	name     string
	price    float64
	duration int
}, userId int, hashed string) (*struct {
	code           string
	discountAmount float64
}, error) {
	// Log bắt đầu quá trình xử lý voucher
	log.Printf("[INFO] Start processing voucher for CampaignID: %d, UserID: %d, Hashed: %s", campaignID, userId, hashed)

	// Get campaign details and decrement voucher count atomically
	var campaign struct {
		discountType  string
		discountValue float64
		maxAmount     sql.NullFloat64
	}
	err := tx.QueryRow(`
		UPDATE promotional_campaigns
		SET remaining_vouchers = remaining_vouchers - 1
		WHERE id = $1 
		AND remaining_vouchers > 0
		AND status = true 
		AND start_date <= CURRENT_TIMESTAMP
		AND end_date > CURRENT_TIMESTAMP
		RETURNING discount_type, discount_value, max_discount_amount`,
		campaignID).Scan(&campaign.discountType, &campaign.discountValue, &campaign.maxAmount)
	if err != nil {
		log.Printf("[ERROR] Failed to update promotional_campaigns for CampaignID: %d, Error: %v", campaignID, err)
		return nil, err
	}
	log.Printf("[INFO] Campaign updated successfully. DiscountType: %s, DiscountValue: %.2f, MaxAmount: %.2f",
		campaign.discountType, campaign.discountValue, campaign.maxAmount.Float64)

	// Update campaign platform limits
	result, err := tx.Exec(`
		UPDATE campaign_platform_limits
		SET used_count = used_count + 1
		WHERE hashed = $1
		AND (voucher_limit IS NULL OR used_count < voucher_limit);`,
		hashed)
	if err != nil {
		log.Printf("[ERROR] Failed to update campaign_platform_limits for Hashed: %s, Error: %v", hashed, err)
		return nil, err
	}

	// Kiểm tra số hàng bị ảnh hưởng
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Failed to get affected rows for campaign_platform_limits. Hashed: %s, Error: %v", hashed, err)
		return nil, err
	}
	if rowsAffected == 0 {
		log.Printf("[ERROR] No rows updated in campaign_platform_limits. Hashed: %s", hashed)
		return nil, errors.New("no platform limit updated or exceeded limit")
	}

	log.Printf("[INFO] Platform limits updated successfully for Hashed: %s, RowsAffected: %d", hashed, rowsAffected)
	// Calculate discount
	var discountAmount float64
	if campaign.discountType == "percentage" {
		discountAmount = pkg.price * campaign.discountValue / 100
		if campaign.maxAmount.Valid && discountAmount > campaign.maxAmount.Float64 {
			discountAmount = campaign.maxAmount.Float64
		}
	} else {
		discountAmount = campaign.discountValue
	}
	log.Printf("[INFO] Discount calculated. DiscountAmount: %.2f, PackagePrice: %.2f", discountAmount, pkg.price)

	// Create voucher
	voucherCode := utils.GenerateMockPaymentToken()
	log.Printf("[INFO] Generated Voucher Code: %s", voucherCode)

	_, err = tx.Exec(`
		INSERT INTO vouchers (
			campaign_id, tracking_id, code, discount_amount, 
			discount_type, max_discount_amount, valid_from, 
			valid_until, status, user_id
		) 
		SELECT 
			$1, $2, $3, $4, pc.discount_type, pc.max_discount_amount, 
			pc.start_date, pc.end_date, 'active', $5
		FROM promotional_campaigns pc
		WHERE pc.id = $1
		`,
		campaignID, trackingID, voucherCode,
		discountAmount, userId)
	if err != nil {
		log.Printf("[ERROR] Failed to create voucher for CampaignID: %d, UserID: %d, Error: %v", campaignID, userId, err)
		return nil, err
	}
	log.Printf("[INFO] Voucher created successfully. CampaignID: %d, VoucherCode: %s", campaignID, voucherCode)

	return &struct {
		code           string
		discountAmount float64
	}{
		code:           voucherCode,
		discountAmount: discountAmount,
	}, nil
}

func (r *SubscriptionRepo) CreateSubscription(userID int, req *model.SubscriptionRequest, userIP, userAgent string) (*model.SubscriptionResponse, error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[ERROR] Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	referrerURL := req.CampaignHash

	// Check if user already has an active subscription
	var existingCount int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM subscriptions 
		WHERE user_id = $1 AND status = 'active'`, userID).Scan(&existingCount)
	if err != nil {
		log.Printf("[ERROR] Failed to check active subscription for user %d: %v", userID, err)
		return nil, err
	}
	if existingCount > 0 {
		log.Printf("[INFO] User %d already has an active subscription", userID)
		return nil, errors.New("user already has an active subscription")
	}

	// Create tracking record
	trackingID, err := r.createTracking(tx, req.CampaignHash, userIP, userAgent, referrerURL)
	if err != nil {
		log.Printf("[ERROR] Failed to create tracking record for user %d: %v", userID, err)
		return nil, err
	}

	// Get package details
	var pkg struct {
		name     string
		price    float64
		duration int
	}
	err = tx.QueryRow(`
		SELECT name, price, duration_months 
		FROM packages 
		WHERE id = $1 AND status = true`, req.PackageID).Scan(&pkg.name, &pkg.price, &pkg.duration)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch package details for package ID %d: %v", req.PackageID, err)
		return nil, err
	}

	// Get campaign ID and process voucher
	var campaignID int
	err = tx.QueryRow(`
		SELECT pc.id
		FROM campaign_platform_limits cpl
		JOIN promotional_campaigns pc ON cpl.campaign_id = pc.id
		WHERE cpl.hashed = $1`,
		req.CampaignHash).Scan(&campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign ID for hash %s: %v", req.CampaignHash, err)
		return nil, err
	}

	voucher, err := r.processVoucher(tx, campaignID, trackingID, pkg, userID, req.CampaignHash)
	if err != nil {
		log.Printf("[ERROR] Failed to process voucher for campaign ID %d: %v", campaignID, err)
		return nil, err
	}

	// Create subscription
	startDate := time.Now()
	endDate := startDate.AddDate(0, pkg.duration, 0)
	var subID int
	err = tx.QueryRow(`
		INSERT INTO subscriptions (
			user_id, package_id, voucher_id, start_date, 
			end_date, original_price, discount_amount, 
			final_price, status
		) VALUES (
			$1, $2, (SELECT id FROM vouchers WHERE code = $3), 
			$4, $5, $6, $7, $8, 'pending'
		) RETURNING id`,
		userID, req.PackageID, voucher.code, startDate, endDate,
		pkg.price, voucher.discountAmount, pkg.price-voucher.discountAmount).Scan(&subID)
	if err != nil {
		log.Printf("[ERROR] Failed to create subscription for user %d: %v", userID, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[ERROR] Failed to commit transaction for user %d: %v", userID, err)
		return nil, err
	}

	log.Printf("[INFO] Subscription created successfully for user %d with subscription ID %d", userID, subID)

	return &model.SubscriptionResponse{
		SubscriptionID: subID,
		PackageName:    pkg.name,
		OriginalPrice:  pkg.price,
		DiscountAmount: voucher.discountAmount,
		FinalPrice:     pkg.price - voucher.discountAmount,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         "pending",
		VoucherCode:    voucher.code,
	}, nil
}

func (r *SubscriptionRepo) CreateUserWithSubscription(req *model.NewUserSubscriptionRequest, userIP, userAgent string) (*model.SubscriptionResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// Check if email exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		req.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var userID int
	err = tx.QueryRow(`
		INSERT INTO users (role_id, email, password, full_name, status)
		VALUES (2, $1, $2, $3, true)
		RETURNING id`,
		req.Email, hashedPassword, req.FullName).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Create subscription using existing method
	subReq := &model.SubscriptionRequest{
		PackageID:    req.PackageID,
		CampaignHash: req.CampaignHash,
	}

	tx.Commit()

	// Create subscription in a new transaction
	subRes, err := r.CreateSubscription(userID, subReq, userIP, userAgent)
	if err != nil {
		return nil, err
	}

	return subRes, nil
}

func (r *SubscriptionRepo) HandlePaymentCallback(callback *model.PaymentCallback) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get subscription details
	var subStatus string
	var voucherID sql.NullInt64
	var trackingID sql.NullInt64
	err = tx.QueryRow(`
		SELECT s.status, s.voucher_id, v.tracking_id
		FROM subscriptions s
		LEFT JOIN vouchers v ON s.voucher_id = v.id
		WHERE s.id = $1`,
		callback.SubscriptionID).Scan(&subStatus, &voucherID, &trackingID)
	if err != nil {
		return err
	}

	if subStatus != "pending" {
		return errors.New("invalid subscription status")
	}

	if callback.Status == "completed" {
		// Update subscription status
		_, err = tx.Exec(`
			UPDATE subscriptions 
			SET status = 'active' 
			WHERE id = $1`, callback.SubscriptionID)
		if err != nil {
			return err
		}

		// If voucher exists, mark it as used and update tracking
		if voucherID.Valid {
			_, err = tx.Exec(`
				UPDATE vouchers 
				SET status = 'used', used_at = CURRENT_TIMESTAMP 
				WHERE id = $1`, voucherID.Int64)
			if err != nil {
				return err
			}

			if trackingID.Valid {
				_, err = tx.Exec(`
					UPDATE campaign_tracking 
					SET converted = true 
					WHERE id = $1`, trackingID.Int64)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Mark subscription as cancelled
		_, err = tx.Exec(`
			UPDATE subscriptions 
			SET status = 'cancelled' 
			WHERE id = $1`, callback.SubscriptionID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
