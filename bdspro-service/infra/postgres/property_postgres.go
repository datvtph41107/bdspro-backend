package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"common/case/crud3"
	_dto "common/domain/dto"
	_domain "common/domain/entity"
	_models "common/models"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"math"
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

// FindDuplicates tìm các property trùng lặp dựa trên tọa độ, bán kính, loại BĐS và diện tích.
// Công thức Haversine được dùng để tính khoảng cách.
func (r *PostgrePropertyRepo) FindDuplicates(
	ctx context.Context,
	lat, lng float64,
	radiusMeters float64,
	propertyTypeID uint64,
	area float64,
	tolerancePercent float64,
) ([]*domain.PropertyLineage, error) {
	const MaxDuplicateReturn = 10
	var lineages []*domain.PropertyLineage

	radiusKm := radiusMeters / 1000.0

	// Tính bounding box để giảm số bản ghi cần tính Haversine
	latDiff := radiusKm / 111.0 // 1 độ ≈ 111km
	lngDiff := radiusKm / (111.0 * math.Cos(lat*math.Pi/180))

	query := `
        SELECT pl.*
        FROM property_lineage pl
        JOIN property_location loc ON pl.location_id = loc.id
        JOIN property_info pi ON pl.property_info_id = pi.id
        LEFT JOIN property_land_info land ON pl.land_info_id = land.id
        LEFT JOIN property_building_info build ON pl.building_info_id = build.id
        WHERE pl.deleted_at IS NULL
          AND pi.property_type_id = $1
          AND COALESCE(land.area_total, build.area_actual) BETWEEN $2 * (1 - $3) AND $2 * (1 + $3)
          -- Bounding box filter
          AND loc.latitude BETWEEN $4 - $5 AND $4 + $5
          AND loc.longitude BETWEEN $6 - $7 AND $6 + $7
          -- Haversine chính xác
          AND (
              6371 * acos( 
                  LEAST(1, GREATEST(-1,
                      cos(radians($4)) * cos(radians(loc.latitude)) * 
                      cos(radians(loc.longitude) - radians($6)) + 
                      sin(radians($4)) * sin(radians(loc.latitude))
                  ))
              ) <= $8
          )
        LIMIT $9
        FOR UPDATE SKIP LOCKED
    `

	err := r.DB.WithContext(ctx).Raw(query,
		propertyTypeID,
		area, tolerancePercent,
		lat, latDiff,
		lng, lngDiff,
		radiusKm,
		MaxDuplicateReturn,
	).Scan(&lineages).Error

	if err != nil {
		return nil, fmt.Errorf("FindDuplicates query failed: %w", err)
	}

	return lineages, nil
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

// chỉ dùng để lấy định danh
func (r *PostgrePropertyRepo) Search(ctx context.Context, req *dto.PropertySearchDTO, userId uint64) ([]*domain.PropertyLineage, int64, error) {
	originId := _utils.GetOriginIdFromContext(ctx)

	whereClauses := []string{"pi.deleted_at IS NULL", "pl.deleted_at IS NULL"}
	args := []any{originId}

	if req != nil {
		// if req.SourceType != nil {
		// 	whereClauses = append(whereClauses, "pif.source_type = ?")
		// 	args = append(args, *req.SourceType)
		// }
		if req.Owned && userId != 0 {
			whereClauses = append(whereClauses, `EXISTS (
				SELECT 1 FROM property_user pu_owned
				WHERE pu_owned.property_lineage_id = pl.id
				  AND pu_owned.owner_origin_id = ?
				  AND pu_owned.deleted_at IS NULL
			)`)
			args = append(args, userId)
		}
		if req.Text != "" {
			text := "%" + strings.ToLower(req.Text) + "%"
			whereClauses = append(whereClauses, `(LOWER(COALESCE(pif.title, '')) LIKE ? OR LOWER(COALESCE(loc.address_detail, '')) LIKE ? OR LOWER(COALESCE(pif.unit_code, '')) LIKE ? OR LOWER(COALESCE(pif.identifier, '')) LIKE ?)`)
			args = append(args, text, text, text, text)
		}
		if req.PropertyTypeID != nil {
			whereClauses = append(whereClauses, "pif.property_type_id = ?")
			args = append(args, *req.PropertyTypeID)
		}
		if req.ProjectID != nil {
			whereClauses = append(whereClauses, "pif.project_id = ?")
			args = append(args, *req.ProjectID)
		}
		if req.ProvinceID != nil {
			whereClauses = append(whereClauses, "loc.province_id = ?")
			args = append(args, *req.ProvinceID)
		}
		if req.DistrictID != nil {
			whereClauses = append(whereClauses, "loc.district_id = ?")
			args = append(args, *req.DistrictID)
		}
		if req.WardID != nil {
			whereClauses = append(whereClauses, "loc.ward_id = ?")
			args = append(args, *req.WardID)
		}
		if req.CreatedFrom != nil {
			whereClauses = append(whereClauses, "pl.created_at >= ?")
			args = append(args, *req.CreatedFrom)
		}
		if req.CreatedTo != nil {
			whereClauses = append(whereClauses, "pl.created_at <= ?")
			args = append(args, *req.CreatedTo)
		}
		if req.UpdatedFrom != nil {
			whereClauses = append(whereClauses, "COALESCE(pl.updated_at, pl.created_at) >= ?")
			args = append(args, *req.UpdatedFrom)
		}
		if req.UpdatedTo != nil {
			whereClauses = append(whereClauses, "COALESCE(pl.updated_at, pl.created_at) <= ?")
			args = append(args, *req.UpdatedTo)
		}
		if req.Scope != nil {
			switch strings.ToLower(*req.Scope) {
			case "private":
				whereClauses = append(whereClauses, "pif.scope = ?")
				args = append(args, enums.PropertyScopePrivate)
			case "shared":
				whereClauses = append(whereClauses, "pif.scope = ?")
				args = append(args, enums.PropertyScopeShared)
			case "public":
				whereClauses = append(whereClauses, "pif.scope = ?")
				args = append(args, enums.PropertyScopePublic)
			}
		}
		if req.Identified != nil && !*req.Identified {
			return []*domain.PropertyLineage{}, 0, nil
		}
		if req.NationalVerified != nil {
			if *req.NationalVerified {
				whereClauses = append(whereClauses, "pl.verified_national_at IS NOT NULL")
			} else {
				whereClauses = append(whereClauses, "pl.verified_national_at IS NULL")
			}
		}
		if req.HasAsset != nil {
			if *req.HasAsset {
				whereClauses = append(whereClauses, `EXISTS (
					SELECT 1 FROM property_relation pr_asset
					WHERE pr_asset.property_id = pl.id
					  AND pr_asset.relation_type = ?
					  AND pr_asset.deleted_at IS NULL
				)`)
			} else {
				whereClauses = append(whereClauses, `NOT EXISTS (
					SELECT 1 FROM property_relation pr_asset
					WHERE pr_asset.property_id = pl.id
					  AND pr_asset.relation_type = ?
					  AND pr_asset.deleted_at IS NULL
				)`)
			}
			args = append(args, enums.ERelationTypeAsset)
		}
		if req.HasProduct != nil {
			if *req.HasProduct {
				whereClauses = append(whereClauses, `EXISTS (
					SELECT 1 FROM property_relation pr_product
					WHERE pr_product.property_id = pl.id
					  AND pr_product.relation_type = ?
					  AND pr_product.deleted_at IS NULL
				)`)
			} else {
				whereClauses = append(whereClauses, `NOT EXISTS (
					SELECT 1 FROM property_relation pr_product
					WHERE pr_product.property_id = pl.id
					  AND pr_product.relation_type = ?
					  AND pr_product.deleted_at IS NULL
				)`)
			}
			args = append(args, enums.ERelationTypeProduct)
		}
		if req.HasListing != nil {
			if *req.HasListing {
				whereClauses = append(whereClauses, `EXISTS (
					SELECT 1
					FROM property_relation pr_listing
					JOIN posts post_listing ON post_listing.product_id = pr_listing.relation_id AND post_listing.deleted_at IS NULL
					WHERE pr_listing.property_id = pl.id
					  AND pr_listing.relation_type = ?
					  AND pr_listing.deleted_at IS NULL
				)`)
			} else {
				whereClauses = append(whereClauses, `NOT EXISTS (
					SELECT 1
					FROM property_relation pr_listing
					JOIN posts post_listing ON post_listing.product_id = pr_listing.relation_id AND post_listing.deleted_at IS NULL
					WHERE pr_listing.property_id = pl.id
					  AND pr_listing.relation_type = ?
					  AND pr_listing.deleted_at IS NULL
				)`)
			}
			args = append(args, enums.ERelationTypeProduct)
		}
		// In the Search method, update the RecordStatus filter section:
		if req.RecordStatus != "" {
			switch strings.ToLower(req.RecordStatus) {
			case "archived":
				whereClauses = append(whereClauses, "pu.archived_at IS NOT NULL")
			case "inactive", "hidden", "hide":
				whereClauses = append(whereClauses, "pu.archived_at IS NULL AND pu.hidden_at IS NOT NULL")
			case "active":
				whereClauses = append(whereClauses, "pu.archived_at IS NULL AND pu.hidden_at IS NULL")
			}
		}
		if req.RequestUserId != nil {
			whereClauses = append(whereClauses, `EXISTS (
				SELECT 1 FROM property_user pu_filter
				WHERE pu_filter.property_lineage_id = pl.id
				  AND pu_filter.owner_origin_id = ?
				  AND pu_filter.deleted_at IS NULL
			)`)
			args = append(args, *req.RequestUserId)
		}
	}

	fromClause := `
FROM property_identify pi
INNER JOIN property_lineage pl ON pl.property_identify_id = pi.id AND pl.deleted_at IS NULL
LEFT JOIN property_info pif ON pl.property_info_id = pif.id AND pif.deleted_at IS NULL
LEFT JOIN property_location loc ON pl.location_id = loc.id
LEFT JOIN property_land_info li ON pl.land_info_id = li.id AND li.deleted_at IS NULL
LEFT JOIN property_user pu ON pu.property_lineage_id = pl.id AND pu.deleted_at IS NULL AND pu.owner_origin_id = ?
LEFT JOIN province_v2 prov ON loc.province_id = prov.id AND prov.deleted_at IS NULL
LEFT JOIN district_v2 dist ON loc.district_id = dist.id AND dist.deleted_at IS NULL
LEFT JOIN ward_v2 ward ON loc.ward_id = ward.id AND ward.deleted_at IS NULL
LEFT JOIN projects proj ON pif.project_id = proj.id AND proj.deleted_at IS NULL
LEFT JOIN property_type pt ON pif.property_type_id = pt.id AND pt.deleted_at IS NULL
LEFT JOIN LATERAL (
	SELECT pm.id, pm.media_url, pm.thumb_url, pm.media_type
	FROM property_media pm
	WHERE (pm.id = pif.avatar_id OR (pif.avatar_id IS NULL AND pm.lineage_id = pl.id))
	  AND pm.deleted_at IS NULL
	ORDER BY CASE WHEN pm.id = pif.avatar_id THEN 0 ELSE 1 END, pm.sort_order
	LIMIT 1
) avatar ON true`

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "\nWHERE " + strings.Join(whereClauses, "\n  AND ") +
			" AND pi.deleted_at IS NULL"
	}

	countQ := "SELECT COUNT(*) " + fromClause + whereSQL

	var total int64
	if err := r.DB.WithContext(ctx).Raw(countQ, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := "ORDER BY COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST"
	if req != nil {
		switch req.Sort {
		case "title_asc":
			orderBy = "ORDER BY pif.title ASC NULLS LAST, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST"
		case "title_desc":
			orderBy = "ORDER BY pif.title DESC NULLS LAST, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST"
		case "updated_at_asc":
			orderBy = "ORDER BY COALESCE(pl.updated_at, pl.created_at) ASC NULLS LAST"
		case "updated_at_desc":
			orderBy = "ORDER BY COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST"
		case "priority_product":
			orderBy = `ORDER BY (
				SELECT COUNT(*)
				FROM property_relation pr_product
				WHERE pr_product.property_id = pl.id
				  AND pr_product.relation_type = 10
				  AND pr_product.deleted_at IS NULL
			) DESC, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST`
		case "priority_asset":
			orderBy = `ORDER BY (
				SELECT COUNT(*)
				FROM property_relation pr_asset
				WHERE pr_asset.property_id = pl.id
				  AND pr_asset.relation_type = 20
				  AND pr_asset.deleted_at IS NULL
			) DESC, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST`
		case "priority_listing":
			orderBy = `ORDER BY (
				SELECT COUNT(*)
				FROM property_relation pr_listing
				JOIN posts post_listing ON post_listing.product_id = pr_listing.relation_id AND post_listing.deleted_at IS NULL
				WHERE pr_listing.property_id = pl.id
				  AND pr_listing.relation_type = 10
				  AND pr_listing.deleted_at IS NULL
			) DESC, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST`
		case "verified_first":
			orderBy = "ORDER BY CASE WHEN pl.verified_national_at IS NOT NULL THEN 0 ELSE 1 END, COALESCE(pl.updated_at, pl.created_at) DESC NULLS LAST"
		}
	}

	listQ := `
SELECT
	pl.id,
	pl.property_identify_id,
	pi.id AS pi_id,
	pi.p_id AS pi_pid,
	pi.version AS pi_version,
	pi.lineage_id AS lineage_id,
	pl.created_at,
	pl.updated_at,
	pl.national_id,
	pl.verified_national_at,
	pif.id AS info_id,
	pif.title,
	pif.note,
	pif.unit_code,
	pif.identifier,
	pif.level,
	pif.source_type,
	pif.legal_status,
	pif.project_id,
	pif.scope,
	pt.id AS property_type_id,
	pt.name AS property_type_name,
	proj.name AS project_name,
	loc.id AS loc_id,
	loc.address_detail,
	loc.region_id,
	loc.province_id,
	loc.district_id,
	loc.ward_id,
	loc.latitude,
	loc.longitude,
	loc.map_url,
	prov.name AS province_name,
	dist.name AS district_name,
	ward.name AS ward_name,
	li.id AS land_id,
	li.document_type,
	li.document_no,
	li.issuring_auth,
	li.area_total AS land_area_total,
	li.plot,
	li.sheet,
	li.note AS land_note,
	li.land_note AS land_note2,
	avatar.id AS avatar_id,
	avatar.media_url AS avatar_url,
	avatar.thumb_url AS avatar_thumb_url,
	avatar.media_type AS avatar_media_type,
	pu.id AS personalization_id,
	pu.owner_origin_id AS personalization_owner_origin_id,
	pu.property_lineage_id AS personalization_property_lineage_id,
	pu.ownered_at AS personalization_owner_at,
	pu.role_id AS personalization_role_id,
	pu.origin_id AS personalization_origin_profile_id,
	pu.archived_at AS personalization_archived_at,
	pu.hidden_at AS personalization_hidden_at
` + fromClause + whereSQL + "\n" + orderBy + "\nLIMIT ? OFFSET ?"

	var rows []struct {
		ID                               uint64
		PropertyIdentifyID               *uint64
		PiID                             *uint64
		PiPid                            *uint64
		PiVersion                        *uint32
		LineageID                        *uint64
		CreatedAt                        time.Time
		UpdatedAt                        *time.Time
		NationalID                       string
		VerifiedNationalAt               *time.Time
		InfoID                           *uint64
		Title                            *string
		Note                             *string
		UnitCode                         *string
		Identifier                       *string
		Level                            *string
		SourceType                       *uint32
		LegalStatus                      *uint32
		ProjectID                        *uint64
		Scope                            *int
		PropertyTypeID                   *uint64
		PropertyTypeName                 *string
		ProjectName                      *string
		LocID                            *uint64
		AddressDetail                    *string
		RegionID                         *uint64
		ProvinceID                       *uint64
		DistrictID                       *uint64
		WardID                           *uint64
		Latitude                         *float64
		Longitude                        *float64
		MapURL                           *string
		ProvinceName                     *string
		DistrictName                     *string
		WardName                         *string
		LandID                           *uint64
		DocumentType                     *uint32
		DocumentNo                       *string
		IssuringAuth                     *string
		LandAreaTotal                    *float64
		Plot                             uint32
		Sheet                            uint32
		LandNote                         *string
		LandNote2                        *string
		AvatarID                         *uint64
		AvatarURL                        *string
		AvatarThumbURL                   *string
		AvatarMediaType                  *string
		PersonalizationID                *uint64
		PersonalizationOwnerOriginID     *uint64
		PersonalizationPropertyLineageID *uint64
		PersonalizationOwnerAt           *time.Time
		PersonalizationRoleID            *uint64
		PersonalizationOriginProfileID   *uint64
		PersonalizationArchivedAt        *time.Time
		PersonalizationHiddenAt          *time.Time
	}

	offset := 0
	limit := 20
	if req != nil {
		offset = req.GetOffset()
		limit = req.GetLimit()
	}
	queryArgs := append(append([]any{}, args...), limit, offset)
	if err := r.DB.WithContext(ctx).Raw(listQ, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	properties := make([]*domain.PropertyLineage, 0, len(rows))
	for _, row := range rows {
		p := &domain.PropertyLineage{
			BaseEntity: _domain.BaseEntity{ID: row.ID},
		}
		p.PropertyIdentifyID = row.PropertyIdentifyID
		if row.PiID != nil {
			pi := &domain.PropertyIdentify{
				BaseEntity: _domain.BaseEntity{ID: *row.PiID},
			}
			if row.PiPid != nil {
				pi.PID = *row.PiPid
			}
			if row.PiVersion != nil {
				pi.Version = *row.PiVersion
			}
			if row.LineageID != nil {
				pi.LineageID = row.LineageID
			}
			p.PropertyIdentify = pi
		}
		p.CreatedAt = &row.CreatedAt
		p.UpdatedAt = row.UpdatedAt
		p.NationalID = row.NationalID
		p.VerifiedNationalAt = row.VerifiedNationalAt

		if row.InfoID != nil {
			info := &domain.PropertyInfo{
				BaseEntity: _domain.BaseEntity{ID: *row.InfoID},
				Title:      ptrStr(row.Title),
				Note:       ptrStr(row.Note),
				UnitCode:   ptrStr(row.UnitCode),
				Identifier: ptrStr(row.Identifier),
				Level:      ptrStr(row.Level),
				ProjectID:  row.ProjectID,
			}
			if row.SourceType != nil {
				info.SourceType = enums.EPropertySourceType(*row.SourceType)
			}
			if row.LegalStatus != nil {
				info.LegalStatus = enums.EHouseCertificate(*row.LegalStatus)
			}
			if row.Scope != nil {
				info.Scope = enums.EPropertyScope(*row.Scope)
			}
			if row.PropertyTypeID != nil {
				info.PropertyTypeID = row.PropertyTypeID
				info.PropertyType = &domain.PropertyType{
					BaseEntity: _models.BaseEntity{ID: *row.PropertyTypeID},
					Name:       ptrStr(row.PropertyTypeName),
				}
			}
			if row.ProjectID != nil {
				info.Project = &domain.Project{
					BaseEntity: _models.BaseEntity{ID: *row.ProjectID},
					Name:       ptrStr(row.ProjectName),
				}
			}
			if row.AvatarURL != nil && *row.AvatarURL != "" {
				mt := "image"
				if row.AvatarMediaType != nil {
					mt = *row.AvatarMediaType
				}
				avatar := &domain.PropertyMedia{
					MediaURL:  *row.AvatarURL,
					ThumbURL:  ptrStr(row.AvatarThumbURL),
					MediaType: mt,
				}
				if row.AvatarID != nil {
					avatar.BaseEntity = _domain.BaseEntity{ID: *row.AvatarID}
				}
				info.Avatar = avatar
			}
			p.PropertyInfo = info
		}

		if row.LocID != nil {
			loc := &domain.PropertyLocation{
				BaseEntity:    _models.BaseEntity{ID: *row.LocID},
				AddressDetail: ptrStr(row.AddressDetail),
				RegionID:      row.RegionID,
				ProvinceID:    row.ProvinceID,
				// DistrictID:    row.DistrictID,
				WardID:    row.WardID,
				Latitude:  row.Latitude,
				Longitude: row.Longitude,
				MapURL:    ptrStr(row.MapURL),
			}
			if row.ProvinceName != nil {
				loc.Province = &domain.ProvinceV2{Name: *row.ProvinceName}
			}
			// if row.DistrictName != nil {
			// 	loc.District = &domain.DistrictV2{Name: *row.DistrictName}
			// }
			if row.WardName != nil {
				loc.Ward = &domain.WardV2{Name: *row.WardName}
			}
			p.Location = loc
		}

		if row.LandID != nil {
			land := &domain.PropertyLandInfo{
				BaseEntity:   _models.BaseEntity{ID: *row.LandID},
				DocumentNo:   ptrStr(row.DocumentNo),
				IssuringAuth: ptrStr(row.IssuringAuth),
				AreaTotal:    row.LandAreaTotal,
				Plot:         row.Plot,
				Sheet:        row.Sheet,
				Note:         ptrStr(row.LandNote),
				LandNote:     ptrStr(row.LandNote2),
			}
			if row.DocumentType != nil {
				land.DocumentType = enums.EDocUnknown.Parse(row.DocumentType)
			}
			p.LandInfo = land
		}

		if row.PersonalizationID != nil {
			p.Personalization = &domain.PropertyUser{
				BaseEntity:        _models.BaseEntity{ID: *row.PersonalizationID},
				PropertyLineageID: row.PersonalizationPropertyLineageID,
				RoleID:            row.PersonalizationRoleID,
				OriginProfileID:   row.PersonalizationOriginProfileID,
				ArchivedAt:        row.PersonalizationArchivedAt,
				HiddenAt:          row.PersonalizationHiddenAt,
			}
			if row.PersonalizationOwnerOriginID != nil {
				p.Personalization.OwnerOriginID = *row.PersonalizationOwnerOriginID
			}
			if row.PersonalizationOwnerAt != nil {
				p.Personalization.OwnerAt = *row.PersonalizationOwnerAt
			}
		}

		properties = append(properties, p)
	}

	return properties, total, nil
}

// ListMe dùng raw SQL lấy danh sách property của user, chỉ trả title, address, project, buildingInfo, avatar
func (r *PostgrePropertyRepo) ListMe(ctx context.Context, originId uint64, pagable _dto.Pagable) ([]*domain.PropertyLineage, int64, error) {
	const countQ = `
SELECT COUNT(*)
FROM property_user pu
INNER JOIN property_lineage pl ON pl.id = pu.property_lineage_id AND pl.deleted_at IS NULL
WHERE pu.owner_origin_id = ? AND pl.deleted_at IS NULL`

	var total int64
	if err := r.DB.WithContext(ctx).Raw(countQ, originId).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	const listQ = `
SELECT
	pl.id,
	pi.id AS pi_id,
	pi.p_id AS pi_pid,
	pi.version AS pi_version,
	pi.lineage_id AS pi_lineage_id,
	pif.title,
	pif.identifier,
	loc.address_detail,
	loc.province_id, loc.district_id, loc.ward_id,
	prov.name AS province_name,
	dist.name AS district_name,
	ward.name AS ward_name,
	pif.project_id AS project_id,
	proj.name AS project_name,
	dev.id AS developer_id, dev.name AS developer_name, dev.slug AS developer_slug, dev.logo_url AS developer_logo_url,
	dev.phone AS developer_phone, dev.email AS developer_email, dev.website AS developer_website,
	bi.id AS bi_id, bi.area_actual, bi.area_floor, bi.area_construction,
	li.area_total AS land_area_total,
	bi.floors, bi.room_number, bi.bedrooms, bi.bathrooms,
	bi.build_status, bi.building_type, bi.direction, bi.balcony_direction, bi.note AS bi_note,
	avatar.id AS avatar_id,
	avatar.media_url AS avatar_url,
	avatar.thumb_url AS avatar_thumb_url,
	avatar.media_type AS avatar_media_type,
	pu.id AS personalization_id,
	pu.owner_origin_id AS personalization_owner_origin_id,
	pu.property_lineage_id AS personalization_property_lineage_id,
	pu.ownered_at AS personalization_owner_at,
	pu.role_id AS personalization_role_id,
	pu.origin_id AS personalization_origin_profile_id,
	pu.archived_at AS personalization_archived_at,
	pu.hidden_at AS personalization_hidden_at,
	ps.id AS statistic_id, ps.product_count, ps.asset_count, ps.listing_count,
	pl.updated_at, pt.id AS property_type_id, pt.name AS property_type_name
FROM property_user pu
INNER JOIN property_lineage pl ON pl.id = pu.property_lineage_id AND pl.deleted_at IS NULL
LEFT JOIN property_info pif ON pl.property_info_id = pif.id AND pif.deleted_at IS NULL
LEFT JOIN property_identify pi ON pl.property_identify_id = pi.id AND pi.deleted_at IS NULL
LEFT JOIN property_location loc ON pl.location_id = loc.id
LEFT JOIN province_v2 prov ON loc.province_id = prov.id AND prov.deleted_at IS NULL
LEFT JOIN district_v2 dist ON loc.district_id = dist.id AND dist.deleted_at IS NULL
LEFT JOIN ward_v2 ward ON loc.ward_id = ward.id AND ward.deleted_at IS NULL
LEFT JOIN projects proj ON pif.project_id = proj.id AND proj.deleted_at IS NULL
LEFT JOIN developer dev ON proj.developer_id = dev.id AND dev.deleted_at IS NULL
LEFT JOIN property_building_info bi ON pl.building_info_id = bi.id AND bi.deleted_at IS NULL
LEFT JOIN property_land_info li ON pl.land_info_id = li.id AND li.deleted_at IS NULL
LEFT JOIN property_statistic ps ON pl.statistic_id = ps.id AND ps.deleted_at IS NULL
LEFT JOIN property_type pt ON pif.property_type_id = pt.id AND pt.deleted_at IS NULL
LEFT JOIN LATERAL (
	SELECT pm.id, pm.media_url, pm.thumb_url, pm.media_type
	FROM property_media pm
	WHERE (pm.id = pif.avatar_id OR (pif.avatar_id IS NULL AND pm.lineage_id = pl.id))
	  AND pm.deleted_at IS NULL
	ORDER BY CASE WHEN pm.id = pif.avatar_id THEN 0 ELSE 1 END, pm.sort_order
	LIMIT 1
) avatar ON true
WHERE pu.owner_origin_id = ? AND pu.deleted_at IS NULL and pl.deleted_at IS NULL
ORDER BY pl.updated_at DESC NULLS LAST`

	limit := pagable.GetLimit()
	offset := pagable.GetOffset()

	var rows []struct {
		ID                               uint64           `gorm:"column:id"`
		PiID                             *uint64          `gorm:"column:pi_id"`
		PiPid                            *uint64          `gorm:"column:pi_pid"`
		PiVersion                        *uint32          `gorm:"column:pi_version"`
		Title                            *string          `gorm:"column:title"`
		Identifier                       *string          `gorm:"column:identifier"`
		LineageID                        *uint64          `gorm:"column:pi_lineage_id"`
		AddressDetail                    *string          `gorm:"column:address_detail"`
		ProvinceID                       *uint64          `gorm:"column:province_id"`
		DistrictID                       *uint64          `gorm:"column:district_id"`
		WardID                           *uint64          `gorm:"column:ward_id"`
		ProvinceName                     *string          `gorm:"column:province_name"`
		DistrictName                     *string          `gorm:"column:district_name"`
		WardName                         *string          `gorm:"column:ward_name"`
		ProjectID                        *uint64          `gorm:"column:project_id"`
		ProjectName                      *string          `gorm:"column:project_name"`
		DeveloperID                      *uint64          `gorm:"column:developer_id"`
		DeveloperName                    *string          `gorm:"column:developer_name"`
		DeveloperSlug                    *string          `gorm:"column:developer_slug"`
		DeveloperLogoURL                 *string          `gorm:"column:developer_logo_url"`
		DeveloperPhone                   *string          `gorm:"column:developer_phone"`
		DeveloperEmail                   *string          `gorm:"column:developer_email"`
		DeveloperWebsite                 *string          `gorm:"column:developer_website"`
		BiID                             *uint64          `gorm:"column:bi_id"`
		AreaActual                       *float64         `gorm:"column:area_actual"`
		AreaFloor                        *float64         `gorm:"column:area_floor"`
		AreaConstruction                 *float64         `gorm:"column:area_construction"`
		Floors                           *uint32          `gorm:"column:floors"`
		RoomNumber                       *uint32          `gorm:"column:room_number"`
		Bedrooms                         *uint32          `gorm:"column:bedrooms"`
		Bathrooms                        *uint32          `gorm:"column:bathrooms"`
		BuildStatus                      *int             `gorm:"column:build_status"`
		BuildingType                     *int             `gorm:"column:building_type"`
		Direction                        *int             `gorm:"column:direction"`
		BalconyDir                       *int             `gorm:"column:balcony_direction"`
		BiNote                           *string          `gorm:"column:bi_note"`
		LandAreaTotal                    *float64         `gorm:"column:land_area_total"`
		AvatarID                         *uint64          `gorm:"column:avatar_id"`
		AvatarURL                        *string          `gorm:"column:avatar_url"`
		AvatarThumbURL                   *string          `gorm:"column:avatar_thumb_url"`
		AvatarMediaType                  *string          `gorm:"column:avatar_media_type"`
		MediaListRaw                     *json.RawMessage `gorm:"column:media_list_raw"`
		PersonalizationID                *uint64          `gorm:"column:personalization_id"`
		PersonalizationOwnerOriginID     *uint64          `gorm:"column:personalization_owner_origin_id"`
		PersonalizationPropertyLineageID *uint64          `gorm:"column:personalization_property_lineage_id"`
		PersonalizationOwnerAt           *time.Time       `gorm:"column:personalization_owner_at"`
		PersonalizationRoleID            *uint64          `gorm:"column:personalization_role_id"`
		PersonalizationOriginProfileID   *uint64          `gorm:"column:personalization_origin_profile_id"`
		PersonalizationArchivedAt        *time.Time       `gorm:"column:personalization_archived_at"`
		PersonalizationHiddenAt          *time.Time       `gorm:"column:personalization_hidden_at"`
		StatisticID                      *uint64          `gorm:"column:statistic_id"`
		ProductCount                     int64            `gorm:"column:product_count"`
		AssetCount                       int64            `gorm:"column:asset_count"`
		ListingCount                     int64            `gorm:"column:listing_count"`
		UpdatedAt                        *time.Time       `gorm:"column:updated_at"`
		PropertyTypeID                   *uint64          `gorm:"column:property_type_id"`
		PropertyTypeName                 *string          `gorm:"column:property_type_name"`
	}

	if err := r.DB.WithContext(ctx).Raw(listQ+" LIMIT ? OFFSET ?", originId, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	list := make([]*domain.PropertyLineage, 0, len(rows))
	for _, row := range rows {
		p := &domain.PropertyLineage{
			BaseEntity: _domain.BaseEntity{ID: row.ID},
		}
		if row.PiID != nil {
			pi := &domain.PropertyIdentify{
				BaseEntity: _domain.BaseEntity{ID: *row.PiID},
			}
			if row.PiPid != nil {
				pi.PID = *row.PiPid
			}
			if row.PiVersion != nil {
				pi.Version = *row.PiVersion
			}
			if row.LineageID != nil {
				pi.LineageID = row.LineageID
			}
			p.PropertyIdentify = pi
			p.PropertyIdentifyID = &pi.ID
		}
		if row.UpdatedAt != nil {
			p.UpdatedAt = row.UpdatedAt
		}
		if row.PersonalizationID != nil {
			p.Personalization = &domain.PropertyUser{
				BaseEntity:        _models.BaseEntity{ID: *row.PersonalizationID},
				PropertyLineageID: row.PersonalizationPropertyLineageID,
				RoleID:            row.PersonalizationRoleID,
				OriginProfileID:   row.PersonalizationOriginProfileID,
			}
			if row.PersonalizationOwnerOriginID != nil {
				p.Personalization.OwnerOriginID = *row.PersonalizationOwnerOriginID
			}
			if row.PersonalizationOwnerAt != nil {
				p.Personalization.OwnerAt = *row.PersonalizationOwnerAt
			}
			p.Personalization.ArchivedAt = row.PersonalizationArchivedAt
			p.Personalization.HiddenAt = row.PersonalizationHiddenAt
		}

		if row.StatisticID != nil {
			p.Statistic = &domain.PropertyStatistic{
				BaseEntity:   _domain.BaseEntity{ID: *row.StatisticID},
				ProductCount: row.ProductCount,
				AssetCount:   row.AssetCount,
				ListingCount: row.ListingCount,
			}
		}

		if row.Title != nil || row.ProjectID != nil {
			info := &domain.PropertyInfo{
				// BaseEntity: _domain.BaseEntity{ID: *row.I},
				Title:      ptrStr(row.Title),
				Identifier: ptrStr(row.Identifier),
				AvatarID:   row.AvatarID,
			}
			if row.ProjectID != nil {
				info.ProjectID = row.ProjectID
				info.Project = &domain.Project{
					BaseEntity: _models.BaseEntity{ID: *row.ProjectID},
					Name:       ptrStr(row.ProjectName),
				}
			}

			if row.AvatarURL != nil && *row.AvatarURL != "" {
				info.Avatar = &domain.PropertyMedia{
					BaseEntity: _domain.BaseEntity{ID: *row.AvatarID},
					MediaURL:   *row.AvatarURL,
					ThumbURL:   ptrStr(row.AvatarThumbURL),
					MediaType:  *row.AvatarMediaType,
				}
			}
			if row.PropertyTypeID != nil {
				info.PropertyType = &domain.PropertyType{
					BaseEntity: _models.BaseEntity{ID: *row.PropertyTypeID},
					Name:       ptrStr(row.PropertyTypeName),
				}
			}
			p.PropertyInfo = info
		}
		if row.DeveloperID != nil && *row.DeveloperID > 0 {
			p.Developer = &domain.Developer{
				BaseEntity: _models.BaseEntity{ID: *row.DeveloperID},
				Name:       ptrStr(row.DeveloperName),
				Slug:       ptrStr(row.DeveloperSlug),
				LogoUrl:    ptrStr(row.DeveloperLogoURL),
				Phone:      ptrStr(row.DeveloperPhone),
				Email:      ptrStr(row.DeveloperEmail),
				Website:    ptrStr(row.DeveloperWebsite),
			}
		}

		if row.AddressDetail != nil || row.ProvinceID != nil || row.ProvinceName != nil || row.DistrictName != nil || row.WardName != nil {
			loc := &domain.PropertyLocation{}
			if row.AddressDetail != nil {
				loc.AddressDetail = *row.AddressDetail
			}
			loc.ProvinceID = row.ProvinceID
			// loc.DistrictID = row.DistrictID
			loc.WardID = row.WardID
			if row.ProvinceName != nil {
				loc.Province = &domain.ProvinceV2{Name: *row.ProvinceName}
			}
			// if row.DistrictName != nil {
			// 	loc.District = &domain.DistrictV2{Name: *row.DistrictName}
			// }
			if row.WardName != nil {
				loc.Ward = &domain.WardV2{Name: *row.WardName}
			}
			p.Location = loc
		}

		if row.LandAreaTotal != nil {
			p.LandInfo = &domain.PropertyLandInfo{AreaTotal: row.LandAreaTotal}
		}

		if row.BiID != nil || row.AreaActual != nil || row.Floors != nil || row.Bedrooms != nil {
			bi := &domain.PropertyBuildingInfo{}
			bi.AreaActual = row.AreaActual
			bi.AreaFloor = row.AreaFloor
			bi.AreaConstruction = row.AreaConstruction
			bi.Floors = row.Floors
			bi.RoomNumber = row.RoomNumber
			bi.Bedrooms = row.Bedrooms
			bi.Bathrooms = row.Bathrooms
			bi.Note = ptrStr(row.BiNote)
			if row.BuildStatus != nil {
				bi.BuildStatus = enums.EBuildStatus(*row.BuildStatus)
			}
			if row.BuildingType != nil {
				bi.BuildingType = enums.EBuildingType(*row.BuildingType)
			}
			if row.Direction != nil {
				bi.Direction = enums.EHouseOrient(*row.Direction)
			}
			if row.BalconyDir != nil {
				bi.BalconyDirection = enums.EHouseOrient(*row.BalconyDir)
			}
			p.BuildingInfo = bi
		}

		// MediaList: ưu tiên danh sách đầy đủ, không có thì dùng avatar
		if row.MediaListRaw != nil {
			var mediaList []struct {
				ID        uint64 `json:"id"`
				MediaURL  string `json:"mediaUrl"`
				ThumbURL  string `json:"thumbUrl"`
				MediaType string `json:"mediaType"`
				SortOrder int32  `json:"sortOrder"`
			}
			if err := json.Unmarshal(*row.MediaListRaw, &mediaList); err == nil && len(mediaList) > 0 {
				p.MediaList = make([]domain.PropertyMedia, len(mediaList))
				for i, m := range mediaList {
					p.MediaList[i] = domain.PropertyMedia{
						BaseEntity: _domain.BaseEntity{ID: m.ID},
						MediaURL:   m.MediaURL,
						ThumbURL:   m.ThumbURL,
						MediaType:  m.MediaType,
						SortOrder:  m.SortOrder,
					}
				}
			}
		}
		if len(p.MediaList) == 0 && row.AvatarURL != nil && *row.AvatarURL != "" {
			mt := "image"
			if row.AvatarMediaType != nil {
				mt = *row.AvatarMediaType
			}
			thumb := ""
			if row.AvatarThumbURL != nil {
				thumb = *row.AvatarThumbURL
			}
			p.MediaList = []domain.PropertyMedia{{
				MediaURL:  *row.AvatarURL,
				ThumbURL:  thumb,
				MediaType: mt,
			}}
		}

		list = append(list, p)
	}

	return list, total, nil
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
	// r.DB.WithContext(ctx).Raw(
	// 	"SELECT COUNT(*) FROM property_relation WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL",
	// 	id, enums.ERelationTypeAsset,
	// ).Scan(&property.AssetCount)
	// r.DB.WithContext(ctx).Raw(
	// 	"SELECT COUNT(*) FROM property_relation WHERE property_id = ? AND relation_type = ? AND deleted_at IS NULL",
	// 	id, enums.ERelationTypeProduct,
	// ).Scan(&property.ProductCount)
	// r.DB.WithContext(ctx).Raw(
	// 	"SELECT COUNT(*) FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id WHERE pr.property_id = ? AND pr.relation_type = ? AND pr.deleted_at IS NULL AND posts.deleted_at IS NULL",
	// 	id, enums.ERelationTypeProduct,
	// ).Scan(&property.ListingCount)

	return &property, nil
}

func (r *PostgrePropertyRepo) GetDetail2(
	ctx context.Context,
	id uint64,
	sourceType uint32,
) (*domain.PropertyLineage, error) {
	originId := _utils.GetOriginIdFromContext(ctx)
	const q = `
SELECT
	pl.id, pl.property_identify_id, pl.location_id, pl.property_info_id,
	pl.land_info_id, pl.edvidence_id, pl.building_info_id, pl.external_ref_id,
	pl.national_id, pl.verified_national_at, pl.created_at, pl.updated_at,
	pl.deleted_at, pl.created_by, pl.updated_by,
	pu.id AS personalization_id, pu.owner_origin_id AS personalization_owner_origin_id,
	pu.property_lineage_id AS personalization_property_lineage_id, pu.ownered_at AS personalization_owner_at,
	pu.role_id AS personalization_role_id, pu.origin_id AS personalization_origin_profile_id,
	pu.archived_at AS personalization_archived_at, pu.hidden_at AS personalization_hidden_at,
	pt.id AS property_type_id, pt.name AS property_type_name,
	pif.id AS info_id, pif.title, pif.unit_code, pif.identifier, pif.level, pif.note, pif.source_type, pif.record_status, pif.visibility, pif.legal_status, pif.privacy_level, pif.project_id,
	pu.owner_origin_id,
	pj.id AS project_id, pj.name AS project_name,
	dev.id AS developer_id, dev.name AS developer_name, dev.slug AS developer_slug, dev.logo_url AS developer_logo_url,
	dev.phone AS developer_phone, dev.email AS developer_email, dev.website AS developer_website,
	avatar.id AS avatar_id, avatar.media_url AS avatar_url, avatar.thumb_url AS avatar_thumb_url, avatar.media_type AS avatar_media_type,
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
	ps.id AS statistic_id, ps.product_count, ps.asset_count, ps.listing_count,
	(SELECT COALESCE(
			json_agg(
				json_build_object(
					'id', pm.id,
					'mediaUrl', pm.media_url,
					'mediaType', pm.media_type,
					'sortOrder', pm.sort_order
				)
				ORDER BY pm.sort_order ASC
			), '[]')
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
	(SELECT COALESCE(json_agg(json_build_object(
		'id', ar.id,
		'name', ar.name
	)), '[]')
	FROM property_area_region par
	LEFT JOIN area_regions ar ON par.area_region_id = ar.id AND ar.deleted_at IS NULL
	WHERE par.property_lineage_id = pl.id
	AND ar.active = true
	) AS area_regions_raw,
	(SELECT COUNT(*) FROM property_relation WHERE property_id = pl.id AND relation_type = ? AND deleted_at IS NULL) AS asset_count,
	(SELECT COUNT(*) FROM property_relation WHERE property_id = pl.id AND relation_type = ? AND deleted_at IS NULL) AS product_count,
	(SELECT COUNT(*) FROM property_relation pr JOIN posts ON pr.relation_id = posts.product_id
	 WHERE pr.property_id = pl.id AND pr.relation_type = ? AND pr.deleted_at IS NULL AND posts.deleted_at IS NULL) AS listing_count
FROM property_lineage pl
LEFT JOIN property_identify pi ON pl.property_identify_id = pi.id AND pi.deleted_at IS NULL
LEFT JOIN property_info pif ON pl.property_info_id = pif.id AND pif.deleted_at IS NULL
LEFT JOIN property_location loc ON pl.location_id = loc.id
LEFT JOIN property_type pt ON pif.property_type_id = pt.id AND pt.deleted_at IS NULL
LEFT JOIN province_v2 prov ON loc.province_id = prov.id AND prov.deleted_at IS NULL
LEFT JOIN district_v2 dist ON loc.district_id = dist.id AND dist.deleted_at IS NULL
LEFT JOIN ward_v2 ward ON loc.ward_id = ward.id AND ward.deleted_at IS NULL
LEFT JOIN property_land_info li ON pl.land_info_id = li.id AND li.deleted_at IS NULL
LEFT JOIN property_building_info bi ON pl.building_info_id = bi.id AND bi.deleted_at IS NULL
LEFT JOIN property_edvidence ev ON pl.edvidence_id = ev.id AND ev.deleted_at IS NULL
LEFT JOIN projects pj ON pif.project_id = pj.id AND pj.deleted_at IS NULL
LEFT JOIN developer dev ON pj.developer_id = dev.id AND dev.deleted_at IS NULL
LEFT JOIN property_statistic ps ON pl.statistic_id = ps.id AND ps.deleted_at IS NULL
LEFT JOIN property_user pu ON pu.property_lineage_id = pl.id AND pu.deleted_at IS NULL AND pu.owner_origin_id = ?
LEFT JOIN LATERAL (
	SELECT pm.id, pm.media_url, pm.thumb_url, pm.media_type
	FROM property_media pm
	WHERE (pm.id = pif.avatar_id OR (pif.avatar_id IS NULL AND pm.lineage_id = pl.id))
	  AND pm.deleted_at IS NULL
	ORDER BY CASE WHEN pm.id = pif.avatar_id THEN 0 ELSE 1 END, pm.sort_order
	LIMIT 1
) avatar ON true
WHERE pl.id = ? AND pl.deleted_at IS NULL`
	var row struct {
		ID                               uint64
		PropertyIdentifyID               *uint64
		LocationID                       *uint64
		PropertyInfoID                   *uint64
		LandInfoID                       *uint64
		EdvidenceID                      *uint64
		BuildingInfoID                   *uint64
		ExternalRefID                    *uint64
		NationalID                       string
		VerifiedNationalAt               *time.Time
		CreatedAt                        time.Time
		UpdatedAt                        *time.Time
		DeletedAt                        *time.Time
		CreatedBy                        *uint64
		UpdatedBy                        *uint64
		PersonalizationID                *uint64
		PersonalizationOwnerOriginID     *uint64
		PersonalizationPropertyLineageID *uint64
		PersonalizationOwnerAt           *time.Time
		PersonalizationRoleID            *uint64
		PersonalizationOriginProfileID   *uint64
		PersonalizationArchivedAt        *time.Time
		PersonalizationHiddenAt          *time.Time
		PropertyTypeID                   *uint64
		PropertyTypeName                 *string
		// Media
		MediaListRaw   *json.RawMessage
		AmenitiesRaw   *json.RawMessage
		AreaRegionsRaw *json.RawMessage
		// PropertyInfo
		InfoID           *uint64
		Title            *string
		UnitCode         *string
		Identifier       *string
		Level            *string
		Note             *string
		SourceType       *uint32
		RecordStatus     *uint32
		Visibility       *uint32
		LegalStatus      *uint32
		PrivacyLevel     *uint32
		OriginProfileId  *uint64
		ProjectID        *uint64
		ProjectName      *string
		DeveloperID      *uint64
		DeveloperName    *string
		DeveloperSlug    *string
		DeveloperLogoURL *string
		DeveloperPhone   *string
		DeveloperEmail   *string
		DeveloperWebsite *string
		// Avatar
		AvatarID        *uint64
		AvatarURL       *string
		AvatarThumbURL  *string
		AvatarMediaType *string
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
		ProvinceCode  *string
		DistrictName  *string
		DistrictCode  *int
		WardName      *string
		WardCode      *string
		// Statistic
		StatisticID  *uint64
		ProductCount int64
		AssetCount   int64
		ListingCount int64
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
	}
	if err := r.DB.WithContext(ctx).Raw(q,
		enums.ERelationTypeAsset, enums.ERelationTypeProduct, enums.ERelationTypeProduct, originId, id,
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
	// p.ExternalRefID = row.ExternalRefID
	p.NationalID = row.NationalID
	p.VerifiedNationalAt = row.VerifiedNationalAt
	p.CreatedAt = &row.CreatedAt
	p.UpdatedAt = row.UpdatedAt
	if row.PersonalizationID != nil {
		p.Personalization = &domain.PropertyUser{
			BaseEntity: _models.BaseEntity{ID: *row.PersonalizationID},
			RoleID:     row.PersonalizationRoleID,
			OwnerAt:    time.Time{},
		}
		if row.PersonalizationOwnerOriginID != nil {
			p.Personalization.OwnerOriginID = *row.PersonalizationOwnerOriginID
		}
		p.Personalization.PropertyLineageID = row.PersonalizationPropertyLineageID
		p.Personalization.OriginProfileID = row.PersonalizationOriginProfileID
		if row.PersonalizationOwnerAt != nil {
			p.Personalization.OwnerAt = *row.PersonalizationOwnerAt
		}
		p.Personalization.ArchivedAt = row.PersonalizationArchivedAt
		p.Personalization.HiddenAt = row.PersonalizationHiddenAt
	}
	// p.AssetCount = row.AssetCount
	// p.ProductCount = row.ProductCount
	// p.ListingCount = row.ListingCount
	if row.StatisticID != nil {
		p.Statistic = &domain.PropertyStatistic{
			BaseEntity:   _domain.BaseEntity{ID: *row.StatisticID},
			ProductCount: row.ProductCount,
			AssetCount:   row.AssetCount,
			ListingCount: row.ListingCount,
		}
	}
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
			SourceType: enums.EPropertySourceType(*row.SourceType),
			// RecordStatus: enums.EPropertyStatus(*row.RecordStatus),
			// Visibility:      enums.EVisibility(*row.Visibility),
			LegalStatus: enums.EHouseCertificate(*row.LegalStatus),
			// PrivacyLevel:    enums.EVisibility(*row.PrivacyLevel),
			ProjectID: row.ProjectID,
			Title:     ptrStr(row.Title),
		}
		if row.ProjectID != nil && *row.ProjectID > 0 {
			info.Project = &domain.Project{
				BaseEntity: _models.BaseEntity{ID: *row.ProjectID},
				Name:       ptrStr(row.ProjectName),
			}
		}
		if row.AvatarURL != nil && *row.AvatarURL != "" {
			mt := "image"
			if row.AvatarMediaType != nil {
				mt = *row.AvatarMediaType
			}
			avatar := &domain.PropertyMedia{
				MediaURL:  *row.AvatarURL,
				ThumbURL:  ptrStr(row.AvatarThumbURL),
				MediaType: mt,
			}
			if row.AvatarID != nil {
				avatar.BaseEntity = _domain.BaseEntity{ID: *row.AvatarID}
			}
			info.Avatar = avatar
		}
		if row.PropertyTypeID != nil {
			info.PropertyType = &domain.PropertyType{
				BaseEntity: _models.BaseEntity{ID: *row.PropertyTypeID},
				Name:       ptrStr(row.PropertyTypeName),
			}
		}
		info.ID = *row.InfoID
		p.PropertyInfo = info
	}
	if row.DeveloperID != nil && *row.DeveloperID > 0 {
		p.Developer = &domain.Developer{
			BaseEntity: _models.BaseEntity{ID: *row.DeveloperID},
			Name:       ptrStr(row.DeveloperName),
			Slug:       ptrStr(row.DeveloperSlug),
			LogoUrl:    ptrStr(row.DeveloperLogoURL),
			Phone:      ptrStr(row.DeveloperPhone),
			Email:      ptrStr(row.DeveloperEmail),
			Website:    ptrStr(row.DeveloperWebsite),
		}
	}
	if row.LocID != nil {
		loc := &domain.PropertyLocation{
			BaseEntity: _models.BaseEntity{
				ID: *row.LocID,
			},
			AddressDetail: ptrStr(row.AddressDetail),
			ProvinceID:    row.ProvinceID,
			// DistrictID:    row.DistrictID,
			WardID:   row.WardID,
			Latitude: row.Latitude, Longitude: row.Longitude, MapURL: ptrStr(row.MapURL), RegionID: row.RegionID,
		}
		if row.ProvinceName != nil {
			loc.Province = &domain.ProvinceV2{Name: *row.ProvinceName}
			if row.ProvinceCode != nil {
				loc.Province.Code = *row.ProvinceCode
			}
		}
		// if row.DistrictName != nil {
		// 	loc.District = &domain.DistrictV2{Name: *row.DistrictName}
		// 	if row.DistrictCode != nil {
		// 		loc.District.Code = *row.DistrictCode
		// 	}
		// }
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
			DocumentType: enums.EDocUnknown.Parse(row.DocumentType),
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
	var areaRegions []domain.AreaRegion
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
	if row.AreaRegionsRaw != nil {
		if err := json.Unmarshal(*row.AreaRegionsRaw, &areaRegions); err != nil {
			return nil, err
		}
	}
	p.MediaList = media
	p.Amenities = amenities
	p.AreaRegions = areaRegions
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

func (r *PostgrePropertyRepo) ArchiveBatchProperties(ctx context.Context, ids []uint64, archived bool, ownerOriginId uint64) error {
	if len(ids) == 0 {
		return nil
	}

	db := GetDB(ctx, r.DB)

	if archived {
		return db.
			Model(&domain.PropertyUser{}).
			Where("property_lineage_id IN (?) AND owner_origin_id = ? AND deleted_at IS NULL", ids, ownerOriginId).
			Update("archived_at", time.Now()).Error
	}

	return db.
		Model(&domain.PropertyUser{}).
		Where("property_lineage_id IN (?) AND owner_origin_id = ? AND deleted_at IS NULL", ids, ownerOriginId).
		Update("archived_at", nil).Error
}

func (r *PostgrePropertyRepo) HideBatchProperties(ctx context.Context, ids []uint64, hidden bool, ownerOriginId uint64) error {
	if len(ids) == 0 {
		return nil
	}

	db := GetDB(ctx, r.DB)

	if hidden {
		return db.
			Model(&domain.PropertyUser{}).
			Where("property_lineage_id IN (?) AND owner_origin_id = ? AND deleted_at IS NULL", ids, ownerOriginId).
			Update("hidden_at", time.Now()).Error
	}

	return db.
		Model(&domain.PropertyUser{}).
		Where("property_lineage_id IN (?) AND owner_origin_id = ? AND deleted_at IS NULL", ids, ownerOriginId).
		Update("hidden_at", nil).Error
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
