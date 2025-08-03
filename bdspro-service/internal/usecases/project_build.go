package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
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

func (u *ProjectBuildUsecase) Search(c context.Context, dto *dto.TextSearchRequest) ([]domain.ProjectBuild, int64, error) {
	return u.Repo.Search(c, dto)
}

func (u *ProjectBuildUsecase) Create(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error) {
	return u.Repo.Create(c, dto)
}

func (u *ProjectBuildUsecase) Update(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error) {
	return u.Repo.Update(c, dto)
}

func (u *ProjectBuildUsecase) Delete(c context.Context, id uint64) error {
	return u.Repo.Delete(c, id)
}