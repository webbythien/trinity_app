package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
)

// RegionController holds the repository
type RegionController struct {
	Repo repository.RegionRepository
}

// NewRegionController creates a new RegionController
func NewRegionController(repo repository.RegionRepository) *RegionController {
	return &RegionController{Repo: repo}
}

func (c *RegionController) GetAllProvinces(ctx *fiber.Ctx) error {
	provinces, err := c.Repo.GetAllProvince()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(model.Response{
			Data:      nil,
			Msg:       err.Error(),
			Success:   false,
			ErrorCode: nil, // Or provide error code
		})
	}
	return ctx.JSON(model.Response{
		Data:      provinces,
		Msg:       "get provinces successfully",
		Success:   true,
		ErrorCode: nil,
	})
}

func (c *RegionController) GetAllDistricts(ctx *fiber.Ctx) error {
	provinceCode := ctx.Params("provinceCode")
	districts, err := c.Repo.GetAllDistrict(provinceCode)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(model.Response{
			Data:      nil,
			Msg:       err.Error(),
			Success:   false,
			ErrorCode: nil,
		})
	}
	return ctx.JSON(model.Response{
		Data:      districts,
		Msg:       "get districts successfully",
		Success:   true,
		ErrorCode: nil,
	})
}

func (c *RegionController) GetAllWards(ctx *fiber.Ctx) error {
	districtCode := ctx.Params("districtCode")
	wards, err := c.Repo.GetAllWard(districtCode)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(model.Response{
			Data:      nil,
			Msg:       err.Error(),
			Success:   false,
			ErrorCode: nil,
		})
	}
	return ctx.JSON(model.Response{
		Data:      wards,
		Msg:       "get wards successfully",
		Success:   true,
		ErrorCode: nil,
	})
}
