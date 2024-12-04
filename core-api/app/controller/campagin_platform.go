package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// Controller
// GetCampaignPlatformLimits godoc
// @Summary Get campaign platform limits
// @Description Get all platform limits for a specific campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign_id path integer true "Campaign ID"
// @Success 200 {array} model.CampaignPlatformLimitOutput
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 404 {object} model.Response "Not Found"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Security BearerAuth
// @Router /campaigns/{campaign_id}/platform-limits [get]
func GetCampaignPlatformLimits(c *fiber.Ctx) error {
	campaignID, err := c.ParamsInt("campaign_id")
	if err != nil {
		return returnBadRequestWithErrorCode(c, "Invalid campaign_id parameter", "BAD_REQUEST")
	}

	repo := repository.NewCampaignPlatformLimitRepo(database.GetDB())
	limits, err := repo.GetByCampaignID(int64(campaignID))

	return handleResultWithErrorCode(c, limits, err, "Failed to fetch campaign platform limits", "Get platform successfully")
}
