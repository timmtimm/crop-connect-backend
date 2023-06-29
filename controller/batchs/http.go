package batchs

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/controller/batchs/request"
	"crop_connect/controller/batchs/response"
	"crop_connect/helper"
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
	regionUC      regions.UseCase
}

func NewController(batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
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

func (bc *Controller) CreateForPerennials(c echo.Context) error {
	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	statusCode, err := bc.batchUC.Create(proposalID, farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil membuat batch",
	})
}

/*
Read
*/

func (bc *Controller) GetByPaginationAndQuery(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "status", "createdAt"})
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

	queryParam, err := request.QueryParamValidationForBuyer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	batchQuery := batchs.Query{
		Skip:        queryPagination.Skip,
		Limit:       queryPagination.Limit,
		Sort:        queryPagination.Sort,
		Order:       queryPagination.Order,
		FarmerID:    queryParam.FarmerID,
		CommodityID: queryParam.CommodityID,
		Name:        queryParam.Name,
		Status:      queryParam.Status,
	}

	if token.Role == constant.RoleFarmer {
		farmerID, err := helper.GetUIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		batchQuery.FarmerID = farmerID
	}

	batchs, totalData, statusCode, err := bc.batchUC.GetByPaginationAndQuery(batchQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomainArray(batchs, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

	batchResponse, statusCode, err := response.FromDomainArray(batchs, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

func (bc *Controller) CountByYear(c echo.Context) error {
	queryYear, err := request.QueryParamValidationYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	count, statusCode, err := bc.batchUC.CountByYear(queryYear)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan jumlah batch",
		Data:    count,
	})
}

func (bc *Controller) GetByID(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	batch, statusCode, err := bc.batchUC.GetByID(batchID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomain(batch, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

func (bc *Controller) GetForTransactionByCommodityID(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	batch, statusCode, err := bc.batchUC.GetForTransactionByCommodityID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomainArray(batch, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

func (bc *Controller) GetForTransactionByID(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	batch, statusCode, err := bc.batchUC.GetForTransactionByID(batchID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomain(batch, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

func (bc *Controller) GetForHarvestByCommmodityID(c echo.Context) error {
	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	batch, statusCode, err := bc.batchUC.GetForHarvestByFarmerID(farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	batchResponse, statusCode, err := response.FromDomainArray(batch, bc.proposalUC, bc.commodityUC, bc.userUC, bc.regionUC)
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

// func (bc *Controller) Cancel(c echo.Context) error {
// 	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "id batch tidak valid",
// 		})
// 	}

// 	farmerID, err := helper.GetUIDFromToken(c)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
// 			Status:  http.StatusUnauthorized,
// 			Message: err.Error(),
// 		})
// 	}

// 	userInput := request.Cancel{}
// 	c.Bind(&userInput)

// 	validationErr := userInput.Validate()
// 	if validationErr != nil {
// 		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "validasi gagal",
// 			Error:   validationErr,
// 		})
// 	}

// 	inputDomain := userInput.ToDomain()
// 	inputDomain.ID = batchID

// 	statusCode, err := bc.batchUC.Cancel(inputDomain, farmerID)
// 	if err != nil {
// 		return c.JSON(statusCode, helper.BaseResponse{
// 			Status:  statusCode,
// 			Message: err.Error(),
// 		})
// 	}

// 	return c.JSON(statusCode, helper.BaseResponse{
// 		Status:  statusCode,
// 		Message: "berhasil membatalkan batch",
// 	})
// }

/*
Delete
*/
