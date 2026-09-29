package _db

import (
	"log/slog"
	"os"
	"reflect"
)

// AutoMigrate creates the table when it does not exist.
// This compatibility helper preserves its historical immediate-exit behavior.
func AutoMigrate(e interface{}) {
	val := reflect.ValueOf(e)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		slog.Error("database migration requires a struct value")
		os.Exit(1)
	}

	if err := DB.AutoMigrate(e); err != nil {
		slog.Error("database migration failed", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("database migration completed")
}
