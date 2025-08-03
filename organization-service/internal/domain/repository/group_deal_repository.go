package repository

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
)

type GroupDealRepository interface {
	Create(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	Update(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.Deal, error)
	DetailByID(ctx context.Context, id uint64) (*entity.Deal, error)
	GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error)
	GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]*entity.Deal, uint32, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.DealStatus) error
	UpdateStatusAndCancelReason(ctx context.Context, id uint64, status enums.DealStatus, cancelReason string) error
	UpdateSetting(ctx context.Context, id uint64, setting *dto.DealSetting) error
}
