package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type CostDocumentUsecase struct {
	Repo repo.CostDocumentRepo
}

func NewCostDocumentUc(repo repo.CostDocumentRepo) CostDocumentUsecase {
	return CostDocumentUsecase{
		Repo: repo,
	}
}

func (uc CostDocumentUsecase) Create(c context.Context, data *domain.CostDocument) (*domain.CostDocument, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc CostDocumentUsecase) Update(c context.Context, id uint64, data *domain.CostDocument) (*domain.CostDocument, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc CostDocumentUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc CostDocumentUsecase) GetData(c context.Context, id uint64) ([]domain.CostDocument, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc CostDocumentUsecase) GetByID(c context.Context, id uint64) (*domain.CostDocument, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
