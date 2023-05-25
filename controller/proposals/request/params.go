package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	CommodityID primitive.ObjectID
	FarmerID    primitive.ObjectID
	Name        string
	Status      string
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Name: c.QueryParam("name"),
	}

	if status := c.QueryParam("status"); status != "" {
		isAvailable := util.CheckStringOnArray([]string{
			constant.ProposalStatusApproved,
			constant.ProposalStatusRejected,
			constant.ProposalStatusPending,
		}, status)

		if !isAvailable {
			return FilterQuery{}, errors.New("status tidak tersedia")
		}

		filter.Status = status
	}

	if commodityID := c.QueryParam("commodityID"); commodityID != "" {
		id, err := primitive.ObjectIDFromHex(commodityID)
		if err != nil {
			return FilterQuery{}, errors.New("commodityID harus berupa hex")
		}

		filter.CommodityID = id
	}

	if farmerID := c.QueryParam("farmerID"); farmerID != "" {
		id, err := primitive.ObjectIDFromHex(farmerID)
		if err != nil {
			return FilterQuery{}, errors.New("farmerID harus berupa hex")
		}

		filter.FarmerID = id
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
