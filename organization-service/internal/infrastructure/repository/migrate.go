package repository

import (
	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/infrastructure/database"
	"github.com/tuandq2112/organization-service/internal/infrastructure/repository/postgres"
)

func AutoMigrate() {
	db := database.NewPostgresDB(appconfig.DB_DNS)
	db.AutoMigrate(&postgres.GroupChatModel{})
	db.AutoMigrate(&postgres.OrganizationLogActivityModel{})
	db.AutoMigrate(&postgres.GroupDocumentModel{})
	db.AutoMigrate(&postgres.GroupDealModel{})
	db.AutoMigrate(&postgres.GroupNotificationModel{})
	db.AutoMigrate(&postgres.GroupLogActivityModel{})
	db.AutoMigrate(&postgres.GroupMemberModel{})
	db.AutoMigrate(&postgres.GroupSettingModel{})
	db.AutoMigrate(&postgres.GroupModel{})
	db.AutoMigrate(&postgres.OrganizationBranchMemberModel{})
	db.AutoMigrate(&postgres.OrganizationBranchModel{})
	db.AutoMigrate(&postgres.OrganizationPermissionModel{})
	db.AutoMigrate(&postgres.OrganizationRoleModel{})
	db.AutoMigrate(&postgres.OrganizationMemberModel{})
	db.AutoMigrate(&postgres.OrganizationModel{})
}
