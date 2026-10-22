package repo

import (
	entity "mdn4/modules/crm/entity"
)

type IBlockRepo interface {
	ListBlock(profileId uint64) []uint64
	BlockUser(profileId, blockId uint64) (*entity.BlockEntity, error)
	UnblockUser(profileId, blockId uint64) (uint64, error)
	IsBlocked(profileId, blockId uint64) bool
}
