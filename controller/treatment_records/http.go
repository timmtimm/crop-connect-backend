package treatment_records

import (
	treatmentRecord "marketplace-backend/business/treatment_records"
	"marketplace-backend/controller/treatment_records/request"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	treatmentRecordUC treatmentRecord.UseCase
}

func NewTreatmentRecordController(treatmentRecordUC treatmentRecord.UseCase) *Controller {
	return &Controller{
		treatmentRecordUC: treatmentRecordUC,
	}
}

/*
Create
*/

func (trc *Controller) RequestToFarmer(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "batch id tidak valid",
		})
	}

	userInput := request.RequestToFarmer{}
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

	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.BatchID = batchID
	inputDomain.RequesterID = userID

	_, statusCode, err := trc.treatmentRecordUC.RequestToFarmer(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "permintaaan pengisian catatan perawatan berhasil dibuat",
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
