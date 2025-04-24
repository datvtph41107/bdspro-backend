package u_transaction

import (
	"bdspro/v1/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	crud.BaseRepo[models.Ward]
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		BaseRepo: crud.BaseRepo[models.Ward]{DB: db},
	}
}
