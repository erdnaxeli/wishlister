// Package migrator implements database migration functionalities.
package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"regexp"
)

// FilenameRgx is the regular expression to match migration filenames.
var FilenameRgx = regexp.MustCompile(`^(\d+)_(.*)\.sql$`)

// Migrator is the interface to manage database migrations.
type Migrator interface {
	// Migrate applies all pending database migrations.
	Migrate() error

	// Version returns the current version of the database schema.
	Version() (int, error)
}

type migrator struct {
	db *sql.DB

	migrations     []Migration
	currentVersion int
	lastVersion    int
}

// Migration represents a database migration.
type Migration struct {
	version int
	name    string
	upSQL   []string
}

// New creates a new Migrator instance.
func New(db *sql.DB, fs fs.FS) (Migrator, error) {
	migrations, err := loadMigrations(fs)
	if err != nil {
		return nil, err
	}

	lastVersion, err := validateMigrations(migrations)
	if err != nil {
		return nil, err
	}

	currentVersion, err := getCurrentDBVersion(db)
	if err != nil {
		return nil, err
	}

	return &migrator{
		db:             db,
		migrations:     migrations,
		lastVersion:    lastVersion,
		currentVersion: currentVersion,
	}, nil
}

func (m *migrator) Version() (int, error) {
	return getCurrentDBVersion(m.db)
}

func getCurrentDBVersion(db *sql.DB) (int, error) {
	// check if the table exists
	_, err := db.Exec(`SELECT 1 FROM schema_migrations`)
	if err != nil {
		log.Print("error checking schema_migrations table: ", err)
		log.Print("assuming schema_migrations table does not exist, try to create it")

		_, err = db.Exec(
			`
			CREATE TABLE schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)
			`,
		)
		if err != nil {
			log.Print("error creating schema_migrations table: ", err)
			return 0, fmt.Errorf("failed to create schema_migrations table: %w", err)
		}
		return 0, err
	}

	var version int
	err = db.QueryRow(
		`SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`,
	).Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return version, nil
}
