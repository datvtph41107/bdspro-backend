package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupDealRepository interface {
	Create(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	Update(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.Deal, error)
	GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.Deal, uint32, error)
	UpdateStatus(ctx context.Context, id uint32, status entity.DealStatus) error
} 