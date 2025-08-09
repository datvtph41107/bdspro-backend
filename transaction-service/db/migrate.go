package db

import (
	_db "common/db"
	"tx/internal/domain"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&domain.Tx{},
		&domain.TxAction{},
	)
}
