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

// GetAllProvinces godoc
// @Summary Get all provinces
// @Description Retrieve all provinces from the database
// @Tags regions
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /region/provinces [get]
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

// GetAllDistricts godoc
// @Summary Get all districts by province code
// @Description Retrieve all districts for a given province code
// @Tags regions
// @Accept  json
// @Produce  json
// @Param  provinceCode  path  string  true  "Province Code"
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /region/districts/{provinceCode} [get]
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

// GetAllWards godoc
// @Summary Get all wards by district code
// @Description Retrieve all wards for a given district code
// @Tags regions
// @Accept  json
// @Produce  json
// @Param  districtCode  path  string  true  "District Code"
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /region/wards/{districtCode} [get]
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
