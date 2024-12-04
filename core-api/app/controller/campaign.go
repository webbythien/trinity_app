package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// GetActiveCampaigns godoc
// @Summary Retrieve active campaigns
// @Description Get a list of active promotional campaigns with optional filters
// @Tags campaigns
// @Accept json
// @Produce json
// @Param discount_type query string false "Discount Type" Enums(percentage, fixed)
// @Param user_type query string false "User Type" Enums(new, existing)
// @Param min_discount query number false "Minimum Discount Value"
// @Param max_discount query number false "Maximum Discount Value"
// @Param entity_type query string false "Entity Type" Enums(package, tour, merchandise)
// @Param sort query string false "Sort Field" Enums(start_date, end_date, discount_value)
// @Param sort_direction query string false "Sort Direction" Enums(asc, desc)
// @Param page query integer false "Page number" default(1)
// @Param limit query integer false "Number of records per page" default(10)
// @Success 200 {object} model.PaginatedCampaignResponse "active campaigns"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /campaigns/active [get]
func GetActiveCampaigns(c *fiber.Ctx) error {
	filter := new(model.CampaignFilter)
	if err := c.QueryParser(filter); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	// Set default values if not provided
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	// Validate filter
	if err := validate.Struct(filter); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	campaignRepo := repository.NewCampaignRepo(database.GetDB())
	paginatedResponse, err := campaignRepo.GetActiveCampaigns(filter)

	return handleResultWithErrorCode(c, paginatedResponse, err, "", "Active campaigns retrieved successfully")
}

// CreateCampaign godoc
// @Summary Create a new or update promotional campaign
// @Description Create a new promotional campaign with entities and platform limits
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body model.CreateCampaignRequest true "Campaign details"
// @Success 201 {object} model.Response "Campaign created successfully"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 401 {object} model.Response "Unauthorized"
// @Failure 422 {object} model.Response "Validation Error"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Security BearerAuth
// @Router /campaigns [post]
func CreateCampaign(c *fiber.Ctx) error {
	campaign := new(model.CreateCampaignRequest)
	if err := c.BodyParser(campaign); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	// Validate request
	if err := validate.Struct(campaign); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "VALIDATION_ERROR")
	}

	// Additional validation
	if campaign.EndDate.Before(campaign.StartDate) {
		return returnBadRequestWithErrorCode(c, "End date must be after start date", "VALIDATION_ERROR")
	}
	msg := "Create campaigns successfully"
	if campaign.ID != 0 {
		msg = "Update campaigns successfully"

	}
	campaignRepo := repository.NewCampaignRepo(database.GetDB())
	err := campaignRepo.CreateOrUpdateCampaign(campaign)
	return handleResultWithErrorCode(c, nil, err, "", msg)

}

// @Summary Get campaign user types
// @Description Get all available campaign user types
// @Tags campaigns
// @Produce json
// @Success 200 {array} model.CampaignUserType "Campaign user types"
// @Router /campaigns/user-types [get]
func GetCampaignUserTypes(c *fiber.Ctx) error {
	campaignRepo := repository.NewCampaignRepo(database.GetDB())
	userTypes := campaignRepo.GetUserTypes()
	return handleResultWithErrorCode(c, userTypes, nil, "", "Get campaigns user type successfully")

}
