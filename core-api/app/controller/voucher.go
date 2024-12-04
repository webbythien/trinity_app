package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/middleware"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// GetAllVouchers godoc
// @Summary Get all vouchers
// @Description Retrieve a list of all vouchers with optional filters
// @Tags vouchers
// @Accept json
// @Produce json
// @Param status query string false "Filter by voucher status (active/used/expired)"
// @Success 200 {array} model.Voucher
// @Failure 401 {object} model.Response "Unauthorized"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Security Bearer
// @Router /vouchers [get]
func GetAllVouchers(c *fiber.Ctx) error {
	status := c.Query("status")
	voucherRepo := repository.NewVoucherRepo(database.GetDB())
	vouchers, err := voucherRepo.GetAllVouchers(status)

	return handleResultWithErrorCode(c, vouchers, err, "", "Vouchers retrieved successfully")
}

// GetUserVouchers godoc
// @Summary Get user's vouchers
// @Description Retrieve all vouchers for a specific user
// @Tags vouchers
// @Accept json
// @Produce json
// @Param status query string false "Filter by voucher status (active/used/expired)"
// @Success 200 {array} model.Voucher
// @Failure 401 {object} model.Response "Unauthorized"
// @Failure 404 {object} model.Response "User not found"
// @Failure 500 {object} model.Response "Internal Server Error"
// @Security Bearer
// @Router /vouchers/users [get]
func GetUserVouchers(c *fiber.Ctx) error {

	user := c.Locals(constants.UserLocal).(middleware.UserJWTProtect)

	status := c.Query("status")
	voucherRepo := repository.NewVoucherRepo(database.GetDB())
	vouchers, err := voucherRepo.GetVouchersByUserId(user.UserID, status)

	return handleResultWithErrorCode(c, vouchers, err, "", "User vouchers retrieved successfully")
}

// CreateGuestVoucher godoc
// @Summary Create a voucher for guest users
// @Description Allows non-registered users to claim a voucher based on their IP and User-Agent information.
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param body body model.CreateGuestVoucher true "Request body to create a guest voucher"
// @Success 200 {object} model.GuestVoucherResponse "Voucher created successfully"
// @Failure 400 {object} model.Response "Invalid request or validation error"
// @Failure 429 {object} model.Response "IP address has already claimed a voucher"
// @Router /vouchers/guest [post]
func CreateGuestVoucher(c *fiber.Ctx) error {
	req := new(model.CreateGuestVoucher)
	if err := c.BodyParser(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "BAD_REQUEST")
	}

	if err := validate.Struct(req); err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), "VALIDATION_ERROR")
	}

	voucherRepo := repository.NewVoucherRepo(database.GetDB())
	userIP := c.IP()
	userAgent := c.Get("User-Agent")

	voucher, err := voucherRepo.GetGuestVoucher(req.CampaignHash, userIP, userAgent)
	return handleResultWithErrorCode(c, voucher, err, "", "Create guest voucher successfully")
}
