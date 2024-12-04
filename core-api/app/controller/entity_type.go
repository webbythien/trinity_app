package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// GetAllEntityTypes godoc
// @Summary Get all entity types
// @Description Retrieve a list of all active entity types
// @Tags entities
// @Accept json
// @Produce json
// @Success 200 {array} model.EntityType
// @Failure 401 {object} model.Response "Unauthorized"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Security BearerAuth
// @Router /campaigns/entity-type [get]
func GetAllEntityTypes(c *fiber.Ctx) error {
	entityRepo := repository.NewEntityRepo(database.GetDB())
	entities, err := entityRepo.GetAllEntityTypes()

	return handleResultWithErrorCode(c, entities, err, "", "Entity types retrieved successfully")
}
