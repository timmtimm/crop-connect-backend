package main

import (
	"fmt"

	_route "marketplace-backend/app/route"
	_driver "marketplace-backend/driver"
	_mongo "marketplace-backend/driver/mongo"
	_util "marketplace-backend/util"

	_commodityUseCase "marketplace-backend/business/commodities"
	_proposalUseCase "marketplace-backend/business/proposals"
	_transactionUseCase "marketplace-backend/business/transactions"
	_userUseCase "marketplace-backend/business/users"

	_commodityController "marketplace-backend/controller/commodities"
	_proposalController "marketplace-backend/controller/proposals"
	_transactionController "marketplace-backend/controller/transactions"
	_userController "marketplace-backend/controller/users"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))

	userRepository := _driver.NewUserRepository(database)
	commodityRepository := _driver.NewCommodityRepository(database)
	proposalRepository := _driver.NewProposalRepository(database)
	transactionRepository := _driver.NewTransactionRepository(database)

	userUseCase := _userUseCase.NewUserUseCase(userRepository)
	commodityUsecase := _commodityUseCase.NewCommodityUseCase(commodityRepository)
	proposalUseCase := _proposalUseCase.NewProposalUseCase(proposalRepository, commodityRepository)
	transactionUseCase := _transactionUseCase.NewTransactionUseCase(transactionRepository, commodityRepository, proposalRepository)

	userController := _userController.NewUserController(userUseCase)
	commodityController := _commodityController.NewCommodityController(commodityUsecase, userUseCase, proposalUseCase)
	proposalController := _proposalController.NewProposalController(proposalUseCase, commodityUsecase)
	transactionController := _transactionController.NewTransactionController(transactionUseCase, proposalUseCase, commodityUsecase, userUseCase)

	routeController := _route.ControllerList{
		UserController:        userController,
		CommodityController:   commodityController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
	}

	routeController.Init(e)

	appPort := fmt.Sprintf(":%s", _util.GetConfig("APP_PORT"))
	e.Logger.Fatal(e.Start(appPort))
}
