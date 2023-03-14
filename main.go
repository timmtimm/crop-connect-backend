package main

import (
	"fmt"

	_route "marketplace-backend/app/route"
	_driver "marketplace-backend/driver"
	_mongo "marketplace-backend/driver/mongo"
	_util "marketplace-backend/util"

	_commodityUseCase "marketplace-backend/business/commodities"
	_proposalUseCase "marketplace-backend/business/proposals"
	_userUseCase "marketplace-backend/business/users"

	_commodityController "marketplace-backend/controller/commodities"
	_proposalController "marketplace-backend/controller/proposals"
	_userController "marketplace-backend/controller/users"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))

	userRepository := _driver.NewUserRepository(database)
	commodityRepository := _driver.NewCommodityRepository(database)
	proposalRepository := _driver.NewProposalRepository(database)

	userUsecase := _userUseCase.NewUserUseCase(userRepository)
	commodityUsecase := _commodityUseCase.NewCommodityUseCase(commodityRepository)
	proposalUsecase := _proposalUseCase.NewProposalUseCase(proposalRepository, commodityRepository)

	userController := _userController.NewUserController(userUsecase)
	commodityController := _commodityController.NewCommodityController(commodityUsecase, userUsecase, proposalUsecase)
	proposalController := _proposalController.NewProposalController(proposalUsecase, commodityUsecase)

	routeController := _route.ControllerList{
		UserController:      userController,
		CommodityController: commodityController,
		ProposalController:  proposalController,
	}

	routeController.Init(e)

	appPort := fmt.Sprintf(":%s", _util.GetConfig("APP_PORT"))
	e.Logger.Fatal(e.Start(appPort))
}
