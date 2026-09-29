package postgres

import (
	_dto "common/domain/dto"
	"context"
	"time"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/repo"
	models "user/internal/models"

	"gorm.io/gorm"
)

type FriendPostgres struct {
	DB *gorm.DB
}

func NewFriendPostgres(DB *gorm.DB) repo.IFriendRepo {
	return &FriendPostgres{
		DB: DB,
	}
}

func (repo *FriendPostgres) Create(ctx context.Context, friend *models.FriendEntity) error {
	return repo.DB.WithContext(ctx).Create(friend).Error
}

func (repo *FriendPostgres) UpdateStatusById(status enums.EFriendStatus, id uint64, respondedAt *time.Time) error {
	return repo.DB.Model(&models.FriendEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"responded_at": respondedAt,
	}).Error
}

func (repo *FriendPostgres) UpdateGroupById(isGroup bool, groupId uint64, id uint64) error {
	tx := repo.DB.Model(&models.FriendEntity{}).Where("id = ?", id)
	if isGroup {
		return tx.Updates(map[string]interface{}{"group_id": groupId}).Error
	}
	return tx.Updates(map[string]interface{}{"group_receiver_id": groupId}).Error
}

func (repo *FriendPostgres) UpdateGroupId(id *uint64, groupId uint64) error {
	return repo.DB.Model(&models.FriendEntity{}).Where("id = ?", id).Update("group_id", groupId).Error
}

func (repo *FriendPostgres) UpdateGroupReceiverId(id uint, groupReceiverId uint) error {
	return repo.DB.Model(&models.FriendEntity{}).Where("id = ?", id).Update("group_receiver_id", groupReceiverId).Error
}

func (repo *FriendPostgres) FindBySenderIdAndReceiverId(sender uint64, receiver uint64) (*models.FriendEntity, error) {
	var request models.FriendEntity
	err := repo.DB.Where("created_by = ? AND receiver_id = ?", sender, receiver).First(&request).Error
	return &request, err
}

func (repo *FriendPostgres) GetById(id uint) (*models.FriendEntity, error) {
	var request models.FriendEntity
	err := repo.DB.Where("id = ?", id).First(&request).Error
	return &request, err
}

func (repo *FriendPostgres) Search(dto dto.FriendDTO, createdBy *uint64, page int, pageSize int) ([]models.FriendEntity, error) {
	var requests []models.FriendEntity
	query := repo.DB.Model(&models.FriendEntity{})

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

func (repo *FriendPostgres) SearchFriend(currentId uint64, dto dto.FriendDTO, page int, pageSize int) ([]models.FriendEntity, error) {
	var friends []models.FriendEntity

	query := repo.DB.
		Model(&models.FriendEntity{}).
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins(`LEFT JOIN profile_transfer ON (friend.created_by = profile_transfer.profile_id OR friend.receiver_id = profile_transfer.profile_id)
				AND profile_transfer.profile_id <> ? `, currentId).
		Where("friend.status = 20 and friend.deleted_at is null")

	if dto.Text != nil && *dto.Text != "" {
		query = query.Where("profile_transfer.full_name ILIKE ? OR profile_transfer.phone ILIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FriendPostgres) ExistsByCreatedByAndReceiverId(createdBy uint64, receiverId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.FriendEntity{}).
		Where("created_by = ? AND receiver_id = ? AND status <> ?", createdBy, receiverId, enums.EFriendStatusReject).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgres) FindFriendRequest(user1 uint, user2 uint) (*models.FriendEntity, error) {
	var request models.FriendEntity
	err := repo.DB.Where("deleted_at is null AND ((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?))", user1, user2, user2, user1).
		First(&request).Error
	return &request, err
}

func (repo *FriendPostgres) GetSumFriend(userId uint) (int64, error) {
	var count int64
	err := repo.DB.Model(&models.FriendEntity{}).
		Where("(created_by = ? OR receiver_id = ?) AND status = 2", userId, userId).
		Count(&count).Error
	return count, err
}

func (repo *FriendPostgres) FindByID(id uint64) (*models.FriendEntity, error) {
	var request models.FriendEntity
	err := repo.DB.First(&request, id).Error
	return &request, err
}

func (repo *FriendPostgres) DeleteByID(id uint64) error {
	return repo.DB.Model(&models.FriendEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgres) BlockFriend(profileId uint64, targetId uint64) error {
	return repo.DB.Model(&models.FriendEntity{}).
		Where("(created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgres) ExistsFriendship(id uint64, receiverID uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.FriendEntity{}).
		Where("id = ? AND receiver_id = ?", id, receiverID).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgres) GetRelationShip(currentId uint64, targetId uint64) (*_dto.FriendItemDTO, error) {
	var status *_dto.FriendItemDTO

	err := repo.DB.
		Model(&_dto.FriendItemDTO{}).
		Select("id, status, created_by, receiver_id").
		Where("((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)) AND deleted_at IS NULL",
			currentId, targetId, targetId, currentId).
		Find(&status).Error

	if err != nil {
		return nil, err
	}

	return status, nil
}
