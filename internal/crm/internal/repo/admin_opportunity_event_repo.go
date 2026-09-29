package repo

import (
	"context"
	"crm/internal/domain"
)

type AdminOpportunityEventRepo interface {
	AddEvent(ctx context.Context, event *domain.AdminOpportunityEvent) error
	ListEvents(ctx context.Context, opportunityID uint64, limit int) ([]domain.AdminOpportunityEvent, error)
}
