package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
)

type ContactRepo interface {
	Search(c context.Context, profileId uint64, ownerType enums.EOwnerType, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error)
	Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	Update(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.ContactEntity, error)
	GetByProfileID(c context.Context, profileId uint64, currentUserId uint64, ownerType enums.EOwnerType) (*domain.ContactEntity, error)
	Existed(c context.Context, followId uint64) error
	Bulk(c context.Context, entities []domain.ContactEntity) ([]domain.ContactEntity, error)
	ExistByProfile(c context.Context, profileId uint64) (bool, error)
	UpdateNote(c context.Context, id uint64, note string) error
}
