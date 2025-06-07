package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type ContactRepo interface {
	Search(c context.Context, profileId uint64, dto dto.ContactSearchDTO) ([]domain.ContactEntity, int64, error)
	Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	Update(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.ContactEntity, error)
}
