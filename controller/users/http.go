package users

import (
	"marketplace-backend/business/users"
	"marketplace-backend/controller/users/request"
	"marketplace-backend/helper"
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
		Data: map[string]interface{}{
			"token": token,
		},
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
		Data: map[string]interface{}{
			"token": token,
		},
	})
}

/*
Update
*/

/*
Delete
*/
