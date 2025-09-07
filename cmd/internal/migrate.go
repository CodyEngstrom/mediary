package internal

import (
	"database/sql"
	"fmt"
	"mediary/config"
	"mediary/internal"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MigratePostgresql(cfg *config.Config, migrationsPath string) error {
	databaseHost := cfg.DBHost
	databasePort := cfg.DBPort
	databaseUsername := cfg.DBUser
	databasePassword := cfg.DBPass
	databaseName := cfg.DBName
	databaseSSLMode := cfg.DBSSLMode

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%d", databaseHost, databasePort),
		Path:   databaseName,
	}
	q := dsn.Query()
	q.Add("sslmode", databaseSSLMode)
	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return internal.WrapDatabaseError(err, "sql.open failed during migration")
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return internal.WrapDatabaseError(err, "failed to create postgres migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		return internal.WrapDatabaseError(err, "failed to create postgres migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return internal.WrapDatabaseError(err, "postgres migration failed")
	}
	return nil
}
