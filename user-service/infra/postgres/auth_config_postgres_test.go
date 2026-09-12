package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type authConfigFakeResult struct {
	rows int64
}

func (r authConfigFakeResult) LastInsertId() (int64, error) { return 0, nil }
func (r authConfigFakeResult) RowsAffected() (int64, error) { return r.rows, nil }

type authConfigFakePool struct {
	rowsAffected int64
	failAt       int
	writeErr     error
	tx           *authConfigFakeTx
}

func (p *authConfigFakePool) BeginTx(context.Context, *sql.TxOptions) (gorm.ConnPool, error) {
	p.tx = &authConfigFakeTx{
		rowsAffected: p.rowsAffected,
		failAt:       p.failAt,
		writeErr:     p.writeErr,
	}
	return p.tx, nil
}

func (p *authConfigFakePool) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, errors.New("unexpected PrepareContext on root pool")
}

func (p *authConfigFakePool) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, errors.New("unexpected ExecContext outside transaction")
}

func (p *authConfigFakePool) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("unexpected QueryContext on root pool")
}

func (p *authConfigFakePool) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return &sql.Row{}
}

type authConfigFakeTx struct {
	rowsAffected int64
	failAt       int
	writeErr     error
	execCalls    int
	committed    bool
	rolledBack   bool
}

func (tx *authConfigFakeTx) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, errors.New("unexpected PrepareContext in transaction")
}

func (tx *authConfigFakeTx) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	tx.execCalls++
	if tx.failAt > 0 && tx.execCalls == tx.failAt {
		return nil, tx.writeErr
	}
	return authConfigFakeResult{rows: tx.rowsAffected}, nil
}

func (tx *authConfigFakeTx) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("unexpected QueryContext in transaction")
}

func (tx *authConfigFakeTx) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (tx *authConfigFakeTx) Commit() error {
	tx.committed = true
	return nil
}

func (tx *authConfigFakeTx) Rollback() error {
	tx.rolledBack = true
	return nil
}

func newAuthConfigTransactionTestDB(t *testing.T, pool *authConfigFakePool) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:             pool,
			WithoutReturning: true,
		}),
		&gorm.Config{
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}
	return db
}

func TestUpdateValuesAtomicallyPropagatesSecondWriteFailureAndRollsBack(t *testing.T) {
	writeErr := errors.New("second update failed")
	pool := &authConfigFakePool{
		rowsAffected: 1,
		failAt:       2,
		writeErr:     writeErr,
	}
	repo := NewAuthConfigRepository(newAuthConfigTransactionTestDB(t, pool))

	err := repo.UpdateValuesAtomically(context.Background(), map[string]string{
		"ZNS_TOKEN":   "new-access",
		"ZNS_REFRESH": "new-refresh",
	})
	if !errors.Is(err, writeErr) {
		t.Fatalf("UpdateValuesAtomically() error = %v, want wrapped %v", err, writeErr)
	}
	if pool.tx == nil || pool.tx.execCalls != 2 {
		t.Fatalf("transaction writes = %v, want 2", pool.tx)
	}
	if !pool.tx.rolledBack || pool.tx.committed {
		t.Fatalf("transaction state committed=%v rolledBack=%v, want rollback only", pool.tx.committed, pool.tx.rolledBack)
	}
}

func TestUpdateValuesAtomicallyRejectsMissingConfigRowAndRollsBack(t *testing.T) {
	pool := &authConfigFakePool{rowsAffected: 0}
	repo := NewAuthConfigRepository(newAuthConfigTransactionTestDB(t, pool))

	err := repo.UpdateValuesAtomically(context.Background(), map[string]string{
		"ZNS_TOKEN": "new-access",
	})
	if err == nil {
		t.Fatal("UpdateValuesAtomically() error = nil, want missing-row failure")
	}
	if pool.tx == nil || pool.tx.execCalls != 1 {
		t.Fatalf("transaction writes = %v, want 1", pool.tx)
	}
	if !pool.tx.rolledBack || pool.tx.committed {
		t.Fatalf("transaction state committed=%v rolledBack=%v, want rollback only", pool.tx.committed, pool.tx.rolledBack)
	}
}

func TestUpdateValuesAtomicallyCommitsCompletePair(t *testing.T) {
	pool := &authConfigFakePool{rowsAffected: 1}
	repo := NewAuthConfigRepository(newAuthConfigTransactionTestDB(t, pool))

	err := repo.UpdateValuesAtomically(context.Background(), map[string]string{
		"ZNS_TOKEN":   "new-access",
		"ZNS_REFRESH": "new-refresh",
	})
	if err != nil {
		t.Fatalf("UpdateValuesAtomically() error = %v", err)
	}
	if pool.tx == nil || pool.tx.execCalls != 2 {
		t.Fatalf("transaction writes = %v, want 2", pool.tx)
	}
	if !pool.tx.committed || pool.tx.rolledBack {
		t.Fatalf("transaction state committed=%v rolledBack=%v, want commit only", pool.tx.committed, pool.tx.rolledBack)
	}
}
