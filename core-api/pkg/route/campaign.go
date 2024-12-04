package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
)

// PublicRoutes func for describe group of public route.
func CampaignRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/campaigns")
	route.Get("/active", controller.GetActiveCampaigns)
	route.Post("", controller.CreateCampaign)
	route.Get("/:campaign_id/platform-limits", controller.GetCampaignPlatformLimits)
	route.Get("/entity-type", controller.GetAllEntityTypes)
	route.Get("/user-types", controller.GetCampaignUserTypes)
}
