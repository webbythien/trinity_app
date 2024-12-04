package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// Controller
// GetPackages godoc
// @Summary Get all packages with active campaigns
// @Description Get a list of all active packages with their associated campaigns and user vouchers
// @Tags packages
// @Accept json
// @Produce json
// @Success 200 {array} model.PackageWithCampaignOutput
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /packages [get]
func GetPackages(c *fiber.Ctx) error {
	var userID *int
	user := c.Locals(constants.UserLocal).(middleware.UserJWTProtect)
	if user.IsExist {
		userID = &user.UserID
	}
	packageRepo := repository.NewPackageRepo(database.GetDB())
	packages, err := packageRepo.GetPackagesWithCampaigns(userID)
	return handleResultWithErrorCode(c, packages, err, "Failed to fetch packages", "Get all packages successfully")

}
