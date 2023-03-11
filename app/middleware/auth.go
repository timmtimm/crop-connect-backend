package middleware

import (
	"marketplace-backend/helper"
	"marketplace-backend/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := util.GetPayloadFromToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, helper.BaseResponse{
					Status:  http.StatusBadRequest,
					Message: "token tidak valid",
					Data:    nil,
				})
			}

			return next(c)
		}
	}
}

func Check(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := util.GetPayloadFromToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, helper.BaseResponse{
					Status:  http.StatusBadRequest,
					Message: "token tidak valid",
					Data:    nil,
				})
			}

			for _, role := range roles {
				if token.Role == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
				Data:    nil,
			})
		}
	}
}
