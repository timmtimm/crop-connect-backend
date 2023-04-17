package request

import (
	"fmt"
	"marketplace-backend/constant"
	"marketplace-backend/util"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Commodity string
	Name      string
	Status    string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}

	commodity := c.QueryParam("commodity")
	name := c.QueryParam("name")
	status := c.QueryParam("status")

	if commodity != "" {
		filter.Commodity = commodity
	}

	if status != "" {
		filter.Status = status
		if !util.CheckStringOnArray([]string{constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel}, status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel)
		}
	}

	if name != "" {
		filter.Name = name
	}

	return filter, nil
}
