package route

import (
	_middleware "marketplace-backend/app/middleware"
	"marketplace-backend/constant"
	"marketplace-backend/controller/commodities"
	"marketplace-backend/controller/proposals"
	"marketplace-backend/controller/transactions"
	"marketplace-backend/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware      echo.MiddlewareFunc
	UserController        *users.Controller
	CommodityController   *commodities.Controller
	ProposalController    *proposals.Controller
	TransactionController *transactions.Controller
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
	user.POST("/register-validator", cl.UserController.RegisterValidator, _middleware.CheckOneRole(constant.RoleAdmin))
	user.POST("/login", cl.UserController.Login)
	user.GET("/profile", cl.UserController.GetProfile, _middleware.Authenticated())
	user.PUT("/profile", cl.UserController.UpdateProfile, _middleware.Authenticated())
	user.GET("/find-farmer/:farmer-name", cl.UserController.GetFarmerByName)

	commodity := apiV1.Group("/commodity")
	commodity.GET("/page/:page", cl.CommodityController.GetForBuyer)
	commodity.POST("", cl.CommodityController.Create, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.GET("/:commodity-id", cl.CommodityController.GetByID)
	commodity.PUT("/:commodity-id", cl.CommodityController.Update, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.DELETE("/:commodity-id", cl.CommodityController.Delete, _middleware.CheckOneRole(constant.RoleFarmer))

	proposal := apiV1.Group("/proposal")
	proposal.GET("/:commodity-id", cl.ProposalController.GetByCommodityIDForBuyer)
	proposal.POST("/:commodity-id", cl.ProposalController.Create, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.PUT("/:proposal-id", cl.ProposalController.Update, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.DELETE("/:proposal-id", cl.ProposalController.Delete, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.PUT("/validate/:proposal-id", cl.ProposalController.ValidateByValidator, _middleware.CheckOneRole(constant.RoleValidator))

	transaction := apiV1.Group("/transaction")
	transaction.POST("/:proposal-id", cl.TransactionController.Create, _middleware.CheckOneRole(constant.RoleBuyer))
	transaction.GET("/:page", cl.TransactionController.GetUserTransaction, _middleware.CheckManyRole([]string{constant.RoleBuyer, constant.RoleFarmer}))
}
