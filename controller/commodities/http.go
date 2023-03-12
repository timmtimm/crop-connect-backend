package commodities

import (
	"fmt"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/users"
	"marketplace-backend/controller/commodities/request"
	"marketplace-backend/controller/commodities/response"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
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
	userInput := request.Create{}
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
	pagination := helper.PaginationParam{
		Page:  c.Param("page"),
		Limit: c.QueryParam("limit"),
		Sort:  c.QueryParam("sort"),
		Order: c.QueryParam("order"),
	}

	queryPagination, err := helper.PaginationToQuery(pagination, []string{"name", "plantingPeriod", "pricePerKg", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	commodities, totalData, statusCode, err := cc.commodityUC.GetByPaginationAndQuery(commodities.Query{
		Skip:  queryPagination.Skip,
		Limit: queryPagination.Limit,
		Sort:  queryPagination.Sort,
		Order: queryPagination.Order,
		Name:  c.QueryParam("name"),
	})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	fmt.Println(commodities)

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

/*
Delete
*/
