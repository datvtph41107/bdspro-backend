package _db

import (
	"bytes"
	"common/logging"
	"common/request"
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestSQLParameterizedAndCorrelated(t *testing.T) {
	configureSQLRecorder()
	var output bytes.Buffer
	old := slog.Default()
	defer slog.SetDefault(old)
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})))
	log := &sqlLogger{level: 4, config: logging.SQLConfig{Enabled: true, SlowThreshold: time.Second}}
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost user=test dbname=test", PreferSimpleProtocol: true}), &gorm.Config{DisableAutomaticPing: true, DryRun: true, Logger: log})
	if err != nil {
		t.Fatal(err)
	}
	ctx := request.WithRequestID(context.Background(), "req_9770ca1e2a97738baee1447ecd49df9f")
	db.WithContext(ctx).Table("profiles").Where("phone = ?", "private-phone").Find(&[]struct{ ID uint64 }{})
	db.WithContext(ctx).Raw("SELECT id FROM profiles WHERE phone = ?", "private-phone").Scan(&[]struct{ ID uint64 }{})
	if strings.Contains(output.String(), "private-phone") || !strings.Contains(output.String(), "$1") || !strings.Contains(output.String(), "request_id") || !strings.Contains(output.String(), "sql_logger_test.go") {
		t.Fatalf("query log: %s", output.String())
	}
}

func TestDisabledSQLDoesNotEvaluateNormalQuery(t *testing.T) {
	log := &sqlLogger{level: 4, config: logging.SQLConfig{SlowThreshold: time.Second}}
	log.Trace(context.Background(), time.Now(), func() (string, int64) { t.Fatal("normal disabled query evaluated"); return "", 0 }, nil)
}
