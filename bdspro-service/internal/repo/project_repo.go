package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_enum "common/domain/enum"
	"context"
)

type ProjectRepo interface {
	SearchItem(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, int64, error)
	// Đếm số dự án theo ownerOf và ownerId
	CountByOwner(ctx context.Context, ownerOf _enum.EOwnerOf, ownerId uint64) (uint32, error)
}
