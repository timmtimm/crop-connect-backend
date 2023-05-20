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

func QueryParamValidationForQuery(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Country:  c.QueryParam("country"),
		Province: c.QueryParam("province"),
		Regency:  c.QueryParam("regency"),
		District: c.QueryParam("district"),
	}

	return filter, nil
}

func QueryParamValidation(c echo.Context, param []string) (map[string]string, error) {
	result := map[string]string{}
	for _, v := range param {
		value := c.QueryParam(v)

		if value == "" {
			return map[string]string{}, errors.New("parameter " + v + " tidak boleh kosong")
		}

		result[v] = value
	}

	return result, nil
}
