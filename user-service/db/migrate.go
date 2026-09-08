// Package db owns User Service's GORM bootstrap registry.
//
// Production schema evolution remains owned by user-service/migrate. The
// registry exists for the explicit QHPRO_USER_DB_SCHEMA_MODE=automigrate mode
// and contains only current, reachable GORM-owned models.
package db

import (
	checkoutpostgres "user/infra/postgres/checkout"
	subscriptionpostgres "user/infra/postgres/subscription"
	access "user/internal/domain/access"
	auth "user/internal/domain/auth"
	"user/internal/models"

	"gorm.io/gorm"
)

// AutoMigrateModels là registry duy nhất giữa bootstrap GORM và contract test
// SQL. Khi thêm model persistence mới phải bổ sung tại đây và trong migrate/.
func AutoMigrateModels() []any {
	registry := []any{
		// Auth identity and session state (owned by User, not auth-service).
		&auth.AuthMethod{},
		&auth.AuthUser{},
		&auth.UserOTPEntity{},
		&auth.UserSessionEntity{},
		&auth.DeviceEntity{},
		&auth.UserStatusEntity{},
		&auth.AuthConfig{},
		&auth.UserPINEntity{},

		// Access-control state. Role and RoleGroup materialize the canonical
		// role_permissions and group_permissions join tables.
		&access.Permission{},
		&access.Color{},
		&access.RoleGroup{},
		&access.Role{},
		&access.RoleProfile{},
		&access.AdminAccessDomain{},
		&access.AdminAccessLogDomain{},

		// Profile and social state.
		&models.UserProfileEntity{},
		&models.ProfileDeleted{},
		&models.AdminProfile{},
		&models.BlockEntity{},
		&models.ContactEntity{},
		&models.FollowEntity{},
		&models.FriendEntity{},
		&models.GroupEntity{},
		&models.BookmarkUserEntity{},
		&models.MainAreaEntity{},
		&models.MainAreaProfileEntity{},
		&models.PurposeUseEntity{},
		&models.PurposeUseProfileEntity{},
		&models.TagEntity{},
		&models.TagUserEntity{},
		&models.KYCEntity{},
		&models.ProfileMediaEntity{},
		&models.ProfessionEntity{},
		&models.CertificationEntity{},
		&models.Profile{},
		&models.PriceTableDomain{},

		// Canonical commercial catalog/subscription read and write models.
		&models.CatalogProduct{},
		&models.CatalogPlan{},
		&models.CatalogPlanVersion{},
		&models.CatalogPlanEntitlement{},
		&models.CatalogPlanOperationPolicy{},
		&models.CatalogPriceItem{},
		&models.CatalogSubscription{},
		&models.CatalogSubscriptionEvent{},
	}
	registry = append(registry, checkoutpostgres.AutoMigrateModels()...)
	registry = append(registry, subscriptionpostgres.AutoMigrateModels()...)
	return registry
}

func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(AutoMigrateModels()...)
}
