package postgresusage

import (
	commonmetering "common/metering"
	"context"
	"errors"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota/reconciliation/usageprojection"
)

const runtimeUsageRetention = 48 * time.Hour

/**
 * ListUsageScopes lấy các counter đã có usage_event để background check Redis.
 */
func (s *Store) ListUsageScopes(
	ctx context.Context,
	offset int,
	limit int,
) ([]usageprojection.Scope, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("usage database is not configured")
	}
	if limit <= 0 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	type row struct {
		SubjectType string
		SubjectID   string
		MeterCode   string
		PeriodStart time.Time
		PeriodEnd   time.Time
	}

	var rows []row
	err := s.db.WithContext(ctx).
		Model(&model{}).
		Select("subject_type, subject_id, meter_code, period_start, period_end").
		Where("period_end > ?", time.Now().UTC().Add(-runtimeUsageRetention)).
		Group("subject_type, subject_id, meter_code, period_start, period_end").
		Order("period_end DESC, subject_type ASC, subject_id ASC, meter_code ASC, period_start ASC").
		Offset(offset).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]usageprojection.Scope, 0, len(rows))
	for _, item := range rows {
		meterCode := commonmetering.Code(item.MeterCode)
		subject := access.Subject{
			Type: access.SubjectType(item.SubjectType),
			ID:   item.SubjectID,
		}
		if !subject.IsValid() || !meterCode.IsValid() {
			continue
		}
		result = append(result, usageprojection.Scope{
			Subject:     subject,
			MeterCode:   meterCode,
			PeriodStart: item.PeriodStart,
			PeriodEnd:   item.PeriodEnd,
		})
	}
	return result, nil
}
