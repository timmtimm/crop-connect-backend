package forgot_password

import (
	forgotPassword "crop_connect/business/forgot_password"
	"crop_connect/controller/forgot_password/request"
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	forgotPasswordUC forgotPassword.UseCase
}

func NewController(forgotPasswordUC forgotPassword.UseCase) *Controller {
	return &Controller{
		forgotPasswordUC: forgotPasswordUC,
	}
}

/*
Create
*/

func (fpc *Controller) Generate(c echo.Context) error {
	userInput := request.Generate{}
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   err,
		})
	}

	statusCode, _ := fpc.forgotPasswordUC.Generate(userInput.Domain, userInput.Email)
	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "jika email terdaftar, maka akan dikirimkan link untuk mereset password",
	})
}

/*
Read
*/

func (fpc *Controller) ValidateToken(c echo.Context) error {
	statusCode, err := fpc.forgotPasswordUC.ValidateToken(c.Param("token"))
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "token dapat digunakan",
	})
}

/*
Update
*/

func (fpc *Controller) ResetPassword(c echo.Context) error {
	userInput := request.UpdatePassword{}
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   err,
		})
	}

	statusCode, err := fpc.forgotPasswordUC.ResetPassword(c.Param("token"), userInput.Password)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "password berhasil diubah",
	})
}

/*
Delete
*/
