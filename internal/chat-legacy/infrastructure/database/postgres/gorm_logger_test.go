package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

func decodeSingleRecord(
	t *testing.T,
	buffer *bytes.Buffer,
) map[string]any {
	t.Helper()

	var record map[string]any

	if err := json.Unmarshal(
		bytes.TrimSpace(buffer.Bytes()),
		&record,
	); err != nil {
		t.Fatalf(
			"decode canonical GORM log: %v",
			err,
		)
	}

	return record
}

func TestNewGormLoggerPreservesBaselineConfiguration(
	t *testing.T,
) {
	logger, ok := newGormLogger().(slogGormLogger)

	if !ok {
		t.Fatal(
			"newGormLogger did not return slogGormLogger",
		)
	}

	if logger.level != gormlogger.Info {
		t.Fatalf(
			"expected Info level, got %v",
			logger.level,
		)
	}

	if logger.slowThreshold != 200*time.Millisecond {
		t.Fatalf(
			"expected 200ms slow threshold, got %v",
			logger.slowThreshold,
		)
	}

	if !logger.ignoreRecordNotFound {
		t.Fatal(
			"expected ignoreRecordNotFound=true",
		)
	}
}

func TestSlogGormLoggerProjectsInfoQuery(
	t *testing.T,
) {
	var buffer bytes.Buffer

	original := slog.Default()

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				&buffer,
				nil,
			),
		),
	)

	defer slog.SetDefault(original)

	logger := newGormLogger()

	logger.Trace(
		context.Background(),
		time.Now().Add(
			-10*time.Millisecond,
		),
		func() (string, int64) {
			return "SELECT 1", 1
		},
		nil,
	)

	record := decodeSingleRecord(
		t,
		&buffer,
	)

	if record["msg"] != "database query" {
		t.Fatalf(
			"unexpected msg: %v",
			record["msg"],
		)
	}

	if record["level"] != "INFO" {
		t.Fatalf(
			"unexpected level: %v",
			record["level"],
		)
	}

	if record["component"] != "postgres" {
		t.Fatalf(
			"unexpected component: %v",
			record["component"],
		)
	}

	if record["sql"] != "SELECT 1" {
		t.Fatalf(
			"unexpected sql: %v",
			record["sql"],
		)
	}

	if record["rows"] != float64(1) {
		t.Fatalf(
			"unexpected rows: %v",
			record["rows"],
		)
	}
}

func TestSlogGormLoggerProjectsSlowQuery(
	t *testing.T,
) {
	var buffer bytes.Buffer

	original := slog.Default()

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				&buffer,
				nil,
			),
		),
	)

	defer slog.SetDefault(original)

	logger := newGormLogger()

	logger.Trace(
		context.Background(),
		time.Now().Add(
			-500*time.Millisecond,
		),
		func() (string, int64) {
			return "SELECT pg_sleep(1)", 2
		},
		nil,
	)

	record := decodeSingleRecord(
		t,
		&buffer,
	)

	if record["msg"] != "slow database query" {
		t.Fatalf(
			"unexpected msg: %v",
			record["msg"],
		)
	}

	if record["level"] != "WARN" {
		t.Fatalf(
			"unexpected level: %v",
			record["level"],
		)
	}

	if record["component"] != "postgres" {
		t.Fatalf(
			"unexpected component: %v",
			record["component"],
		)
	}

	if record["sql"] != "SELECT pg_sleep(1)" {
		t.Fatalf(
			"unexpected sql: %v",
			record["sql"],
		)
	}

	if record["rows"] != float64(2) {
		t.Fatalf(
			"unexpected rows: %v",
			record["rows"],
		)
	}
}
