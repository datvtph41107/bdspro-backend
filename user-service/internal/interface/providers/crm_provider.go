package providers

import (
	_dto "common/domain/dto"
	"context"
	models "user/internal/models"
)

type CrmProvider interface {
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*models.FriendItem, error)
	GetRelationShipV3(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendV3DTO, error)
	GetUserRateStats(ctx context.Context, userId uint64) (*models.RateStats, error)
	GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error)
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
}
