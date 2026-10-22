package maindb

import (
	entity "mdn4/modules/crm/entity"

	"gorm.io/gorm"
)

// @bind: mdn4/internal/interface/repo.IBlockRepo
type BlockPostgres struct {
	DB *gorm.DB
}

func NewBlockPostgres(DB *gorm.DB) *BlockPostgres {
	return &BlockPostgres{DB: DB}
}

func (repo BlockPostgres) ListBlock(profileId uint64) []uint64 {
	var blockedUsers []uint64
	repo.DB.Model(&entity.BlockEntity{}).Where("profile_id = ?", profileId).Pluck("blocked_id", &blockedUsers)
	return blockedUsers
}

func (repo BlockPostgres) BlockUser(profileId, blockId uint64) (*entity.BlockEntity, error) {
	block := entity.BlockEntity{ProfileId: profileId, BlockedId: blockId}
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		FirstOrCreate(&block).Error
	return &block, err
}

func (repo BlockPostgres) UnblockUser(profileId, blockId uint64) (uint64, error) {
	err := repo.DB.
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Delete(&entity.BlockEntity{}).Error

	return blockId, err
}

func (repo BlockPostgres) IsBlocked(profileId, blockId uint64) bool {
	var count int64
	repo.DB.Model(&entity.BlockEntity{}).
		Where("profile_id = ? AND blocked_id = ?", profileId, blockId).
		Count(&count)

	return count > 0
}
