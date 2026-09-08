package postgres

import (
	"fmt"
	paymentpg "payment/infra/postgres/payment"
	walletpg "payment/infra/postgres/wallet"
	bankdomain "payment/internal/domain/bank"

	"gorm.io/gorm"
)

// AutoMigrateModels là source-of-truth giữa GORM bootstrap và SQL contract.
func AutoMigrateModels() []any {
	models := []any{&bankdomain.Bank{}}
	models = append(models, walletpg.AutoMigrateModels()...)
	models = append(models, paymentpg.AutoMigrateModels()...)
	return models
}

// AutoMigrate is Payment's complete GORM bootstrap registry. Canonical
// production evolution remains in payment-service/migrate.
func AutoMigrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(hashtext(?))",
			"qhpro:payment:gorm-schema",
		).Error; err != nil {
			return fmt.Errorf("acquire Payment schema lock: %w", err)
		}
		if err := tx.AutoMigrate(AutoMigrateModels()...); err != nil {
			return fmt.Errorf("migrate Payment models: %w", err)
		}
		return nil
	})
}
