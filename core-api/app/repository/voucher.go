package repository

import (
	"database/sql"
	"log"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type VoucherRepo struct {
	db *database.DB
}

func NewVoucherRepo(db *database.DB) VoucherRepository {
	return &VoucherRepo{db}
}

func (r *VoucherRepo) GetAllVouchers(status string) ([]model.Voucher, error) {
	query := `
		SELECT 
			v.id, v.campaign_id, v.user_id, v.tracking_id, v.code,
			v.discount_amount, v.discount_type, v.max_discount_amount,
			v.valid_from, v.valid_until, v.used_at, v.status,
			v.created_at, v.updated_at,
			c.name as campaign_name, c.description as campaign_description
		FROM vouchers v
		LEFT JOIN promotional_campaigns c ON v.campaign_id = c.id
		WHERE 1=1`

	if status != "" {
		query += ` AND v.status = $1`
		rows, err := r.db.Query(query, status)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanVouchers(rows)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVouchers(rows)
}

func (r *VoucherRepo) GetVouchersByUserId(userId int, status string) ([]model.Voucher, error) {
	query := `
		SELECT 
			v.id, v.campaign_id, v.user_id, v.tracking_id, v.code,
			v.discount_amount, v.discount_type, v.max_discount_amount,
			v.valid_from, v.valid_until, v.used_at, v.status,
			v.created_at, v.updated_at,
			c.name as campaign_name, c.description as campaign_description
		FROM vouchers v
		LEFT JOIN promotional_campaigns c ON v.campaign_id = c.id
		WHERE v.user_id = $1`

	if status != "" {
		query += ` AND v.status = $2`
		rows, err := r.db.Query(query, userId, status)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanVouchers(rows)
	}

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVouchers(rows)
}

func scanVouchers(rows *sql.Rows) ([]model.Voucher, error) {
	var vouchers []model.Voucher
	for rows.Next() {
		var v model.Voucher
		var campaignName, campaignDescription string

		err := rows.Scan(
			&v.ID, &v.CampaignID, &v.UserID, &v.TrackingID, &v.Code,
			&v.DiscountAmount, &v.DiscountType, &v.MaxDiscountAmount,
			&v.ValidFrom, &v.ValidUntil, &v.UsedAt, &v.Status,
			&v.CreatedAt, &v.UpdatedAt,
			&campaignName, &campaignDescription,
		)
		if err != nil {
			return nil, err
		}

		// v.Campaign = &model.Campaign{
		// 	ID:          v.CampaignID,
		// 	Name:        campaignName,
		// 	Description: campaignDescription,
		// }

		vouchers = append(vouchers, v)
	}
	return vouchers, nil
}

func (r *VoucherRepo) GetGuestVoucher(campaignHash, userIP, userAgent string) (*model.GuestVoucherResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	var existingVoucher model.GuestVoucherResponse
	err = tx.QueryRow(`
        SELECT v.code, v.discount_amount, v.discount_type, 
               v.valid_from, v.valid_until
        FROM vouchers v
        JOIN campaign_tracking ct ON v.tracking_id = ct.id
        WHERE ct.ip_address = $1 and v.user_id is null
    `, userIP).Scan(
		&existingVoucher.Code,
		&existingVoucher.DiscountAmount,
		&existingVoucher.DiscountType,
		&existingVoucher.ValidFrom,
		&existingVoucher.ValidUntil,
	)
	if err == nil {
		return &existingVoucher, nil
	} else {
		log.Printf("No existing voucher found for IP: %s, CampaignHash: %s. Error: %v", userIP, campaignHash, err)
	}

	var campaignID int
	var remainingVouchers int
	var platformLimitID int
	err = tx.QueryRow(`
        UPDATE promotional_campaigns pc
        SET remaining_vouchers = remaining_vouchers - 1
        FROM campaign_platform_limits cpl
        WHERE cpl.campaign_id = pc.id
        AND cpl.hashed = $1
        AND pc.remaining_vouchers > 0
        AND pc.status = true
        AND pc.start_date <= CURRENT_TIMESTAMP
        AND pc.end_date > CURRENT_TIMESTAMP
        AND (cpl.voucher_limit IS NULL OR cpl.used_count < cpl.voucher_limit)
        RETURNING pc.id, pc.remaining_vouchers, cpl.id
    `, campaignHash).Scan(&campaignID, &remainingVouchers, &platformLimitID)
	if err != nil {
		log.Printf("Failed to update campaign promotional limits for CampaignHash: %s. Error: %v", campaignHash, err)
		return nil, err
	}

	_, err = tx.Exec(`
        UPDATE campaign_platform_limits
        SET used_count = used_count + 1
        WHERE hashed = $1
    `, campaignHash)
	if err != nil {
		log.Printf("Failed to increment used_count for CampaignHash: %s. Error: %v", campaignHash, err)
		return nil, err
	}

	var trackingID int
	err = tx.QueryRow(`
        INSERT INTO campaign_tracking (
             platform_limit_id , ip_address, user_agent
        ) VALUES ($1, $2, $3)
        RETURNING id
    `, platformLimitID, userIP, userAgent).Scan(&trackingID)
	if err != nil {
		log.Printf("Failed to insert tracking record for CampaignHash: %s, IP: %s. Error: %v", campaignHash, userIP, err)
		return nil, err
	}

	voucherCode := utils.GenerateMockPaymentToken()
	var voucher model.GuestVoucherResponse
	err = tx.QueryRow(`
        INSERT INTO vouchers (
            campaign_id, tracking_id, code, discount_amount,
            discount_type, max_discount_amount, valid_from,
            valid_until, status
        )
        SELECT 
            $1, $2, $3, pc.discount_value,
            pc.discount_type, pc.max_discount_amount,
            pc.start_date, pc.end_date, 'active'
        FROM promotional_campaigns pc
        WHERE pc.id = $1
        RETURNING code, discount_amount, discount_type,
                valid_from, valid_until
    `, campaignID, trackingID, voucherCode).Scan(
		&voucher.Code,
		&voucher.DiscountAmount,
		&voucher.DiscountType,
		&voucher.ValidFrom,
		&voucher.ValidUntil,
	)
	if err != nil {
		log.Printf("Failed to create voucher for CampaignID: %d, TrackingID: %d. Error: %v", campaignID, trackingID, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction for CampaignHash: %s. Error: %v", campaignHash, err)
		return nil, err
	}

	log.Printf("Voucher created successfully for CampaignHash: %s, VoucherCode: %s", campaignHash, voucher.Code)
	return &voucher, nil
}
