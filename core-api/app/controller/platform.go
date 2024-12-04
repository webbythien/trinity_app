package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// GetAllPlatforms godoc
// @Summary Get all platforms
// @Description Retrieve a list of all platforms
// @Tags platforms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Platform
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /platforms [get]
func GetAllPlatforms(c *fiber.Ctx) error {
	platformRepo := repository.NewPlatformRepo(database.GetDB())
	platforms, err := platformRepo.GetAllPlatforms()

	return handleResultWithErrorCode(c, platforms, err, "", "Platforms retrieved successfully")
}
