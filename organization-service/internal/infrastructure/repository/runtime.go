package repository

import (
	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/internal/infrastructure/database"
	"github.com/tuandq2112/organization-service/internal/infrastructure/repository/postgres"
)

func NewOrganizationRepository() repository.OrganizationRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationPostgresRepository(db)
	}
	return nil
}

func NewOrganizationMemberRepository() repository.OrganizationMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationMemberPostgresRepository(db)
	}
	return nil
}

func NewOrganizationRoleRepository() repository.OrganizationRoleRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationRolePostgresRepository(db)
	}
	return nil
}

func NewOrganizationPermissionRepository() repository.OrganizationPermissionRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationPermissionPostgresRepository(db)
	}
	return nil
}

func NewOrganizationBranchRepository() repository.OrganizationBranchRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationBranchPostgresRepository(db)
	}
	return nil
}

func NewOrganizationBranchMemberRepository() repository.OrganizationBranchMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationBranchMemberPostgresRepository(db)
	}
	return nil
}

func NewGroupRepository() repository.GroupRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupPostgresRepository(db)
	}
	return nil
}

func NewGroupSettingRepository() repository.GroupSettingRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupSettingPostgresRepository(db)
	}
	return nil
}

func NewGroupMemberRepository() repository.GroupMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupMemberPostgresRepository(db)
	}
	return nil
}

func NewGroupLogActivityRepository() repository.GroupLogActivityRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupLogActivityPostgresRepository(db)
	}
	return nil
}

func NewGroupNotificationRepository() repository.GroupNotificationRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupNotificationPostgresRepository(db)
	}
	return nil
}

func NewGroupDealRepository() repository.GroupDealRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupDealPostgresRepository(db)
	}
	return nil
}

func NewGroupDocumentRepository() repository.GroupDocumentRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewGroupDocumentPostgresRepository(db)
	}
	return nil
}

func NewOrganizationLogActivityRepository() repository.OrganizationLogActivityRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)

		return postgres.NewOrganizationLogActivityPostgresRepository(db)
	}
	return nil
}

func NewGroupChatRepository() repository.GroupChatRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db := database.NewPostgresDB(appconfig.DB_DNS)
		return postgres.NewGroupChatPostgresRepository(db)
	}
	return nil
}
