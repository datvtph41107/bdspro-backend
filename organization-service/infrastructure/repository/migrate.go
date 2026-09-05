package repository

import (
	"fmt"
	"organization/infrastructure/repository/postgres"
	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

// AutoMigrateModels là registry bootstrap GORM duy nhất của Organization.
// SQL production vẫn do organization-service/migrate sở hữu; registry này
// chỉ chạy khi process owner chọn schema mode automigrate.
func AutoMigrateModels() []any {
	return []any{
		&entity.Color{},
		&entity.BusinessDomain{},
		&entity.Organization{},
		&postgres.OrganizationBusinessDomainModel{},
		&postgres.OrganizationPermissionModel{},
		&postgres.OrganizationRoleModel{},
		&postgres.OrganizationMemberModel{},
		&postgres.OrganizationBranchModel{},
		&postgres.OrganizationBranchMemberModel{},
		&postgres.OrganizationLogActivityModel{},
		&postgres.GroupModel{},
		&postgres.GroupSettingModel{},
		&postgres.GroupMemberModel{},
		&postgres.GroupLogActivityModel{},
		&postgres.GroupNotificationModel{},
		&postgres.GroupDocumentModel{},
		&postgres.GroupChatModel{},
		&entity.BankAccount{},
		&entity.Deal{},
		&entity.DealProduct{},
		&entity.DealMember{},
		&entity.DealMilestone{},
		&entity.Investment{},
		&entity.AttachDocument{},
		&entity.InternalNote{},
		&entity.DealOfOrganization{},
		&entity.DealOfBranch{},
		&entity.DealOfGroup{},
	}
}

func AutoMigrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(hashtext(?))",
			"qhpro:organization:gorm-schema",
		).Error; err != nil {
			return fmt.Errorf("acquire Organization schema lock: %w", err)
		}
		if err := tx.AutoMigrate(AutoMigrateModels()...); err != nil {
			return fmt.Errorf("migrate Organization models: %w", err)
		}
		return nil
	})
}
