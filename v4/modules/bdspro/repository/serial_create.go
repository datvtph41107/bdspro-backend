package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// createSerialEntity inserts a row with BIGSERIAL id, omitting explicit zero IDs.
// After dev/import with hand-assigned ids, Postgres sequences can lag MAX(id); resync once on pkey conflict.
func createSerialEntity(tx *gorm.DB, table string, row any) error {
	sp := fmt.Sprintf("serial_%s", table)
	_ = tx.SavePoint(sp)
	err := tx.Omit("ID").Create(row).Error
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "duplicate key") {
		_ = tx.RollbackTo(sp)
		return err
	}
	// Failed INSERT aborts the outer tx unless we roll back to a savepoint first.
	_ = tx.RollbackTo(sp)
	if syncErr := resyncSerialAfterPKConflict(tx, table); syncErr != nil {
		return err
	}
	return tx.Omit("ID").Create(row).Error
}

func resyncSerialAfterPKConflict(tx *gorm.DB, table string) error {
	q := fmt.Sprintf(
		`SELECT setval(pg_get_serial_sequence('%s', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM %s), 0), 1), (SELECT COALESCE(MAX(id), 0) FROM %s) > 0)`,
		table, table, table,
	)
	return tx.Exec(q).Error
}
