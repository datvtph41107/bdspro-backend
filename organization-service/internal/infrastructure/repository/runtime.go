package repository

import (
	"log"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/internal/infra/database"
	"github.com/tuandq2112/organization-service/internal/infra/repository/postgres"
)

func NewOrganizationRepository() repository.OrganizationRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationModel{}, &postgres.OrganizationMemberModel{})
		return postgres.NewOrganizationPostgresRepository(db)
	}
	return nil
}

func NewOrganizationMemberRepository() repository.OrganizationMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationMemberModel{})
		return postgres.NewOrganizationMemberPostgresRepository(db)
	}
	return nil
}

func NewOrganizationRoleRepository() repository.OrganizationRoleRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationRoleModel{})
		return postgres.NewOrganizationRolePostgresRepository(db)
	}
	return nil
}

func NewOrganizationPermissionRepository() repository.OrganizationPermissionRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationPermissionModel{})
		return postgres.NewOrganizationPermissionPostgresRepository(db)
	}
	return nil
}

func NewOrganizationBranchRepository() repository.OrganizationBranchRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationBranchModel{})
		return postgres.NewOrganizationBranchPostgresRepository(db)
	}
	return nil
}

func NewOrganizationBranchMemberRepository() repository.OrganizationBranchMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.OrganizationBranchMemberModel{})
		return postgres.NewOrganizationBranchMemberPostgresRepository(db)
	}
	return nil
}

func NewGroupRepository() repository.GroupRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.GroupModel{})
		return postgres.NewGroupPostgresRepository(db)
	}
	return nil
}

func NewGroupSettingRepository() repository.GroupSettingRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.GroupSettingModel{})

		return postgres.NewGroupSettingPostgresRepository(db)
	}
	return nil
}

func NewGroupMemberRepository() repository.GroupMemberRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.GroupMemberModel{})
		return postgres.NewGroupMemberPostgresRepository(db)
	}
	return nil
}

func NewGroupLogActivityRepository() repository.GroupLogActivityRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.GroupLogActivityModel{})
		return postgres.NewGroupLogActivityPostgresRepository(db)
	}
	return nil
}

func NewGroupNotificationRepository() repository.GroupNotificationRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db.AutoMigrate(&postgres.GroupNotificationModel{})
		return postgres.NewGroupNotificationPostgresRepository(db)
	}
	return nil
}

func NewGroupDealRepository() repository.GroupDealRepository {
	switch appconfig.PRIMARY_DB_TYPE {
	case appconfig.DBTypePostgres:
		db, err := database.NewPostgresDB(appconfig.DB_DNS)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}	
		db.AutoMigrate(&postgres.GroupDealModel{})
		return postgres.NewGroupDealPostgresRepository(db)
	}
	return nil
}

