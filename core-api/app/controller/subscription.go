package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// CreateSubscription godoc
// @Summary Create a new subscription for existing user
// @Description Create a subscription with campaign tracking and voucher
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.SubscriptionRequest true "Subscription details"
// @Success 201 {object} model.SubscriptionResponse "Subscription created"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 401 {object} model.Response "Unauthorized"
// @Failure 422 {object} model.Response "Validation Error"
// @Security Bearer
// @Router /subscriptions [post]
func CreateSubscription(c *fiber.Ctx) error {
	req := new(model.SubscriptionRequest)
	if err := c.BodyParser(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	if err := validate.Struct(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "VALIDATION_ERROR")
	}
	user := c.Locals(constants.UserLocal).(middleware.UserJWTProtect)
	userID := user.UserID
	userIP := utils.GetIPFromContext(c)
	userAgent := utils.GetUserAgentFromContext(c)
	subRepo := repository.NewSubscriptionRepo(database.GetDB())
	result, err := subRepo.CreateSubscription(userID, req, userIP, userAgent)

	return handleResultWithErrorCode(c, result, err, "", "Subscription created successfully")
}

func CreateUserWithSubscription(c *fiber.Ctx) error {
	req := new(model.NewUserSubscriptionRequest)
	if err := c.BodyParser(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	if err := validate.Struct(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "VALIDATION_ERROR")
	}
	userIP := utils.GetIPFromContext(c)
	userAgent := utils.GetUserAgentFromContext(c)

	subRepo := repository.NewSubscriptionRepo(database.GetDB())
	result, err := subRepo.CreateUserWithSubscription(req, userIP, userAgent)

	return handleResultWithErrorCode(c, result, err, "", "User registered and subscription created successfully")
}

// HandlePaymentCallback godoc
// @Summary Handle payment callback
// @Description Update subscription and voucher status based on payment result
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param callback body model.PaymentCallback true "Payment callback details"
// @Success 200 {object} model.Response "Payment processed successfully"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 422 {object} model.Response "Validation Error"
// @Router /subscriptions/payment-callback [post]
func HandlePaymentCallback(c *fiber.Ctx) error {
	callback := new(model.PaymentCallback)
	if err := c.BodyParser(callback); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	if err := validate.Struct(callback); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "VALIDATION_ERROR")
	}

	subRepo := repository.NewSubscriptionRepo(database.GetDB())
	err := subRepo.HandlePaymentCallback(callback)

	return handleResultWithErrorCode(c, nil, err, "", "Payment processed successfully")
}
