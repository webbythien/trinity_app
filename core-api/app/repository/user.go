package repository

import (
	"database/sql"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}

func (r *UserRepo) GetUserInfo(userID int) (*model.UserSubscriptionInfo, error) {
	query := `
        WITH active_subscription AS (
            SELECT 
                s.*, 
                p.name as package_name,
                p.package_type
            FROM subscriptions s
            JOIN packages p ON s.package_id = p.id
            WHERE s.user_id = $1 
            AND s.status = 'active'
            AND s.end_date > CURRENT_TIMESTAMP
            ORDER BY s.end_date DESC
            LIMIT 1
        )
        SELECT 
            u.id, u.email, u.full_name, u.status, u.role_id,
            u.created_at, u.updated_at,
            CASE WHEN as2.id IS NOT NULL THEN true ELSE false END as is_subscribed,
            as2.package_id, as2.package_name, as2.package_type,
            as2.start_date, as2.end_date,
            as2.original_price, as2.discount_amount, as2.final_price,
            as2.status as subscription_status
        FROM users u
        LEFT JOIN active_subscription as2 ON true
        WHERE u.id = $1`

	var user model.UserSubscriptionInfo
	var sub model.SubscriptionInfo

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Status, &user.RoleID,
		&user.CreatedAt, &user.UpdatedAt, &user.IsSubscribed,
		&sub.PackageID, &sub.PackageName, &sub.PackageType,
		&sub.StartDate, &sub.EndDate,
		&sub.OriginalPrice, &sub.DiscountAmount, &sub.FinalPrice,
		&sub.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	if user.IsSubscribed {
		user.Subscription = &sub
	}

	return &user, nil
}
