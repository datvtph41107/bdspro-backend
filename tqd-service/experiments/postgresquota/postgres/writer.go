package postgresquota

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tqd/experiments/postgresquota/quota"

	"gorm.io/gorm"
)

// Writer owns only the legal PostgreSQL mutation of a resolved quota pool.
// The Report acceptance adapter will own the surrounding transaction when P2.4 wires B+.
type Writer struct{}

func NewWriter() *Writer { return &Writer{} }

type consumeRow struct {
	PeriodEnd time.Time `gorm:"column:period_end"`
	Used      int64     `gorm:"column:used"`
}

func (w *Writer) ConsumeTx(
	ctx context.Context,
	tx *gorm.DB,
	input quota.Consumption,
) (int64, error) {
	if tx == nil {
		return 0, errors.New("quota transaction is not configured")
	}
	if !input.IsValid() {
		return 0, errors.New("quota consumption is invalid")
	}

	var row consumeRow
	result := tx.WithContext(ctx).Raw(`
INSERT INTO quota_pools (
    subject_type,
    subject_id,
    meter_code,
    period_start,
    period_end,
    used
)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (subject_type, subject_id, meter_code, period_start)
DO UPDATE SET
    used = quota_pools.used + EXCLUDED.used,
    updated_at = NOW()
WHERE quota_pools.period_end = EXCLUDED.period_end
  AND quota_pools.used <= ? - EXCLUDED.used
RETURNING period_end, used
`,
		string(input.Subject.Type),
		input.Subject.ID,
		string(input.Meter),
		input.PeriodStart,
		input.PeriodEnd,
		input.Amount,
		input.Limit,
	).Scan(&row)
	if result.Error != nil {
		return 0, fmt.Errorf("consume quota pool: %w", result.Error)
	}
	if result.RowsAffected == 1 {
		return row.Used, nil
	}

	var existing consumeRow
	err := tx.WithContext(ctx).Raw(`
SELECT period_end, used
  FROM quota_pools
 WHERE subject_type = ?
   AND subject_id = ?
   AND meter_code = ?
   AND period_start = ?
`,
		string(input.Subject.Type),
		input.Subject.ID,
		string(input.Meter),
		input.PeriodStart,
	).Scan(&existing).Error
	if err != nil {
		return 0, fmt.Errorf("classify quota rejection: %w", err)
	}
	if !existing.PeriodEnd.Equal(input.PeriodEnd) {
		return existing.Used, quota.ErrPeriodConflict
	}
	return existing.Used, quota.ErrExceeded
}
