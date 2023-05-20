package regions

import (
	"crop_connect/business/regions"
	"crop_connect/controller/regions/request"
	"crop_connect/controller/regions/response"
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	regionUC regions.UseCase
}

func NewController(regionUC regions.UseCase) *Controller {
	return &Controller{
		regionUC: regionUC,
	}
}

/*
Create
*/

/*
Read
*/

func (rc *Controller) GetByCountry(c echo.Context) error {
	statusCode := http.StatusBadRequest

	queryParam, err := request.QueryParamValidation(c, []string{"country"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	province, statusCode, err := rc.regionUC.GetByCountry(queryParam["country"])
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data provinsi berdasarkan negara",
		Data:    province,
	})
}

func (rc *Controller) GetByProvince(c echo.Context) error {
	statusCode := http.StatusBadRequest

	queryParam, err := request.QueryParamValidation(c, []string{"country", "province"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	city, statusCode, err := rc.regionUC.GetByProvince(queryParam["country"], queryParam["province"])
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data kabupaten berdasarkan provinsi",
		Data:    city,
	})
}

func (rc *Controller) GetByRegency(c echo.Context) error {
	statusCode := http.StatusBadRequest

	queryParam, err := request.QueryParamValidation(c, []string{"country", "province", "regency"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	district, statusCode, err := rc.regionUC.GetByRegency(queryParam["country"], queryParam["province"], queryParam["regency"])
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data kecamatan berdasarkan kabupaten",
		Data:    district,
	})
}

func (rc *Controller) GetByDistrict(c echo.Context) error {
	statusCode := http.StatusBadRequest

	queryParam, err := request.QueryParamValidation(c, []string{"country", "province", "regency", "district"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	subdistrict, statusCode, err := rc.regionUC.GetByDistrict(queryParam["country"], queryParam["province"], queryParam["regency"], queryParam["district"])
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data kelurahan berdasarkan kecamatan",
		Data:    response.FromDomainArray(subdistrict),
	})
}

/*
Update
*/

/*
Delete
*/
