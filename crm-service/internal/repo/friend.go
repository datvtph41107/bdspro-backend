package repo

import (
	"context"
	"crm/internal/domain"
)

type FriendRepo interface {
	GetAll(c context.Context) ([]domain.FriendEntity, error)
	GetByID(c context.Context, id uint64) (*domain.FriendEntity, error)
	Create(c context.Context, entity *domain.FriendEntity) error
	Update(c context.Context, id uint64, entity *domain.FriendEntity) error
	Delete(c context.Context, id uint64) error
	ExistByProfile(c context.Context, profileId uint64) (bool, error)
	ExistByProfileAndFriend(c context.Context, profileId uint64, friendId uint64) (bool, error)
}
