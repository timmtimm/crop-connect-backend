package request

import (
	"errors"
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
	filter := FilterQuery{
		Commodity: c.QueryParam("commodity"),
		Status:    c.QueryParam("status"),
	}

	if startDate := c.QueryParam("startDate"); startDate != "" {
		date, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return FilterQuery{}, errors.New("startDate harus berupa tanggal")
		}

		filter.StartDate = primitive.NewDateTimeFromTime(date)
	}

	if endDate := c.QueryParam("endDate"); endDate != "" {
		date, err := time.Parse("2006-01-02", c.QueryParam("endDate"))
		if err != nil {
			return FilterQuery{}, errors.New("endDate harus berupa tanggal")
		}

		filter.EndDate = primitive.NewDateTimeFromTime(date)
	}

	if filter.StartDate > filter.EndDate {
		return FilterQuery{}, errors.New("startDate tidak boleh dari endDate")
	}

	return filter, nil
}
