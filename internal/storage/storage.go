package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Close() error
}

type PgxConnInterface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close(ctx context.Context) error
}
