package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
)

// PublicRoutes func for describe group of public route.
func SubscriptionRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/subscriptions")
	route.Post("", middleware.JWTOptional(), controller.CreateSubscription)
	route.Post("/register", controller.CreateUserWithSubscription)
	route.Post("/payment-callback", controller.HandlePaymentCallback)
}
