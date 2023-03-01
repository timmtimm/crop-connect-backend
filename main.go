package main

import (
	"fmt"
	_route "marketplace-backend/app/route"
	"marketplace-backend/util"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	_route.Init(e)

	appPort := fmt.Sprintf(":%s", util.GetConfig("APP_PORT"))
	e.Logger.Fatal(e.Start(appPort))
}
