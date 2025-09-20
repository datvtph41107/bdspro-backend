package repo

import (
	base_enum "base/enum"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ProjectRepo interface {
	SearchItem(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, int64, error)
	// Đếm số dự án theo ownerOf và ownerId
	CountByOwner(ctx context.Context, ownerOf base_enum.EOwnerOf, ownerId uint64) (uint32, error)
}
