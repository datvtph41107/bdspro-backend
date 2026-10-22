package repo

import (
	_dto "common/domain/dto"
	"context"
	"mdn4/modules/user/domain/user"
	"mdn4/modules/user/enums"
)

type ITagRepo interface {
	ListByProfileID(ctx context.Context, profileID uint64) ([]user.TagEntity, error)
	ReplaceProfileTags(ctx context.Context, profileID *uint64, tags []user.TagEntity) error
	Create(ctx context.Context, tag *user.TagEntity) error
	Update(ctx context.Context, id uint64, tag *user.TagEntity) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*user.TagEntity, error)
	ListForUser(ctx context.Context, profileID uint64, tagType enums.TagType, pagable _dto.IPagable, searchText string) ([]user.TagEntity, uint32, error)
	ListAdmin(ctx context.Context, tagTypes []enums.TagType, pagable _dto.IPagable, searchText string, isDefault *bool) ([]user.TagEntity, uint32, error)
}
