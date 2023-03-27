package request

import (
	"errors"
	"fmt"
	"marketplace-backend/constant"
	"marketplace-backend/util"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	Commodity string
	FarmerID  primitive.ObjectID
	BuyerID   primitive.ObjectID
	Status    string
	StartDate primitive.DateTime
	EndDate   primitive.DateTime
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}

	commodity := c.QueryParam("commodity")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	if commodity != "" {
		filter.Commodity = commodity
	}

	if status != "" {
		filter.Status = status
		if !util.CheckStringOnArray([]string{constant.TransactionStatusPending, constant.TransactionStatusAccepted, constant.ProposalStatusRejected}, status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.TransactionStatusPending, constant.TransactionStatusAccepted, constant.ProposalStatusRejected)
		}
	}

	if startDate != "" {
		date, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return FilterQuery{}, errors.New("startDate harus berupa tanggal")
		}

		filter.StartDate = primitive.NewDateTimeFromTime(date)
	}

	if endDate != "" {
		date, err := time.Parse("2006-01-02", c.QueryParam("endDate"))
		if err != nil {
			return FilterQuery{}, errors.New("endDate harus berupa tanggal")
		}

		filter.EndDate = primitive.NewDateTimeFromTime(date)
	}

	if startDate != "" && endDate != "" {
		if filter.StartDate > filter.EndDate {
			return FilterQuery{}, errors.New("startDate tidak boleh dari endDate")
		}
	}

	return filter, nil
}
