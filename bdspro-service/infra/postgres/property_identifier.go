package postgres

import (
	"context"
	"errors"
	"fmt"

	"bdspro/internal/domain"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type PropertyIdentifierPostgres struct {
	db *gorm.DB
}

func NewIdentifierRepository(db *gorm.DB) repo.PropertyIdentifierRepo {
	return &PropertyIdentifierPostgres{
		db: db,
	}
}

func (r *PropertyIdentifierPostgres) Create(ctx context.Context, identifier *domain.CountryIdentifier) error {
	result := r.db.WithContext(ctx).Create(identifier)
	if result.Error != nil {
		return fmt.Errorf("db create failed: %w", result.Error)
	}
	return nil
}

func (r *PropertyIdentifierPostgres) GetByID(ctx context.Context, id string) (*domain.CountryIdentifier, error) {
	var identifier domain.CountryIdentifier
	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&identifier)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found (let usecase handle)
		}
		return nil, fmt.Errorf("db get by id failed: %w", result.Error)
	}
	return &identifier, nil
}

// Update - ONLY updates data
func (r *PropertyIdentifierPostgres) Update(ctx context.Context, identifier *domain.CountryIdentifier) error {
	result := r.db.WithContext(ctx).
		Model(&domain.CountryIdentifier{}).
		Where("id = ?", identifier.ID).
		Updates(map[string]interface{}{
			"land_parcel_code": identifier.LandParcelCode,
			"project_code":     identifier.ProjectCode,
			"province_id":      identifier.ProvinceID,
			"district_id":      identifier.DistrictID,
			"ward_id":          identifier.WardID,
			"latitude":         identifier.Latitude,
			"longitude":        identifier.Longitude,
			"type":             identifier.Type,
			"legal_status":     identifier.LegalStatus,
			"current_owner_id": identifier.CurrentOwnerID,
			"updated_at":       identifier.UpdatedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("db update failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete - ONLY soft deletes
func (r *PropertyIdentifierPostgres) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.CountryIdentifier{})

	if result.Error != nil {
		return fmt.Errorf("db delete failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Search - ONLY builds and executes search query
func (r *PropertyIdentifierPostgres) Search(ctx context.Context, filters map[string]interface{}, page, limit int32) ([]*domain.CountryIdentifier, int64, error) {
	var identifiers []*domain.CountryIdentifier
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.CountryIdentifier{}).Where("deleted_at IS NULL")

	// Apply filters (purely mechanical, no logic)
	for field, value := range filters {
		if value != nil && value != "" {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("db count failed: %w", err)
	}

	// Apply pagination
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(int(offset)).Limit(int(limit))
	}

	// Execute
	if err := query.Order("created_at DESC").Find(&identifiers).Error; err != nil {
		return nil, 0, fmt.Errorf("db search failed: %w", err)
	}

	return identifiers, total, nil
}

// GetByLocation - ONLY spatial query
func (r *PropertyIdentifierPostgres) GetByLocation(ctx context.Context, lat, lng float64, radiusKm float64) ([]*domain.CountryIdentifier, error) {
	var identifiers []*domain.CountryIdentifier

	// Pure spatial query - no business logic
	point := fmt.Sprintf("POINT(%f %f)", lng, lat)
	err := r.db.WithContext(ctx).
		Where("ST_DWithin(geo_location, ST_GeographyFromText(?), ?) AND deleted_at IS NULL",
			point, radiusKm*1000).
		Find(&identifiers).Error

	if err != nil {
		return nil, fmt.Errorf("db location query failed: %w", err)
	}
	return identifiers, nil
}

// GetByOwner - ONLY query by owner
func (r *PropertyIdentifierPostgres) GetByOwner(ctx context.Context, ownerID string) ([]*domain.CountryIdentifier, error) {
	var identifiers []*domain.CountryIdentifier

	err := r.db.WithContext(ctx).
		Where("current_owner_id = ? AND deleted_at IS NULL", ownerID).
		Order("created_at DESC").
		Find(&identifiers).Error

	if err != nil {
		return nil, fmt.Errorf("db owner query failed: %w", err)
	}
	return identifiers, nil
}

// CheckExists - ONLY existence check
func (r *PropertyIdentifierPostgres) CheckExists(ctx context.Context, landParcelCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.CountryIdentifier{}).
		Where("land_parcel_code = ? AND deleted_at IS NULL", landParcelCode).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("db check exists failed: %w", err)
	}
	return count > 0, nil
}
