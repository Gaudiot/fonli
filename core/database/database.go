package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(databaseURL string) (*DB, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool == nil {
		db.Pool.Close()
	}
}
