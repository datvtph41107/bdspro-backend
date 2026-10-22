package postgres

import (
	"context"
	domain "mdn4/modules/crm/entity"
	repo "mdn4/modules/crm/repository"

	"gorm.io/gorm"
)

type AdminOpportunityEventPostgres struct {
	db *gorm.DB
}

func NewAdminOpportunityEventPostgres(db *gorm.DB) repo.AdminOpportunityEventRepo {
	return &AdminOpportunityEventPostgres{db: db}
}

func (r *AdminOpportunityEventPostgres) AddEvent(ctx context.Context, event *domain.AdminOpportunityEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *AdminOpportunityEventPostgres) ListEvents(ctx context.Context, opportunityID uint64, limit int) ([]domain.AdminOpportunityEvent, error) {
	var events []domain.AdminOpportunityEvent
	q := r.db.WithContext(ctx).
		Where("opportunity_id = ? AND deleted_at IS NULL", opportunityID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&events).Error
	return events, err
}
