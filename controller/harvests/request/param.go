package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"fmt"

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
