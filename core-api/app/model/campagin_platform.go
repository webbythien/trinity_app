package model

type CampaignPlatformLimit struct {
	ID           int    `json:"id"`
	CampaignID   int    `json:"campaign_id"`
	Platform     string `json:"platform"`
	VoucherLimit *int   `json:"voucher_limit"`
	UsedCount    int    `json:"used_count"`
}
