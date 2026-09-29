package quota

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	quotadomain "tqd/internal/domain/quota"
	quotausecase "tqd/internal/usecase/quota"
)

type AdminUsageStore struct{ db *gorm.DB }

func NewAdminUsageStore(db *gorm.DB) *AdminUsageStore { return &AdminUsageStore{db: db} }

func (s *AdminUsageStore) List(ctx context.Context, query quotausecase.AdminUsageQuery) (quotausecase.AdminUsagePage, error) {
	if s == nil || s.db == nil {
		return quotausecase.AdminUsagePage{}, fmt.Errorf("admin usage database is not configured")
	}
	base := s.db.WithContext(ctx).Table("quota_usage_events").Where("subject_type = ?", "profile")
	if query.ProfileID != 0 {
		base = base.Where("subject_id = ?", fmt.Sprint(query.ProfileID))
	}
	if query.Operation != "" {
		base = base.Where("operation = ?", query.Operation)
	}
	if query.MeterCode != "" {
		base = base.Where("meter_code = ?", query.MeterCode)
	}
	grouped := base.Session(&gorm.Session{}).Select("1").Group("subject_type, subject_id, operation, meter_code, period_start, period_end")
	var total int64
	if err := s.db.WithContext(ctx).Table("(?) AS grouped", grouped).Count(&total).Error; err != nil {
		return quotausecase.AdminUsagePage{}, fmt.Errorf("count usage summaries: %w", err)
	}
	rows := make([]quotadomain.AdminUsageProjection, 0)
	offset := int((query.Page - 1) * query.PageSize)
	selectSQL := `subject_type, subject_id, operation, meter_code,
		COALESCE(SUM(amount), 0) AS durable_used,
		BOOL_OR(limit_snapshot IS NOT NULL) AS limit_known,
		COALESCE(MAX(limit_snapshot), 0) AS limit_snapshot,
		period_start, period_end,
		COALESCE(MAX(subscription_id), 0) AS subscription_id,
		COALESCE(MAX(plan_code), '') AS plan_code,
		COALESCE(MAX(plan_version), '') AS plan_version,
		COALESCE(MAX(policy_version), '') AS policy_version,
		MAX(created_at) AS last_usage_at`
	if err := base.Select(selectSQL).
		Group("subject_type, subject_id, operation, meter_code, period_start, period_end").
		Order("period_end DESC, subject_id ASC, meter_code ASC").Offset(offset).Limit(int(query.PageSize)).Scan(&rows).Error; err != nil {
		return quotausecase.AdminUsagePage{}, fmt.Errorf("list usage summaries: %w", err)
	}
	return quotausecase.AdminUsagePage{Usage: rows, Total: uint64(total), Page: query.Page, PageSize: query.PageSize}, nil
}
