// Implements a web server that expose a web app to manage wishlists.
package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/erdnaxeli/migrator"

	"github.com/erdnaxeli/wishlister"
	"github.com/erdnaxeli/wishlister/pkg/email"
	"github.com/erdnaxeli/wishlister/pkg/server"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

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
		mailSender, err = email.NewSMTPSender(
			"contact@malistedevoeux.fr",
			cfg.EmailPassword,
			"mail.infomaniak.fr",
			465,
			"\"Ma liste de voeux.fr\" <contact@malistedevoeux.fr>",
		)
		if err != nil {
			log.Fatalf("Error while creating mail client: %s", err)
		}
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	err = runServer(db, mailSender)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	log.Print("Opening database")

	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %w", err)
	}

	log.Print("Applying migrations")

	subFS, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error while applying migrations: %w", err)
	}
	migrator, err := migrator.New(db, subFS)
	if err != nil {
		return nil, fmt.Errorf("error while applying migrations: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return nil, fmt.Errorf("error while applying migrations: %w", err)
	}

	return db, nil
}

func runServer(db *sql.DB, mailSender email.Sender) error {
	defer func() { _ = db.Close() }()
	log.Print("Starting application")

	app, err := wishlister.NewWithConfig(db, mailSender)
	if err != nil {
		return err
	}

	server.New(server.Config{Wishlister: app}).Run()

	return nil
}
