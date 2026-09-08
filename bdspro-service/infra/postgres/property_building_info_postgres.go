package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyBuildingInfoRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyBuildingInfoRepo(db *gorm.DB) repo.PropertyBuildingInfoRepo {
	return &PostgrePropertyBuildingInfoRepo{
		DB: db,
	}
}

func (r *PostgrePropertyBuildingInfoRepo) Create(ctx context.Context, buildingInfo *domain.PropertyBuildingInfo) error {
	return GetDB(ctx, r.DB).Create(buildingInfo).Error
}

func (r *PostgrePropertyBuildingInfoRepo) Update(ctx context.Context, buildingInfo *domain.PropertyBuildingInfo) error {
	err := GetDB(ctx, r.DB).
		Model(&domain.PropertyBuildingInfo{}).
		Where("property_identify_id = ? AND deleted_at IS NULL", buildingInfo.PropertyIdentifyID).
		Updates(map[string]interface{}{
			"building_type": buildingInfo.BuildingType,
			// "construction_area": buildingInfo.ConstructionArea,
			// "floor_area":        buildingInfo.FloorArea,
			"floors":            buildingInfo.Floors,
			"bedrooms":          buildingInfo.Bedrooms,
			"bathrooms":         buildingInfo.Bathrooms,
			"direction":         buildingInfo.Direction,
			"balcony_direction": buildingInfo.BalconyDirection,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgrePropertyBuildingInfoRepo) GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyBuildingInfo, error) {
	var buildingInfo domain.PropertyBuildingInfo
	err := GetDB(ctx, r.DB).
		Where("property_id = ? AND deleted_at IS NULL", propertyID).
		First(&buildingInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &buildingInfo, nil
}
