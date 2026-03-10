package database

import (
	"context"
	"database/sql"
)

type Database interface {
	DB() *sql.DB
	Ping(ctx context.Context) error
	Close() error
}
