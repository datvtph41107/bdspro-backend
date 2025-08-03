package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ProjectBuildRepo interface {
	GetAll(c context.Context) ([]domain.ProjectBuild, error)
	GetByID(c context.Context, id uint64) (*domain.ProjectBuild, error)
	Search(c context.Context, dto *dto.TextSearchRequest) ([]domain.ProjectBuild, int64, error)
	Create(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error)
	Update(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error)
	Delete(c context.Context, id uint64) error
}
