package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
)

// PublicRoutes func for describe group of public route.
func PlatformRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/platforms")
	route.Get("", controller.GetAllPlatforms)
}
