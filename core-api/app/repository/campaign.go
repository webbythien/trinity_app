package repository

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type CampaignRepo struct {
	db *database.DB
}

func NewCampaignRepo(db *database.DB) CampaignRepository {
	return &CampaignRepo{db}
}

func (r *CampaignRepo) GetActiveCampaigns(filter *model.CampaignFilter) (*model.PaginatedCampaignResponse, error) {
	// First, build the count query
	countQuery := `
        SELECT COUNT(DISTINCT pc.id)
        FROM promotional_campaigns pc
        WHERE pc.status = true
        AND pc.start_date <= CURRENT_TIMESTAMP
        AND pc.end_date > CURRENT_TIMESTAMP
        AND pc.remaining_vouchers > 0`

	// Initialize params slice with some capacity
	params := make([]interface{}, 0, 8)
	paramCount := 1

	// Add filters to count query
	if filter.DiscountType != nil {
		countQuery += fmt.Sprintf(" AND pc.discount_type = $%d", paramCount)
		params = append(params, *filter.DiscountType)
		paramCount++
	}

	if filter.UserType != nil {
		countQuery += fmt.Sprintf(" AND pc.user_type = $%d", paramCount)
		params = append(params, *filter.UserType)
		paramCount++
	}

	if filter.MinDiscount != nil {
		countQuery += fmt.Sprintf(" AND pc.discount_value >= $%d", paramCount)
		params = append(params, *filter.MinDiscount)
		paramCount++
	}

	if filter.MaxDiscount != nil {
		countQuery += fmt.Sprintf(" AND pc.discount_value <= $%d", paramCount)
		params = append(params, *filter.MaxDiscount)
		paramCount++
	}

	if filter.EntityType != nil {
		countQuery += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM campaign_entities WHERE campaign_id = pc.id AND entity_type = $%d)", paramCount)
		params = append(params, *filter.EntityType)
		paramCount++
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, params...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	// Now build the main query
	baseQuery := `
        WITH campaign_entities AS (
            SELECT 
                ce.campaign_id,
                ce.entity_type,
                ce.entity_id,
                CASE 
                    WHEN ce.entity_type = 'package' THEN p.name
                    -- WHEN ce.entity_type = 'tour' THEN t.name
                    -- WHEN ce.entity_type = 'merchandise' THEN m.name
                END as entity_name,
                CASE 
                    WHEN ce.entity_type = 'package' THEN p.price
                    -- WHEN ce.entity_type = 'tour' THEN t.price
                    -- WHEN ce.entity_type = 'merchandise' THEN m.price
                END as entity_price
            FROM campaign_entities ce
            LEFT JOIN packages p ON ce.entity_type = 'package' AND ce.entity_id = p.id
            -- LEFT JOIN tours t ON ce.entity_type = 'tour' AND ce.entity_id = t.id
            -- LEFT JOIN merchandise m ON ce.entity_type = 'merchandise' AND ce.entity_id = m.id
        )
        SELECT 
            pc.id, pc.name, pc.description, pc.start_date, pc.end_date,
            pc.discount_type, pc.discount_value, pc.max_discount_amount,
            pc.user_type, pc.max_vouchers, pc.remaining_vouchers,
            COALESCE(json_agg(
                json_build_object(
                    'entity_type', ce.entity_type,
                    'entity_id', ce.entity_id,
                    'entity_name', ce.entity_name,
                    'entity_price', ce.entity_price
                )
            ) FILTER (WHERE ce.entity_type IS NOT NULL), '[]') as entities
        FROM promotional_campaigns pc
        LEFT JOIN campaign_entities ce ON pc.id = ce.campaign_id
        WHERE pc.status = true
        AND pc.start_date <= CURRENT_TIMESTAMP
        AND pc.end_date > CURRENT_TIMESTAMP
        AND pc.remaining_vouchers > 0`

	// Reset params and paramCount for main query
	params = make([]interface{}, 0, 8)
	paramCount = 1

	// Add filters to main query (same as before)
	if filter.DiscountType != nil {
		baseQuery += fmt.Sprintf(" AND pc.discount_type = $%d", paramCount)
		params = append(params, *filter.DiscountType)
		paramCount++
	}
	// ... (add other filters similarly)

	// Add GROUP BY
	baseQuery += " GROUP BY pc.id"

	// Validate and sanitize sort field
	allowedSortFields := map[string]bool{
		"start_date":         true,
		"end_date":           true,
		"discount_value":     true,
		"remaining_vouchers": true,
		"created_at":         true,
	}

	sortField := "start_date"
	if filter.Sort != "" {
		if allowed := allowedSortFields[filter.Sort]; allowed {
			sortField = filter.Sort
		}
	}

	sortDirection := "DESC"
	if filter.SortDirection != "" && (strings.ToUpper(filter.SortDirection) == "ASC" || strings.ToUpper(filter.SortDirection) == "DESC") {
		sortDirection = strings.ToUpper(filter.SortDirection)
	}

	baseQuery += fmt.Sprintf(" ORDER BY pc.%s %s", sortField, sortDirection)
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, filter.Limit, (filter.Page-1)*filter.Limit)

	// Execute main query
	rows, err := r.db.Query(baseQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []model.CampaignResponse
	for rows.Next() {
		var campaign model.CampaignResponse
		var entitiesJSON []byte

		err := rows.Scan(
			&campaign.ID, &campaign.Name, &campaign.Description,
			&campaign.StartDate, &campaign.EndDate, &campaign.DiscountType,
			&campaign.DiscountValue, &campaign.MaxDiscountAmount,
			&campaign.UserType, &campaign.MaxVouchers, &campaign.RemainingVouchers,
			&entitiesJSON,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(entitiesJSON, &campaign.Entities)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, campaign)
	}

	return &model.PaginatedCampaignResponse{
		Data:        campaigns,
		Total:       total,
		CurrentPage: filter.Page,
		PerPage:     filter.Limit,
		TotalPages:  totalPages,
	}, nil
}

func (r *CampaignRepo) CreateOrUpdateCampaign(campaign *model.CreateCampaignRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var campaignID int
	if campaign.ID > 0 {
		// Check if campaign exists
		var exists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM promotional_campaigns WHERE id = $1)", campaign.ID).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("campaign with ID %d not found", campaign.ID)
		}

		// Update existing campaign
		_, err = tx.Exec(`
            UPDATE promotional_campaigns SET
                name = $1, description = $2, start_date = $3, end_date = $4,
                discount_type = $5, discount_value = $6, max_discount_amount = $7,
                user_type = $8, max_vouchers = $9, updated_at = CURRENT_TIMESTAMP
            WHERE id = $10`,
			campaign.Name, campaign.Description, campaign.StartDate, campaign.EndDate,
			campaign.DiscountType, campaign.DiscountValue, campaign.MaxDiscountAmount,
			campaign.UserType, campaign.MaxVouchers, campaign.ID,
		)
		if err != nil {
			return err
		}
		campaignID = campaign.ID

		// Delete existing entities and platform limits
		_, err = tx.Exec("DELETE FROM campaign_entities WHERE campaign_id = $1", campaignID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM campaign_platform_limits WHERE campaign_id = $1", campaignID)
		if err != nil {
			return err
		}
	} else {
		// Insert new campaign
		err = tx.QueryRow(`
            INSERT INTO promotional_campaigns (
                name, description, start_date, end_date,
                discount_type, discount_value, max_discount_amount,
                user_type, max_vouchers, remaining_vouchers, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9, true)
            RETURNING id`,
			campaign.Name, campaign.Description, campaign.StartDate, campaign.EndDate,
			campaign.DiscountType, campaign.DiscountValue, campaign.MaxDiscountAmount,
			campaign.UserType, campaign.MaxVouchers,
		).Scan(&campaignID)
		if err != nil {
			return err
		}
	}

	// Insert campaign entities
	for _, entity := range campaign.Entities {
		_, err = tx.Exec(`
            INSERT INTO campaign_entities (
                campaign_id, entity_type, entity_id
            ) VALUES ($1, $2, $3)`,
			campaignID, entity.EntityType, entity.EntityID,
		)
		if err != nil {
			return err
		}
	}

	// Insert platform limits
	for _, platform := range campaign.PlatformLimits {
		hash := utils.GenerateHash(campaignID, platform.PlatformID)
		_, err = tx.Exec(`
            INSERT INTO campaign_platform_limits (
                campaign_id, platform_id, voucher_limit, hashed
            ) VALUES ($1, $2, $3, $4)`,
			campaignID, platform.PlatformID, platform.VoucherLimit, hash,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CampaignRepo) GetUserTypes() []model.CampaignUserType {
	return []model.CampaignUserType{
		{Value: "internal", Description: "For internal users only"},
		{Value: "external", Description: "For external users only"},
		{Value: "both", Description: "For both internal and external users"},
	}
}
