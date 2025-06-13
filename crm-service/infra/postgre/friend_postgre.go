package postgre

import (
	"context"
	"crm/data"
	"crm/internal/domain"
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.FriendRepo
type FriendPostgre struct {
	DB *gorm.DB
}

func NewFriendRepo(DB *gorm.DB) *FriendPostgre {
	return &FriendPostgre{
		DB: DB,
	}
}

func (repo *FriendPostgre) UpdateStatusById(c context.Context, status enums.FriendStatus, id uint64, respondedAt *time.Time) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"responded_at": respondedAt,
	}).Error
}

func (repo *FriendPostgre) UpdateGroupById(c context.Context, isGroup bool, groupId uint64, id uint64) error {
	tx := repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id)
	if isGroup {
		return tx.Updates(map[string]interface{}{"group_id": groupId}).Error
	}
	return tx.Updates(map[string]interface{}{"group_receiver_id": groupId}).Error
}

func (repo *FriendPostgre) UpdateGroupId(id *uint64, groupId uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Update("group_id", groupId).Error
}

func (repo *FriendPostgre) UpdateGroupReceiverId(id uint, groupReceiverId uint) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Update("group_receiver_id", groupReceiverId).Error
}

func (repo *FriendPostgre) FindBySenderIdAndReceiverId(sender uint64, receiver uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("created_by = ? AND receiver_id = ?", sender, receiver).First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) GetById(id uint) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("id = ?", id).First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) Search(dto data.FriendDTO, createdBy *uint64, page int, pageSize int) ([]domain.FriendEntity, error) {
	var requests []domain.FriendEntity
	query := repo.DB.Model(&domain.FriendEntity{})

	if dto.ReceiverID != 0 {
		query = query.Where("receiver_id = ?", dto.ReceiverID)
	}
	if createdBy != nil {
		query = query.Where("created_by = ?", createdBy)
	}
	if dto.Status != 0 {
		query = query.Where("status = ?", dto.Status)
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&requests).Error
	return requests, err
}

func (repo *FriendPostgre) SearchFriend(currentId uint64, dto data.FriendDTO, page int, pageSize int) ([]domain.FriendEntity, error) {
	var friends []domain.FriendEntity

	query := repo.DB.
		Model(&domain.FriendEntity{}).
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins("LEFT JOIN profile_transfer ON (friend.created_by = profile_transfer.profile_id OR friend.receiver_id = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("friend.status = 20 and friend.deleted_at is null")

	if dto.Text != nil && *dto.Text != "" {
		query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FriendPostgre) ExistsByCreatedByAndReceiverId(createdBy uint64, receiverId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("created_by = ? AND receiver_id = ? AND status <> ?", createdBy, receiverId, enums.Reject).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgre) FindFriendRequest(user1 uint, user2 uint) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("deleted_at is null AND ((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?))", user1, user2, user2, user1).
		First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) GetSumFriend(userId uint) (int64, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("(created_by = ? OR receiver_id = ?) AND status = 2", userId, userId).
		Count(&count).Error
	return count, err
}

func (repo *FriendPostgre) FindByID(id uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.First(&request, id).Error
	return &request, err
}

func (repo *FriendPostgre) DeleteByID(id uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgre) BlockFriend(c context.Context, profileId uint64, targetId uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).
		Where("(created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgre) ExistsFriendship(id uint64, receiverID uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("id = ? AND receiver_id = ?", id, receiverID).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgre) GetByUserID(c context.Context, userId uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("created_by = ? OR receiver_id = ?", userId, userId).First(&request).Error
	return &request, err
}
