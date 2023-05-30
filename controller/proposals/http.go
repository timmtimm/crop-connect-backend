package proposals

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	"crop_connect/constant"
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

func (pc *Controller) StatisticByYear(c echo.Context) error {
	year, err := request.QueryParamValidationYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	proposals, statusCode, err := pc.proposalUC.StatisticByYear(year)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil didapatkan",
		Data:    proposals,
	})
}

func (pc *Controller) CountTotalProposalByFarmer(c echo.Context) error {
	farmerID, err := primitive.ObjectIDFromHex(c.Param("farmer-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id petani tidak valid",
		})
	}

	user, _, err := pc.userUC.GetByID(farmerID)
	if user.Role != constant.RoleFarmer || err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "petani tidak ditemukan",
		})
	}

	totalProposal, statusCode, err := pc.proposalUC.CountTotalProposalByFarmer(farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "proposal berhasil didapatkan",
		Data:    totalProposal,
	})
}

func (pc *Controller) GetByPaginationAndQuery(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "status", "plantingArea", "estimatedTotalHarvest", "isAvailable", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	queryparam, err := request.QueryParamValidation(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	query := proposals.Query{
		Skip:        queryPagination.Skip,
		Limit:       queryPagination.Limit,
		Sort:        queryPagination.Sort,
		Order:       queryPagination.Order,
		FarmerID:    userID,
		CommodityID: queryparam.CommodityID,
		Name:        queryparam.Name,
		Status:      queryparam.Status,
	}

	proposals, totalData, statusCode, err := pc.proposalUC.GetByPaginationAndQuery(query)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	proposalResponse, statusCode, err := response.FromDomainArrayToProposalWithCommodity(proposals, pc.userUC, pc.commodityUC, pc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:     http.StatusOK,
		Message:    "berhasil mendapatkan proposal",
		Data:       proposalResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (pc *Controller) GetByID(c echo.Context) error {
	proposalID, err := primitive.ObjectIDFromHex(c.Param("proposal-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id proposal tidak valid",
		})
	}

	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	farmerID := primitive.NilObjectID
	if token.Role == constant.RoleFarmer {
		farmerID, err = primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "id petani tidak valid",
			})
		}
	}

	proposal, statusCode, err := pc.proposalUC.GetByIDAndFarmerID(proposalID, farmerID)
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
		Message: "berhasil mendapatkan proposal",
		Data:    proposalResponse,
	})
}

func (pc *Controller) GetForPerennials(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "token tidak valid",
		})
	}

	proposals, statusCode, err := pc.proposalUC.GetForPerennials(userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	proposalResponse, statusCode, err := response.FromDomainArrayToProposalWithCommodity(proposals, pc.userUC, pc.commodityUC, pc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "berhasil mendapatkan proposal",
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
