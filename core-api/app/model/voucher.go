package model

import "time"

type Voucher struct {
	ID                int        `json:"id"`
	CampaignID        int        `json:"campaign_id"`
	UserID            *int       `json:"user_id"`
	TrackingID        *int       `json:"tracking_id"`
	Code              string     `json:"code"`
	DiscountAmount    float64    `json:"discount_amount"`
	DiscountType      string     `json:"discount_type"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	ValidFrom         time.Time  `json:"valid_from"`
	ValidUntil        time.Time  `json:"valid_until"`
	UsedAt            *time.Time `json:"used_at,omitempty"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Additional fields for related data
	// Campaign *Campaign `json:"campaign,omitempty"`
}

type VoucherInfo struct {
	ID             int64   `json:"id"`
	Code           string  `json:"code"`
	DiscountType   string  `json:"discount_type"`
	DiscountAmount float64 `json:"discount_amount"`
	ValidUntil     string  `json:"valid_until"`
}

type GuestVoucherResponse struct {
	Code           string  `json:"code"`
	DiscountAmount float64 `json:"discount_amount"`
	DiscountType   string  `json:"discount_type"`
	ValidFrom      string  `json:"valid_from"`
	ValidUntil     string  `json:"valid_until"`
}

type CreateGuestVoucher struct {
	CampaignHash string `json:"campaign_hash" validate:"required"`
}
