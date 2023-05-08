package route

import (
	_middleware "crop_connect/app/middleware"
	"crop_connect/constant"
	"crop_connect/controller/batchs"
	"crop_connect/controller/commodities"
	forgotPassword "crop_connect/controller/forgot_password"
	"crop_connect/controller/harvests"
	"crop_connect/controller/proposals"
	"crop_connect/controller/regions"
	"crop_connect/controller/transactions"
	treatmentRecords "crop_connect/controller/treatment_records"
	"crop_connect/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	UserController            *users.Controller
	CommodityController       *commodities.Controller
	ProposalController        *proposals.Controller
	TransactionController     *transactions.Controller
	BatchController           *batchs.Controller
	TreatmentRecordController *treatmentRecords.Controller
	HarvestController         *harvests.Controller
	RegionController          *regions.Controller
	ForgotPasswordController  *forgotPassword.Controller
}

func (ctrl *ControllerList) Init(e *echo.Echo) {
	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World!",
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome to Echo!</h1>
			<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
		`)
	})

	// API V1
	apiV1 := e.Group("/api/v1")

	user := apiV1.Group("/user")
	user.POST("/register", ctrl.UserController.Register)
	user.POST("/register-validator", ctrl.UserController.RegisterValidator, _middleware.CheckOneRole(constant.RoleAdmin))
	user.POST("/login", ctrl.UserController.Login)
	user.GET("/profile", ctrl.UserController.GetProfile, _middleware.Authenticated())
	user.PUT("/profile", ctrl.UserController.UpdateProfile, _middleware.Authenticated())
	user.GET("", ctrl.UserController.GetByPaginationAndQueryForAdmin, _middleware.CheckOneRole(constant.RoleAdmin))
	user.GET("/farmer", ctrl.UserController.GetFarmerByPaginationAndQueryForBuyer)
	user.GET("/farmer/:farmer-id", ctrl.UserController.GetFarmerByIDForBuyer)
	user.PUT("/change-password", ctrl.UserController.UpdatePassword, _middleware.Authenticated())

	forgotPassword := user.Group("/forgot-password")
	forgotPassword.POST("", ctrl.ForgotPasswordController.Generate)
	forgotPassword.GET("/:token", ctrl.ForgotPasswordController.ValidateToken)
	forgotPassword.PUT("/:token", ctrl.ForgotPasswordController.ResetPassword)

	commodity := apiV1.Group("/commodity")
	commodity.GET("", ctrl.CommodityController.GetForBuyer)
	commodity.GET("/farmer", ctrl.CommodityController.GetForFarmer, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.POST("", ctrl.CommodityController.Create, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.GET("/:commodity-id", ctrl.CommodityController.GetByID)
	commodity.PUT("/:commodity-id", ctrl.CommodityController.Update, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.DELETE("/:commodity-id", ctrl.CommodityController.Delete, _middleware.CheckOneRole(constant.RoleFarmer))
	commodity.GET("/statistic-total", ctrl.CommodityController.CountTotalCommodity, _middleware.CheckOneRole(constant.RoleAdmin))

	proposal := apiV1.Group("/proposal")
	proposal.GET("/:commodity-id", ctrl.ProposalController.GetByCommodityIDForBuyer)
	proposal.POST("/:commodity-id", ctrl.ProposalController.Create, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.PUT("/:proposal-id", ctrl.ProposalController.Update, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.DELETE("/:proposal-id", ctrl.ProposalController.Delete, _middleware.CheckOneRole(constant.RoleFarmer))
	proposal.PUT("/validate/:proposal-id", ctrl.ProposalController.ValidateByValidator, _middleware.CheckOneRole(constant.RoleValidator))

	transaction := apiV1.Group("/transaction")
	transaction.GET("", ctrl.TransactionController.GetUserTransactionWithPagination, _middleware.CheckManyRole([]string{constant.RoleBuyer, constant.RoleFarmer}))
	transaction.POST("/:proposal-id", ctrl.TransactionController.Create, _middleware.CheckOneRole(constant.RoleBuyer))
	transaction.PUT("/:transaction-id", ctrl.TransactionController.MakeDecision, _middleware.CheckOneRole(constant.RoleFarmer))
	transaction.GET("/statistic", ctrl.TransactionController.StatisticByYear, _middleware.CheckManyRole([]string{constant.RoleAdmin, constant.RoleFarmer}))
	transaction.GET("/statistic-province", ctrl.TransactionController.StatisticTopProvince, _middleware.CheckOneRole(constant.RoleAdmin))
	transaction.GET("/statistic-commodity", ctrl.TransactionController.StatisticTopCommodity, _middleware.CheckManyRole([]string{constant.RoleAdmin, constant.RoleFarmer}))

	batch := apiV1.Group("/batch")
	batch.GET("", ctrl.BatchController.GetFarmerBatch, _middleware.CheckOneRole(constant.RoleFarmer))
	batch.GET("/:commodity-id", ctrl.BatchController.GetByCommodityID)
	batch.GET("/statistic-total", ctrl.BatchController.CountByYear, _middleware.CheckOneRole(constant.RoleAdmin))
	// batch.PUT("/cancel/:batch-id", ctrl.BatchController.Cancel, _middleware.CheckOneRole(constant.RoleFarmer))

	treatmentRecord := apiV1.Group("/treatment-record")
	treatmentRecord.GET("", ctrl.TreatmentRecordController.GetByPaginationAndQuery, _middleware.CheckManyRole([]string{constant.RoleFarmer, constant.RoleValidator}))
	treatmentRecord.POST("/:batch-id", ctrl.TreatmentRecordController.RequestToFarmer, _middleware.CheckOneRole(constant.RoleValidator))
	treatmentRecord.PUT("/:treatment-record-id", ctrl.TreatmentRecordController.FillTreatmentRecord, _middleware.CheckOneRole(constant.RoleFarmer))
	treatmentRecord.PUT("/validate/:treatment-record-id", ctrl.TreatmentRecordController.Validate, _middleware.CheckOneRole(constant.RoleValidator))
	treatmentRecord.PUT("/note/:treatment-record-id", ctrl.TreatmentRecordController.UpdateNotes, _middleware.CheckOneRole(constant.RoleValidator))
	treatmentRecord.GET("/statistic-total", ctrl.TreatmentRecordController.CountByYear, _middleware.CheckOneRole(constant.RoleValidator))

	harvest := apiV1.Group("/harvest")
	harvest.GET("", ctrl.HarvestController.GetByPaginationAndQuery, _middleware.CheckManyRole([]string{constant.RoleFarmer, constant.RoleValidator}))
	harvest.GET("/:batch-id", ctrl.HarvestController.GetByBatchID)
	harvest.POST("/:batch-id", ctrl.HarvestController.SubmitHarvest, _middleware.CheckOneRole(constant.RoleFarmer))
	harvest.PUT("/validate/:harvest-id", ctrl.HarvestController.Validate, _middleware.CheckOneRole(constant.RoleValidator))

	region := apiV1.Group("/region")
	region.GET("/province", ctrl.RegionController.GetByCountry)
	region.GET("/regency", ctrl.RegionController.GetByProvince)
	region.GET("/district", ctrl.RegionController.GetByRegency)
	region.GET("/sub-district", ctrl.RegionController.GetByDistrict)

}
