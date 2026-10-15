package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"userv4/internal/interface/repo"
	"userv4/modules/user/entity"
)

type PriceTableRepo struct {
	_provider.CrudRepo[entity.PriceTableDomain]
}

// @bind: user/internal/interface/repo.IPriceTableRepo
func NewPriceTableRepo(db *_db.TransactionRepo) repo.IPriceTableRepo {
	repo := &PriceTableRepo{}
	repo.Init(repo, db)
	return repo
}
