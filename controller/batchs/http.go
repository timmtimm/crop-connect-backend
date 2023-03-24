package batchs

import (
	"marketplace-backend/business/batchs"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/transactions"
	"marketplace-backend/business/users"
	"marketplace-backend/controller/batchs/request"
	"marketplace-backend/controller/batchs/response"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	batchUC       batchs.UseCase
	transactionUC transactions.UseCase
	proposalUC    proposals.UseCase
	commodityUC   commodities.UseCase
	userUC        users.UseCase
}

func NewBatchController(batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase) *Controller {
	return &Controller{
		batchUC:       batchUC,
		transactionUC: transactionUC,
		proposalUC:    proposalUC,
		commodityUC:   commodityUC,
		userUC:        userUC,
	}
}

/*
Create
*/

/*
Read
*/

func (bc *Controller) GetFarmerBatch(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "status", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidationForBuyer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	batchQuery := batchs.Query{
		Skip:      queryPagination.Skip,
		Limit:     queryPagination.Limit,
		Sort:      queryPagination.Sort,
		Order:     queryPagination.Order,
		FarmerID:  farmerID,
		Commodity: queryParam.Commodity,
		Name:      queryParam.Name,
		Status:    queryParam.Status,
	}

	batchs, totalData, statusCode, err := bc.batchUC.GetByPaginationAndQuery(batchQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomainArray(batchs, bc.transactionUC, bc.proposalUC, bc.commodityUC, bc.userUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan batch petani",
		Data:       batchResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (bc *Controller) GetByCommodityID(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	batchs, statusCode, err := bc.batchUC.GetByCommodityID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomainArray(batchs, bc.transactionUC, bc.proposalUC, bc.commodityUC, bc.userUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan batch",
		Data:    batchResponse,
	})
}

/*
Update
*/

func (bc *Controller) Cancel(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userInput := request.Cancel{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	inputDomain := userInput.ToDomain()
	inputDomain.ID = batchID

	statusCode, err := bc.batchUC.Cancel(inputDomain, farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil membatalkan batch",
	})
}

/*
Delete
*/
