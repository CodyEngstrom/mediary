package internal

import (
	"context"
	"fmt"
	"mediary/config"
	"mediary/internal"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL(cfg *config.Config) (*pgxpool.Pool, error) {
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

	pool, err := pgxpool.New(context.Background(), dsn.String())
	if err != nil {
		return nil, internal.WrapDatabaseError(err, "pgxpool.Connect failed")
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, internal.WrapDatabaseError(err, "db.Ping failed")
	}
	return pool, nil
}
