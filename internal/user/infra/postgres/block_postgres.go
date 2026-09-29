package postgres

import (
	models "user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IBlockRepo
type BlockPostgres struct {
	DB *gorm.DB
}

func NewBlockPostgres(DB *gorm.DB) *BlockPostgres {
	return &BlockPostgres{DB: DB}
}

func (repo BlockPostgres) ListBlock(profileId uint64) []uint64 {
	var blockedUsers []uint64
	repo.DB.Model(&models.BlockEntity{}).Where("profile_id = ?", profileId).Pluck("blocked_id", &blockedUsers)
	return blockedUsers
}

func (repo BlockPostgres) BlockUser(profileId, blockId uint64) (*models.BlockEntity, error) {
	block := models.BlockEntity{ProfileId: profileId, BlockedId: blockId}
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		FirstOrCreate(&block).Error
	return &block, err
}

func (repo BlockPostgres) UnblockUser(profileId, blockId uint64) (uint64, error) {
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Delete(&models.BlockEntity{}).Error

	return blockId, err
}

func (repo BlockPostgres) IsBlocked(profileId, blockId uint64) bool {
	var count int64
	repo.DB.Model(&models.BlockEntity{}).
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Count(&count)

	return count > 0
}
