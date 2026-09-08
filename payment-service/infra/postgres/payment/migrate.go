package postgres

import "gorm.io/gorm"

// AutoMigrateModels công bố registry persistence private cho schema owner của
// Payment và contract test; package khác không được tự gọi AutoMigrate.
func AutoMigrateModels() []any {
	return []any{
		&orderModel{},
		&paymentAttemptModel{},
		&settlementModel{},
		&fulfillmentModel{},
		&commandEffectModel{},
		&outboxEventModel{},
	}
}

// AutoMigrate keeps private persistence rows private while exposing one schema
// bootstrap operation to Payment's process-owned registry.
func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(AutoMigrateModels()...)
}
