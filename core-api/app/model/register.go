package model

import "time"

type SubscriptionRequest struct {
	PackageID    int    `json:"package_id" validate:"required"`
	CampaignHash string `json:"campaign_hash" validate:"required"`
	// TrackingID   int    `json:"tracking_id" validate:"required"`
}

type NewUserSubscriptionRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	PackageID    int    `json:"package_id" validate:"required"`
	CampaignHash string `json:"campaign_hash" validate:"required"`
	// TrackingID   int    `json:"tracking_id" validate:"required"`
}

type SubscriptionResponse struct {
	SubscriptionID int       `json:"subscription_id"`
	PackageName    string    `json:"package_name"`
	OriginalPrice  float64   `json:"original_price"`
	DiscountAmount float64   `json:"discount_amount"`
	FinalPrice     float64   `json:"final_price"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	VoucherCode    string    `json:"voucher_code,omitempty"`
}

type PaymentCallback struct {
	SubscriptionID int `json:"subscription_id" validate:"required"`
	// PaymentID      string  `json:"payment_id" validate:"required"`
	Status string  `json:"status" validate:"required,oneof=completed failed"`
	Amount float64 `json:"amount" validate:"required"`
}
