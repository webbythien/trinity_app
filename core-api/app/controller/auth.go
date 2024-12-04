package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	repo "github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with the provided details
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  data  body  model.RegisterInput  true  "User Registration Data"
// @Success 200 {object} model.AuthResponse "Registration successful"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	var req *model.RegisterInput

	// Parse the request body and handle errors
	if err := parseRequestBody(c, &req); err != nil {
		return err
	}

	// Validate the common fields (step and optional id)
	if err := validate.Struct(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}
	authRepo := repo.NewAuthRepo(database.GetDB())

	result, err := authRepo.Register(req)

	return handleResultWithErrorCode(c, result, err, "", "Register successfully")
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  data  body  model.LoginInput  true  "User Login Data"
// @Success 200 {object} model.AuthResponse "Login successful"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var req *model.LoginInput

	// Parse the request body and handle errors
	if err := parseRequestBody(c, &req); err != nil {
		return err
	}
	// Validate the common fields (step and optional id)
	if err := validate.Struct(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}
	authRepo := repo.NewAuthRepo(database.GetDB())

	result, err := authRepo.Login(req)

	return handleResultWithErrorCode(c, result, err, "", "Login successfully")
}
