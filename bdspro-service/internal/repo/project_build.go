package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ProjectBuildRepo interface {
	GetAll(c context.Context) ([]domain.ProjectBuild, error)
	GetByID(c context.Context, id uint64) (*domain.ProjectBuild, error)
}
