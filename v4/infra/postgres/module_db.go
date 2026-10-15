package postgres

import "gorm.io/gorm"

// CrmDB, PaymentDB, and TqdDB are wire disambiguators over the process-wide
// PostgreSQL pool (common/db.NewDB → database.dsn).

type CrmDB struct {
	DB *gorm.DB
}

func NewCrmDB(db *gorm.DB) *CrmDB {
	return &CrmDB{DB: db}
}

type PaymentDB struct {
	DB *gorm.DB
}

func NewPaymentDB(db *gorm.DB) *PaymentDB {
	return &PaymentDB{DB: db}
}

type TqdDB struct {
	DB *gorm.DB
}

func NewTqdDB(db *gorm.DB) (*TqdDB, error) {
	if err := ensureTqdLocationCatalogSchema(db); err != nil {
		return nil, err
	}
	return &TqdDB{DB: db}, nil
}
