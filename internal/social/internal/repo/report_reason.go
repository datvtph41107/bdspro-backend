package repo

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
)

type ReportReasonRepo interface {
	Create(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error)
	Update(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error)
	Delete(ctx context.Context, id uint64) error
	GetList(ctx context.Context, dto *dto.ReportReasonRequest) ([]domain.ReportReason, int64, error)
	GetByID(ctx context.Context, id uint64) (*domain.ReportReason, error)
	// -- public --
	GetPublic(ctx context.Context) ([]domain.ReportReason, error)
	ExistedByID(ctx context.Context, id uint64) (bool, error)
}
