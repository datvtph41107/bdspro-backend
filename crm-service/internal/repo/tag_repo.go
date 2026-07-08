package repo

import (
	_dto "common/domain/dto"
	"context"
	"crm/internal/domain"
)

type TagRepo interface {
	ListByContactID(ctx context.Context, contactID uint64) ([]domain.TagEntity, error)
	ReplaceContactTags(ctx context.Context, contactID uint64, tags []domain.TagEntity) error
	Create(ctx context.Context, tag *domain.TagEntity) error
	Update(ctx context.Context, id uint64, tag *domain.TagEntity) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.TagEntity, error)
	ListForContact(ctx context.Context, ownerID uint64, ownerOf int32, pagable _dto.IPagable, searchText string) ([]domain.TagEntity, uint32, error)
	ListAdmin(ctx context.Context, pagable _dto.IPagable, searchText string) ([]domain.TagEntity, uint32, error)
}