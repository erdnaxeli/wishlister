// Implements a web server that expose a web app to manage wishlists.
package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/erdnaxeli/wishlister"
	"github.com/erdnaxeli/wishlister/pkg/email"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type config struct {
	Email         string `env:"EMAIL"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Error while reading configuration: %s", err)
	}

	if cfg.Email == "" && cfg.EmailPassword == "" {
		log.Fatalf(
			"You must provide an email password with the env var EMAIL_PASSWORD or disable emailing with EMAIL=off.",
		)
	}

	e := echo.New()
	e.Debug = true
	e.Pre(
		middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{RedirectCode: 308}),
	)

	var mailSender email.Sender
	if cfg.Email == "off" {
		mailSender = email.NoMailer{}
	} else {
		mailSender = email.NewSMTPSender(
			"contact@malistedevoeux.fr",
			cfg.EmailPassword,
			"mail.infomaniak.fr",
			465,
			"contact@malistedevoeux.fr",
		)
	}

	app, err := wishlister.New(mailSender)
	if err != nil {
		e.Logger.Fatal(err)
	}

	setRoutes(e, app)
	setStatics(e)
	loadTemplates(e)

	e.Logger.Fatal(e.Start(":8080"))
}
