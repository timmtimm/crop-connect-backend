package transactions

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/controller/transactions/request"
	"crop_connect/controller/transactions/response"
	"crop_connect/helper"
	"crop_connect/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	transactionUC transactions.UseCase
	proposalUC    proposals.UseCase
	commodityUC   commodities.UseCase
	userUC        users.UseCase
	batchUC       batchs.UseCase
	regionUC      regions.UseCase
}

func NewController(transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, batchUC batchs.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		transactionUC: transactionUC,
		proposalUC:    proposalUC,
		commodityUC:   commodityUC,
		userUC:        userUC,
		batchUC:       batchUC,
		regionUC:      regionUC,
	}
}

/*
Create
*/

func (tc *Controller) Create(c echo.Context) error {
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
			Message: "token tidak valid",
		})
	}
	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.BuyerID = userID

	statusCode, err := tc.transactionUC.Create(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}
	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "transaksi berhasil dibuat",
	})
}

/*
Read
*/

func (tc *Controller) GetUserTransactionWithPagination(c echo.Context) error {
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

	queryParam, err := request.QueryParamValidationForBuyer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	transactionQuery := transactions.Query{
		Skip:      queryPagination.Skip,
		Limit:     queryPagination.Limit,
		Sort:      queryPagination.Sort,
		Order:     queryPagination.Order,
		Commodity: queryParam.Commodity,
		Proposal:  queryParam.Proposal,
		Batch:     queryParam.Batch,
		Status:    queryParam.Status,
		StartDate: queryParam.StartDate,
		EndDate:   queryParam.EndDate,
	}

	if token.Role == constant.RoleBuyer {
		transactionQuery.BuyerID, err = primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}
	} else if token.Role == constant.RoleFarmer {
		transactionQuery.FarmerID, err = primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}
	}

	transactions, totalData, statusCode, err := tc.transactionUC.GetByPaginationAndQuery(transactionQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	transactionResponse, statusCode, err := response.FromArrayToResponseArray(transactions, tc.batchUC, tc.proposalUC, tc.commodityUC, tc.userUC, tc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan transaksi",
		Data:       transactionResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (tc *Controller) GetByID(c echo.Context) error {
	transactionID, err := primitive.ObjectIDFromHex(c.Param("transaction-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "transaction id tidak valid",
		})
	}

	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	transaction := transactions.Domain{}
	statusCode := http.StatusInternalServerError

	userID, err := primitive.ObjectIDFromHex(token.UID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "token tidak valid",
		})
	}

	buyerID := primitive.NilObjectID
	farmerID := primitive.NilObjectID

	if token.Role == constant.RoleBuyer {
		buyerID = userID
	} else {
		farmerID = userID
	}

	transaction, statusCode, err = tc.transactionUC.GetByIDAndBuyerIDOrFarmerID(transactionID, buyerID, farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	transactionResponse, statusCode, err := response.ConvertToTransactionResponse(&transaction, tc.batchUC, tc.proposalUC, tc.commodityUC, tc.userUC, tc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan transaksi",
		Data:    transactionResponse,
	})
}

func (tc *Controller) StatisticByYear(c echo.Context) error {
	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamStatistic(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	farmerID := primitive.NilObjectID
	if token.Role == constant.RoleFarmer {
		farmerID, err = primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}
	}

	transactionStatistic, statusCode, err := tc.transactionUC.StatisticByYear(farmerID, queryParam.Year)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan statistik",
		Data:    response.FromDomainArrayToStatistic(transactionStatistic),
	})
}

func (tc *Controller) StatisticTopProvince(c echo.Context) error {
	queryParam, err := request.QueryParamLimitAndYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	transactionStatistic, statusCode, err := tc.transactionUC.StatisticTopProvince(queryParam.Year, queryParam.Limit)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan statistik",
		Data:    response.FromDomainArrayToStatisticProvince(transactionStatistic),
	})
}

func (tc *Controller) StatisticTopCommodity(c echo.Context) error {
	queryParam, err := request.QueryParamLimitAndYear(c)
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

	farmerID := primitive.NilObjectID
	if token.Role == constant.RoleFarmer {
		farmerID, err = primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}
	}

	transactionStatistic, statusCode, err := tc.transactionUC.StatisticTopCommodity(farmerID, queryParam.Year, queryParam.Limit)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	responseStatistic, statusCode, err := response.FromDomainArrayToStatisticTopCommodity(transactionStatistic, tc.userUC, tc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan statistik",
		Data:    responseStatistic,
	})
}

func (tc *Controller) CountByCommodityID(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	totalTransaction, totalWeight, statusCode, err := tc.transactionUC.CountByCommodityID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan total transaksi berdasarkan komoditas",
		Data:    response.FromDomainToTransactionStatisticForCommodityPage(totalTransaction, totalWeight),
	})
}

/*
Update
*/

func (tc *Controller) MakeDecision(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	transactionID, err := primitive.ObjectIDFromHex(c.Param("transaction-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "transaction id tidak valid",
		})
	}

	userInput := request.Decision{}
	c.Bind(&userInput)

	isStatusValid := util.CheckStringOnArray([]string{constant.TransactionStatusAccepted, constant.TransactionStatusPending, constant.TransactionStatusRejected}, userInput.Decision)
	if !isStatusValid {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "keputusan hanya tersedia untuk accepted, pending, dan rejected",
		})
	}

	inputDomain := userInput.ToDomain()
	inputDomain.ID = transactionID

	statusCode, err := tc.transactionUC.MakeDecision(inputDomain, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "transaksi berhasil dibuat keputusan",
	})
}

func (tc *Controller) CancelOnPending(c echo.Context) error {
	transactionID, err := primitive.ObjectIDFromHex(c.Param("transaction-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "transaction id tidak valid",
		})
	}

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	statusCode, err := tc.transactionUC.CancelOnPending(transactionID, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "transaksi berhasil dibatalkan",
	})
}

/*
Delete
*/
