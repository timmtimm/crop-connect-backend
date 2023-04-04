package harvests

import (
	"marketplace-backend/business/batchs"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/harvests"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/regions"
	"marketplace-backend/business/transactions"
	"marketplace-backend/business/users"
	"marketplace-backend/constant"
	"marketplace-backend/controller/harvests/request"
	"marketplace-backend/controller/harvests/response"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	harvestUC     harvests.UseCase
	batchUC       batchs.UseCase
	transactionUC transactions.UseCase
	proposalUC    proposals.UseCase
	commodityUC   commodities.UseCase
	userUC        users.UseCase
	regionUC      regions.UseCase
}

func NewHarvestController(harvestUC harvests.UseCase, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		harvestUC:     harvestUC,
		batchUC:       batchUC,
		transactionUC: transactionUC,
		proposalUC:    proposalUC,
		commodityUC:   commodityUC,
		userUC:        userUC,
		regionUC:      regionUC,
	}
}

/*
Create
*/

func (hc *Controller) SubmitHarvest(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	userInput := request.SubmitHarvest{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	images, statusCode, err := helper.GetCreateImageRequest(c, []string{"image1", "image2", "image3", "image4", "image5"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	inputDomain := userInput.ToDomain()
	inputDomain.BatchID = batchID

	_, statusCode, err = hc.harvestUC.SubmitHarvest(inputDomain, userID, images, userInput.Notes)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mengajukan hasil panen",
	})
}

/*
Read
*/

func (hc *Controller) GetByPaginationAndQuery(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"status", "totalPrice", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidation(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	harvestQuery := harvests.Query{
		Skip:      queryPagination.Skip,
		Limit:     queryPagination.Limit,
		Sort:      queryPagination.Sort,
		Order:     queryPagination.Order,
		Commodity: queryParam.Commodity,
		Batch:     queryParam.Batch,
		Status:    queryParam.Status,
	}

	if token.Role == constant.RoleFarmer {
		FarmerID, err := primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}

		harvestQuery.FarmerID = FarmerID
	}

	harvests, totalData, statusCode, err := hc.harvestUC.GetByPaginationAndQuery(harvestQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	harvestResponses, statusCode, err := response.FromDomainArrayToResponse(harvests, hc.batchUC, hc.transactionUC, hc.proposalUC, hc.commodityUC, hc.userUC, hc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan panen",
		Data:       harvestResponses,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

/*
Update
*/

/*
Delete
*/
