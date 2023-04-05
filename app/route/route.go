package route

import (
	_middleware "marketplace-backend/app/middleware"
	"marketplace-backend/constant"
	"marketplace-backend/controller/batchs"
	"marketplace-backend/controller/commodities"
	"marketplace-backend/controller/harvests"
	"marketplace-backend/controller/proposals"
	"marketplace-backend/controller/regions"
	"marketplace-backend/controller/transactions"
	treatmentRecords "marketplace-backend/controller/treatment_records"
	"marketplace-backend/controller/users"
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
	commodity.GET("/farmer", cl.CommodityController.GetForFarmer, _middleware.CheckOneRole(constant.RoleFarmer))
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
	transaction.GET("/page/:page", cl.TransactionController.GetUserTransactionWithPagination, _middleware.CheckManyRole([]string{constant.RoleBuyer, constant.RoleFarmer}))
	transaction.POST("/:proposal-id", cl.TransactionController.Create, _middleware.CheckOneRole(constant.RoleBuyer))
	transaction.PUT("/:transaction-id", cl.TransactionController.MakeDecision, _middleware.CheckOneRole(constant.RoleFarmer))

	batch := apiV1.Group("/batch")
	batch.GET("/page/:page", cl.BatchController.GetFarmerBatch, _middleware.CheckOneRole(constant.RoleFarmer))
	batch.GET("/commodity/:commodity-id", cl.BatchController.GetByCommodityID)
	// batch.PUT("/cancel/:batch-id", cl.BatchController.Cancel, _middleware.CheckOneRole(constant.RoleFarmer))

	treatmentRecord := apiV1.Group("/treatment-record")
	treatmentRecord.GET("/page/:page", cl.TreatmentRecordController.GetByPaginationAndQuery, _middleware.CheckManyRole([]string{constant.RoleFarmer, constant.RoleValidator}))
	treatmentRecord.POST("/:batch-id", cl.TreatmentRecordController.RequestToFarmer, _middleware.CheckOneRole(constant.RoleValidator))
	treatmentRecord.PUT("/:treatment-record-id", cl.TreatmentRecordController.FillTreatmentRecord, _middleware.CheckOneRole(constant.RoleFarmer))
	treatmentRecord.PUT("/validate/:treatment-record-id", cl.TreatmentRecordController.Validate, _middleware.CheckOneRole(constant.RoleValidator))
	treatmentRecord.PUT("/note/:treatment-record-id", cl.TreatmentRecordController.UpdateNotes, _middleware.CheckOneRole(constant.RoleValidator))

	harvest := apiV1.Group("/harvest")
	harvest.GET("/page/:page", cl.HarvestController.GetByPaginationAndQuery, _middleware.CheckManyRole([]string{constant.RoleFarmer, constant.RoleValidator}))
	harvest.GET("/:batch-id", cl.HarvestController.GetByBatchID)
	harvest.POST("/:batch-id", cl.HarvestController.SubmitHarvest, _middleware.CheckOneRole(constant.RoleFarmer))
	harvest.PUT("/validate/:harvest-id", cl.HarvestController.Validate, _middleware.CheckOneRole(constant.RoleValidator))
}
