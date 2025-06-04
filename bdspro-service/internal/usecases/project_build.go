package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type ProjectBuildUsecase struct {
	Repo repo.ProjectBuildRepo
}

func NewProjectBuildUsecase(repo repo.ProjectBuildRepo) *ProjectBuildUsecase {
	return &ProjectBuildUsecase{
		Repo: repo,
	}
}

func (u *ProjectBuildUsecase) GetAll(c context.Context) ([]domain.ProjectBuild, error) {
	return u.Repo.GetAll(c)
}

func (u *ProjectBuildUsecase) GetByID(c context.Context, id uint64) (*domain.ProjectBuild, error) {
	return u.Repo.GetByID(c, id)
}
