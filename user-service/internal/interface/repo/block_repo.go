package repo

import (
	models "user/internal/models"
)

type IBlockRepo interface {
	ListBlock(profileId uint64) []uint64
	BlockUser(profileId, blockId uint64) (*models.BlockEntity, error)
	UnblockUser(profileId, blockId uint64) (uint64, error)
	IsBlocked(profileId, blockId uint64) bool
}
