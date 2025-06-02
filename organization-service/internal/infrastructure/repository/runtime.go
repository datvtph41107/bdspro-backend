package repository

import (
	"log"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/internal/infrastructure/database"
	"github.com/tuandq2112/organization-service/internal/infrastructure/repository/postgres"
)

func NewRepository() repository.OrganizationRepository {
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
