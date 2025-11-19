package providers

import (
	_dto "common/domain/dto"
	"context"
	models "user/internal/models"
)

type FriendProvider interface {
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*models.FriendItem, error)
	GetRelationShipV3(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendV3DTO, error)
}
