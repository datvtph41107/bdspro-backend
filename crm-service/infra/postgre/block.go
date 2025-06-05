package postgre

import (
	"crm/internal/domain"

	"gorm.io/gorm"
)

type BlockRepo struct {
	DB *gorm.DB
}

func NewBlockRepo(DB *gorm.DB) *BlockRepo {
	return &BlockRepo{DB: DB}
}

func (repo BlockRepo) ListBlock(profileId uint64) []uint64 {
	var blockedUsers []uint64
	repo.DB.Model(&domain.BlockEntity{}).Where("profile_id = ?", profileId).Pluck("blocked_id", &blockedUsers)
	return blockedUsers
}

func (repo BlockRepo) BlockUser(profileId, blockId uint64) (*domain.BlockEntity, error) {
	block := domain.BlockEntity{ProfileId: profileId, BlockedId: blockId}
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		FirstOrCreate(&block).Error
	return &block, err
}

func (repo BlockRepo) UnblockUser(profileId, blockId uint64) (uint64, error) {
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Delete(&domain.BlockEntity{}).Error

	return blockId, err
}

func (repo BlockRepo) IsBlocked(profileId, blockId uint64) bool {
	var count int64
	repo.DB.Model(&domain.BlockEntity{}).
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Count(&count)

	return count > 0
}
