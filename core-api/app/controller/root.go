package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/task"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/workers"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/logger"
)

var logr = logger.GetLogger()

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
