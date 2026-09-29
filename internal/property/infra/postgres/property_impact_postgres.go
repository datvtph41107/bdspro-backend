package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	property_repo "bdspro/internal/repo/property"
	"context"
	"time"

	"gorm.io/gorm"
)

type PropertyImpactRepoImpl struct {
	db *gorm.DB
}

func NewSnapshotRepository(db *gorm.DB) property_repo.PropertyImpactRepository {
	return &PropertyImpactRepoImpl{db: db}
}

func (r *PropertyImpactRepoImpl) GetCurrent(ctx context.Context, entityType string, entityID uint64) (*domain.PropertyImpactVersion, error) {
	var snapshot domain.PropertyImpactVersion
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND valid_to IS NULL", entityType, entityID).
		Order("version DESC").
		First(&snapshot).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &snapshot, err
}

func (r *PropertyImpactRepoImpl) GetByProperty(ctx context.Context, propertyID uint64, entityTypes []string) ([]*domain.PropertyImpactVersion, error) {
	var snapshots []*domain.PropertyImpactVersion
	query := r.db.WithContext(ctx).Where("property_id = ? AND valid_to IS NULL", propertyID)
	if len(entityTypes) > 0 {
		query = query.Where("entity_type IN ?", entityTypes)
	}
	err := query.Order("created_at DESC").Find(&snapshots).Error
	return snapshots, err
}

func (r *PropertyImpactRepoImpl) Create(ctx context.Context, snapshot *domain.PropertyImpactVersion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Model(&domain.PropertyImpactVersion{}).
			Where("entity_type = ? AND entity_id = ? AND valid_to IS NULL", snapshot.EntityType, snapshot.EntityID).
			Update("valid_to", now).Error; err != nil {
			return err
		}
		return tx.Create(snapshot).Error
	})
}

func (r *PropertyImpactRepoImpl) Invalidate(ctx context.Context, snapshotID uint64, validTo time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.PropertyImpactVersion{}).
		Where("id = ?", snapshotID).
		Update("valid_to", validTo).Error
}

func (r *PropertyImpactRepoImpl) GetAffectedEntities(ctx context.Context, propertyID uint64) (*dto.AffectedEntities, error) {
	result := &dto.AffectedEntities{}

	var entities []struct {
		EntityType string
		EntityID   uint64
	}

	err := r.db.WithContext(ctx).
		Model(&domain.PropertyImpactVersion{}).
		Select("DISTINCT entity_type, entity_id").
		Where("property_id = ? AND valid_to IS NULL", propertyID).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	for _, e := range entities {
		switch e.EntityType {
		case "product":
			result.Products = append(result.Products, e.EntityID)
		case "listing":
			result.Listings = append(result.Listings, e.EntityID)
		case "asset":
			result.Assets = append(result.Assets, e.EntityID)
		case "deal":
			result.Deals = append(result.Deals, e.EntityID)
		case "crm_note":
			result.CrmNotes = append(result.CrmNotes, e.EntityID)
		}
	}

	return result, nil
}
