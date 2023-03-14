package route

import (
	_middleware "marketplace-backend/app/middleware"
	"marketplace-backend/controller/commodities"
	"marketplace-backend/controller/proposals"
	"marketplace-backend/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware    echo.MiddlewareFunc
	UserController      *users.Controller
	CommodityController *commodities.Controller
	ProposalController  *proposals.Controller
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
	commodity.GET("/page/:page", cl.CommodityController.GetForBuyer)
	commodity.POST("", cl.CommodityController.Create, _middleware.CheckOneRole("farmer"))
	commodity.GET("/:commodity-id", cl.CommodityController.GetByID)
	commodity.PUT("/:commodity-id", cl.CommodityController.Update, _middleware.CheckOneRole("farmer"))
	commodity.DELETE("/:commodity-id", cl.CommodityController.Delete, _middleware.CheckOneRole("farmer"))

	proposal := apiV1.Group("/proposal")
	proposal.GET("/:commodity-id", cl.ProposalController.GetByCommodityIDForBuyer)
	proposal.POST("/:commodity-id", cl.ProposalController.Create, _middleware.CheckOneRole("farmer"))
	proposal.PUT("/:proposal-id", cl.ProposalController.Update, _middleware.CheckOneRole("farmer"))
	proposal.DELETE("/:proposal-id", cl.ProposalController.Delete, _middleware.CheckOneRole("farmer"))
}
