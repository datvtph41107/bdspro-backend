package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"userv4/internal/interface/repo"
	"userv4/internal/models"
)

type PriceTableRepo struct {
	_provider.CrudRepo[models.PriceTableDomain]
}

// @bind: user/internal/interface/repo.IPriceTableRepo
func NewPriceTableRepo(db *_db.TransactionRepo) repo.IPriceTableRepo {
	repo := &PriceTableRepo{}
	repo.Init(repo, db)
	return repo
}
