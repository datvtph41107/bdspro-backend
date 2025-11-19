package providers

import (
	_dto "common/domain/dto"
	"context"
	models "user/internal/models"
)

type FeedbackProvider interface {
	GetUserRateStats(ctx context.Context, userId uint64) (*models.RateStats, error)
	GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error)
}
