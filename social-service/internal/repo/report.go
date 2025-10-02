package repo

import (
	_dto "common/domain/dto"
	"context"
	"social/internal/domain"
	"social/internal/enums"
)

type ReportRepo interface {
	Create(ctx context.Context, report *domain.Report) (*domain.Report, error)
	ExistByTargetIdAndUserIdAndTargetType(ctx context.Context, targetId uint64, userId uint64, targetType enums.TargetType) (bool, error)
	GetByCreatedBy(ctx context.Context, createdBy uint64, dto *_dto.Pagable) ([]domain.Report, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.ReportStatus) error
}
