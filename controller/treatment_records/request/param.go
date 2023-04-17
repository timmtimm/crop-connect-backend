package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Commodity string
	Batch     string
	Number    int
	Status    string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}

	commodity := c.QueryParam("commodity")
	batch := c.QueryParam("batch")
	number := c.QueryParam("number")
	status := c.QueryParam("status")

	if commodity != "" {
		filter.Commodity = commodity
	}

	if status != "" {
		filter.Status = status
		if !util.CheckStringOnArray([]string{constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse}, status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, %s, dan %s", constant.TreatmentRecordStatusApproved, constant.TreatmentRecordStatusPending, constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusWaitingResponse)
		}
	}

	if batch != "" {
		filter.Batch = batch
	}

	if number != "" {
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			return FilterQuery{}, errors.New("number harus berupa angka")
		}

		filter.Number = numberInt
	}

	return filter, nil
}
