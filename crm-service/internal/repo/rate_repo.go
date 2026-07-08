package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type RateRepo interface {
	Create(ctx context.Context, rate *domain.Rate) (*domain.Rate, error)
	GetByID(ctx context.Context, id uint64) (*domain.Rate, error)
	DetailRate(ctx context.Context, id uint64) (*domain.Rate, error)
	GetByProfileIdAndOwnerIdAndOwnerOf(ctx context.Context, profileId uint64, ownerId uint64, ownerOf uint8) (*domain.Rate, error)
	GetRateList(ctx context.Context, id uint64, ownerOf uint32, req *dto.RateSearchDTO) ([]domain.Rate, int64, error)
	Update(ctx context.Context, id uint64, rate *domain.Rate) error
	Delete(ctx context.Context, id uint64) error
	UpdateHidden(ctx context.Context, id uint64, hidden bool) error
	CountInfo(ctx context.Context, ownerID uint64, ownerOf uint32) (domain.RateStats, error)
	GetHistoryRate(ctx context.Context, rateId uint64, page, size uint32) ([]domain.Rate, int64, error)
}