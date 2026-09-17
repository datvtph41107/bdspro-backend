package walletpostgres

import (
	"context"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DryRun kiểm tra SQL của adapter thật; test usecase mô phỏng lock chưa đủ
// nếu query production vô tình bỏ FOR UPDATE.
func TestDepositLookupLocksTransactionRow(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost user=test dbname=test sslmode=disable"}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	var query string
	err = db.Callback().Query().After("gorm:query").Register("capture_deposit_lock", func(tx *gorm.DB) { query = tx.Statement.SQL.String() })
	if err != nil {
		t.Fatal(err)
	}
	repo := &walletTransactionRepository{db: db}
	if _, err := repo.GetTransactionByCodeForUpdate(context.Background(), "DEP123"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(query, "FOR UPDATE") || !strings.Contains(query, "transaction_code") {
		t.Fatalf("deposit query does not lock by command: %s", query)
	}
}
