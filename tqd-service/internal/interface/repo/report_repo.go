package repo

import (
	"context"
	"tqd/internal/domain"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	Update(ctx context.Context, report *domain.Report) error
	UpdateStatus(ctx context.Context, id uint64, status uint32, errorMsg *string) error
	UpdateFileInfo(ctx context.Context, id uint64, fileURL *string, fileSize *int64, fileHash *string, completedAt *string) error
	GetByID(ctx context.Context, id uint64) (*domain.Report, error)
	ListByUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
	AdminList(ctx context.Context, filter AdminReportFilter) ([]domain.Report, int64, error)
	AdminSummary(ctx context.Context) (map[uint32]int64, error)
	Delete(ctx context.Context, id uint64) error
	DeleteExpired(ctx context.Context) error
}

type AdminReportFilter struct {
	UserID        *uint64
	ReportType    *uint32
	Status        *uint32
	ProblemReport *uint32
	QaStatus      *uint32
	Severity      *uint32
	AssigneeID    *uint64
	Q             string
	Page          int
	Limit         int
}
