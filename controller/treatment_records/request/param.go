package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Commodity string
	Batch     string
	Number    int
	Status    string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Commodity: c.QueryParam("commodity"),
		Batch:     c.QueryParam("batch"),
		Status:    c.QueryParam("status"),
	}

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, %s, dan %s", constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse)
		}
	}

	if number := c.QueryParam("number"); number != "" {
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			return FilterQuery{}, errors.New("number harus berupa angka")
		}

		filter.Number = numberInt
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
