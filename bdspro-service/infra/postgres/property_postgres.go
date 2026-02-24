package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
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

func (r *PostgrePropertyRepo) GetRefreshSyncIds(c context.Context, req *dto.CheckVersionSyncRequest) ([]uint64, error) {
	// profileId := _utils.GetProfileIdWithContext(c)
	// originId := _utils.GetOriginIdFromContext(c)
	var propertyIds []uint64
	err := r.DB.WithContext(c).
		Model(&domain.Property{}).
		// Joins(`LEFT JOIN product_user pu ON (pu.product_id = products.id AND pu.deleted_at IS NULL)`).
		Where(`property.deleted_at is null `).
		Order("property.id desc").
		Pluck("property.id", &propertyIds).Error

	if err != nil {
		return nil, err
	}

	return propertyIds, nil
}

func (r *PostgrePropertyRepo) Search(ctx context.Context, req *dto.PropertySearchDTO) ([]dto.PropertyWithCounts, int64, error) {
	var properties []dto.PropertyWithCounts
	var total int64

	db := r.DB.WithContext(ctx).Table("property").
		Select(`
			property.*,
			COALESCE(prov.name, '') as province_name,
            COALESCE(dist.name, '') as district_name,
            COALESCE(ward.name, '') as ward_name,
            COALESCE(proj.name, '') as project_name,
            COALESCE(blk.name, '') as block_name,
            COALESCE(bld.name, '') as building_name,
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
            avatar.sort_order as avatar_sort_order,
            (SELECT COUNT(*) FROM property_relation pr WHERE pr.property_id = property.id AND pr.relation_type = ? AND pr.deleted_at IS NULL) as asset_count,
            (SELECT COUNT(*) FROM property_relation pr WHERE pr.property_id = property.id AND pr.relation_type = ? AND pr.deleted_at IS NULL) as product_count,
            (SELECT COUNT(*) FROM property_relation pr 
                LEFT JOIN posts ON pr.relation_id = posts.product_id 
                WHERE pr.property_id = property.id AND pr.relation_type = ? AND posts.deleted_at IS NULL) as post_count
        `, enums.ERelationTypeAsset, enums.ERelationTypeProduct, enums.ERelationTypeProduct).
		Joins("LEFT JOIN province_v2 prov ON property.province_id = prov.id AND prov.deleted_at IS NULL").
		Joins("LEFT JOIN district_v2 dist ON property.district_id = dist.id AND dist.deleted_at IS NULL").
		Joins("LEFT JOIN ward_v2 ward ON property.ward_id = ward.id AND ward.deleted_at IS NULL").
		Joins("LEFT JOIN projects proj ON property.project_id = proj.id AND proj.deleted_at IS NULL").
		Joins("LEFT JOIN blocks blk ON property.block_id = blk.id AND blk.deleted_at IS NULL").
		Joins("LEFT JOIN buildings bld ON property.building_id = bld.id AND bld.deleted_at IS NULL").
		Joins("LEFT JOIN property_type pt ON property.property_type_id = pt.id AND pt.deleted_at IS NULL").
		Joins("LEFT JOIN property_building_info pbi ON property.id = pbi.property_id AND pbi.deleted_at IS NULL").
		Joins("LEFT JOIN property_media avatar ON property.avatar_id = avatar.id AND avatar.deleted_at IS NULL").
		Where("property.deleted_at IS NULL")

	db = r.applyFilters(db, req)

	var countDB *gorm.DB
	countDB = r.DB.WithContext(ctx).Table("property").
		Where("property.deleted_at IS NULL")
	countDB = r.applyFiltersForCount(countDB, req)
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	db = r.applySorting(db, req.Sort)

	// Pagination
	offset := req.GetOffset()
	limit := req.GetLimit()
	if limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

	// Execute query
	err := db.Scan(&properties).Error
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *PostgrePropertyRepo) applyFilters(db *gorm.DB, p *dto.PropertySearchDTO) *gorm.DB {
	// Text search
	if p.Text != "" {
		text := "%" + strings.ToLower(p.Text) + "%"
		db = db.Where("(LOWER(property.title) LIKE ? OR LOWER(property.address_detail) LIKE ? OR LOWER(property.unit_code) LIKE ?)", text, text, text)
	}

	// Property type
	if p.PropertyTypeID != nil {
		db = db.Where("property.property_type_id = ?", *p.PropertyTypeID)
	}

	// Project
	if p.ProjectID != nil {
		db = db.Where("property.project_id = ?", *p.ProjectID)
	}

	// Record status
	if p.RecordStatus != "" {
		var statusVal int
		switch p.RecordStatus {
		case "active":
			statusVal = int(enums.EPropertyStatusActive)
		case "inactive":
			statusVal = int(enums.EPropertyStatusInactive)
		case "archived":
			statusVal = int(enums.EPropertyStatusArchived)
		}
		if statusVal != 0 {
			db = db.Where("property.record_status = ?", statusVal)
		}
	}

	// Scope
	if p.Scope != nil {
		var scopeVal int
		switch *p.Scope {
		case "private":
			scopeVal = int(enums.PropertyScopePrivate)
		case "shared":
			scopeVal = int(enums.PropertyScopeShared)
		case "public":
			scopeVal = int(enums.PropertyScopePublic)
		}
		if scopeVal != 0 {
			db = db.Where("property.scope = ?", scopeVal)
		}
	}

	// Identified (có identifier hoặc national_id)
	if p.Identified != nil {
		if *p.Identified {
			db = db.Where("(property.identifier != '' OR property.national_id != '')")
		} else {
			db = db.Where("(property.identifier = '' OR property.identifier IS NULL) AND (property.national_id = '' OR property.national_id IS NULL)")
		}
	}

	// National verified
	if p.NationalVerified != nil {
		db = db.Where("property.national_id_verified = ?", *p.NationalVerified)
	}

	// Địa chỉ hành chính
	if p.ProvinceID != nil {
		db = db.Where("property.province_id = ?", *p.ProvinceID)
	}
	if p.DistrictID != nil {
		db = db.Where("property.district_id = ?", *p.DistrictID)
	}
	if p.WardID != nil {
		db = db.Where("property.ward_id = ?", *p.WardID)
	}

	// HasAsset, HasProduct, HasListing - sử dụng subquery
	if p.HasAsset != nil {
		if *p.HasAsset {
			db = db.Where("EXISTS (SELECT 1 FROM property_relation WHERE property_id = property.id AND relation_type = ? AND deleted_at IS NULL)", enums.ERelationTypeAsset)
		} else {
			db = db.Where("NOT EXISTS (SELECT 1 FROM property_relation WHERE property_id = property.id AND relation_type = ? AND deleted_at IS NULL)", enums.ERelationTypeAsset)
		}
	}
	if p.HasProduct != nil {
		if *p.HasProduct {
			db = db.Where("EXISTS (SELECT 1 FROM property_relation WHERE property_id = property.id AND relation_type = ? AND deleted_at IS NULL)", enums.ERelationTypeProduct)
		} else {
			db = db.Where("NOT EXISTS (SELECT 1 FROM property_relation WHERE property_id = property.id AND relation_type = ? AND deleted_at IS NULL)", enums.ERelationTypeProduct)
		}
	}
	if p.HasListing != nil {
		if *p.HasListing {
			db = db.Where("EXISTS (SELECT 1 FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id WHERE pr.property_id = property.id AND pr.relation_type = ? AND posts.deleted_at IS NULL)", enums.ERelationTypeProduct)
		} else {
			db = db.Where("NOT EXISTS (SELECT 1 FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id WHERE pr.property_id = property.id AND pr.relation_type = ? AND posts.deleted_at IS NULL)", enums.ERelationTypeProduct)
		}
	}

	// Thời gian
	if p.CreatedFrom != nil {
		db = db.Where("property.created_at >= ?", *p.CreatedFrom)
	}
	if p.CreatedTo != nil {
		db = db.Where("property.created_at <= ?", *p.CreatedTo)
	}
	if p.UpdatedFrom != nil {
		db = db.Where("property.updated_at >= ?", *p.UpdatedFrom)
	}
	if p.UpdatedTo != nil {
		db = db.Where("property.updated_at <= ?", *p.UpdatedTo)
	}

	return db
}

// applyFiltersForCount phiên bản đơn giản để đếm, chỉ áp dụng filter không liên quan đến JOIN phức tạp
func (r *PostgrePropertyRepo) applyFiltersForCount(db *gorm.DB, p *dto.PropertySearchDTO) *gorm.DB {
	// Chỉ cần lọc trên bảng property, nhưng vẫn cần các subquery cho hasAsset, hasProduct, hasListing
	// Nên ta dùng cùng logic applyFilters, nhưng với Table("property")
	return r.applyFilters(db, p)
}

// applySorting xử lý sort
func (r *PostgrePropertyRepo) applySorting(db *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "title_asc":
		db = db.Order("property.title ASC")
	case "title_desc":
		db = db.Order("property.title DESC")
	case "updated_at_asc":
		db = db.Order("property.updated_at ASC")
	case "updated_at_desc":
		db = db.Order("property.updated_at DESC")
	case "priority_product":
		// Ưu tiên property có product trước
		db = db.Order(
			gorm.Expr(
				`(SELECT COUNT(*) 
			  FROM property_relation 
			  WHERE property_id = property.id 
			    AND relation_type = ?) DESC,
			  property.updated_at DESC`,
				enums.ERelationTypeProduct,
			),
		)
	case "priority_asset":
		db = db.Order(
			gorm.Expr(
				`(SELECT COUNT(*) 
			  FROM property_relation 
			  WHERE property_id = property.id 
			    AND relation_type = ?) DESC,
			  property.updated_at DESC`,
				enums.ERelationTypeAsset,
			),
		)

	case "priority_listing":
		db = db.Order(
			gorm.Expr(
				`(SELECT COUNT(*) 
			  FROM property_relation pr
			  JOIN posts 
			    ON pr.relation_id = posts.product_id
			  WHERE pr.property_id = property.id 
			    AND pr.relation_type = ?
			    AND posts.deleted_at IS NULL) DESC,
			  property.updated_at DESC`,
				enums.ERelationTypeProduct,
			),
		)
	case "identified_first":
		db = db.Order("(CASE WHEN property.identifier != '' OR property.national_id != '' THEN 0 ELSE 1 END), property.updated_at DESC")
	case "verified_first":
		db = db.Order("(CASE WHEN property.national_id_verified = true THEN 0 ELSE 1 END), property.updated_at DESC")
	default:
		db = db.Order("property.updated_at DESC")
	}
	return db
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

func (r *PostgrePropertyRepo) Update(ctx context.Context, id uint64, entity *domain.Property) error {
	updates := map[string]interface{}{}

	if entity.PropertyTypeID != nil {
		updates["property_type_id"] = entity.PropertyTypeID
	}
	if entity.Title != "" {
		updates["title"] = entity.Title
	}
	if entity.LocationID != nil {
		updates["location_id"] = entity.LocationID
	}
	if entity.AddressDetail != "" {
		updates["address_detail"] = entity.AddressDetail
	}
	if entity.Latitude != nil {
		updates["latitude"] = entity.Latitude
	}
	if entity.Longitude != nil {
		updates["longitude"] = entity.Longitude
	}
	if entity.MapURL != "" {
		updates["map_url"] = entity.MapURL
	}
	if entity.ProjectID != nil {
		updates["project_id"] = entity.ProjectID
	}
	if entity.BlockID != nil {
		updates["block_id"] = entity.BlockID
	}
	if entity.BuildingID != nil {
		updates["building_id"] = entity.BuildingID
	}
	if entity.Level != "" {
		updates["level"] = entity.Level
	}
	if entity.UnitCode != "" {
		updates["unit_code"] = entity.UnitCode
	}
	if entity.AreaTotal != nil {
		updates["area_total"] = entity.AreaTotal
	}
	if entity.AreaLand != nil {
		updates["area_land"] = entity.AreaLand
	}
	if entity.AreaResidential != nil {
		updates["area_residential"] = entity.AreaResidential
	}
	if entity.LegalNote != "" {
		updates["legal_note"] = entity.LegalNote
	}
	if entity.BuildingInfoID != nil {
		updates["building_info_id"] = entity.BuildingInfoID
	} else {
		// Cho phép set NULL
		updates["building_info_id"] = nil
	}
	if entity.LegalStatus != 0 {
		updates["legal_status"] = entity.LegalStatus
	}
	if entity.RecordStatus != 0 {
		updates["record_status"] = entity.RecordStatus
	}
	if entity.AvatarID != nil {
		updates["avatar_id"] = entity.AvatarID
	} else {
		updates["avatar_id"] = nil
	}
	if entity.ProvinceID != nil {
		updates["province_id"] = entity.ProvinceID
	}
	if entity.DistrictID != nil {
		updates["district_id"] = entity.DistrictID
	}
	if entity.WardID != nil {
		updates["ward_id"] = entity.WardID
	}
	if entity.Scope != 0 {
		updates["scope"] = entity.Scope
	}
	if entity.NationalID != "" {
		updates["national_id"] = entity.NationalID
	}
	if entity.NationalIDVerified {
		updates["national_id_verified"] = entity.NationalIDVerified
	}
	if entity.Identifier != "" {
		updates["identifier"] = entity.Identifier
	}

	if len(updates) == 0 {
		return nil
	}

	return GetDB(ctx, r.DB).
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

func (r *PostgrePropertyRepo) ArchiveBatchProperties(ctx context.Context, ids []uint64, archived bool) error {
	if len(ids) == 0 {
		return nil
	}
	status := enums.EPropertyStatusActive
	if archived {
		status = enums.EPropertyStatusArchived
	}
	return GetDB(ctx, r.DB).
		Model(&domain.Property{}).
		Where("id IN (?) AND deleted_at IS NULL", ids).
		Update("record_status", status).Error
}
