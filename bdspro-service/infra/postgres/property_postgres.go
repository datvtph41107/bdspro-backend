package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"common/case/crud3"
	_models "common/models"
	_utils "common/utils"
	"context"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostgrePropertyRepo struct {
	crud3.BaseRepo[domain.PropertyLineage]
}

func NewPostgrePropertyRepo(db *gorm.DB) repo.PropertyRepo {
	return &PostgrePropertyRepo{
		BaseRepo: crud3.BaseRepo[domain.PropertyLineage]{DB: db},
	}
}

func (r *PostgrePropertyRepo) GetRefreshSyncIds(c context.Context, req *dto.CheckVersionSyncRequest) ([]uint64, error) {
	// profileId := _utils.GetProfileIdWithContext(c)
	// originId := _utils.GetOriginIdFromContext(c)
	var propertyIds []uint64
	err := r.DB.WithContext(c).
		Model(&domain.PropertyLineage{}).
		// Joins(`LEFT JOIN product_user pu ON (pu.product_id = products.id AND pu.deleted_at IS NULL)`).
		Where(`property.deleted_at is null `).
		Order("property.id desc").
		Pluck("property.id", &propertyIds).Error

	if err != nil {
		return nil, err
	}

	return propertyIds, nil
}

func (r *PostgrePropertyRepo) Search(ctx context.Context, req *dto.PropertySearchDTO, userId uint64) ([]*domain.PropertyLineage, int64, error) {
	var rows []dto.PropertyWithCounts
	var total int64

	baseDB := r.DB.WithContext(ctx).
		Table("property")

	if !req.Sync {
		baseDB = baseDB.Where("property.deleted_at IS NULL")
	}

	baseDB = r.applyFilters(baseDB, req, userId)

	if err := baseDB.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db := baseDB.Select(`
		property.id,
		property.title,
		property.address_detail,
		property.area_total,
		property.record_status,
		property.created_at,
		property.updated_at,
		property.created_by,
		property.province_id,
		property.district_id,
		property.ward_id,
		property.project_id,
		property.property_type_id,
		property.scope,
		property.source_type,
		prov.name as province_name,
		dist.name as district_name,
		ward.name as ward_name,
		proj.name as project_name,
		pt.name as property_type_name,
		avatar.media_url as avatar_url,
		avatar.thumb_url as avatar_thumb_url,
		0 as asset_count,
		0 as product_count,
		0 as listing_count
	`)
	originId := _utils.GetOriginIdFromContext(ctx)
	db = db.
		Joins("LEFT JOIN property_user pu ON pu.property_lineage_id = property.id AND pu.deleted_at IS NULL AND pu.owner_origin_id = ?", originId).
		Joins("LEFT JOIN province_v2 prov ON property.province_id = prov.id AND prov.deleted_at IS NULL").
		Joins("LEFT JOIN district_v2 dist ON property.district_id = dist.id AND dist.deleted_at IS NULL").
		Joins("LEFT JOIN ward_v2 ward ON property.ward_id = ward.id AND ward.deleted_at IS NULL").
		Joins("LEFT JOIN projects proj ON property.project_id = proj.id AND proj.deleted_at IS NULL").
		Joins("LEFT JOIN property_type pt ON property.property_type_id = pt.id AND pt.deleted_at IS NULL").
		Joins("LEFT JOIN property_media avatar ON property.avatar_id = avatar.id AND avatar.deleted_at IS NULL")

	if req.RequestUserId != nil {
		db.Where("(pu.owner_origin_id = ?)", req.RequestUserId)
	} else {
		// db.Where("(property_lineage.source_type = 10)")
	}
	// if len(whereClauses) > 0 {
	// 	db.Where(strings.Join(whereClauses, " AND "))
	// }

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
	properties := make([]*domain.PropertyLineage, 0, len(rows))

	// for _, row := range rows {
	// 	p := &domain.PropertyBase{
	// 		ID:             row.ID,
	// 		Title:          row.Title,
	// 		AddressDetail:  row.AddressDetail,
	// 		ProvinceID:     row.ProvinceID,
	// 		DistrictID:     row.DistrictID,
	// 		WardID:         row.WardID,
	// 		ProjectID:      row.ProjectID,
	// 		PropertyTypeID: row.PropertyTypeID,
	// 		Scope:          enums.EPropertyScope(row.Scope),
	// 		CreatedAt:      row.CreatedAt,
	// 		UpdatedAt:      &row.UpdatedAt,
	// 		ProductCount:   row.ProductCount,
	// 		AssetCount:     row.AssetCount,
	// 		ListingCount:   row.ListingCount,
	// 		AreaTotal:      &row.AreaTotal,
	// 		SourceType:     enums.EPropertySourceType(row.SourceType),
	// 		RecordStatus:   enums.EPropertyStatus(row.StatusRecord),
	// 		Visibility:     enums.EVisibility(row.Visibility),
	// 		CreatedBy:      row.CreatedBy,
	// 	}

	// 	// ===== ADDRESS REL =====
	// 	if row.ProvinceName != "" {
	// 		p.Province = &domain.ProvinceV2{Name: row.ProvinceName}
	// 	}
	// 	if row.DistrictName != "" {
	// 		p.District = &domain.DistrictV2{Name: row.DistrictName}
	// 	}
	// 	if row.WardName != "" {
	// 		p.Ward = &domain.WardV2{Name: row.WardName}
	// 	}

	// 	// ===== PROPERTY TYPE =====
	// 	if row.PropertyTypeName != "" {
	// 		p.PropertyType = &domain.PropertyTypeItem{
	// 			Name: row.PropertyTypeName,
	// 		}
	// 	}

	// 	// ===== AVATAR =====
	// 	// if row.AvatarURL != "" {
	// 	// 	p.Avatar = &domain.PropertyMedia{
	// 	// 		URL: row.AvatarURL,
	// 	// 	}
	// 	// }

	// 	properties = append(properties, p)
	// }

	return properties, total, nil
}

func (r *PostgrePropertyRepo) applyFilters(db *gorm.DB, p *dto.PropertySearchDTO, userId uint64) *gorm.DB {
	if p.SourceType != nil {
		db = db.Where("property.source_type = ?", *p.SourceType)
	}

	if p.Owned {
		db = db.Where("property.created_by = ?", userId)
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
) (*domain.PropertyLineage, error) {
	var property domain.PropertyLineage

	db := r.DB.WithContext(ctx).
		Preload("PropertyIdentify").
		Preload("Location", "deleted_at IS NULL").
		Preload("Location.Province", "deleted_at IS NULL").
		Preload("Location.District", "deleted_at IS NULL").
		Preload("Location.Ward", "deleted_at IS NULL").
		Preload("LandInfo", "deleted_at IS NULL").
		Preload("BuildingInfo", "deleted_at IS NULL").
		Preload("Edvidence", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", id)

	if err := db.First(&property).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Load counts
	r.DB.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM property_relation WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL",
		id, enums.ERelationTypeAsset,
	).Scan(&property.AssetCount)
	r.DB.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM property_relation WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL",
		id, enums.ERelationTypeProduct,
	).Scan(&property.ProductCount)
	r.DB.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id WHERE pr.property_id = ? AND pr.relation_type = ? AND pr.deleted_at IS NULL AND posts.deleted_at IS NULL",
		id, enums.ERelationTypeProduct,
	).Scan(&property.ListingCount)

	return &property, nil
}

func (r *PostgrePropertyRepo) GetDetail2(
	ctx context.Context,
	id uint64,
	sourceType uint32,
) (*domain.PropertyLineage, error) {
	const q = `
SELECT
	pl.id, pl.property_identify_id, pl.location_id, pl.property_info_id,
	pl.land_info_id, pl.edvidence_id, pl.building_info_id, pl.external_ref_id,
	pl.national_id, pl.verified_national_at, pl.created_at, pl.updated_at,
	pl.deleted_at, pl.created_by, pl.updated_by,
	pif.id AS info_id, pif.title, pif.unit_code, pif.identifier, pif.level, pif.note, pif.source_type, pif.record_status, pif.visibility, pif.legal_status, pif.privacy_level, pif.origin_profile_id, pif.project_id,
	COALESCE(pt.id, 0) AS project_id, COALESCE(pt.name, '') AS project_name,
	COALESCE(pi.id, 0) AS pi_id, COALESCE(pi.p_id, 0) AS pi_pid, COALESCE(pi.version, 0) AS pi_version, COALESCE(pi.created_at, CURRENT_TIMESTAMP) AS pi_created_at, COALESCE(pi.updated_at, CURRENT_TIMESTAMP) AS pi_updated_at,
	loc.id AS loc_id, loc.address_detail, loc.province_id, loc.district_id, loc.ward_id,
	loc.latitude, loc.longitude, loc.map_url, loc.region_id,
	prov.name AS province_name, prov.code AS province_code,
	dist.name AS district_name, dist.code AS district_code,
	ward.name AS ward_name, ward.code AS ward_code,
	li.id AS land_id, li.document_type, li.document_no, li.issuring_auth, li.plot, li.sheet,
	li.area_total AS land_area_total, li.area_land, li.area_plant, li.purpose_used, li.note AS land_note,
	li.front_width, li.depth, li.street_width, li.land_note AS land_note2, li.expired_land, li.expired_plant,
	bi.id AS bi_id, bi.area_actual, bi.area_floor, bi.area_construction, bi.floors,
	bi.room_number, bi.bedrooms, bi.bathrooms, bi.build_status, bi.building_type,
	bi.direction, bi.balcony_direction, bi.note AS bi_note,
	ev.id AS ev_id, ev.title AS ev_title, ev.file_id, ev.description AS ev_description,
	(SELECT COALESCE(json_agg(json_build_object(
			'id', pm.id,
			'mediaUrl', pm.media_url,
			'mediaType', pm.media_type,
			'sortOrder', pm.sort_order
		)), '[]')
		FROM property_media pm
		WHERE pm.lineage_id = pl.id
		AND pm.deleted_at IS NULL
	) AS media_list_raw,
	 (SELECT COALESCE(json_agg(json_build_object(
			'id', a.id,
			'name', a.name
		)), '[]')
		FROM property_amenity pa
		LEFT JOIN amenity a ON pa.amenity_item_id = a.id AND a.deleted_at IS NULL
		WHERE pa.property_lineage_id = pl.id
		AND a.active = true
	) AS amenities_raw,
	(SELECT COUNT(*) FROM property_relation WHERE property_id = pl.id AND relation_type = ? AND deleted_at IS NULL) AS asset_count,
	(SELECT COUNT(*) FROM property_relation WHERE property_id = pl.id AND relation_type = ? AND deleted_at IS NULL) AS product_count,
	(SELECT COUNT(*) FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id
	 WHERE pr.property_id = pl.id AND pr.relation_type = ? AND pr.deleted_at IS NULL AND posts.deleted_at IS NULL) AS listing_count
FROM property_lineage pl
LEFT JOIN property_identify pi ON pl.property_identify_id = pi.id AND pi.deleted_at IS NULL
LEFT JOIN property_info pif ON pl.property_info_id = pif.id AND pif.deleted_at IS NULL
LEFT JOIN property_location loc ON pl.location_id = loc.id
LEFT JOIN province_v2 prov ON loc.province_id = prov.id AND prov.deleted_at IS NULL
LEFT JOIN district_v2 dist ON loc.district_id = dist.id AND dist.deleted_at IS NULL
LEFT JOIN ward_v2 ward ON loc.ward_id = ward.id AND ward.deleted_at IS NULL
LEFT JOIN property_land_info li ON pl.land_info_id = li.id AND li.deleted_at IS NULL
LEFT JOIN property_building_info bi ON pl.building_info_id = bi.id AND bi.deleted_at IS NULL
LEFT JOIN property_edvidence ev ON pl.edvidence_id = ev.id AND ev.deleted_at IS NULL
LEFT JOIN projects pt ON pif.project_id = pt.id AND pt.deleted_at IS NULL
WHERE pl.id = ? AND pl.deleted_at IS NULL`
	var row struct {
		ID                 uint64
		PropertyIdentifyID *uint64
		LocationID         *uint64
		PropertyInfoID     *uint64
		LandInfoID         *uint64
		EdvidenceID        *uint64
		BuildingInfoID     *uint64
		ExternalRefID      *uint64
		NationalID         string
		VerifiedNationalAt *time.Time
		CreatedAt          time.Time
		UpdatedAt          *time.Time
		DeletedAt          *time.Time
		CreatedBy          *uint64
		UpdatedBy          *uint64
		// Media
		MediaListRaw *json.RawMessage
		AmenitiesRaw *json.RawMessage
		// PropertyInfo
		InfoID          *uint64
		Title           *string
		UnitCode        *string
		Identifier      *string
		Level           *string
		Note            *string
		SourceType      *uint32
		RecordStatus    *uint32
		Visibility      *uint32
		LegalStatus     *uint32
		PrivacyLevel    *uint32
		OriginProfileId *uint64
		ProjectID       *uint64
		ProjectName     *string
		// PropertyIdentify
		PiID        *uint64
		PiPid       *uint64
		PiVersion   *uint32
		PiCreatedAt *time.Time
		PiUpdatedAt *time.Time
		// Location
		LocID         *uint64
		AddressDetail *string
		ProvinceID    *uint64
		DistrictID    *uint64
		WardID        *uint64
		Latitude      *float64
		Longitude     *float64
		MapURL        *string
		RegionID      *uint64
		ProvinceName  *string
		ProvinceCode  *int
		DistrictName  *string
		DistrictCode  *int
		WardName      *string
		WardCode      *int
		// LandInfo
		LandID        *uint64
		DocumentType  *uint32
		DocumentNo    *string
		IssuringAuth  *string
		Plot          *uint32
		Sheet         *uint32
		LandAreaTotal *float64
		AreaLand      *float64
		AreaPlant     *float64
		PurposeUsed   *int
		LandNote      *string
		FrontWidth    *float64
		Depth         *float64
		StreetWidth   *float64
		LandNote2     *string
		ExpiredLand   *time.Time
		ExpiredPlant  *time.Time
		// BuildingInfo
		BiID             *uint64
		AreaActual       *float64
		AreaFloor        *float64
		AreaConstruction *float64
		Floors           *uint32
		RoomNumber       *uint32
		Bedrooms         *uint32
		Bathrooms        *uint32
		BuildStatus      *int
		BuildingType     *int
		Direction        *int
		BalconyDirection *int
		BiNote           *string
		// Edvidence
		EvID          *uint64
		EvTitle       *string
		EvFileID      *uint64
		EvDescription *string
		AssetCount    int64
		ProductCount  int64
		ListingCount  int64
	}
	if err := r.DB.WithContext(ctx).Raw(q,
		enums.ERelationTypeAsset, enums.ERelationTypeProduct, enums.ERelationTypeProduct, id,
	).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	p := &domain.PropertyLineage{}
	p.ID = row.ID
	p.PropertyIdentifyID = row.PropertyIdentifyID
	p.LocationID = row.LocationID
	p.PropertyInfoID = row.PropertyInfoID
	p.LandInfoID = row.LandInfoID
	p.EdvidenceID = row.EdvidenceID
	p.BuildingInfoID = row.BuildingInfoID
	p.ExternalRefID = row.ExternalRefID
	p.NationalID = row.NationalID
	p.VerifiedNationalAt = row.VerifiedNationalAt
	p.AssetCount = row.AssetCount
	p.ProductCount = row.ProductCount
	p.ListingCount = row.ListingCount
	if row.PiID != nil {
		pi := &domain.PropertyIdentify{}
		pi.ID = *row.PiID
		if row.PiPid != nil {
			pi.PID = *row.PiPid
		}
		if row.PiVersion != nil {
			pi.Version = *row.PiVersion
		}
		if row.PiCreatedAt != nil {
			pi.CreatedAt = row.PiCreatedAt
		}
		if row.PiUpdatedAt != nil {
			pi.UpdatedAt = row.PiUpdatedAt
		}
		p.PropertyIdentify = pi
	}
	if row.InfoID != nil {
		info := &domain.PropertyInfo{
			UnitCode:   ptrStr(row.UnitCode),
			Identifier: ptrStr(row.Identifier),
			Level:      ptrStr(row.Level), Note: ptrStr(row.Note),
			SourceType:      enums.EPropertySourceType(*row.SourceType),
			RecordStatus:    enums.EPropertyStatus(*row.RecordStatus),
			Visibility:      enums.EVisibility(*row.Visibility),
			LegalStatus:     enums.EHouseCertificate(*row.LegalStatus),
			PrivacyLevel:    enums.EVisibility(*row.PrivacyLevel),
			OriginProfileId: row.OriginProfileId,
			ProjectID:       row.ProjectID,
			Title:           ptrStr(row.Title),
		}
		if row.ProjectID != nil {
			info.Project = &domain.Project{
				BaseEntity: _models.BaseEntity{
					ID: *row.ProjectID,
				},
				Name: ptrStr(row.ProjectName),
			}
		}
		info.ID = *row.InfoID
		p.PropertyInfo = info
	}
	if row.LocID != nil {
		loc := &domain.PropertyLocation{
			ID: *row.LocID, AddressDetail: ptrStr(row.AddressDetail),
			ProvinceID: row.ProvinceID, DistrictID: row.DistrictID, WardID: row.WardID,
			Latitude: row.Latitude, Longitude: row.Longitude, MapURL: ptrStr(row.MapURL), RegionID: row.RegionID,
		}
		if row.ProvinceName != nil {
			loc.Province = &domain.ProvinceV2{Name: *row.ProvinceName}
			if row.ProvinceCode != nil {
				loc.Province.Code = *row.ProvinceCode
			}
		}
		if row.DistrictName != nil {
			loc.District = &domain.DistrictV2{Name: *row.DistrictName}
			if row.DistrictCode != nil {
				loc.District.Code = *row.DistrictCode
			}
		}
		if row.WardName != nil {
			loc.Ward = &domain.WardV2{Name: *row.WardName}
			if row.WardCode != nil {
				loc.Ward.Code = *row.WardCode
			}
		}
		p.Location = loc
	}
	if row.LandID != nil {
		li := &domain.PropertyLandInfo{
			DocumentType: ptrU32(row.DocumentType),
			DocumentNo:   ptrStr(row.DocumentNo),
			IssuringAuth: ptrStr(row.IssuringAuth),
			Plot:         ptrU32(row.Plot),
			Sheet:        ptrU32(row.Sheet),
			AreaTotal:    row.LandAreaTotal,
			AreaLand:     row.AreaLand,
			AreaPlant:    row.AreaPlant,
			Note:         ptrStr(row.LandNote),
			FrontWidth:   row.FrontWidth,
			Depth:        row.Depth,
			StreetWidth:  row.StreetWidth,
			LandNote:     ptrStr(row.LandNote2),
			ExpiredLand:  row.ExpiredLand,
			ExpiredPlant: row.ExpiredPlant,
		}
		li.ID = *row.LandID
		if row.PurposeUsed != nil {
			li.PurposeUsed = enums.EPurposeUsed(*row.PurposeUsed)
		}
		p.LandInfo = li
	}
	if row.BiID != nil {
		bi := &domain.PropertyBuildingInfo{
			AreaActual: row.AreaActual, AreaFloor: row.AreaFloor, AreaConstruction: row.AreaConstruction,
			Floors: row.Floors, RoomNumber: row.RoomNumber, Bedrooms: row.Bedrooms, Bathrooms: row.Bathrooms,
			Note: ptrStr(row.BiNote),
		}
		bi.ID = *row.BiID
		if row.BuildStatus != nil {
			bi.BuildStatus = enums.EBuildStatus(*row.BuildStatus)
		}
		if row.BuildingType != nil {
			bi.BuildingType = enums.EBuildingType(*row.BuildingType)
		}
		if row.Direction != nil {
			bi.Direction = enums.EHouseOrient(*row.Direction)
		}
		if row.BalconyDirection != nil {
			bi.BalconyDirection = enums.EHouseOrient(*row.BalconyDirection)
		}
		p.BuildingInfo = bi
	}
	if row.EvID != nil {
		ev := &domain.PropertyEdvidence{Title: ptrStr(row.EvTitle), FileID: row.EvFileID, Description: ptrStr(row.EvDescription)}
		ev.ID = *row.EvID
		p.Edvidence = ev
	}
	// MediaList + Amenities: 2 query nhỏ
	var media []domain.PropertyMedia
	var amenities []domain.AmenityItem
	if row.MediaListRaw != nil {
		if err := json.Unmarshal(*row.MediaListRaw, &media); err != nil {
			return nil, err
		}
	}
	if row.AmenitiesRaw != nil {
		if err := json.Unmarshal(*row.AmenitiesRaw, &amenities); err != nil {
			return nil, err
		}
	}
	p.MediaList = media
	p.Amenities = amenities
	return p, nil
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func ptrU32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

func (r *PostgrePropertyRepo) Update(ctx context.Context, id uint64, entity *domain.PropertyLineage) error {
	updates := map[string]interface{}{}

	// if entity.PropertyTypeID != nil {
	// 	updates["property_type_id"] = entity.PropertyTypeID
	// }
	// if entity.Title != "" {
	// 	updates["title"] = entity.Title
	// }
	// if entity.LocationID != nil {
	// 	updates["location_id"] = entity.LocationID
	// }
	// if entity.AddressDetail != "" {
	// 	updates["address_detail"] = entity.AddressDetail
	// }
	// if entity.Latitude != nil {
	// 	updates["latitude"] = entity.Latitude
	// }
	// if entity.Longitude != nil {
	// 	updates["longitude"] = entity.Longitude
	// }
	// if entity.MapURL != "" {
	// 	updates["map_url"] = entity.MapURL
	// }
	// if entity.ProjectID != nil {
	// 	updates["project_id"] = entity.ProjectID
	// }
	// if entity.Level != "" {
	// 	updates["level"] = entity.Level
	// }
	// if entity.UnitCode != "" {
	// 	updates["unit_code"] = entity.UnitCode
	// }
	// if entity.AreaTotal != nil {
	// 	updates["area_total"] = entity.AreaTotal
	// }
	// // if entity.AreaLand != nil {
	// // 	updates["area_land"] = entity.AreaLand
	// // }
	// // if entity.AreaResidential != nil {
	// // 	updates["area_residential"] = entity.AreaResidential
	// // }
	// // if entity.LegalNote != "" {
	// // 	updates["legal_note"] = entity.LegalNote
	// // }
	// if entity.BuildingInfoID != nil {
	// 	updates["building_info_id"] = entity.BuildingInfoID
	// } else {
	// 	// Cho phép set NULL
	// 	updates["building_info_id"] = nil
	// }
	// // if entity.LegalStatus != 0 {
	// // 	updates["legal_status"] = entity.LegalStatus
	// // }
	// if entity.RecordStatus != 0 {
	// 	updates["record_status"] = entity.RecordStatus
	// }
	// if entity.AvatarID != nil {
	// 	updates["avatar_id"] = entity.AvatarID
	// } else {
	// 	updates["avatar_id"] = nil
	// }
	// if entity.ProvinceID != nil {
	// 	updates["province_id"] = entity.ProvinceID
	// }
	// if entity.DistrictID != nil {
	// 	updates["district_id"] = entity.DistrictID
	// }
	// if entity.WardID != nil {
	// 	updates["ward_id"] = entity.WardID
	// }
	// if entity.Scope != 0 {
	// 	updates["scope"] = entity.Scope
	// }
	// if entity.NationalID != "" {
	// 	updates["national_id"] = entity.NationalID
	// }
	// if entity.NationalIDVerified {
	// 	updates["national_id_verified"] = entity.NationalIDVerified
	// }
	// if entity.Identifier != "" {
	// 	updates["identifier"] = entity.Identifier
	// }

	if len(updates) == 0 {
		return nil
	}

	return GetDB(ctx, r.DB).
		Model(&domain.PropertyLineage{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *PostgrePropertyRepo) Create(c context.Context, entity *domain.PropertyLineage) error {
	return GetDB(c, r.DB).Create(entity).Error
}

func (r *PostgrePropertyRepo) UpdateOption(c context.Context, id uint64, property *domain.PropertyLineage) error {
	// Chỉ update các field cụ thể: area_land, area_total, project_id, property_type_id
	updates := map[string]interface{}{}

	// if property.AreaLand != nil {
	// 	updates["area_land"] = property.AreaLand
	// }
	// if property.AreaTotal != nil {
	// 	updates["area_total"] = property.AreaTotal
	// }
	// if property.ProjectID != nil {
	// 	updates["project_id"] = property.ProjectID
	// }
	// if property.PropertyTypeID != nil {
	// 	updates["property_type_id"] = property.PropertyTypeID
	// }

	// Nếu không có field nào cần update, return nil
	if len(updates) == 0 {
		return nil
	}

	return GetDB(c, r.DB).
		Model(&domain.PropertyLineage{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *PostgrePropertyRepo) GetOneByID(c context.Context, id uint64) (*domain.PropertyLineage, error) {
	var property domain.PropertyLineage
	if err := GetDB(c, r.DB).
		Model(&domain.PropertyLineage{}).
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
		Delete(&domain.PropertyLineage{}).Error
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
		Model(&domain.PropertyLineage{}).
		Where("id IN (?) AND deleted_at IS NULL", ids).
		Update("record_status", status).Error
}

func (r *PostgrePropertyRepo) GetUpdatedAt(c context.Context, id uint64) (int64, error) {
	var result time.Time

	err := GetDB(c, r.DB).
		Model(&domain.PropertyLineage{}).
		Select("updated_at").
		Where("id = ?", id).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Unix(), nil
}
