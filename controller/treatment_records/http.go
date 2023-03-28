package treatment_records

import (
	"marketplace-backend/business/batchs"
	treatmentRecord "marketplace-backend/business/treatment_records"
	"marketplace-backend/business/users"
	"marketplace-backend/constant"
	"marketplace-backend/controller/treatment_records/request"
	"marketplace-backend/controller/treatment_records/response"
	"marketplace-backend/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	treatmentRecordUC treatmentRecord.UseCase
	batchUC           batchs.UseCase
	userUC            users.UseCase
}

func NewTreatmentRecordController(treatmentRecordUC treatmentRecord.UseCase, batchUC batchs.UseCase, userUC users.UseCase) *Controller {
	return &Controller{
		treatmentRecordUC: treatmentRecordUC,
		batchUC:           batchUC,
		userUC:            userUC,
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

func (trc *Controller) GetByPaginationAndQuery(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"number", "date", "status", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidationForBuyer(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	treatmentRecordQuery := treatmentRecord.Query{
		Skip:      queryPagination.Skip,
		Limit:     queryPagination.Limit,
		Sort:      queryPagination.Sort,
		Order:     queryPagination.Order,
		Commodity: queryParam.Commodity,
		Batch:     queryParam.Batch,
		Number:    queryParam.Number,
		Status:    queryParam.Status,
	}

	if token.Role == constant.RoleFarmer {
		farmerID, err := primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}

		treatmentRecordQuery.FarmerID = farmerID
	}

	treatmentRecords, totalData, statusCode, err := trc.treatmentRecordUC.GetByPaginationAndQuery(treatmentRecordQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	transactionResponse, statusCode, err := response.FromDomainArray(treatmentRecords, trc.batchUC, trc.userUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan riwayat perawatan",
		Data:       transactionResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

/*
Update
*/

func (trc *Controller) FillTreatmentRecord(c echo.Context) error {
	treatmentRecordID, err := primitive.ObjectIDFromHex(c.Param("treatment-record-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "treatment record id tidak valid",
		})
	}

	userInput := request.FillTreatmentRecord{}
	c.Bind(&userInput)

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	inputDomain := treatmentRecord.Domain{
		ID: treatmentRecordID,
	}

	images, statusCode, err := helper.GetCreateImageRequest(c, []string{"image1", "image2", "image3", "image4", "image5"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	if len(images) == 0 {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "gambar tidak boleh kosong",
		})
	}

	_, statusCode, err = trc.treatmentRecordUC.FillTreatmentRecord(&inputDomain, userID, images, userInput.Notes)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "catatan perawatan berhasil diisi",
	})
}

/*
Delete
*/
