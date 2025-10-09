package providers

import (
	"context"
	models "user/internal/models"
)

type FriendProvider interface {
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*models.FriendItem, error)
}
