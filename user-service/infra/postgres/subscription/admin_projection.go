package subscriptionpostgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	subscriptiondomain "user/internal/domain/subscription"
	subscriptionusecase "user/internal/usecase/subscription"
)

type AdminProjectionStore struct{ db *gorm.DB }

func NewAdminProjectionStore(db *gorm.DB) *AdminProjectionStore { return &AdminProjectionStore{db: db} }

func (s *AdminProjectionStore) List(ctx context.Context, query subscriptionusecase.AdminQuery) (subscriptionusecase.AdminPage, error) {
	if s == nil || s.db == nil {
		return subscriptionusecase.AdminPage{}, fmt.Errorf("subscription projection database is not configured")
	}
	db := s.filtered(ctx, query)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return subscriptionusecase.AdminPage{}, fmt.Errorf("count subscriptions: %w", err)
	}
	rows := make([]subscriptiondomain.AdminProjection, 0)
	offset := int((query.Page - 1) * query.PageSize)
	if err := projectionSelect(db).Order("s.updated_at DESC, s.id DESC").Offset(offset).Limit(int(query.PageSize)).Scan(&rows).Error; err != nil {
		return subscriptionusecase.AdminPage{}, fmt.Errorf("list subscriptions: %w", err)
	}
	return subscriptionusecase.AdminPage{Subscriptions: rows, Total: uint64(total), Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *AdminProjectionStore) ListEffectiveByProfiles(ctx context.Context, profileIDs []uint64) ([]subscriptiondomain.AdminUserProjection, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("subscription projection database is not configured")
	}
	if len(profileIDs) == 0 {
		return []subscriptiondomain.AdminUserProjection{}, nil
	}
	rows := make([]subscriptiondomain.AdminProjection, 0)
	base := projectionBase(s.db.WithContext(ctx).Table("catalog_subscriptions AS s")).
		Where("s.subject_kind = ? AND s.subject_id IN ?", "profile", uint64Strings(profileIDs)).
		Where("s.status IN ?", []string{"pending", "active", "past_due"})
	if err := projectionSelect(base).
		Order("s.subject_id ASC, CASE s.status WHEN 'active' THEN 0 WHEN 'past_due' THEN 1 ELSE 2 END, s.updated_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list effective subscriptions: %w", err)
	}
	result := make([]subscriptiondomain.AdminUserProjection, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for i := range rows {
		profileID, err := subscriptionusecase.ParseProfileID(rows[i].SubjectID)
		if err != nil {
			continue
		}
		if _, ok := seen[profileID]; ok {
			continue
		}
		seen[profileID] = struct{}{}
		projection := rows[i]
		result = append(result, subscriptiondomain.AdminUserProjection{ProfileID: profileID, EffectiveSubscription: &projection})
	}
	return result, nil
}

func (s *AdminProjectionStore) filtered(ctx context.Context, query subscriptionusecase.AdminQuery) *gorm.DB {
	db := projectionBase(s.db.WithContext(ctx).Table("catalog_subscriptions AS s"))
	if query.ProfileID != 0 {
		db = db.Where("s.subject_kind = ? AND s.subject_id = ?", "profile", fmt.Sprint(query.ProfileID))
	}
	if query.Status != "" {
		db = db.Where("s.status = ?", query.Status)
	}
	if query.ProductCode != "" {
		db = db.Where("product.code = ?", query.ProductCode)
	}
	if query.PlanCode != "" {
		db = db.Where("plan.code = ?", query.PlanCode)
	}
	return db
}

func projectionSelect(db *gorm.DB) *gorm.DB {
	return db.Select(`s.id AS subscription_id, s.subject_kind, s.subject_id,
			product.code AS product_code, product.display_name AS product_display_name,
			plan.code AS plan_code, pv.display_name AS plan_display_name,
			s.plan_version_id, pv.version AS plan_version, s.status, s.started_at,
			s.current_period_start, s.current_period_end, s.access_until,
			s.pending_plan_version_id, s.pending_effective_at, s.order_reference, s.updated_at`)
}

func projectionBase(db *gorm.DB) *gorm.DB {
	return db.
		Joins("JOIN catalog_products product ON product.id = s.product_id").
		Joins("JOIN catalog_plan_versions pv ON pv.id = s.plan_version_id").
		Joins("JOIN catalog_plans plan ON plan.id = pv.plan_id")
}

func uint64Strings(ids []uint64) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = fmt.Sprint(id)
	}
	return result
}
