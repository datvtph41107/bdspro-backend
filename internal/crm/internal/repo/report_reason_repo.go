package repo

import (
	"context"
	"crm/internal/domain"
)

// ReportReasonRepo defines the interface for report reason repository operations
type ReportReasonRepo interface {
	Create(ctx context.Context, reason *domain.ReportReason) (*domain.ReportReason, error)

	// GetByID gets a report reason by ID
	GetByID(ctx context.Context, id uint64) (*domain.ReportReason, error)

	// GetList gets a list of report reasons with pagination
	GetList(ctx context.Context, req *domain.ReportReasonListRequest) ([]domain.ReportReason, int64, error)

	// Update updates a report reason
	Update(ctx context.Context, reason *domain.ReportReason) error

	// Delete deletes a report reason by ID
	Delete(ctx context.Context, id uint64) error

	// GetAll gets all report reasons
	GetAll(ctx context.Context) ([]domain.ReportReason, error)
}