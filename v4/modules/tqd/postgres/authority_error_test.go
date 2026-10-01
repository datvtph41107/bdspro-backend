package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"testing"
	"userv4/modules/tqd/services"
)

func TestAuthorityConstraintMapping(t *testing.T) {
	conflict := &pgconn.PgError{Code: "23505", ConstraintName: "idx_qh_authority_issuring_code_active"}
	if !errors.Is(authorityWriteError(conflict), usecase.ErrAuthorityCodeConflict) {
		t.Fatal("code constraint must be a business conflict")
	}
	other := &pgconn.PgError{Code: "23505", ConstraintName: "qh_authority_issuring_pkey"}
	if authorityWriteError(other) != other {
		t.Fatal("unrelated constraint must preserve original failure")
	}
}
