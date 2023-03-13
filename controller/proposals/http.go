package proposals

import (
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/controller/proposals/request"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	proposalUC  proposals.UseCase
	commodityUC commodities.UseCase
}

func NewProposalController(proposalUC proposals.UseCase, commodityUC commodities.UseCase) *Controller {
	return &Controller{
		proposalUC:  proposalUC,
		commodityUC: commodityUC,
	}
}

/*
Create
*/

func (pc *Controller) Create(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
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
	inputDomain.CommodityID = commodityID
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Error:   validationErr,
		})
	}

	_, statusCode, err := pc.commodityUC.GetByID(inputDomain.CommodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: "komoditas tidak ditemukan",
			Error:   validationErr,
		})
	}

	statusCode, err = pc.proposalUC.Create(inputDomain, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "proposal berhasil dibuat",
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
