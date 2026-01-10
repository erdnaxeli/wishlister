// Implements a web server that expose a web app to manage wishlists.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"

	"github.com/erdnaxeli/wishlister"
	"github.com/erdnaxeli/wishlister/migrator"
	"github.com/erdnaxeli/wishlister/pkg/email"
	"github.com/erdnaxeli/wishlister/pkg/server"
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
			"\"Ma liste de voeux.fr\" <contact@malistedevoeux.fr>",
		)
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting application")

	app, err := wishlister.NewWithConfig(db, mailSender)
	if err != nil {
		log.Fatal(err)
	}

	server.New(server.Config{Wishlister: app}).Run()
}

func openDB() (*sql.DB, error) {
	log.Print("Opening database")

	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %w", err)
	}

	defer func() { _ = db.Close() }()
	log.Print("Applying migrations")

	directory := os.DirFS("db/migrations")
	migrator, err := migrator.New(db, directory)
	if err != nil {
		return nil, fmt.Errorf("error while applying migrations: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return nil, fmt.Errorf("error while applying migrations: %w", err)
	}

	return db, nil
}
