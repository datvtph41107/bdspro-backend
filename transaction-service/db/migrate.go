package db

import (
	_db "common/db"
	"tx/internal/domain"
	"tx/internal/domain/entity"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&domain.Tx{},
		&domain.TxAction{},
		&entity.ContractDeal{},
		&entity.DealCost{},
		&entity.CostType{},
	)
}
