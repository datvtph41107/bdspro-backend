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

func (repo *FriendPostgre) UpdateStatusById(ctx context.Context, status enums.FriendStatus, id uint64, respondedAt *time.Time) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"responded_at": respondedAt,
	}).Error
}

func (repo *FriendPostgre) UpdateGroupById(ctx context.Context, isGroup bool, groupId uint64, id uint64) error {
	tx := repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id)
	if isGroup {
		return tx.Updates(map[string]interface{}{"group_id": groupId}).Error
	}
	return tx.Updates(map[string]interface{}{"group_receiver_id": groupId}).Error
}

func (repo *FriendPostgre) UpdateGroupId(ctx context.Context, id *uint64, groupId uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Update("group_id", groupId).Error
}

func (repo *FriendPostgre) UpdateGroupReceiverId(ctx context.Context, id uint, groupReceiverId uint) error {
	return repo.DB.Model(&domain.FriendEntity{}).Where("id = ?", id).Update("group_receiver_id", groupReceiverId).Error
}

func (repo *FriendPostgre) FindBySenderIdAndReceiverId(ctx context.Context, sender uint64, receiver uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("created_by = ? AND receiver_id = ? and status = ? and deleted_at is null", sender, receiver, enums.FriendStatusPending).First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) GetByID(ctx context.Context, id uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("id = ?", id).First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) Search(ctx context.Context, dto data.FriendDTO, createdBy *uint64, page int, pageSize int) ([]domain.FriendEntity, error) {
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

func (repo *FriendPostgre) SearchFriend(ctx context.Context, currentId uint64, dto data.FriendRequest) ([]domain.FriendEntity, int64, error) {
	var friends []domain.FriendEntity

	query := repo.DB.
		Model(&domain.FriendEntity{}).
		Debug().
		Preload("GroupUser").
		Preload("GroupReceiver").
		// Joins(`LEFT JOIN friend fr ON (
		// 	(fr.created_by = tb_contact.profile_id and fr.receiver_id = ?)
		// 	OR (fr.receiver_id = tb_contact.profile_id and fr.created_by = ?))
		// 	and fr.deleted_at is null`,
		// 	currentId, currentId).
		// Joins("LEFT JOIN profile_transfer ON (friend.created_by = profile_transfer.profile_id OR friend.receiver_id = profile_transfer.profile_id) "+
		// 	"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("status = 20 and deleted_at is null").
		Where(`((created_by = ? and receiver_id <> ?)
		 	OR (receiver_id = ? and created_by <> ?))`, currentId, currentId, currentId, currentId)

	// if dto.Text != "" {
	// 	query = query.Where("lower(full_name) LIKE ? OR phone LIKE ?", "%"+dto.Text+"%", "%"+dto.Text+"%")
	// }

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(dto.GetOffset()).Limit(dto.GetLimit()).Find(&friends).Error
	return friends, total, err
}

func (repo *FriendPostgre) GetFriendStatusByProfileIds(ctx context.Context, currentId uint64, profileIds []uint64) ([]domain.FriendItem, error) {
	var friends []domain.FriendItem

	query := repo.DB.
		Debug().
		Model(&domain.FriendEntity{}).
		Where("deleted_at is null").
		Where(`((created_by = ? and receiver_id in (?))
		 	OR (receiver_id = ? and created_by in (?)))`, currentId, profileIds, currentId, profileIds)

	err := query.Select("id, status, created_by, receiver_id").Find(&friends).Error
	return friends, err
}

func (repo *FriendPostgre) ExistsByCreatedByAndReceiverId(ctx context.Context, createdBy uint64, receiverId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("created_by = ? AND receiver_id = ? AND status <> ? and deleted_at is null", createdBy, receiverId, enums.FriendStatusReject).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgre) FindFriendRequest(ctx context.Context, user1 uint, user2 uint) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("deleted_at is null AND ((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?))", user1, user2, user2, user1).
		First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) GetSumFriend(ctx context.Context, userId uint) (int64, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("(created_by = ? OR receiver_id = ?) AND status = 2", userId, userId).
		Count(&count).Error
	return count, err
}

func (repo *FriendPostgre) FindByID(ctx context.Context, id uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.First(&request, id).Error
	return &request, err
}

func (repo *FriendPostgre) DeleteByID(ctx context.Context, id uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgre) BlockFriend(ctx context.Context, profileId uint64, targetId uint64) error {
	return repo.DB.Model(&domain.FriendEntity{}).
		Where("(created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *FriendPostgre) ExistsFriendship(ctx context.Context, id uint64, receiverID uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.FriendEntity{}).
		Where("id = ? AND receiver_id = ?", id, receiverID).
		Count(&count).Error
	return count > 0, err
}

func (repo *FriendPostgre) GetByUserID(ctx context.Context, userId uint64) (*domain.FriendEntity, error) {
	var request domain.FriendEntity
	err := repo.DB.Where("created_by = ? OR receiver_id = ?", userId, userId).First(&request).Error
	return &request, err
}

func (repo *FriendPostgre) GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*domain.FriendItem, error) {
	var status *domain.FriendItem

	err := repo.DB.
		Model(&domain.FriendItem{}).
		Select("id, status, created_by, receiver_id").
		Where("((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)) AND deleted_at IS NULL",
			currentId, targetId, targetId, currentId).
		Find(&status).Error

	if err != nil {
		return nil, err
	}

	return status, nil
}

func (repo *FriendPostgre) CancelRequest(ctx context.Context, id uint64) error {
	return repo.DB.
		Model(&domain.FriendEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}
