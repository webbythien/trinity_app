package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
)

// PublicRoutes func for describe group of public route.
func VoucherRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/vouchers")
	route.Get("", middleware.JWTProtected(), middleware.IsAdmin, controller.GetAllVouchers)
	route.Get("/users", middleware.JWTProtected(), controller.GetUserVouchers)
	route.Post("/guest", controller.CreateGuestVoucher)
}
