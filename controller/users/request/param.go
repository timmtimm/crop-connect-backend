package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"fmt"

	"github.com/labstack/echo/v4"
)

type FilterQuery struct {
	Name        string
	Email       string
	PhoneNumber string
	Role        string
}

func QueryParamValidationForAdmin(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Name:        c.QueryParam("name"),
		Email:       c.QueryParam("email"),
		PhoneNumber: c.QueryParam("phoneNumber"),
		Role:        c.QueryParam("role"),
	}

	if filter.Role != "" {
		if !util.CheckStringOnArray([]string{constant.RoleAdmin, constant.RoleValidator, constant.RoleFarmer, constant.RoleBuyer}, filter.Role) {
			return FilterQuery{}, fmt.Errorf("role tersedia hanya %s, %s, %s, %s", constant.RoleAdmin, constant.RoleValidator, constant.RoleFarmer, constant.RoleBuyer)
		}
	}

	return filter, nil
}

func QueryParamValidationForSearchFarmer(c echo.Context) (FilterQuery, error) {
	filter := FilterQuery{
		Name: c.QueryParam("name"),
	}

	return filter, nil
}
