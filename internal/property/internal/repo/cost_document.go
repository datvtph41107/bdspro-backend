package repo

import (
	"bdspro/internal/domain"
	"context"
)

type CostDocumentRepo interface {
	Create(c context.Context, entity *domain.CostDocument) error
	Update(c context.Context, id uint64, entity *domain.CostDocument) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.CostDocument, error)
	GetByID(c context.Context, id uint64) (*domain.CostDocument, error)
}
