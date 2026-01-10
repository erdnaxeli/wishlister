package migrator

import (
	"fmt"
	"log"
)

func (m *migrator) Migrate() error {
	if len(m.migrations) == 0 {
		log.Print("No migrations to apply.")
		return nil
	}

	if m.currentVersion == m.lastVersion {
		log.Print("Database is already up to date.")
		return nil
	}

	for i := m.currentVersion + 1; i <= m.lastVersion; i++ {
		err := m.applyMigration(i)
		if err != nil {
			return err
		}
	}

	log.Print("All migrations applied successfully.")

	return nil
}

func (m *migrator) applyMigration(version int) error {
	migration := m.migrations[version-1]
	log.Printf("Applying migration %d: %s.", migration.version, migration.name)

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf(
			"failed to begin transaction for migration %d: %w", migration.version, err,
		)
	}

	defer func() { _ = tx.Rollback() }()

	for _, sql := range migration.upSQL {
		_, err = tx.Exec(sql)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.version, err)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version) VALUES (?)`,
		migration.version,
	)
	if err != nil {
		return fmt.Errorf("failed to record migration %d: %w", migration.version, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit migration %d: %w", migration.version, err)
	}

	m.currentVersion = migration.version
	return nil
}
