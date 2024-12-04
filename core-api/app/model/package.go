package model

type PackageWithCampaignOutput struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	PackageType    string         `json:"package_type"`
	Price          float64        `json:"price"`
	Description    string         `json:"description"`
	DurationMonths int            `json:"duration_months"`
	Status         bool           `json:"status"`
	Campaigns      []CampaignInfo `json:"campaigns"`
	UserVouchers   []VoucherInfo  `json:"user_vouchers"`
}

type CampaignInfo struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	DiscountType         string   `json:"discount_type"`
	DiscountValue        float64  `json:"discount_value"`
	MaxDiscountAmount    *float64 `json:"max_discount_amount"`
	EndDate              string   `json:"end_date"`
	PricePackageDiscount float64  `json:"price_package_discount"`
}
