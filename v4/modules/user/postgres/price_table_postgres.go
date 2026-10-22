package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"mdn4/internal/interface/repo"
	"mdn4/modules/user/entity"
)

type PriceTableRepo struct {
	_provider.CrudRepo[entity.PriceTableDomain]
}

func NewPriceTableRepo(db *_db.TransactionRepo) repo.IPriceTableRepo {
	repo := &PriceTableRepo{}
	repo.Init(repo, db)
	return repo
}
