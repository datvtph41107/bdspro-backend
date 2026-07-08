package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type ReportRepo interface {
	Create(ctx context.Context, report *domain.Report) (*domain.Report, error)
	GetByID(ctx context.Context, id uint64) (*domain.Report, error)
	GetByCreatedBy(ctx context.Context, createdBy uint64, req *dto.ReportListRequest) ([]domain.Report, int64, error)
	GetList(ctx context.Context, req *dto.ReportListRequest) ([]domain.Report, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	ExistByOwnerIdAndUserIdAndOwnerOf(ctx context.Context, ownerId uint64, userId uint64, ownerOf uint32) (bool, error)
	GetIDByOwnerIdAndUserIdAndOwnerOf(ctx context.Context, ownerId uint64, userId uint64, ownerOf uint32) (*uint64, error)
}