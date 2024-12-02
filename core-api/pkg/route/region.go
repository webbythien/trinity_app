package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

func RegionRoutes(a *fiber.App) {
	route := a.Group("/region")
	regionRepo := repository.NewRegionRepo(database.GetDB())

	regionController := controller.NewRegionController(regionRepo)

	route.Get("/provinces", regionController.GetAllProvinces)
	route.Get("/districts/:provinceCode", regionController.GetAllDistricts)
	route.Get("/wards/:districtCode", regionController.GetAllWards)
}
