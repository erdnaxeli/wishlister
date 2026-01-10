package migrator

import "fmt"

// InvalidMigrationFilenameError is returned when a migration filename does not match the expected pattern.
type InvalidMigrationFilenameError struct {
	Filename string
}

func (e InvalidMigrationFilenameError) Error() string {
	return "invalid migration filename: " + e.Filename
}

// EmptyMigrationError is returned when a migration file is empty.
type EmptyMigrationError struct {
	Filename string
}

func (e EmptyMigrationError) Error() string {
	return "empty migration file: " + e.Filename
}

// InvalidMigrationFileError is returned when a migration file foramt is invalid.
type InvalidMigrationFileError struct {
	Filename string
}

func (e InvalidMigrationFileError) Error() string {
	return "invalid migration file: " + e.Filename + ", first line must be \"-- +migrate Up\""
}

// DuplicateMigrationVersionError is returned when there are multiple migrations with the same version.
type DuplicateMigrationVersionError struct {
	Version int
}

func (e DuplicateMigrationVersionError) Error() string {
	return fmt.Sprintf("duplicate migration version: %d", e.Version)
}

// MissingMigrationVersionError is returned when a migration version is missing in the sequence.
type MissingMigrationVersionError struct {
	Version int
}

func (e MissingMigrationVersionError) Error() string {
	return fmt.Sprintf("missing migration version: %d", e.Version)
}
