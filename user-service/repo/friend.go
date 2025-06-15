package repo

import (
	"time"
	"user/data"
	"user/enums"
	user_model "user/models/user"

	"gorm.io/gorm"
)

type FriendRepo struct {
	DB *gorm.DB
}

func NewFriendRepo(DB *gorm.DB) *FriendRepo {
	return &FriendRepo{
		DB: DB,
	}
}

func (repo *FriendRepo) UpdateStatusById(status enums.EFriendStatus, id uint64, respondedAt *time.Time) error {
	return repo.DB.Model(&user_model.FriendEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"responded_at": respondedAt,
	}).Error
}

func (repo *FriendRepo) UpdateGroupById(isGroup bool, groupId uint64, id uint64) error {
	tx := repo.DB.Model(&user_model.FriendEntity{}).Where("id = ?", id)
	if isGroup {
		return tx.Updates(map[string]interface{}{"group_id": groupId}).Error
	}
	return tx.Updates(map[string]interface{}{"group_receiver_id": groupId}).Error
}

func (repo *FriendRepo) UpdateGroupId(id *uint64, groupId uint64) error {
	return repo.DB.Model(&user_model.FriendEntity{}).Where("id = ?", id).Update("group_id", groupId).Error
}

func (repo *FriendRepo) UpdateGroupReceiverId(id uint, groupReceiverId uint) error {
	return repo.DB.Model(&user_model.FriendEntity{}).Where("id = ?", id).Update("group_receiver_id", groupReceiverId).Error
}

func (repo *FriendRepo) FindBySenderIdAndReceiverId(sender uint64, receiver uint64) (*user_model.FriendEntity, error) {
	var request user_model.FriendEntity
	err := repo.DB.Where("created_by = ? AND receiver_id = ?", sender, receiver).First(&request).Error
	return &request, err
}

func (repo *FriendRepo) GetById(id uint) (*user_model.FriendEntity, error) {
	var request user_model.FriendEntity
	err := repo.DB.Where("id = ?", id).First(&request).Error
	return &request, err
}

func (repo *FriendRepo) Search(dto data.FriendDTO, createdBy *uint64, page int, pageSize int) ([]user_model.FriendEntity, error) {
	var requests []user_model.FriendEntity
	query := repo.DB.Model(&user_model.FriendEntity{})

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

func (repo *FriendRepo) SearchFriend(currentId uint64, dto data.FriendDTO, page int, pageSize int) ([]user_model.FriendEntity, error) {
	var friends []user_model.FriendEntity

	query := repo.DB.
		Model(&user_model.FriendEntity{}).
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins(`LEFT JOIN profile_transfer ON (friend.created_by = profile_transfer.profile_id OR friend.receiver_id = profile_transfer.profile_id)
				AND profile_transfer.profile_id <> ? `, currentId).
		Where("friend.status = 20 and friend.deleted_at is null")

	if dto.Text != nil && *dto.Text != "" {
		query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FriendRepo) ExistsByCreatedByAndReceiverId(createdBy uint64, receiverId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&user_model.FriendEntity{}).
		Where("created_by = ? AND receiver_id = ? AND status <> ?", createdBy, receiverId, enums.EFriendStatusReject).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendRepo) FindFriendRequest(user1 uint, user2 uint) (*user_model.FriendEntity, error) {
	var request user_model.FriendEntity
	err := repo.DB.Where("deleted_at is null AND ((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?))", user1, user2, user2, user1).
		First(&request).Error
	return &request, err
}

func (repo *FriendRepo) GetSumFriend(userId uint) (int64, error) {
	var count int64
	err := repo.DB.Model(&user_model.FriendEntity{}).
		Where("(created_by = ? OR receiver_id = ?) AND status = 2", userId, userId).
		Count(&count).Error
	return count, err
}

func (repo *FriendRepo) FindByID(id uint64) (*user_model.FriendEntity, error) {
	var request user_model.FriendEntity
	err := repo.DB.First(&request, id).Error
	return &request, err
}

func (repo *FriendRepo) DeleteByID(id uint64) error {
	return repo.DB.Model(&user_model.FriendEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendRepo) BlockFriend(profileId uint64, targetId uint64) error {
	return repo.DB.Model(&user_model.FriendEntity{}).
		Where("(created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendRepo) ExistsFriendship(id uint64, receiverID uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&user_model.FriendEntity{}).
		Where("id = ? AND receiver_id = ?", id, receiverID).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendRepo) GetRelationShip(currentId uint64, targetId uint64) (*user_model.FriendItem, error) {
	var status *user_model.FriendItem

	err := repo.DB.
		Model(&user_model.FriendItem{}).
		Select("id, status, created_by, receiver_id").
		Where("((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)) AND deleted_at IS NULL",
			currentId, targetId, targetId, currentId).
		Find(&status).Error

	if err != nil {
		return nil, err
	}

	return status, nil
}
