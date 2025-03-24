package models

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&Province{},
		&District{},
		&Ward{},
		&Product{},
		&ProductAccess{},
		&DocType{},
		&PropertyType{},
	)
}
