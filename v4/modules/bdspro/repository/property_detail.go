package repository

import (
	"context"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"

	"gorm.io/gorm"
)

func (r *propertyPostgres) GetDetail(ctx context.Context, id, ownerOriginID uint64) (*entity.PropertyLineageDetail, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var p entity.PropertyLineage
	err := r.db.WithContext(ctx).
		Preload("Info", "deleted_at IS NULL").
		Preload("Info.PropertyType", "deleted_at IS NULL").
		Preload("Info.Project").
		Preload("Info.Avatar", "deleted_at IS NULL").
		Preload("Location", "deleted_at IS NULL").
		Preload("Location.Province").
		Preload("Location.Ward").
		Preload("LandInfo", "deleted_at IS NULL").
		Preload("BuildingInfo", "deleted_at IS NULL").
		Preload("Statistic", "deleted_at IS NULL").
		Preload("Identify", "deleted_at IS NULL").
		Preload("Edvidence", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	out := &entity.PropertyLineageDetail{PropertyLineage: &p, VerifiedNationalAt: p.VerifiedNationalAt}

	if ownerOriginID > 0 {
		var pu entity.PropertyUser
		if err := r.db.WithContext(ctx).
			Where("property_lineage_id = ? AND owner_origin_id = ? AND deleted_at IS NULL", id, ownerOriginID).
			First(&pu).Error; err == nil {
			out.Personalization = &pu
		}
	}

	if p.Info != nil && p.Info.ProjectID != nil && *p.Info.ProjectID > 0 {
		var dev entity.PropertyDeveloper
		_ = r.db.WithContext(ctx).Raw(`
			SELECT d.id, d.name, d.slug, d.logo_url, d.phone, d.email, d.website
			FROM projects pj
			JOIN developer d ON d.id = pj.developer_id AND d.deleted_at IS NULL
			WHERE pj.id = ? AND pj.deleted_at IS NULL
			LIMIT 1`, *p.Info.ProjectID).Scan(&dev).Error
		if dev.ID > 0 {
			out.Developer = &dev
		}
	}

	_ = r.db.WithContext(ctx).Raw(`
		SELECT a.id, a.name
		FROM property_amenity pa
		JOIN amenity a ON pa.amenity_item_id = a.id AND a.deleted_at IS NULL
		WHERE pa.property_lineage_id = ? AND a.active = true`, id).Scan(&out.Amenities).Error

	_ = r.db.WithContext(ctx).Raw(`
		SELECT ar.id, ar.name
		FROM property_area_region par
		JOIN area_regions ar ON par.area_region_id = ar.id AND ar.deleted_at IS NULL
		WHERE par.property_lineage_id = ? AND ar.active = true`, id).Scan(&out.AreaRegions).Error

	if out.PropertyLineage.Statistic != nil {
		stat := out.PropertyLineage.Statistic
		var assetCount, productCount int64
		r.db.WithContext(ctx).Raw(`
			SELECT COUNT(*) FROM property_relation
			WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL`,
			id, enums.ERelationTypeAsset).Scan(&assetCount)
		r.db.WithContext(ctx).Raw(`
			SELECT COUNT(*) FROM property_relation
			WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL`,
			id, enums.ERelationTypeProduct).Scan(&productCount)
		stat.AssetCount = assetCount
		stat.ProductCount = productCount
	}

	return out, nil
}
