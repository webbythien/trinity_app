package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

func GetUserInfo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	userRepo := repository.NewUserRepo(database.GetDB())
	userInfo, err := userRepo.GetUserInfo(userID)

	return handleResultWithErrorCode(c, userInfo, err, "", "User info retrieved successfully")
}
