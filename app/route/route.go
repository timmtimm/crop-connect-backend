package route

import (
	_middleware "marketplace-backend/app/middleware"
	"marketplace-backend/controller/commodities"
	"marketplace-backend/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware    echo.MiddlewareFunc
	UserController      *users.Controller
	CommodityController *commodities.Controller
}

func (cl *ControllerList) Init(e *echo.Echo) {
	_middleware.InitLogger(e)

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World!",
		})
	})

	// API V1
	apiV1 := e.Group("/api/v1")

	user := apiV1.Group("/user")
	user.POST("/register", cl.UserController.Register)
	user.POST("/login", cl.UserController.Login)
	user.GET("/profile", cl.UserController.GetProfile, _middleware.Authenticated())
	user.PUT("/profile", cl.UserController.UpdateProfile, _middleware.Authenticated())

	commodity := apiV1.Group("/commodity")
	commodity.GET("/:page", cl.CommodityController.GetForBuyer)
	commodity.POST("", cl.CommodityController.Create, _middleware.CheckOneRole("farmer"))
}
