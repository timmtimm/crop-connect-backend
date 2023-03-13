package request

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Name     string
	Farmer   string
	MinPrice int
	MaxPrice int
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}
	var err error

	if c.QueryParam("name") != "" {
		filter.Name = c.QueryParam("name")
	}

	if c.QueryParam("farmer") != "" {
		filter.Farmer = c.QueryParam("farmer")
	}

	if c.QueryParam("minPrice") != "" {
		filter.MinPrice, err = strconv.Atoi(c.QueryParam("minPrice"))
		if err != nil {
			return FilterQuery{}, errors.New("harga minimal harus berupa angka")
		}
	}

	if c.QueryParam("maxPrice") != "" {
		filter.MaxPrice, err = strconv.Atoi(c.QueryParam("maxPrice"))
		if err != nil {
			return FilterQuery{}, errors.New("harga maksimal harus berupa angka")
		}
	}

	return filter, nil
}
