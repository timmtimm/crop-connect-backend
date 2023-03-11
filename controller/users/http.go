package users

import (
	"marketplace-backend/business/users"
	"marketplace-backend/controller/users/request"
	"marketplace-backend/helper"
	"marketplace-backend/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUC users.UseCase
}

func NewUserController(userUC users.UseCase) *UserController {
	return &UserController{
		userUC: userUC,
	}
}

/*
Create
*/

func (uc *UserController) Register(c echo.Context) error {
	userInput := request.Register{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	token, statusCode, err := uc.userUC.Register(userInput.ToDomain())
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

/*
Read
*/

func (uc *UserController) Login(c echo.Context) error {
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

func (uc *UserController) GetProfile(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	user, err := uc.userUC.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "berhasil mendapatkan data user",
		Data:    user,
	})
}

/*
Update
*/

/*
Delete
*/
