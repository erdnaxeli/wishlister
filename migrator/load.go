package migrator

import (
	"bufio"
	"cmp"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

func loadMigrations(directory fs.FS) ([]Migration, error) {
	matches, err := fs.Glob(directory, "*.sql")
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, filename := range matches {
		migration, err := loadMigration(directory, filename)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func loadMigration(directory fs.FS, filename string) (Migration, error) {
	submatches := FilenameRgx.FindSubmatch([]byte(filename))
	if submatches == nil {
		return Migration{}, InvalidMigrationFilenameError{Filename: filename}
	}

	versionStr := string(submatches[1])
	name := string(submatches[2])

	var version int
	_, err := fmt.Sscanf(versionStr, "%d", &version)
	if err != nil {
		return Migration{}, fmt.Errorf("error while parsing version: %s, %w", versionStr, err)
	}

	upSQL, err := readMigrationSQL(directory, filename)
	if err != nil {
		return Migration{}, err
	}

	migration := Migration{
		version: version,
		name:    name,
		upSQL:   upSQL,
	}
	return migration, nil
}

func readMigrationSQL(directory fs.FS, filename string) ([]string, error) {
	file, err := directory.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, EmptyMigrationError{Filename: filename}
	}

	if scanner.Text() != "-- +migrate Up" {
		return nil, InvalidMigrationFileError{Filename: filename}
	}

	var b strings.Builder
	var upSQL []string

	for scanner.Scan() {
		bytes := scanner.Bytes()
		b.Write(bytes)

		if len(bytes) > 0 && bytes[len(bytes)-1] == ';' {
			upSQL = append(upSQL, b.String())
			b.Reset()
		} else {
			b.WriteString("\n")
		}
	}

	if b.Len() > 0 {
		upSQL = append(upSQL, b.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return upSQL, nil
}

func validateMigrations(migrations []Migration) (int, error) {
	seenVersions := make(map[int]bool)
	maxVersion := 0
	for _, m := range migrations {
		if seenVersions[m.version] {
			return 0, DuplicateMigrationVersionError{Version: m.version}
		}

		if m.version > maxVersion {
			maxVersion = m.version
		}

		seenVersions[m.version] = true
	}

	for v := 1; v <= maxVersion; v++ {
		if !seenVersions[v] {
			return 0, MissingMigrationVersionError{Version: v}
		}
	}

	slices.SortFunc(
		migrations,
		func(a Migration, b Migration) int { return cmp.Compare(a.version, b.version) },
	)

	return maxVersion, nil
}
