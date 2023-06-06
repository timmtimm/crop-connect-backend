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
	CommodityID primitive.ObjectID
	BatchID     primitive.ObjectID
	Status      string
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Status: c.QueryParam("status"),
	}

	var err error

	if filter.Status != "" {
		if !util.CheckStringOnArray([]string{constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision}, filter.Status) {
			return FilterQuery{}, fmt.Errorf("status tersedia hanya %s, %s, dan %s", constant.HarvestStatusApproved, constant.HarvestStatusPending, constant.HarvestStatusRevision)
		}
	}

	if commodityID := c.QueryParam("commodityID"); commodityID != "" {
		filter.CommodityID, err = primitive.ObjectIDFromHex(commodityID)
		if err != nil {
			return FilterQuery{}, errors.New("commodityID harus berupa hex")
		}
	}

	if batchID := c.QueryParam("batchID"); batchID != "" {
		filter.BatchID, err = primitive.ObjectIDFromHex(batchID)
		if err != nil {
			return FilterQuery{}, errors.New("batchID harus berupa hex")
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
