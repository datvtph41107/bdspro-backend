package repo

import (
	base_enum "base/enum"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type ContactRepo interface {
	Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error)
	Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	UpdateContactID(c context.Context, contactId uint64, leadId uint64) (*domain.ContactEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.ContactEntity, error)
	GetByProfileID(c context.Context, profileId uint64, currentUserId uint64, ownerType base_enum.EOwnerOf) (*domain.ContactEntity, error)
	Existed(c context.Context, followId uint64) error
	Bulk(c context.Context, entities []domain.ContactEntity) ([]domain.ContactEntity, error)
	ExistByProfile(c context.Context, profileId uint64) (bool, error)
	UpdateNote(c context.Context, id uint64, note string) error
	GetByPhone(c context.Context, profileId uint64, ownerType base_enum.EOwnerOf, phone string) (*domain.ContactEntity, error)
}
