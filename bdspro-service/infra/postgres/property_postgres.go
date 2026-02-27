package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"common/case/crud3"
	"context"
	"strings"
	"time"

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

func (r *PostgrePropertyRepo) Search(
	ctx context.Context,
	req *dto.PropertySearchDTO,
) ([]*domain.Property, int64, error) {
	var rows []dto.PropertyWithCounts
	var total int64

	baseDB := r.DB.WithContext(ctx).
		Table("property").
		Where("property.deleted_at IS NULL")

	baseDB = r.applyFilters(baseDB, req)

	// ===== COUNT =====
	if err := baseDB.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ===== SELECT =====
	db := baseDB.Select(`
		property.*,

		COALESCE(prov.name, '') as province_name,
		COALESCE(dist.name, '') as district_name,
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

		(
			SELECT COUNT(*) 
			FROM property_relation pr 
			WHERE pr.property_id = property.id 
			  AND pr.relation_type = ?
			  AND pr.deleted_at IS NULL
		) as asset_count,

		(
			SELECT COUNT(*) 
			FROM property_relation pr 
			WHERE pr.property_id = property.id 
			  AND pr.relation_type = ?
			  AND pr.deleted_at IS NULL
		) as product_count,

		(
			SELECT COUNT(*) 
			FROM property_relation pr
			JOIN posts ON pr.relation_id = posts.product_id
			WHERE pr.property_id = property.id
			  AND pr.relation_type = ?
			  AND posts.deleted_at IS NULL
		) as listing_count
	`,
		enums.ERelationTypeAsset,
		enums.ERelationTypeProduct,
		enums.ERelationTypeProduct,
	)

	db = db.
		Joins("LEFT JOIN province_v2 prov ON property.province_id = prov.id AND prov.deleted_at IS NULL").
		Joins("LEFT JOIN district_v2 dist ON property.district_id = dist.id AND dist.deleted_at IS NULL").
		Joins("LEFT JOIN ward_v2 ward ON property.ward_id = ward.id AND ward.deleted_at IS NULL").
		Joins("LEFT JOIN projects proj ON property.project_id = proj.id AND proj.deleted_at IS NULL").
		Joins("LEFT JOIN property_type pt ON property.property_type_id = pt.id AND pt.deleted_at IS NULL").
		Joins("LEFT JOIN property_building_info pbi ON property.id = pbi.property_id AND pbi.deleted_at IS NULL").
		Joins("LEFT JOIN property_media avatar ON property.avatar_id = avatar.id AND avatar.deleted_at IS NULL")

	db = r.applySorting(db, req.Sort)

	// pagination
	offset := req.GetOffset()
	limit := req.GetLimit()
	if limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	// ===== MAP → DOMAIN =====
	properties := make([]*domain.Property, 0, len(rows))

	for _, row := range rows {
		p := &domain.Property{
			ID:             row.ID,
			Title:          row.Title,
			AddressDetail:  row.AddressDetail,
			ProvinceID:     row.ProvinceID,
			DistrictID:     row.DistrictID,
			WardID:         row.WardID,
			ProjectID:      row.ProjectID,
			PropertyTypeID: row.PropertyTypeID,
			Scope:          enums.EPropertyScope(row.Scope),
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
			ProductCount:   row.ProductCount,
			AssetCount:     row.AssetCount,
			ListingCount:   row.ListingCount,
		}

		// ===== ADDRESS REL =====
		if row.ProvinceName != "" {
			p.Province = &domain.ProvinceV2{Name: row.ProvinceName}
		}
		if row.DistrictName != "" {
			p.District = &domain.DistrictV2{Name: row.DistrictName}
		}
		if row.WardName != "" {
			p.Ward = &domain.WardV2{Name: row.WardName}
		}

		// ===== PROPERTY TYPE =====
		if row.PropertyTypeName != "" {
			p.PropertyType = &domain.PropertyTypeItem{
				Name: row.PropertyTypeName,
			}
		}

		// ===== AVATAR =====
		if row.AvatarURL != "" {
			p.Avatar = &domain.PropertyMedia{
				URL: row.AvatarURL,
			}
		}

		properties = append(properties, p)
	}

	return properties, total, nil
}

func (r *PostgrePropertyRepo) applyFilters(db *gorm.DB, p *dto.PropertySearchDTO) *gorm.DB {
	if p.SourceType != nil {
		db = db.Where("property.source_type = ?", *p.SourceType)
	}

	if p.Text != "" {
		text := "%" + strings.ToLower(p.Text) + "%"
		db = db.Where(`
			LOWER(property.title) LIKE ? OR 
			LOWER(property.address_detail) LIKE ? OR 
			LOWER(property.unit_code) LIKE ?`,
			text, text, text)
	}

	if p.PropertyTypeID != nil {
		db = db.Where("property.property_type_id = ?", *p.PropertyTypeID)
	}

	if p.ProjectID != nil {
		db = db.Where("property.project_id = ?", *p.ProjectID)
	}

	if p.ProvinceID != nil {
		db = db.Where("property.province_id = ?", *p.ProvinceID)
	}
	if p.DistrictID != nil {
		db = db.Where("property.district_id = ?", *p.DistrictID)
	}
	if p.WardID != nil {
		db = db.Where("property.ward_id = ?", *p.WardID)
	}

	// Time filter
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

func (r *PostgrePropertyRepo) GetDetail(
	ctx context.Context,
	id uint64,
	sourceType uint32,
) (*domain.Property, error) {
	var property domain.Property

	db := r.DB.WithContext(ctx).
		Preload("PropertyType").
		Preload("Project").
		Preload("Province").
		Preload("MediaList").
		Preload("Ward").
		Preload("BuildingInfo").
		Preload("LandInfo").
		Preload("Avatar").
		Preload("Amenities").
		Where("id = ? AND deleted_at IS NULL", id).
		Where("source_type = ?", sourceType)

	if err := db.First(&property).Error; err != nil {
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

func (r *PostgrePropertyRepo) GetUpdatedAt(c context.Context, id uint64) (int64, error) {
	var result time.Time

	err := GetDB(c, r.DB).
		Model(&domain.Property{}).
		Select("updated_at").
		Where("id = ?", id).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Unix(), nil
}
