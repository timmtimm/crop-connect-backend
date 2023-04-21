package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"fmt"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Commodity string
	Name      string
	Status    string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Commodity: c.QueryParam("commodity"),
		Name:      c.QueryParam("name"),
		Status:    c.QueryParam("status"),
	}

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel)
		}
	}

	return filter, nil
}
