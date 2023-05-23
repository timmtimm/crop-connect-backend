package commodities

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	"crop_connect/controller/commodities/request"
	"crop_connect/controller/commodities/response"
	"crop_connect/helper"
	"crop_connect/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	commodityUC commodities.UseCase
	userUC      users.UseCase
	proposalUC  proposals.UseCase
	regionUC    regions.UseCase
}

func NewController(commodityUC commodities.UseCase, userUC users.UseCase, proposalUC proposals.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		commodityUC: commodityUC,
		userUC:      userUC,
		proposalUC:  proposalUC,
		regionUC:    regionUC,
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

	userDomain := userInput.ToDomain()
	userDomain.FarmerID = userID

	statusCode, err = cc.commodityUC.Create(userDomain, images)
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
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "plantingPeriod", "pricePerKg", "isAvailable", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
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

	queryRegion, err := request.QueryValidationForRegion(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	commodities, totalData, statusCode, err := cc.commodityUC.GetByPaginationAndQuery(commodities.Query{
		Skip:     queryPagination.Skip,
		Limit:    queryPagination.Limit,
		Sort:     queryPagination.Sort,
		Order:    queryPagination.Order,
		Name:     queryParam.Name,
		Farmer:   queryParam.Farmer,
		MinPrice: queryParam.MinPrice,
		MaxPrice: queryParam.MaxPrice,
		FarmerID: queryParam.FarmerID,
		Province: queryRegion.Province,
		Regency:  queryRegion.Regency,
		District: queryRegion.District,
		RegionID: queryRegion.RegionID,
	})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	commodityResponse, statusCode, err := response.FromDomainArray(commodities, cc.userUC, cc.regionUC)
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

func (cc *Controller) GetByID(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	commodity, statusCode, err := cc.commodityUC.GetByID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	commodityResponse, statusCode, err := response.FromDomain(commodity, cc.userUC, cc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan komoditas",
		Data:    commodityResponse,
	})
}

func (cc *Controller) GetForFarmer(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	commodities, statusCode, err := cc.commodityUC.GetByFarmerID(userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	commodityResponse, statusCode, err := response.FromDomainArray(commodities, cc.userUC, cc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan komoditas",
		Data:    commodityResponse,
	})
}

func (cc *Controller) CountTotalCommodity(c echo.Context) error {
	year, err := request.QueryParamValidationYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	totalCommodity, statusCode, err := cc.commodityUC.CountTotalCommodity(year)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan total komoditas",
		Data:    totalCommodity,
	})
}

func (cc *Controller) CountTotalCommodityByFarmer(c echo.Context) error {
	farmerID, err := primitive.ObjectIDFromHex(c.Param("farmer-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id petani tidak valid",
		})
	}

	totalCommodity, statusCode, err := cc.commodityUC.CountTotalCommodityByFarmer(farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan total komoditas petani",
		Data:    totalCommodity,
	})
}

/*
Update
*/

func (cc *Controller) Update(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	userInput := request.Update{}
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

	commodity, statusCode, err := cc.commodityUC.GetByIDAndFarmerID(commodityID, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	updateImage, statusCode, err := helper.GetUpdateImageRequest(c, []string{"image1", "image2", "image3", "image4", "image5"}, commodity.ImageURLs, util.ConvertArrayStringToBool(userInput.IsChange), util.ConvertArrayStringToBool(userInput.IsDelete))
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.ID = commodityID
	userDomain.FarmerID = userID

	_, statusCode, err = cc.commodityUC.Update(userDomain, updateImage)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	statusCode, err = cc.proposalUC.UpdateCommodityID(commodityID, commodity.ID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mengubah komoditas",
	})
}

/*
Delete
*/

func (cc *Controller) Delete(c echo.Context) error {
	commodityID, err := primitive.ObjectIDFromHex(c.Param("commodity-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	statusCode, err := cc.commodityUC.Delete(commodityID, farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	statusCode, err = cc.proposalUC.DeleteByCommodityID(commodityID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil menghapus komoditas",
	})
}
