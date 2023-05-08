package users

import (
	"crop_connect/business/regions"
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/controller/users/request"
	"crop_connect/controller/users/response"
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	userUC   users.UseCase
	regionUC regions.UseCase
}

func NewController(userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		userUC:   userUC,
		regionUC: regionUC,
	}
}

/*
Create
*/

func (uc *Controller) Register(c echo.Context) error {
	userInput := request.RegisterUser{}
	c.Bind(&userInput)

	if validationErr := userInput.Validate(); validationErr != nil {
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

	token, statusCode, err := uc.userUC.Register(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "registrasi sukses",
		Data:    token,
	})
}

func (uc *Controller) RegisterValidator(c echo.Context) error {
	userInput := request.RegisterValidator{}
	c.Bind(&userInput)

	if validationErr := userInput.Validate(); validationErr != nil {
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

	token, statusCode, err := uc.userUC.RegisterValidator(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "registrasi validasi sukses",
		Data:    token,
	})
}

/*
Read
*/

func (uc *Controller) Login(c echo.Context) error {
	userInput := request.Login{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	token, statusCode, err := uc.userUC.Login(userInput.ToDomain())
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "login sukses",
		Data:    token,
	})
}

func (uc *Controller) GetProfile(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	user, statusCode, err := uc.userUC.GetByID(userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	userResponse, statusCode, err := response.FromDomain(user, uc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data user",
		Data:    userResponse,
	})
}

func (uc *Controller) GetFarmerByPaginationAndQueryForBuyer(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidationForSearchFarmer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	userQuery := users.Query{
		Skip:  queryPagination.Skip,
		Limit: queryPagination.Limit,
		Sort:  queryPagination.Sort,
		Order: queryPagination.Order,
		Name:  queryParam.Name,
		Role:  constant.RoleFarmer,
	}

	users, totalData, statusCode, err := uc.userUC.GetByPaginationAndQuery(userQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	usersReponse, statusCode, err := response.FromDomainArray(users, uc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan data user",
		Data:       usersReponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (uc *Controller) GetFarmerByIDForBuyer(c echo.Context) error {
	farmerID, err := primitive.ObjectIDFromHex(c.Param("farmer-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id komoditas tidak valid",
		})
	}

	farmer, statusCode, err := uc.userUC.GetFarmerByID(farmerID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	farmerResponse, statusCode, err := response.FromDomain(farmer, uc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan data petani",
		Data:    farmerResponse,
	})
}

func (uc *Controller) GetByPaginationAndQueryForAdmin(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"name", "email", "role", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidationForSearchFarmer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	userQuery := users.Query{
		Skip:        queryPagination.Skip,
		Limit:       queryPagination.Limit,
		Sort:        queryPagination.Sort,
		Order:       queryPagination.Order,
		Name:        queryParam.Name,
		Email:       queryParam.Email,
		PhoneNumber: queryParam.PhoneNumber,
		Role:        queryParam.Role,
	}

	users, totalData, statusCode, err := uc.userUC.GetByPaginationAndQuery(userQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	usersReponse, statusCode, err := response.FromDomainArray(users, uc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan data user",
		Data:       usersReponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

/*
Update
*/

func (uc *Controller) UpdateProfile(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userInput := request.Update{}
	c.Bind(&userInput)

	if validationErr := userInput.Validate(); validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	userDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	userDomain.ID = userID

	_, statusCode, err := uc.userUC.UpdateProfile(userDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil update data user",
	})
}

func (uc *Controller) UpdatePassword(c echo.Context) error {
	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userInput := request.ChangePassword{}
	c.Bind(&userInput)

	if validationErr := userInput.Validate(); validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.ID = userID

	_, statusCode, err := uc.userUC.UpdatePassword(userDomain, userInput.NewPassword)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil update password user",
	})
}

/*
Delete
*/
