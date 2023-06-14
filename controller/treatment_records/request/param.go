package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	Farmer    string
	FarmerID  primitive.ObjectID
	Commodity string
	BatchID   primitive.ObjectID
	Batch     string
	Number    int
	Status    string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Commodity: c.QueryParam("commodity"),
		Status:    c.QueryParam("status"),
		Farmer:    c.QueryParam("farmer"),
		Batch:     c.QueryParam("batch"),
	}

	var err error

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, %s, dan %s", constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse)
		}
	}

	if batch := c.QueryParam("batchID"); batch != "" {
		filter.BatchID, err = primitive.ObjectIDFromHex(batch)
		if err != nil {
			return FilterQuery{}, errors.New("batchID harus berupa hex")
		}
	}

	if farmer := c.QueryParam("farmerID"); farmer != "" {
		filter.BatchID, err = primitive.ObjectIDFromHex(farmer)
		if err != nil {
			return FilterQuery{}, errors.New("farmerID harus berupa hex")
		}
	}

	if number := c.QueryParam("number"); number != "" {
		filter.Number, err = strconv.Atoi(number)
		if err != nil {
			return FilterQuery{}, errors.New("number harus berupa angka")
		}
	}

	return filter, nil
}

func QueryParamValidationYear(c echo.Context) (int, error) {
	if year := c.QueryParam("year"); year != "" {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return 0, errors.New("year harus berupa angka")
		}

		if year > time.Now().Format("2006") {
			return 0, errors.New("year tidak boleh lebih dari tahun sekarang")
		}

		return yearInt, nil
	}

	return time.Now().Year(), nil
}
