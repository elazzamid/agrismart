package platform

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DatabaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}
	return "postgres://agrismart:agrismart@localhost:5432/agrismart?sslmode=disable"
}

func OpenDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}
	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
