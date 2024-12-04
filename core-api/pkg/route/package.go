package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
)

// PublicRoutes func for describe group of public route.
func PackageRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/packages")
	route.Get("", middleware.JWTOptional(), controller.GetPackages)
}
