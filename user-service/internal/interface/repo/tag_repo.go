package repo

import (
	_dto "common/domain/dto"
	"context"
	"user/enums"
	"user/internal/models"
)

type ITagRepo interface {
	ListByProfileID(ctx context.Context, profileID uint64) ([]models.TagEntity, error)
	ReplaceProfileTags(ctx context.Context, profileID *uint64, tags []models.TagEntity) error
	Create(ctx context.Context, tag *models.TagEntity) error
	Update(ctx context.Context, id uint64, tag *models.TagEntity) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.TagEntity, error)
	ListForUser(ctx context.Context, profileID uint64, tagType enums.TagType, pagable _dto.IPagable, searchText string) ([]models.TagEntity, uint32, error)
	ListAdmin(ctx context.Context, tagTypes []enums.TagType, pagable _dto.IPagable, searchText string, isDefault *bool) ([]models.TagEntity, uint32, error)
}
