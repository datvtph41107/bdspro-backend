package repo

import (
	_dto "common/dto"
	"context"
	"user/internal/dto"
	user_model "user/internal/models/user"
)

type RateRepo interface {
	CreateRate(ctx context.Context, rate *user_model.Rate) (*user_model.Rate, error)
	GetRateList(ctx context.Context, targetId uint64, targetType uint8, dto *dto.RateSearchDTO) ([]user_model.RateItem, int64, error)
	UpdateRate(ctx context.Context, rate *user_model.Rate, oldRate *user_model.Rate) (*user_model.Rate, error)
	DeleteRate(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*user_model.Rate, error)
	GetByProfileIdAndTargetIdAndTargetType(ctx context.Context, profileId uint64, targetId uint64, targetType uint8) (*user_model.Rate, error)
	GetHistoryRate(ctx context.Context, rateId uint64, dto *_dto.Pagable) ([]user_model.RateItem, int64, error)
	UpdateHidden(ctx context.Context, id uint64, hidden bool) error
	CountInfo(ctx context.Context, targetId uint64, targetType uint8) (user_model.RateStats, error)
}
