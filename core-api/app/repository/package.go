package repository

import (
	"encoding/json"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type PackageRepo struct {
	db *database.DB
}

func NewPackageRepo(db *database.DB) *PackageRepo {
	return &PackageRepo{db}
}

func (r *PackageRepo) GetPackagesWithCampaigns(userID *int) ([]model.PackageWithCampaignOutput, error) {
	query := `
        WITH active_campaigns AS (
            SELECT 
                ce.entity_id as package_id,
                pc.id as campaign_id,
                pc.name as campaign_name,
                pc.discount_type,
                pc.discount_value,
                pc.max_discount_amount,
                pc.end_date
            FROM promotional_campaigns pc
            JOIN campaign_entities ce ON pc.id = ce.campaign_id
            WHERE pc.status = true 
                AND ce.entity_type = 'package'
                AND pc.end_date > CURRENT_TIMESTAMP
                AND pc.remaining_vouchers > 0
        ),
        user_vouchers AS (
            SELECT 
                ce.entity_id as package_id,
                v.id as voucher_id,
                v.code,
                v.discount_type,
                v.discount_amount,
                v.valid_until
            FROM vouchers v
            JOIN promotional_campaigns pc ON v.campaign_id = pc.id
            JOIN campaign_entities ce ON pc.id = ce.campaign_id
            WHERE ce.entity_type = 'package'
                AND v.status = 'active'
                AND v.valid_until > CURRENT_TIMESTAMP
                AND v.user_id = $1
        )
        SELECT 
            p.id,
            p.name,
            p.package_type,
            p.price,
            p.description,
            p.duration_months,
            p.status,
            COALESCE(
                json_agg(
                    DISTINCT jsonb_build_object(
                        'id', ac.campaign_id,
                        'name', ac.campaign_name,
                        'discount_type', ac.discount_type,
                        'discount_value', ac.discount_value,
                        'max_discount_amount', ac.max_discount_amount,
                        'end_date', ac.end_date,
                        'price_package_discount',
                        CASE 
                            WHEN ac.discount_type = 'percentage' THEN
                                CASE 
                                    WHEN ac.max_discount_amount IS NOT NULL AND ac.max_discount_amount > 0 THEN 
                                        LEAST(p.price * ac.discount_value / 100, ac.max_discount_amount)
                                    ELSE 
                                        p.price * ac.discount_value / 100
                                END
                            WHEN ac.discount_type = 'fixed' THEN
                                CASE 
                                    WHEN ac.max_discount_amount IS NOT NULL AND ac.max_discount_amount > 0 THEN 
                                        LEAST(ac.discount_value, ac.max_discount_amount)
                                    ELSE 
                                        ac.discount_value
                                END
                            ELSE 0
                        END
                    )
                ) FILTER (WHERE ac.campaign_id IS NOT NULL),
                '[]'
            ) as campaigns,
            COALESCE(
                json_agg(
                    DISTINCT jsonb_build_object(
                        'id', uv.voucher_id,
                        'code', uv.code,
                        'discount_type', uv.discount_type,
                        'discount_amount', uv.discount_amount,
                        'valid_until', uv.valid_until
                    )
                ) FILTER (WHERE uv.voucher_id IS NOT NULL),
                '[]'
            ) as user_vouchers
        FROM packages p
        LEFT JOIN active_campaigns ac ON p.id = ac.package_id
        LEFT JOIN user_vouchers uv ON p.id = uv.package_id
        WHERE p.status = true
        GROUP BY p.id, p.name, p.package_type, p.price, p.description, p.duration_months, p.status
        ORDER BY p.id`

	var packages []model.PackageWithCampaignOutput
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pkg model.PackageWithCampaignOutput
		var campaignsJSON, vouchersJSON []byte

		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.PackageType,
			&pkg.Price,
			&pkg.Description,
			&pkg.DurationMonths,
			&pkg.Status,
			&campaignsJSON,
			&vouchersJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse campaigns JSON
		if err := json.Unmarshal(campaignsJSON, &pkg.Campaigns); err != nil {
			return nil, err
		}
		// Parse vouchers JSON if user_id was provided
		if userID != nil {
			if err := json.Unmarshal(vouchersJSON, &pkg.UserVouchers); err != nil {
				return nil, err
			}
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}
