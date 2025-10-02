package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"tqd/internal/interface/repo"

	"tqd/internal/domain"
)

type AmenityRepo struct {
	_provider.CrudRepo[domain.Amenity]
}

// NewAmenityRepo creates a new amenity repository
func NewAmenityRepo(db *_db.TransactionRepo) repo.IAmenityRepo {
	repo := &AmenityRepo{}
	repo.Init(repo, db)
	return repo
}
