package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"common/case/crud3"
	"context"
	"strings"

	"gorm.io/gorm"
)

type PostgrePropertyRepo struct {
	crud3.BaseRepo[domain.Property]
}

func NewPostgrePropertyRepo(db *gorm.DB) repo.PropertyRepo {
	return &PostgrePropertyRepo{
		BaseRepo: crud3.BaseRepo[domain.Property]{DB: db},
	}
}

func (r *PostgrePropertyRepo) Search(ctx context.Context, payload *dto.PropertySearchDTO) ([]dto.PropertyWithCounts, int64, error) {
	var properties []dto.PropertyWithCounts
	var total int64

	baseQuery := r.DB.WithContext(ctx).
		Table("property").
		Select(`property.*,
			CASE WHEN pr.primary_product_id IS NOT NULL THEN 1 ELSE 0 END as product_count,
			CASE WHEN pr.asset_id IS NOT NULL THEN 1 ELSE 0 END as asset_count,
			COALESCE((SELECT COUNT(*) FROM posts WHERE posts.product_id = pr.primary_product_id AND posts.deleted_at IS NULL), 0) as post_count,
			COALESCE(prov.name, '') as province_name,
			COALESCE(ward.name, '') as ward_name,
			COALESCE(proj.name, '') as project_name,
			COALESCE(pt.name, '') as property_type_name,
			pbi.building_type,
			pbi.construction_area,
			pbi.floor_area,
			pbi.floors,
			pbi.bedrooms,
			pbi.bathrooms,
			pbi.direction,
			pbi.balcony_direction,
			avatar.url as avatar_url,
			avatar.thumb_url as avatar_thumb_url,
			avatar.media_type as avatar_media_type,
			avatar.sort_order as avatar_sort_order`).
		Joins("LEFT JOIN property_relation pr ON property.id = pr.property_id").
		Joins("LEFT JOIN province_v2 prov ON property.province_id = prov.id AND prov.deleted_at IS NULL").
		Joins("LEFT JOIN ward_v2 ward ON property.ward_id = ward.id AND ward.deleted_at IS NULL").
		Joins("LEFT JOIN projects proj ON property.project_id = proj.id AND proj.deleted_at IS NULL").
		Joins("LEFT JOIN property_type pt ON property.property_type_id = pt.id AND pt.deleted_at IS NULL").
		Joins("LEFT JOIN property_building_info pbi ON property.id = pbi.property_id AND pbi.deleted_at IS NULL").
		Joins("LEFT JOIN property_media avatar ON property.id = avatar.property_id AND avatar.deleted_at IS NULL").
		Where("property.deleted_at IS NULL")

	// Filter by text
	if payload.Text != "" {
		text := "%" + strings.ToLower(payload.Text) + "%"
		baseQuery = baseQuery.Where("LOWER(property.title) LIKE ? OR LOWER(property.address_text) LIKE ?", text, text)
	}

	// Filter by propertyTypeId
	if payload.PropertyTypeID != nil {
		baseQuery = baseQuery.Where("property.property_type_id = ?", *payload.PropertyTypeID)
	}

	// Filter by locationId
	if payload.LocationID != nil {
		baseQuery = baseQuery.Where("property.location_id = ?", *payload.LocationID)
	}

	// Filter by projectId
	if payload.ProjectID != nil {
		baseQuery = baseQuery.Where("property.project_id = ?", *payload.ProjectID)
	}

	// Filter by recordStatus
	if payload.RecordStatus != "" {
		baseQuery = baseQuery.Where("property.record_status = ?", payload.RecordStatus)
	}

	// Count total (clone query để count)
	countQuery := baseQuery
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := payload.GetOffset()
	limit := payload.GetLimit()

	// Default sort by updated_at DESC
	sort := payload.Sort
	if sort == "" {
		sort = "property.updated_at DESC"
	}

	err := baseQuery.
		Order(sort).
		Offset(offset).
		Limit(limit).
		Scan(&properties).Error

	return properties, total, err
}

func (r *PostgrePropertyRepo) GetDetail(ctx context.Context, id uint64) (*domain.Property, error) {
	var property domain.Property
	if err := r.DB.WithContext(ctx).
		Preload("PropertyType").
		Preload("Project").
		Preload("Province").
		Preload("MediaList").
		Preload("Ward").
		Preload("BuildingInfo").
		Preload("LandInfo").
		Preload("Avatar").
		Preload("Amenities").
		Where("id = ? and deleted_at is null", id).First(&property).Error; err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *PostgrePropertyRepo) Update(c context.Context, id uint64, entity *domain.Property) error {
	// Sử dụng map để update các trường, bao gồm cả khi giá trị là nil (để set về NULL trong DB)
	updates := map[string]interface{}{
		"property_type_id": entity.PropertyTypeID,
		"title":            entity.Title,
		"location_id":      entity.LocationID,
		"address_detail":   entity.AddressDetail,
		"latitude":         entity.Latitude,
		"longitude":        entity.Longitude,
		"map_url":          entity.MapURL,
		"project_id":       entity.ProjectID,
		"block_id":         entity.BlockID,
		"level":            entity.Level,
		"unit_code":        entity.UnitCode,
		"area_total":       entity.AreaTotal,
		"area_land":        entity.AreaLand,
		"area_residential": entity.AreaResidential,
		"legal_note":       entity.LegalNote,
		"building_info_id": entity.BuildingInfoID, // Có thể là nil để set về NULL
		"legal_status":     entity.LegalStatus,
		"record_status":    entity.RecordStatus,
		"avatar_id":        entity.AvatarID,
		"province_id":      entity.ProvinceID,
		"ward_id":          entity.WardID,
	}

	return GetDB(c, r.DB).
		Model(&domain.Property{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *PostgrePropertyRepo) Create(c context.Context, entity *domain.Property) error {
	return GetDB(c, r.DB).Create(entity).Error
}

func (r *PostgrePropertyRepo) UpdateOption(c context.Context, id uint64, property *domain.Property) error {
	// Chỉ update các field cụ thể: area_land, area_total, project_id, property_type_id
	updates := map[string]interface{}{}

	if property.AreaLand != nil {
		updates["area_land"] = property.AreaLand
	}
	if property.AreaTotal != nil {
		updates["area_total"] = property.AreaTotal
	}
	if property.ProjectID != nil {
		updates["project_id"] = property.ProjectID
	}
	if property.PropertyTypeID != nil {
		updates["property_type_id"] = property.PropertyTypeID
	}

	// Nếu không có field nào cần update, return nil
	if len(updates) == 0 {
		return nil
	}

	return GetDB(c, r.DB).
		Model(&domain.Property{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *PostgrePropertyRepo) GetOneByID(c context.Context, id uint64) (*domain.Property, error) {
	var property domain.Property
	if err := GetDB(c, r.DB).
		Model(&domain.Property{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&property).Error; err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *PostgrePropertyRepo) DeleteBatchProperties(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return GetDB(ctx, r.DB).
		Where("id IN (?) AND deleted_at IS NULL", ids).
		Delete(&domain.Property{}).Error
}
