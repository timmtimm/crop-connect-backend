package treatment_records

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	treatmentRecord "crop_connect/business/treatment_records"
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/controller/treatment_records/request"
	"crop_connect/controller/treatment_records/response"
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	treatmentRecordUC treatmentRecord.UseCase
	batchUC           batchs.UseCase
	transactionUC     transactions.UseCase
	proposalUC        proposals.UseCase
	commodityUC       commodities.UseCase
	userUC            users.UseCase
	regionUC          regions.UseCase
}

func NewController(treatmentRecordUC treatmentRecord.UseCase, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		treatmentRecordUC: treatmentRecordUC,
		batchUC:           batchUC,
		transactionUC:     transactionUC,
		proposalUC:        proposalUC,
		commodityUC:       commodityUC,
		userUC:            userUC,
		regionUC:          regionUC,
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

	treatmentRecordsResponse, statusCode, err := response.FromDomainArray(treatmentRecords, trc.batchUC, trc.transactionUC, trc.proposalUC, trc.commodityUC, trc.userUC, trc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan riwayat perawatan",
		Data:       treatmentRecordsResponse,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (trc *Controller) GetByBatchID(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.QueryParam("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "batch id tidak valid",
		})
	}

	treatmentRecords, statusCode, err := trc.treatmentRecordUC.GetByBatchID(batchID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	treatmentRecordsResponse, statusCode, err := response.FromDomainArray(treatmentRecords, trc.batchUC, trc.transactionUC, trc.proposalUC, trc.commodityUC, trc.userUC, trc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan riwayat perawatan",
		Data:    treatmentRecordsResponse,
	})
}

func (trc *Controller) CountByYear(c echo.Context) error {
	queryYear, err := request.QueryParamValidationYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	count, statusCode, err := trc.treatmentRecordUC.CountByYear(queryYear)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan jumlah riwayat perawatan",
		Data:    count,
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

func (trc *Controller) Validate(c echo.Context) error {
	treatmentRecordID, err := primitive.ObjectIDFromHex(c.Param("treatment-record-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "treatment record id tidak valid",
		})
	}

	userInput := request.Validate{}
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

	inputDomain := userInput.ToDomain()
	inputDomain.ID = treatmentRecordID

	_, statusCode, err := trc.treatmentRecordUC.Validate(inputDomain, userID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "catatan perawatan berhasil divalidasi",
	})
}

func (trc *Controller) UpdateNotes(c echo.Context) error {
	treatmentRecordID, err := primitive.ObjectIDFromHex(c.Param("treatment-record-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "treatment record id tidak valid",
		})
	}

	userInput := request.UpdateNotes{}
	c.Bind(&userInput)

	inputDomain := userInput.ToDomain()
	inputDomain.ID = treatmentRecordID

	_, statusCode, err := trc.treatmentRecordUC.UpdateNotes(inputDomain)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "catatan perawatan berhasil diupdate",
	})
}

/*
Delete
*/
