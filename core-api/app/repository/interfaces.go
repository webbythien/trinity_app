package repository

import (
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
)

type UserRepository interface {
	GetUserInfo(userID int) (*model.UserSubscriptionInfo, error)
}

type BookRepository interface {
	Create(b *model.Book) error
	All(limit int, offset uint) ([]*model.Book, error)
	Get(ID uuid.UUID) (*model.Book, error)
	Update(ID uuid.UUID, b *model.Book) error
	Delete(ID uuid.UUID) error
}

type RegionRepository interface {
	GetAllProvince() ([]*model.ProvineUnitName, error)
	GetAllDistrict(provinceCode string) ([]*model.DistrictProvince, error)
	GetAllWard(districtCode string) ([]*model.WardDistrictProvince, error)
}

type AuthRepository interface {
	Register(input *model.RegisterInput) (*model.AuthResponse, error)
	Login(input *model.LoginInput) (*model.AuthResponse, error)
	GetUserByEmail(email string) (*model.AuthOutput, error)
}

type CampaignRepository interface {
	GetActiveCampaigns(filter *model.CampaignFilter) (*model.PaginatedCampaignResponse, error)
	CreateOrUpdateCampaign(campaign *model.CreateCampaignRequest) error
	GetUserTypes() []model.CampaignUserType
}

type CampaignPlatformLimitRepository interface {
	GetByCampaignID(campaignID int64) ([]model.CampaignPlatformLimitOutput, error)
}

type PlatformRepository interface {
	GetAllPlatforms() ([]model.Platform, error)
}

type SubscriptionRepository interface {
	CreateSubscription(userID int, req *model.SubscriptionRequest) (*model.SubscriptionResponse, error)
	CreateUserWithSubscription(req *model.NewUserSubscriptionRequest) (*model.SubscriptionResponse, error)
	HandlePaymentCallback(callback *model.PaymentCallback) error
}
type EntityRepository interface {
	GetAllEntityTypes() ([]model.EntityType, error)
}

type VoucherRepository interface {
	GetAllVouchers(status string) ([]model.Voucher, error)
	GetVouchersByUserId(userId int, status string) ([]model.Voucher, error)
	GetGuestVoucher(campaignHash, userIP, userAgent string) (*model.GuestVoucherResponse, error)
}
