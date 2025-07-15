package repository

import (
	_db "common/db"
	"organization/infrastructure/repository/postgres"
	"organization/internal/domain/entity"
)

func AutoMigrate() {
	// db := database.NewPostgresDB(config.DB_DNS)
	db := _db.DB
	db.AutoMigrate(&postgres.GroupChatModel{})
	db.AutoMigrate(&postgres.OrganizationLogActivityModel{})
	db.AutoMigrate(&postgres.GroupDocumentModel{})
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
	db.AutoMigrate(&entity.Deal{})
	db.AutoMigrate(&entity.DealMember{})
	db.AutoMigrate(&entity.DealProduct{})
	db.AutoMigrate(&entity.AttachDocument{})
	db.AutoMigrate(&entity.BankAccount{})
	db.AutoMigrate(&entity.DealInvitation{}) // Thêm migration cho deal_invitations
	err := db.AutoMigrate(&entity.Investment{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.DealMember{})
	db.AutoMigrate(&entity.DealTransaction{})
	
	// Thêm migration cho DealMemberModel với các trường hoa hồng
	db.AutoMigrate(&postgres.DealMemberModel{})
}
