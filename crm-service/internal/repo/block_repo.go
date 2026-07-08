package repo

import (
	"context"
	"crm/internal/domain"
)

type BlockRepo interface {
	// GetAll(c context.Context) ([]domain.BlockEntity, error)
	ListBlock(ctx context.Context, profileId uint64) []uint64
	BlockUser(ctx context.Context, profileId, blockId uint64) (*domain.BlockEntity, error)
	UnblockUser(ctx context.Context, profileId, blockId uint64) (uint64, error)
	IsBlocked(ctx context.Context, profileId, blockId uint64) bool
}