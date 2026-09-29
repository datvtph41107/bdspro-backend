package repo

import (
	"bdspro/internal/domain"
	"context"
)

type IncomeDocumentRepo interface {
	Create(c context.Context, entity *domain.IncomeDocument) error
	Update(c context.Context, id uint64, entity *domain.IncomeDocument) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.IncomeDocument, error)
	GetByID(c context.Context, id uint64) (*domain.IncomeDocument, error)
}
