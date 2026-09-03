// Package db owns TQD's optional GORM bootstrap registry.
//
// Ordered production evolution is owned exclusively by tqd-service/migrate.
// This registry is selected only through QHPRO_TQD_DB_SCHEMA_MODE=automigrate
// and restores the development bootstrap capability of the legacy source.
package db

import (
	"fmt"
	"tqd/infra/postgres/generatedreport/job"
	"tqd/infra/postgres/usage"
	"tqd/internal/domain"
	qh "tqd/internal/domain/qh"

	"gorm.io/gorm"
)

func AutoMigrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		// Serialize bootstrap across concurrently starting local replicas. SQL
		// migration mode never enters this transaction.
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(hashtext(?))",
			"qhpro:tqd:gorm-schema",
		).Error; err != nil {
			return fmt.Errorf("acquire TQD schema lock: %w", err)
		}
		for _, extension := range []string{"postgis", "unaccent"} {
			if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS " + extension).Error; err != nil {
				return fmt.Errorf("ensure TQD extension %s: %w", extension, err)
			}
		}

		if err := tx.AutoMigrate(
			// Directory and POI state.
			&domain.PoiCategory{},
			&domain.POI{},
			&domain.POIMedia{},
			&domain.OpenHour{},
			&domain.Amenity{},
			&domain.ContactLabel{},
			&domain.DirectoryCategory{},
			&domain.DirectorySource{},
			&domain.DirectorySupplier{},
			&domain.Feature{},

			// Administrative and planning state.
			&domain.Province{},
			&domain.Ward{},
			&qh.QHLayerFamily{},
			&qh.QHAuthorityIssuring{},
			&qh.QHLayer{},
			&qh.QHJurisdiction{},
			&qh.QHLabel{},
			&qh.QHLandUse{},
			&qh.QHLayerLegend{},
			&qh.QHRegion{},
			&qh.QHRegionExtend{},
			&qh.QHLabelLayer{},
			&qh.QHLayerLandUse{},
			&qh.QHLayerLegal{},
			&qh.QHLayerResolverConfig{},
			&qh.QHShape{},
			&qh.Parcel{},
			&qh.QHPlanningUse{},
			&qh.QHParcelInfo{},
			&qh.QHPlanningProject{},
			&qh.QHPlanningDocument{},
			&qh.QHPlanningEvent{},
			&qh.QHPlanningRelation{},
			&qh.QHLegalDocument{},
			&qh.ProAIJob{},

			// Audit, user interaction, report, and subscription state.
			&qh.QHAuditEntry{},
			&qh.QHRegionImportErrorLog{},
			&qh.QHUserFollowedParcel{},
			&qh.QHUserFollowedPlanningProject{},
			&qh.QHUserReported{},
			&qh.QHUserViewHistory{},
			&qh.QHUserViewEvent{},
			&domain.UserNotification{},
			&domain.Report{},
			&domain.ReportEvent{},
			&domain.UserSubscription{},
		); err != nil {
			return fmt.Errorf("migrate TQD domain models: %w", err)
		}
		if err := postgresjob.AutoMigrate(tx); err != nil {
			return fmt.Errorf("migrate generated-report jobs: %w", err)
		}
		if err := postgresusage.AutoMigrate(tx); err != nil {
			return fmt.Errorf("migrate quota usage: %w", err)
		}
		return nil
	})
}
