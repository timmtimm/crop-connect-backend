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

	name := c.QueryParam("name")
	farmer := c.QueryParam("farmer")
	minPrice := c.QueryParam("minPrice")
	maxPrice := c.QueryParam("maxPrice")

	if name != "" {
		filter.Name = name
	}

	if farmer != "" {
		filter.Farmer = farmer
	}

	if minPrice != "" {
		filter.MinPrice, err = strconv.Atoi(minPrice)
		if err != nil {
			return FilterQuery{}, errors.New("harga minimal harus berupa angka")
		}
	}

	if maxPrice != "" {
		filter.MaxPrice, err = strconv.Atoi(maxPrice)
		if err != nil {
			return FilterQuery{}, errors.New("harga maksimal harus berupa angka")
		}
	}

	return filter, nil
}
