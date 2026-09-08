package usecases

import (
	_crud "common/domain/crud"
	"user/internal/interface/repo"
	"user/internal/models"
)

type PriceTableUsecase struct {
	_crud.BaseUsecase[models.PriceTableDomain, repo.IPriceTableRepo]
	priceTableRepo repo.IPriceTableRepo
}

type IPriceTableUsecase interface {
	_crud.IBaseUsecase[models.PriceTableDomain]
}

func NewPriceTableUsecase(priceTableRepo repo.IPriceTableRepo) IPriceTableUsecase {
	return &PriceTableUsecase{
		BaseUsecase: _crud.BaseUsecase[models.PriceTableDomain, repo.IPriceTableRepo]{
			Repo: priceTableRepo,
		},
		priceTableRepo: priceTableRepo,
	}
}
