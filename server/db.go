package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func migrateDatabase() {
	m, err := migrate.New("file://./migrations", "pgx5"+getDbUrl()[8:])

	if err != nil {
		log.Fatalf("failed to create migration: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to apply migration: %v", err)
	}

	defer func(m *migrate.Migrate) {
		if err, _ = m.Close(); err != nil {
			log.Fatalf("failed to close migration: %v", err)
		}
	}(m)
}

func getDbUrl() string {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalf("database connection string not set")
	}

	return dbUrl
}

func getDb(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, getDbUrl())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func closeDb(conn *pgx.Conn, ctx context.Context) {
	_ = conn.Close(ctx)
}
