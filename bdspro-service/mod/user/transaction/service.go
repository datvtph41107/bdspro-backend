package u_transaction

import (
	"bdspro/models"
	"common/case/crud"
)

type TransactionService struct {
	crud.BaseService[models.Ward]
}

func NewTransactionService(repo *TransactionRepo) *TransactionService {
	return &TransactionService{
		BaseService: crud.BaseService[models.Ward]{Repo: repo.BaseRepo},
	}
}
