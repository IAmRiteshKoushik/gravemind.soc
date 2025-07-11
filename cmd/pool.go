package cmd

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func InitDB() (*pgxpool.Pool, error) {

	connString := App.DatabaseUrl
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		Log.Fatal(fmt.Sprintf("Failed to parse database config: %v", err))
		return nil, err
	}

	config.MinConns = 1
	config.MaxConns = 5
	config.MaxConnLifetime = 3600
	config.MaxConnIdleTime = 1800
	config.HealthCheckPeriod = 60
	config.MaxConnLifetimeJitter = 0

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		Log.Fatal(fmt.Sprintf("Failed to create connection pool: %v", err))
		return nil, err
	}

	// Verify connection
	err = pool.Ping(context.Background())
	if err != nil {
		Log.Fatal(fmt.Sprintf("Failed to ping database: %v", err))
		return nil, err
	}

	return pool, nil
}
