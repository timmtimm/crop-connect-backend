package commodities

import (
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/users"
	"marketplace-backend/controller/commodities/request"
	"marketplace-backend/controller/commodities/response"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	commodityUC commodities.UseCase
	userUC      users.UseCase
}

func NewCommodityController(commodityUC commodities.UseCase, userUC users.UseCase) *Controller {
	return &Controller{
		commodityUC: commodityUC,
		userUC:      userUC,
	}
}

/*
Create
*/

func (cc *Controller) Create(c echo.Context) error {
	userInput := request.Commodity{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.FarmerID = userID

	statusCode, err := cc.commodityUC.Create(userDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil membuat komoditas",
	})
}

/*
Read
*/

func (cc *Controller) GetForBuyer(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "plantingPeriod", "pricePerKg", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	QueryParam, err := request.QueryParamValidation(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	commodities, totalData, statusCode, err := cc.commodityUC.GetByPaginationAndQuery(commodities.Query{
		Skip:     queryPagination.Skip,
		Limit:    queryPagination.Limit,
		Sort:     queryPagination.Sort,
		Order:    queryPagination.Order,
		Name:     QueryParam.Name,
		MinPrice: QueryParam.MinPrice,
		MaxPrice: QueryParam.MaxPrice,
	})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	commodityResponse, statusCode, err := response.FromDomainArray(commodities, cc.userUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan komoditas",
		Data:       commodityResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

/*
Update
*/

func (cc *Controller) Update(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	userInput := request.Commodity{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.ID = commodityID
	userDomain.FarmerID = userID

	statusCode, err := cc.commodityUC.Update(userDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mengubah komoditas",
	})
}

/*
Delete
*/

func (cc *Controller) Delete(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	statusCode, err := cc.commodityUC.Delete(commodityID, farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil menghapus komoditas",
	})
}
