package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(
	ctx context.Context,
	databaseURL string,
) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}
	if schema := os.Getenv("DB_SCHEMA"); schema != "" {
		config.ConnConfig.RuntimeParams["search_path"] = pgx.Identifier{schema}.Sanitize() + ",public"
	}
	// Leave connection capacity for the other portfolio apps and migrations.
	config.MaxConns = 3
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf(
			"create database pool: %w",
			err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"connect to database: %w",
			err,
		)
	}

	return pool, nil
}
