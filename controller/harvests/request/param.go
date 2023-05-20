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
	Proposal  string
	Batch     string
	Status    string
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Commodity: c.QueryParam("commodity"),
		Proposal:  c.QueryParam("proposal"),
		Batch:     c.QueryParam("batch"),
		Status:    c.QueryParam("status"),
	}

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision)
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
