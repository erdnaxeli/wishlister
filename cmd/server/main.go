package main

import (
	"log"

	"github.com/erdnaxeli/wishlister"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{RedirectCode: 308}))

	app, err := wishlister.New()
	if err != nil {
		log.Fatal(err)
	}

	setRoutes(e, app)
	setStatics(e)
	loadTemplates(e)

	e.Logger.Fatal(e.Start(":8080"))
}
