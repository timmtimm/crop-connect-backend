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

	return filter, nil
}
