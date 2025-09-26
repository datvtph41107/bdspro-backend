package providers

import (
	"context"
	user_model "user/internal/models/user"
)

type FriendProvider interface {
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*user_model.FriendItem, error)
}
