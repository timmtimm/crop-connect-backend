package request

import (
	"fmt"
	"marketplace-backend/constant"
	"marketplace-backend/util"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Commodity string
	Proposal  string
	Batch     string
	Status    string
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}

	commodity := c.QueryParam("commodity")
	proposal := c.QueryParam("proposal")
	batch := c.QueryParam("batch")
	status := c.QueryParam("status")

	if commodity != "" {
		filter.Commodity = commodity
	}

	if status != "" {
		filter.Status = status
		if !util.CheckStringOnArray([]string{constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision}, status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision)
		}
	}

	if proposal != "" {
		filter.Proposal = proposal
	}

	if batch != "" {
		filter.Batch = batch
	}

	return filter, nil
}
