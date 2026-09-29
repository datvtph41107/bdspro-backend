package repo

import (
	"context"
)

type FriendTagRepo interface {
	UpdateFriendTag(ctx context.Context, newsFeedID uint64, friendTagIds []uint64) error
}
