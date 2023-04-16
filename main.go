package main

import (
	"context"
	"fmt"
	"time"

	_middleware "marketplace-backend/app/middleware"
	_route "marketplace-backend/app/route"
	_driver "marketplace-backend/driver"
	_mongo "marketplace-backend/driver/mongo"
	"marketplace-backend/helper/cloudinary"
	"marketplace-backend/seeds"
	_util "marketplace-backend/util"

	_batchUseCase "marketplace-backend/business/batchs"
	_commodityUseCase "marketplace-backend/business/commodities"
	_harvestUseCase "marketplace-backend/business/harvests"
	_proposalUseCase "marketplace-backend/business/proposals"
	_regionUseCase "marketplace-backend/business/regions"
	_transactionUseCase "marketplace-backend/business/transactions"
	_treatmentRecordUseCase "marketplace-backend/business/treatment_records"
	_userUseCase "marketplace-backend/business/users"

	_batchController "marketplace-backend/controller/batchs"
	_commodityController "marketplace-backend/controller/commodities"
	_harvestController "marketplace-backend/controller/harvests"
	_proposalController "marketplace-backend/controller/proposals"
	_regionController "marketplace-backend/controller/regions"
	_transactionController "marketplace-backend/controller/transactions"
	_treatmentRecordController "marketplace-backend/controller/treatment_records"
	_userController "marketplace-backend/controller/users"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	fmt.Println("Initializing echo...")
	e := echo.New()

	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	fmt.Println("Initializing database...")
	database := _mongo.Init(_util.GetConfig("DB_NAME"))
	cloudinary := cloudinary.Init(_util.GetConfig("CLOUDINARY_UPLOAD_FOLDER"))

	fmt.Println("Initializing repositories...")
	userRepository := _driver.NewUserRepository(database)
	commodityRepository := _driver.NewCommodityRepository(database)
	proposalRepository := _driver.NewProposalRepository(database)
	transactionRepository := _driver.NewTransactionRepository(database)
	batchRepository := _driver.NewBatchRepository(database)
	treatmentRecordRepository := _driver.NewTreatmentRecordRepository(database)
	harvestRepository := _driver.NewHarvestRepository(database)
	regionRepository := _driver.NewRegionRepository(database)

	fmt.Println("Initializing usecases...")
	userUseCase := _userUseCase.NewUserUseCase(userRepository, regionRepository)
	commodityUsecase := _commodityUseCase.NewCommodityUseCase(commodityRepository, cloudinary)
	proposalUseCase := _proposalUseCase.NewProposalUseCase(proposalRepository, commodityRepository, regionRepository)
	transactionUseCase := _transactionUseCase.NewTransactionUseCase(transactionRepository, commodityRepository, proposalRepository)
	batchUseCase := _batchUseCase.NewBatchUseCase(batchRepository, transactionRepository, proposalRepository, commodityRepository)
	treatmentRecordUseCase := _treatmentRecordUseCase.NewTreatmentRecordUseCase(treatmentRecordRepository, batchRepository, transactionRepository, proposalRepository, commodityRepository, cloudinary)
	harvestUseCase := _harvestUseCase.NewHarvestUseCase(harvestRepository, batchRepository, treatmentRecordRepository, transactionRepository, proposalRepository, commodityRepository, cloudinary)
	regionUseCase := _regionUseCase.NewRegionUseCase(regionRepository)

	fmt.Println("Initializing controllers...")
	userController := _userController.NewUserController(userUseCase, regionUseCase)
	commodityController := _commodityController.NewCommodityController(commodityUsecase, userUseCase, proposalUseCase, regionUseCase)
	proposalController := _proposalController.NewProposalController(proposalUseCase, commodityUsecase)
	transactionController := _transactionController.NewTransactionController(transactionUseCase, proposalUseCase, commodityUsecase, userUseCase, batchUseCase, regionUseCase)
	batchController := _batchController.NewBatchController(batchUseCase, transactionUseCase, proposalUseCase, commodityUsecase, userUseCase, regionUseCase)
	treatmentRecordController := _treatmentRecordController.NewTreatmentRecordController(treatmentRecordUseCase, batchUseCase, transactionUseCase, proposalUseCase, commodityUsecase, userUseCase, regionUseCase)
	harvestController := _harvestController.NewHarvestController(harvestUseCase, batchUseCase, transactionUseCase, proposalUseCase, commodityUsecase, userUseCase, regionUseCase)
	regionController := _regionController.NewRegionController(regionUseCase)

	seeds.SeedDatabase(database, regionUseCase)

	fmt.Println("Initializing middlewares...")
	_middleware.InitLogger(e)
	_middleware.InitCORS(e)

	fmt.Println("Initializing routes...")
	routeController := _route.ControllerList{
		UserController:            userController,
		CommodityController:       commodityController,
		ProposalController:        proposalController,
		TransactionController:     transactionController,
		BatchController:           batchController,
		TreatmentRecordController: treatmentRecordController,
		HarvestController:         harvestController,
		RegionController:          regionController,
	}
	routeController.Init(e)

	fmt.Println("Starting server...")

	go func() {
		appPort := fmt.Sprintf(":%s", _util.GetConfig("APP_PORT"))
		if _util.GetConfig("APP_ENV") == "development" {
			e.Logger.Fatal(e.Start(appPort))
		} else {
			e.Logger.Fatal(e.StartAutoTLS(appPort))
		}
	}()

	wait := _util.GracefulShutdown(context.Background(), 2*time.Second, map[string]_util.Operation{
		"database": func(ctx context.Context) error {
			return _mongo.Close(database)
		},
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(context.Background())
		},
	})

	<-wait
}
