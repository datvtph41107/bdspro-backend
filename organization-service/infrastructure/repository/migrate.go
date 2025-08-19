package repository

import (
	_db "common/db"

	"organization/internal/domain/entity"
)

func AutoMigrate() {
	// db := database.NewPostgresDB(config.DB_DNS)
	db := _db.DB
	// Migrate các models hiện có
	err := db.AutoMigrate(
		// &postgres.GroupChatModel{},
		// &postgres.OrganizationLogActivityModel{},
		// &postgres.GroupDocumentModel{},
		// &postgres.GroupNotificationModel{},
		// &postgres.GroupLogActivityModel{},
		// &postgres.GroupMemberModel{},
		// &postgres.GroupSettingModel{},
		// &postgres.GroupModel{},
		// &postgres.OrganizationBranchMemberModel{},
		// &postgres.OrganizationBranchModel{},
		// &postgres.OrganizationPermissionModel{},
		// &postgres.OrganizationRoleModel{},
		// &postgres.OrganizationMemberModel{},
		// &postgres.BusinessDomainModel{},
		// &postgres.OrganizationBusinessDomainModel{},
		// &postgres.OrganizationModel{}, &entity.Deal{},
		// &entity.DealProduct{},
		// &entity.AttachDocument{},
		// &entity.BankAccount{},
		// &entity.DealInvitation{},
		// &entity.Investment{},
		// &entity.DealMember{},
		// &entity.DealTransaction{},
		// &postgres.DealMemberModel{},
		// &postgres.InternalNoteModel{},
		// &entity.Color{},
		// &entity.DealCost{},
		// &entity.CostType{},
		// &entity.DealAction{},
		&entity.Deal{},
	)
	if err != nil {
		panic(err)
	}
}
