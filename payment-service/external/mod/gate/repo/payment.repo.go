package repo

import (
	cmodels "common/models"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	cmodels.BaseRepository
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}
