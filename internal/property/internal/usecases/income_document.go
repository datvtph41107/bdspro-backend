package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type IncomeDocumentUsecase struct {
	Repo repo.IncomeDocumentRepo
}

func NewIncomeDocumentUc(repo repo.IncomeDocumentRepo) IncomeDocumentUsecase {
	return IncomeDocumentUsecase{
		Repo: repo,
	}
}

func (uc IncomeDocumentUsecase) Create(c context.Context, data *domain.IncomeDocument) (*domain.IncomeDocument, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc IncomeDocumentUsecase) Update(c context.Context, id uint64, data *domain.IncomeDocument) (*domain.IncomeDocument, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc IncomeDocumentUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc IncomeDocumentUsecase) GetData(c context.Context, id uint64) ([]domain.IncomeDocument, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc IncomeDocumentUsecase) GetByID(c context.Context, id uint64) (*domain.IncomeDocument, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
