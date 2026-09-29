package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"user/internal/interface/repo"
	"user/internal/models"
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
