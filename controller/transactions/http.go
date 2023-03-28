package transactions

import (
	"marketplace-backend/business/batchs"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/transactions"
	"marketplace-backend/business/users"
	"marketplace-backend/constant"
	"marketplace-backend/controller/transactions/request"
	"marketplace-backend/controller/transactions/response"
	"marketplace-backend/helper"
	"marketplace-backend/util"
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
}

func NewTransactionController(transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, batchUC batchs.UseCase) *Controller {
	return &Controller{
		transactionUC: transactionUC,
		proposalUC:    proposalUC,
		commodityUC:   commodityUC,
		userUC:        userUC,
		batchUC:       batchUC,
	}
}

/*
Create
*/

func (tc *Controller) Create(c echo.Context) error {
	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "proposal id tidak valid",
		})
	}

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

	inputDomain := userInput.ToDomain()
	inputDomain.ProposalID = proposalID
	inputDomain.BuyerID = userID

	statusCode, err := tc.transactionUC.Create(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
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
		Status:    queryParam.Status,
		StartDate: queryParam.StartDate,
		EndDate:   queryParam.EndDate,
	}

	if token.Role == constant.RoleBuyer {
		buyerID, err := primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}

		transactionQuery.BuyerID = buyerID
	} else if token.Role == constant.RoleFarmer {
		FarmerID, err := primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}

		transactionQuery.FarmerID = FarmerID
	}

	transactions, totalData, statusCode, err := tc.transactionUC.GetByPaginationAndQuery(transactionQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	transactionResponse, statusCode, err := response.FromDomainArrayToBuyer(transactions, tc.proposalUC, tc.commodityUC, tc.userUC)
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
	inputDomain.BuyerID = userID

	statusCode, err := tc.transactionUC.MakeDecision(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	statusCode, err = tc.batchUC.Create(inputDomain.ID)
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

/*
Delete
*/
