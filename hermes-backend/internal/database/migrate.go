package database

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// FormatDSNToURL converts a GORM DSN to a URL format required by golang-migrate
func FormatDSNToURL(dsn string) string {
	// If DSN already has the postgres:// scheme, return it
	if strings.HasPrefix(dsn, "postgres://") {
		return dsn
	}

	// Otherwise, try to convert it
	// Example DSN: "host=localhost user=postgres password=postgres dbname=hermes port=5432 sslmode=disable"
	// Should become: "postgres://postgres:postgres@localhost:5432/hermes?sslmode=disable"

	// Parse DSN parts
	var host, user, password, dbname, port, sslmode string

	parts := strings.Split(dsn, " ")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		value := kv[1]

		switch key {
		case "host":
			host = value
		case "user":
			user = value
		case "password":
			password = value
		case "dbname":
			dbname = value
		case "port":
			port = value
		case "sslmode":
			sslmode = value
		}
	}

	// Build URL
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	if sslmode != "" {
		url += fmt.Sprintf("?sslmode=%s", sslmode)
	}

	return url
}

// RunMigrations applies all pending database migrations
func RunMigrations(dsn string, migrationsPath string) error {
	dbURL := FormatDSNToURL(dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// RollbackLastMigration rolls back the most recently applied migration
func RollbackLastMigration(dsn string, migrationsPath string) error {
	dbURL := FormatDSNToURL(dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// RollbackAllMigrations rolls back all migrations
func RollbackAllMigrations(dsn string, migrationsPath string) error {
	dbURL := FormatDSNToURL(dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback all migrations: %w", err)
	}

	return nil
}
