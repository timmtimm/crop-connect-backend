package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	FarmerID    primitive.ObjectID
	CommodityID primitive.ObjectID
	Name        string
	Status      string
}

func QueryParamValidationForBuyer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Name:   c.QueryParam("name"),
		Status: c.QueryParam("status"),
	}

	if farmerID := c.QueryParam("farmerID"); farmerID != "" {
		farmerID, err := primitive.ObjectIDFromHex(farmerID)
		if err != nil {
			return FilterQuery{}, errors.New("farmerID harus berupa hex")
		}

		filter.FarmerID = farmerID
	}

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.BatchStatusPlanting, constant.BatchStatusHarvest, constant.BatchStatusCancel)
		}
	}

	if commodity := c.QueryParam("commodityID"); commodity != "" {
		commodityID, err := primitive.ObjectIDFromHex(commodity)
		if err != nil {
			return FilterQuery{}, errors.New("commodityID harus berupa hex")
		}

		filter.CommodityID = commodityID
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
