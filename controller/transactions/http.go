package transactions

import (
	"marketplace-backend/business/transactions"
	"marketplace-backend/controller/transactions/request"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	transactionUC transactions.UseCase
}

func NewTransactionController(transactionUC transactions.UseCase) *Controller {
	return &Controller{
		transactionUC: transactionUC,
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

/*
Update
*/

/*
Delete
*/
