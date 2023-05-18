package proposals

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	"crop_connect/controller/proposals/request"
	"crop_connect/controller/proposals/response"
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	proposalUC  proposals.UseCase
	commodityUC commodities.UseCase
	userUC      users.UseCase
	regionUC    regions.UseCase
}

func NewController(proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		proposalUC:  proposalUC,
		commodityUC: commodityUC,
		userUC:      userUC,
		regionUC:    regionUC,
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

	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.CommodityID = commodityID
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	_, statusCode, err := pc.commodityUC.GetByID(inputDomain.CommodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: "komoditas tidak ditemukan",
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

func (pc *Controller) GetByCommodityIDForBuyer(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	// _, statusCode, err := pc.commodityUC.GetByID(commodityID)
	// if err != nil {
	// 	return c.JSON(statusCode, helper.BaseResponse{
	// 		Status:  statusCode,
	// 		Message: "komoditas tidak ditemukan",
	// 	})
	// }

	proposals, statusCode, err := pc.proposalUC.GetByCommodityID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil didapatkan",
		Data:    response.FromDomainArrayToBuyer(proposals),
	})
}

func (pc *Controller) GetByIDAccepted(c echo.Context) error {
	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
		})
	}

	proposal, statusCode, err := pc.proposalUC.GetByIDAccepted(proposalID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	proposalResponse, statusCode, err := response.FromDomainToProposalWithCommodity(&proposal, pc.userUC, pc.commodityUC, pc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil didapatkan",
		Data:    proposalResponse,
	})
}

/*
Update
*/

func (pc *Controller) Update(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
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

	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.ID = proposalID

	statusCode, err := pc.proposalUC.Update(inputDomain, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil diubah",
	})
}

func (pc *Controller) ValidateByValidator(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	id, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
		})
	}

	userInput := request.Validate{}
	c.Bind(&userInput)

	inputDomain := userInput.ToDomain()
	inputDomain.ID = id

	statusCode, err := pc.proposalUC.ValidateProposal(inputDomain, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil divalidasi",
	})
}

/*
Delete
*/

func (pc *Controller) Delete(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
		})
	}

	statusCode, err := pc.proposalUC.Delete(proposalID, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil dihapus",
	})
}
