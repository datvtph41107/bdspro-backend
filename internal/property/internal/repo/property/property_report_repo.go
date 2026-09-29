package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyReportRepository interface {
	Create(ctx context.Context, report *domain.PropertyReport) error
	GetByID(ctx context.Context, id uint64) (*domain.PropertyReport, error)
	GetByLineageAndReporter(ctx context.Context, lineageID, reporterOriginID uint64) (*domain.PropertyReport, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
	CountPendingByLineage(ctx context.Context, lineageID uint64) (int64, error)
}
