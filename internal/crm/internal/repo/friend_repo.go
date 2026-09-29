package repo

import (
	"context"
	"crm/data"
	"crm/internal/domain"
	"crm/internal/enums"
	"time"
)

type FriendRepo interface {
	GetByID(ctx context.Context, id uint64) (*domain.FriendEntity, error)
	UpdateStatusById(ctx context.Context, status enums.FriendStatus, id uint64, respondedAt *time.Time) error
	UpdateGroupById(ctx context.Context, isGroup bool, groupId uint64, id uint64) error
	BlockFriend(ctx context.Context, profileId uint64, targetId uint64) error
	GetByUserID(ctx context.Context, userId uint64) (*domain.FriendEntity, error)
	UpdateGroupId(ctx context.Context, id *uint64, groupId uint64) error
	UpdateGroupReceiverId(ctx context.Context, id uint, groupReceiverId uint) error
	FindBySenderIdAndReceiverId(ctx context.Context, sender uint64, receiver uint64) (*domain.FriendEntity, error)
	Search(ctx context.Context, dto data.FriendDTO, createdBy *uint64, page int, pageSize int) ([]domain.FriendEntity, error)
	SearchFriend(ctx context.Context, currentId uint64, dto data.FriendRequest) ([]domain.FriendEntity, int64, error)
	GetFriendStatusByProfileIds(ctx context.Context, currentId uint64, profileIds []uint64) ([]domain.FriendItem, error)
	ExistsByCreatedByAndReceiverId(ctx context.Context, createdBy uint64, receiverId uint64) (bool, error)
	FindFriendRequest(ctx context.Context, user1 uint, user2 uint) (*domain.FriendEntity, error)
	GetSumFriend(ctx context.Context, userId uint) (int64, error)
	FindByID(ctx context.Context, id uint64) (*domain.FriendEntity, error)
	DeleteByID(ctx context.Context, id uint64) error
	ExistsFriendship(ctx context.Context, id uint64, receiverID uint64) (bool, error)
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*domain.FriendItem, error)
	CancelRequest(ctx context.Context, id uint64) error
	GetCommonFriends(ctx context.Context, user1Id uint64, user2Id uint64, limit int) ([]uint64, int64, error)
	GetFriendship(ctx context.Context, userID1, userID2 uint64) (*domain.FriendEntity, error)
}