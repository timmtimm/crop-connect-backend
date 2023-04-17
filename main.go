package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	_middleware "crop_connect/app/middleware"
	_route "crop_connect/app/route"
	_driver "crop_connect/driver"
	_mongo "crop_connect/driver/mongo"
	"crop_connect/helper/cloudinary"
	"crop_connect/seeds"
	_util "crop_connect/util"

	_batchUseCase "crop_connect/business/batchs"
	_commodityUseCase "crop_connect/business/commodities"
	_harvestUseCase "crop_connect/business/harvests"
	_proposalUseCase "crop_connect/business/proposals"
	_regionUseCase "crop_connect/business/regions"
	_transactionUseCase "crop_connect/business/transactions"
	_treatmentRecordUseCase "crop_connect/business/treatment_records"
	_userUseCase "crop_connect/business/users"

	_batchController "crop_connect/controller/batchs"
	_commodityController "crop_connect/controller/commodities"
	_harvestController "crop_connect/controller/harvests"
	_proposalController "crop_connect/controller/proposals"
	_regionController "crop_connect/controller/regions"
	_transactionController "crop_connect/controller/transactions"
	_treatmentRecordController "crop_connect/controller/treatment_records"
	_userController "crop_connect/controller/users"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	fmt.Println("Initializing echo...")
	e := echo.New()

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
			autoTLSManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      autocert.DirCache("/var/www/.cache"),
				HostPolicy: autocert.HostWhitelist(_util.ResontructeDomainName()...),
			}

			s := http.Server{
				Addr:    appPort,
				Handler: e,
				TLSConfig: &tls.Config{
					GetCertificate: autoTLSManager.GetCertificate,
					NextProtos:     []string{acme.ALPNProto},
				},
			}

			e.Logger.Fatal(s.ListenAndServeTLS("", ""))
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
