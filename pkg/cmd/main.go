// Implements a web server that expose a web app to manage wishlists.
package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/erdnaxeli/wishlister/pkg/email"
	"github.com/erdnaxeli/wishlister/pkg/server"

	"github.com/erdnaxeli/wishlister"
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
		log.Fatal(err)
	}

	server.New(server.Config{Wishlister: app}).Run()
}
