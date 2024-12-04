package repository

import (
	"fmt"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type CampaignPlatformLimitRepo struct {
	db *database.DB
}

func NewCampaignPlatformLimitRepo(db *database.DB) CampaignPlatformLimitRepository {
	return &CampaignPlatformLimitRepo{db}
}

func (r *CampaignPlatformLimitRepo) GetByCampaignID(campaignID int64) ([]model.CampaignPlatformLimitOutput, error) {
	query := `
        SELECT 
            cpl.id,
            cpl.campaign_id,
            cpl.platform_id,
            p.name as platform_name,
            cpl.voucher_limit,
            cpl.used_count,
            cpl.hashed,
            cpl.created_at,
            cpl.updated_at
        FROM campaign_platform_limits cpl
        JOIN platforms p ON cpl.platform_id = p.id
        WHERE cpl.campaign_id = $1
        ORDER BY p.name ASC`

	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []model.CampaignPlatformLimitOutput
	for rows.Next() {
		var limit model.CampaignPlatformLimitOutput
		err := rows.Scan(
			&limit.ID,
			&limit.CampaignID,
			&limit.PlatformID,
			&limit.PlatformName,
			&limit.VoucherLimit,
			&limit.UsedCount,
			&limit.Hashed,
			&limit.CreatedAt,
			&limit.UpdatedAt,
		)
		limit.URL = fmt.Sprintf("https://trinity.app/?platform=%d?campaign=%d?hashed=%s", limit.PlatformID, limit.CampaignID, limit.Hashed)
		if err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If no records found, return an empty slice with a specific error
	if len(limits) == 0 {
		return nil, constants.ErrRecordNotFound
	}

	return limits, nil
}
