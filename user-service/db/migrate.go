// Package db owns User Service's GORM bootstrap registry.
//
// Production schema evolution remains owned by user-service/migrate. The
// registry exists for the explicit QHPRO_USER_DB_SCHEMA_MODE=automigrate mode
// and contains only current, reachable GORM-owned models.
package db

import (
	checkoutpostgres "user/infra/postgres/checkout"
	subscriptionpostgres "user/infra/postgres/subscription"
	"user/internal/domain"
	access "user/internal/domain/access"
	auth "user/internal/domain/auth"
	"user/internal/domain/user"
	"user/internal/models"

	"gorm.io/gorm"
)

// AutoMigrateModels là registry duy nhất giữa bootstrap GORM và contract test
// SQL. Khi thêm model persistence mới phải bổ sung tại đây và trong migrate/.
func AutoMigrateModels() []any {
	registry := []any{
		// Auth identity and session state (owned by User, not auth-service).
		&auth.AuthMethod{},
		&auth.AuthStatus{},
		&auth.UserOTPEntity{},
		&auth.UserSessionEntity{},
		&auth.DeviceEntity{},
		&auth.AuthConfig{},
		&auth.UserPINEntity{},
		&domain.UserActivityLog{},

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
		&user.UserProfileEntity{},
		&user.ProfileDeleted{},
		&models.AdminProfile{},
		&models.BlockEntity{},
		&models.ContactEntity{},
		&models.FollowEntity{},
		&models.FriendEntity{},
		&models.GroupEntity{},
		&user.BookmarkUserEntity{},
		&models.MainAreaEntity{},
		&models.MainAreaProfileEntity{},
		&user.PurposeUseEntity{},
		&user.PurposeUseProfileEntity{},
		&user.TagEntity{},
		&user.TagUserEntity{},
		&user.KYCEntity{},
		&user.ProfileMediaEntity{},
		&user.ProfessionEntity{},
		&user.CertificationEntity{},
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
	if err := database.AutoMigrate(AutoMigrateModels()...); err != nil {
		return err
	}
	return ensureAuthStatusFromLegacy(database)
}

// ensureAuthStatusFromLegacy copies registration lock rows from user_status when
// that legacy table is still present (pre-000021 DBs under automigrate bootstrap).
func ensureAuthStatusFromLegacy(database *gorm.DB) error {
	if database == nil {
		return nil
	}
	var exists bool
	if err := database.Raw(`
SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = current_schema()
      AND table_name = 'user_status'
)`).Scan(&exists).Error; err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return database.Exec(`
INSERT INTO auth_status (auth_id, active, verified, locked_until, updated_at)
SELECT
    us.auth_id,
    us.active,
    us.verified,
    us.locked_until,
    COALESCE(us.updated_at, NOW())
FROM user_status us
ON CONFLICT (auth_id) DO NOTHING
`).Error
}
