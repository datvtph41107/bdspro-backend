package postgres

import (
	_db "common/db"

	"marketing/internal/domain/model"
)

func Migrate() {
	db := _db.DB
	db.AutoMigrate(&model.AdvertisingSettings{},
		&model.BudgetConfig{},
		&model.Campaign{},
		&model.Package{},
		&model.BudgetUsage{},
		&model.BudgetAlert{},
	)
}
