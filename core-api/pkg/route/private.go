package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
)

// PrivateRoutes func for describe group of private route.
func PrivateRoutes(a *fiber.App) {
	// Admin route group
	// adminRoute := a.Group("/api/v1/users", middleware.JWTProtected(), middleware.IsAdmin)
	// User

	// Book
	route := a.Group("/api/v1/books", middleware.JWTProtected())
	route.Post("/", controller.CreateBook)
	route.Put("/:id", controller.UpdateBook)
	route.Delete("/:id", controller.DeleteBook)

}
