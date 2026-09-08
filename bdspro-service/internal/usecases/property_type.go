package usecases

import "bdspro/internal/repo"

type PropertyTypeUsecase struct {
	Repo *repo.PropertyTypeRepo
}

func NewPropertyTypeUsecase(repo *repo.PropertyTypeRepo) *PropertyTypeUsecase {
	return &PropertyTypeUsecase{
		Repo: repo,
	}
}

// todo: crud
