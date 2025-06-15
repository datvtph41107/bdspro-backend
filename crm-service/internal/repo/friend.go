package repo

import (
	"context"
	"crm/internal/enums"
	"time"
)

type FriendRepo interface {
	// GetAll(c context.Context) ([]domain.FriendEntity, error)
	// GetByID(c context.Context, id uint64) (*domain.FriendEntity, error)
	// Create(c context.Context, entity *domain.FriendEntity) error
	// Update(c context.Context, id uint64, entity *domain.FriendEntity) error
	// Delete(c context.Context, id uint64) error
	UpdateStatusById(c context.Context, status enums.FriendStatus, id uint64, respondedAt *time.Time) error
	UpdateGroupById(c context.Context, isGroup bool, groupId uint64, id uint64) error
	BlockFriend(c context.Context, profileId uint64, targetId uint64) error
	// ExistByProfile(c context.Context, profileId uint64) (bool, error)
	// ExistByProfileAndFriend(c context.Context, profileId uint64, friendId uint64) (bool, error)
}
