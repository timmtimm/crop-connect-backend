package request

import (
	"errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	Name     string
	Farmer   string
	FarmerID primitive.ObjectID
	MinPrice int
	MaxPrice int
}

var err error

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Name:   c.QueryParam("name"),
		Farmer: c.QueryParam("farmer"),
	}

	if minPrice := c.QueryParam("minPrice"); minPrice != "" {
		filter.MinPrice, err = strconv.Atoi(minPrice)
		if err != nil {
			return FilterQuery{}, errors.New("harga minimal harus berupa angka")
		}
	}

	if maxPrice := c.QueryParam("maxPrice"); maxPrice != "" {
		filter.MaxPrice, err = strconv.Atoi(maxPrice)
		if err != nil {
			return FilterQuery{}, errors.New("harga maksimal harus berupa angka")
		}
	}

	if farmerID := c.QueryParam("farmerID"); farmerID != "" {
		filter.FarmerID, err = primitive.ObjectIDFromHex(farmerID)
		if err != nil {
			return FilterQuery{}, errors.New("farmerID harus berupa hex")
		}
	}

	return filter, nil
}

type QueryRegion struct {
	Province string
	Regency  string
	District string
	RegionID primitive.ObjectID // for subdistrict cases, you can get it using regionID
}

func QueryValidationForRegion(c echo.Context) (QueryRegion, error) {
	query := QueryRegion{
		Province: c.QueryParam("province"),
		Regency:  c.QueryParam("regency"),
		District: c.QueryParam("district"),
	}

	if query.District != "" {
		if query.Province == "" || query.Regency == "" {
			return QueryRegion{}, errors.New("harus menyertakan parameter province dan regency")
		}
	} else if query.Regency != "" {
		if query.Province == "" {
			return QueryRegion{}, errors.New("harus menyertakan parameter province")
		}
	}

	if regionID := c.QueryParam("regionID"); regionID != "" {
		query.RegionID, err = primitive.ObjectIDFromHex(regionID)
		if err != nil {
			return QueryRegion{}, errors.New("regionID harus berupa hex")
		}
	}

	return query, nil
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
