package repo

import (
	_dto "common/domain/dto"
	"context"
	"time"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/enums"
	entity "userv4/modules/crm/entity"
)

type IFriendRepo interface {
	Create(ctx context.Context, friend *entity.FriendEntity) error
	UpdateStatusById(status enums.FriendStatus, id uint64, respondedAt *time.Time) error
	UpdateGroupById(isGroup bool, groupId uint64, id uint64) error
	UpdateGroupId(id *uint64, groupId uint64) error
	UpdateGroupReceiverId(id uint, groupReceiverId uint) error
	FindBySenderIdAndReceiverId(sender uint64, receiver uint64) (*entity.FriendEntity, error)
	GetById(id uint) (*entity.FriendEntity, error)
	Search(dto dto.FriendDTO, createdBy *uint64, page int, pageSize int) ([]entity.FriendEntity, error)
	SearchFriend(currentId uint64, dto dto.FriendDTO, page int, pageSize int) ([]entity.FriendEntity, error)
	ExistsByCreatedByAndReceiverId(createdBy uint64, receiverId uint64) (bool, error)
	FindFriendRequest(user1 uint, user2 uint) (*entity.FriendEntity, error)
	GetSumFriend(userId uint) (int64, error)
	FindByID(id uint64) (*entity.FriendEntity, error)
	DeleteByID(id uint64) error
	BlockFriend(profileId uint64, targetId uint64) error
	ExistsFriendship(id uint64, receiverID uint64) (bool, error)
	GetRelationShip(currentId uint64, targetId uint64) (*_dto.FriendItemDTO, error)
	// CancelRequest(id uint64) error
}
