// Package db owns Notification's optional GORM bootstrap registry.
package db

import (
	"fmt"
	eventstore "notification/infra/postgres/eventing"
	"notification/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(hashtext(?))",
			"qhpro:notification:gorm-schema",
		).Error; err != nil {
			return fmt.Errorf("acquire Notification schema lock: %w", err)
		}
		if err := tx.AutoMigrate(
			&domain.NotificationEntity{},
			&domain.NotificationHistoryEntity{},
			&domain.HistoryEntity{},
			&domain.AdminHistoryEntity{},
			&domain.DealHistoryEntity{},
			&domain.HistoryAuthEntity{},
			&domain.HistoryFcmEntity{},
			&domain.PersonConfigEntity{},
			&domain.PropertyHistory{},
			&domain.AccountWarningTemplateEntity{},
			&domain.AccountWarningEntity{},
			&domain.AccountWarningLogEntity{},
		); err != nil {
			return fmt.Errorf("migrate Notification domain models: %w", err)
		}
		if err := eventstore.AutoMigrate(tx); err != nil {
			return fmt.Errorf("migrate Notification eventing models: %w", err)
		}
		return nil
	})
}
