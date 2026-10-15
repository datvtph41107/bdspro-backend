package services

import (
	_crud "common/domain/crud"
	"userv4/internal/interface/repo"
	"userv4/modules/user/entity"
)

type PriceTableUsecase struct {
	_crud.BaseUsecase[entity.PriceTableDomain, repo.IPriceTableRepo]
	priceTableRepo repo.IPriceTableRepo
}

type IPriceTableUsecase interface {
	_crud.IBaseUsecase[entity.PriceTableDomain]
}

func NewPriceTableUsecase(priceTableRepo repo.IPriceTableRepo) IPriceTableUsecase {
	return &PriceTableUsecase{
		BaseUsecase: _crud.BaseUsecase[entity.PriceTableDomain, repo.IPriceTableRepo]{
			Repo: priceTableRepo,
		},
		priceTableRepo: priceTableRepo,
	}
}
