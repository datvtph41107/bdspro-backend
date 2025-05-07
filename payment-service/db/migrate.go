package migrate

import (
	_db "common/db"
	"payment/internal/wallet/domain"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		&domain.Account{},
		&domain.Transaction{},
	)
}
