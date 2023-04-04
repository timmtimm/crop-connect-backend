package request

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Country     string
	Province    string
	Regency     string
	District    string
	Subdistrict string
}

func QueryParamValidation(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{}

	country := c.QueryParam("country")
	province := c.QueryParam("province")
	regency := c.QueryParam("regency")
	district := c.QueryParam("district")

	if country == "" {
		return FilterQuery{}, errors.New("country wajib diisi")
	}

	filter.Country = country
	if province != "" {
		filter.Province = province
	}

	if regency != "" {
		filter.Regency = regency
	}

	if district != "" {
		filter.District = district
	}

	return filter, nil
}
