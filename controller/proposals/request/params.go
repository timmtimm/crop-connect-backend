package request

import (
	"errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

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
