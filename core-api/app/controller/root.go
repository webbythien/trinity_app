package controller

import (
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/task"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/workers"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/logger"
)

var logr = logger.GetLogger()
var validate = validator.New()

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Msg         string `json:"msg"`
}

type ErrorResponse struct {
	Msg string `json:"msg"`
}

func HealthCheckWorker(c *fiber.Ctx) error {
	return workers.Delay("core_api_queue", "Worker.HealthCheck", task.HealthCheck, int64(1))
}

func GetPagination(c *fiber.Ctx) (pageNo, pageSize int) {
	ps := c.Query("page_size")
	pn := c.Query("page")
	pageSize, pageNo = 10, 1

	if len(ps) > 0 {
		psInt, err := strconv.Atoi(ps)
		if err != nil {
			logr.Error(err)
		} else {
			pageSize = psInt
		}
	}

	if len(pn) > 0 {
		pnInt, err := strconv.Atoi(pn)
		if err != nil {
			logr.Error(err)
		} else {
			pageNo = pnInt
		}
	}

	return pageNo, pageSize
}

func returnBadRequestWithErrorCode(c *fiber.Ctx, message string, errorCode string) error {
	return c.Status(fiber.StatusBadRequest).JSON(model.Response{
		Data:      nil,
		Msg:       message,
		Success:   false,
		ErrorCode: &errorCode,
	})
}
func handleResultWithErrorCode(c *fiber.Ctx, result interface{}, err error, errorCode string, successMessage string) error {
	if err != nil {
		return returnBadRequestWithErrorCode(c, err.Error(), errorCode)
	}
	return c.JSON(model.Response{
		Data:      result,
		Msg:       successMessage,
		Success:   true,
		ErrorCode: nil,
	})
}

func parseRequestBody(c *fiber.Ctx, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return returnBadRequestWithErrorCode(c, "Invalid request format", "BAD_REQUEST")
	}
	return nil
}
