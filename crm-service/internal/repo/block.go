package repo

import (
	"context"
	"crm/internal/domain"
)

type BlockRepo interface {
	GetAll(c context.Context) ([]domain.BlockEntity, error)
}
