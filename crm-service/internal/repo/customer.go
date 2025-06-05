package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type CustomerRepo interface {
	Search(c context.Context, dto dto.CustomerSearchDTO) ([]domain.CustomerEntity, int64, error)
	Create(c context.Context, entity *domain.CustomerEntity) (*domain.CustomerEntity, error)
	Update(c context.Context, id uint64, entity *domain.CustomerEntity) (*domain.CustomerEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.CustomerEntity, error)
	Existed(c context.Context, entity *domain.CustomerEntity) (bool, error)
	UpdateNote(c context.Context, id uint64, note string) (*domain.CustomerEntity, error)
	Assign(c context.Context, id uint64, entity *domain.CustomerEntity) (*domain.CustomerEntity, error)
	SwitchStage(c context.Context, id uint64, stageID *uint64) (*domain.CustomerEntity, error)
}
