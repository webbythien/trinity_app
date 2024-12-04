package model

import "time"

type CampaignFilter struct {
	DiscountType  *string  `json:"discount_type,omitempty" query:"discount_type"`
	UserType      *string  `json:"user_type,omitempty" query:"user_type"`
	MinDiscount   *float64 `json:"min_discount,omitempty" query:"min_discount"`
	MaxDiscount   *float64 `json:"max_discount,omitempty" query:"max_discount"`
	EntityType    *string  `json:"entity_type,omitempty" query:"entity_type"`
	Sort          string   `json:"sort,omitempty" query:"sort"`
	SortDirection string   `json:"sort_direction,omitempty" query:"sort_direction"`
	Page          int      `json:"page,omitempty" query:"page" validate:"min=1"`
	Limit         int      `json:"limit,omitempty" query:"limit" validate:"min=1"`
}
type CampaignResponse struct {
	ID                int              `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	StartDate         time.Time        `json:"start_date"`
	EndDate           time.Time        `json:"end_date"`
	DiscountType      string           `json:"discount_type"`
	DiscountValue     float64          `json:"discount_value"`
	MaxDiscountAmount *float64         `json:"max_discount_amount,omitempty"`
	UserType          string           `json:"user_type"`
	MaxVouchers       int              `json:"max_vouchers"`
	RemainingVouchers int              `json:"remaining_vouchers"`
	Entities          []CampaignEntity `json:"entities"`
}

type CampaignEntity struct {
	EntityType  string   `json:"entity_type"`
	EntityID    int      `json:"entity_id"`
	EntityName  string   `json:"entity_name"`
	EntityPrice *float64 `json:"entity_price,omitempty"`
}

type PaginatedCampaignResponse struct {
	Data        []CampaignResponse `json:"data"`
	Total       int                `json:"total"`
	CurrentPage int                `json:"current_page"`
	PerPage     int                `json:"per_page"`
	TotalPages  int                `json:"total_pages"`
}

type CampaignPlatformLimitOutput struct {
	ID           int64     `json:"id"`
	CampaignID   int64     `json:"campaign_id"`
	PlatformID   int64     `json:"platform_id"`
	URL          string    `json:"url"`
	PlatformName string    `json:"platform_name"`
	VoucherLimit *int      `json:"voucher_limit"`
	UsedCount    int       `json:"used_count"`
	Hashed       string    `json:"hashed"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EntityType struct {
	ID        int    `json:"id"`
	TypeName  string `json:"type_name"`
	TableName string `json:"table_name"`
	Status    bool   `json:"status"`
}

type CreateCampaignRequest struct {
	ID                int                       `json:"id,omitempty" example:"1"`
	Name              string                    `json:"name" example:"Black Friday" validate:"required"`
	Description       string                    `json:"description" example:"Huge discounts for Black Friday!"`
	StartDate         time.Time                 `json:"start_date" example:"2024-12-01T00:00:00Z" validate:"required"`
	EndDate           time.Time                 `json:"end_date" example:"2024-12-10T23:59:59Z" validate:"required"`
	DiscountType      string                    `json:"discount_type" example:"percentage" validate:"required,oneof=percentage fixed"`
	DiscountValue     float64                   `json:"discount_value" example:"15.5" validate:"required,gt=0"`
	MaxDiscountAmount *float64                  `json:"max_discount_amount" example:"100.0"`
	UserType          string                    `json:"user_type" example:"external" validate:"required,oneof=internal external both"`
	MaxVouchers       int                       `json:"max_vouchers" example:"1000" validate:"required,gt=0"`
	Entities          []Entity                  `json:"entities" validate:"required,min=1,dive"`
	PlatformLimits    []PlatformCampaginRequest `json:"platform_limits" validate:"required,min=1,dive"`
}

// Entity represents an entity associated with a campaign
type Entity struct {
	EntityType string `json:"entity_type" example:"product"`
	EntityID   int    `json:"entity_id" example:"123"`
}

// PlatformCampaignRequest represents platform-specific campaign limits
type PlatformCampaginRequest struct {
	PlatformID   int `json:"platform_id" example:"1"`
	VoucherLimit int `json:"voucher_limit" example:"100"`
}

type CampaignUserType struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Campaign struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	DiscountType      string    `json:"discount_type"`
	DiscountValue     float64   `json:"discount_value"`
	MaxDiscountAmount *float64  `json:"max_discount_amount"`
	UserType          string    `json:"user_type"`
	MaxVouchers       int       `json:"max_vouchers"`
	RemainingVouchers int       `json:"remaining_vouchers"`
	Status            bool      `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CampaignTracking struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}
