package request

import (
	"crop_connect/constant"
	"crop_connect/util"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterQuery struct {
	Name        string
	Email       string
	PhoneNumber string
	Role        string
	Province    string
	Regency     string
	District    string
	RegionID    primitive.ObjectID
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
		Name:     c.QueryParam("name"),
		Province: c.QueryParam("province"),
		Regency:  c.QueryParam("regency"),
		District: c.QueryParam("district"),
	}

	var err error

	if regionID := c.QueryParam("regionID"); regionID != "" {
		filter.RegionID, err = primitive.ObjectIDFromHex(regionID)
		if err != nil {
			return FilterQuery{}, fmt.Errorf("regionID tidak valid")
		}
	}

	return filter, nil
}
