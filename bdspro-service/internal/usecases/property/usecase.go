package property_usecases

import (
	"common/pkg/fieldmask"
	"common/pkg/patch"
	"context"
	"fmt"
	"log"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	property_repo "bdspro/internal/repo/property"
)

// ─────────────────────────────────────────────────────────────────────────────
// InfoHandler
// ─────────────────────────────────────────────────────────────────────────────

type InfoHandler struct {
	Patch *dto.UpdateInfoDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyInfoRepository
}

func (h *InfoHandler) Name() string     { return "info" }
func (h *InfoHandler) FKColumn() string { return "property_info_id" }
func (h *InfoHandler) Present() bool    { return h.Patch != nil }

// Insert: tạo record mới khi lineage chưa có info.
// OrElse(zero) vì Insert cần value, không cần phân biệt absent/present.
func (h *InfoHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyInfo{
		Title:          derefStr(h.Patch.Title),
		Note:           derefStr(h.Patch.Note),
		UnitCode:       derefStr(h.Patch.UnitCode),
		Identifier:     derefStr(h.Patch.Identifier),
		Level:          derefStr(h.Patch.Level),
		Scope:          enums.EPropertyScope(derefUint32(h.Patch.Scope)),
		SourceType:     enums.EPropertySourceType(derefUint32(h.Patch.SourceType)),
		LegalStatus:    enums.EHouseCertificate(derefUint32(h.Patch.LegalStatus)),
		AvatarID:       derefUint64Ptr(h.Patch.AvatarID),
		PropertyTypeID: derefUint64Ptr(h.Patch.PropertyTypeID),
		ProjectID:      derefUint64Ptr(h.Patch.ProjectID),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

// Update: partial update — chỉ field present + allowed by mask + diff.
func (h *InfoHandler) Update(ctx context.Context, id uint64) error {
	fields := patch.WithMask(h.Mask).
		String("info.title", "title", h.Patch.Title, "").
		String("info.note", "note", h.Patch.Note, "").
		String("info.unit_code", "unit_code", h.Patch.UnitCode, "").
		String("info.identifier", "identifier", h.Patch.Identifier, "").
		String("info.level", "level", h.Patch.Level, "").
		Uint32("info.scope", "scope", h.Patch.Scope, 0).
		Uint32("info.source_type", "source_type", h.Patch.SourceType, 0).
		Uint32("info.legalStatus", "legal_status", h.Patch.LegalStatus, 0).
		Uint32("info.record_status", "record_status", h.Patch.RecordStatus, 0).
		ForeignKey("info.avatarId", "avatar_id", h.Patch.AvatarID, nil).
		ForeignKey("info.propertyTypeId", "property_type_id", h.Patch.PropertyTypeID, nil).
		ForeignKey("info.projectId", "project_id", h.Patch.ProjectID, nil).
		ForeignKey("info.origin_profile_id", "origin_profile_id", h.Patch.OriginProfileID, nil).
		Build()

	if len(fields) == 0 {
		return nil // không có gì thay đổi → skip write
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// LocationHandler
// ─────────────────────────────────────────────────────────────────────────────

// type LocationHandler struct {
// 	Patch *dto.UpdateLocationDTO
// 	Mask  fieldmask.FieldMask
// 	Repo  property_repo.PropertyLocationRepository
// }

// func (h *LocationHandler) Name() string     { return "location" }
// func (h *LocationHandler) FKColumn() string { return "location_id" }
// func (h *LocationHandler) Present() bool    { return h.Patch != nil }

// func (h *LocationHandler) Insert(ctx context.Context) (uint64, error) {
// 	e := &domain.PropertyLocation{
// 		AddressDetail: derefStr(h.Patch.AddressDetail),
// 		MapURL:        derefStr(h.Patch.MapURL),
// 		Latitude:      h.Patch.Latitude,
// 		Longitude:     h.Patch.Longitude,
// 		RegionID:      derefUint64Ptr(h.Patch.RegionID),
// 		ProvinceID:    derefUint64Ptr(h.Patch.ProvinceID),
// 		WardID:        derefUint64Ptr(h.Patch.WardID),
// 	}
// 	if err := h.Repo.Create(ctx, e); err != nil {
// 		return 0, err
// 	}
// 	return e.ID, nil
// }

// func (h *LocationHandler) Update(ctx context.Context, id uint64) error {
// 	existing, err := h.Repo.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if existing == nil {
// 		return fmt.Errorf("location id=%d not found", id)
// 	}

// 	fields, err := patch.WithMask(h.Mask).
// 		String("location.address_detail", "address_detail", h.Patch.AddressDetail, existing.AddressDetail).
// 		String("location.map_url", "map_url", h.Patch.MapURL, existing.MapURL).
// 		NullableFloat("location.latitude", "latitude", h.Patch.Latitude, existing.Latitude).
// 		NullableFloat("location.longitude", "longitude", h.Patch.Longitude, existing.Longitude).
// 		ForeignKey("location.region_id", "region_id", h.Patch.RegionID, existing.RegionID).
// 		ForeignKey("location.province_id", "province_id", h.Patch.ProvinceID, existing.ProvinceID).
// 		ForeignKey("location.ward_id", "ward_id", h.Patch.WardID, existing.WardID).
// 		Result()

// 	if err != nil {
// 		return err
// 	}

// 	if len(fields) == 0 {
// 		return nil
// 	}

// 	return h.Repo.UpdateFields(ctx, id, fields)
// }

// ─────────────────────────────────────────────────────────────────────────────
// LandInfoHandler
// ─────────────────────────────────────────────────────────────────────────────

type LandInfoHandler struct {
	Patch *dto.UpdateLandInfoDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyLandInfoRepository
}

func (h *LandInfoHandler) Name() string     { return "land_info" }
func (h *LandInfoHandler) FKColumn() string { return "land_info_id" }
func (h *LandInfoHandler) Present() bool    { return h.Patch != nil }

func (h *LandInfoHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyLandInfo{
		Note:     derefStr(h.Patch.Note),
		LandNote: derefStr(h.Patch.LandNote),
		PurposeUsed: enums.EPurposeUsed(func() int32 {
			if h.Patch.PurposeUsed == nil {
				return 0
			}
			return *h.Patch.PurposeUsed
		}()),
	}
	e.DocumentNo = derefStr(h.Patch.DocumentNo)
	e.IssuringAuth = derefStr(h.Patch.IssuringAuth)
	e.DocumentType = enums.EDocType(derefUint32(h.Patch.DocumentType))
	e.Plot = derefUint32(h.Patch.Plot)
	e.Sheet = derefUint32(h.Patch.Sheet)
	e.AreaTotal = h.Patch.AreaTotal
	e.AreaLand = h.Patch.AreaLand
	e.AreaPlant = h.Patch.AreaPlant
	e.FrontWidth = h.Patch.FrontWidth
	e.Depth = h.Patch.Depth
	e.StreetWidth = h.Patch.StreetWidth
	e.ExpiredLand = ParseDate(h.Patch.ExpiredLand)
	e.ExpiredPlant = ParseDate(h.Patch.ExpiredPlant)

	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *LandInfoHandler) Update(ctx context.Context, id uint64) error {
	b := patch.WithMask(h.Mask).
		String("landInfo.note", "note", h.Patch.Note, "").
		String("landInfo.landNote", "land_note", h.Patch.LandNote, "").
		Int32("landInfo.purposeUsed", "purpose_used", h.Patch.PurposeUsed, 0).
		String("landInfo.documentNo", "document_no", h.Patch.DocumentNo, "").
		String("landInfo.issuringAuth", "issuring_auth", h.Patch.IssuringAuth, "").
		Uint32("landInfo.documentType", "document_type", h.Patch.DocumentType, 0).
		Uint32("landInfo.plot", "plot", h.Patch.Plot, 0).
		Uint32("landInfo.sheet", "sheet", h.Patch.Sheet, 0).
		NullableFloat("landInfo.areaTotal", "area_total", h.Patch.AreaTotal, nil).
		NullableFloat("landInfo.areaLand", "area_land", h.Patch.AreaLand, nil).
		NullableFloat("landInfo.areaPlant", "area_plant", h.Patch.AreaPlant, nil).
		NullableFloat("landInfo.frontWidth", "front_width", h.Patch.FrontWidth, nil).
		NullableFloat("landInfo.depth", "depth", h.Patch.Depth, nil).
		NullableFloat("landInfo.streetWidth", "street_width", h.Patch.StreetWidth, nil).
		DateString("landInfo.expiredLand", "expired_land", h.Patch.ExpiredLand, nil).
		DateString("landInfo.expiredPlant", "expired_plant", h.Patch.ExpiredPlant, nil)

	fields, err := b.Result()
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// BuildingHandler
// ─────────────────────────────────────────────────────────────────────────────

type BuildingHandler struct {
	Patch *dto.UpdateBuildingDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyBuildingInfoRepository
}

func (h *BuildingHandler) Name() string     { return "building_info" }
func (h *BuildingHandler) FKColumn() string { return "building_info_id" }
func (h *BuildingHandler) Present() bool    { return h.Patch != nil }

func (h *BuildingHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyBuildingInfo{
		Note:             derefStr(h.Patch.Note),
		BuildStatus:      enums.EBuildStatus(derefUint32(h.Patch.BuildStatus)),
		BuildingType:     enums.EBuildingType(derefUint32(h.Patch.BuildingType)),
		Direction:        enums.EHouseOrient(derefUint32(h.Patch.Direction)),
		BalconyDirection: enums.EHouseOrient(derefUint32(h.Patch.BalconyDir)),
		AreaActual:       h.Patch.AreaActual,
		AreaFloor:        h.Patch.AreaFloor,
		AreaConstruction: h.Patch.AreaConstruction,
		Floors:           h.Patch.Floors,
		RoomNumber:       h.Patch.RoomNumber,
		Bedrooms:         h.Patch.Bedrooms,
		Bathrooms:        h.Patch.Bathrooms,
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *BuildingHandler) Update(ctx context.Context, id uint64) error {
	fields := patch.WithMask(h.Mask).
		String("buildingInfo.note", "note", h.Patch.Note, "").
		Uint32("buildingInfo.buildStatus", "build_status", h.Patch.BuildStatus, 0).
		Uint32("buildingInfo.building_type", "building_type", h.Patch.BuildingType, 0).
		Uint32("buildingInfo.direction", "direction", h.Patch.Direction, 0).
		Uint32("buildingInfo.balconyDirection", "balcony_direction", h.Patch.BalconyDir, 0).
		NullableFloat("buildingInfo.areaActual", "area_actual", h.Patch.AreaActual, nil).
		NullableFloat("buildingInfo.areaFloor", "area_floor", h.Patch.AreaFloor, nil).
		NullableFloat("buildingInfo.areaConstruction", "area_construction", h.Patch.AreaConstruction, nil).
		NullableUint32("buildingInfo.floors", "floors", h.Patch.Floors, nil).
		NullableUint32("buildingInfo.roomNumber", "room_number", h.Patch.RoomNumber, nil).
		NullableUint32("buildingInfo.bedrooms", "bedrooms", h.Patch.Bedrooms, nil).
		NullableUint32("buildingInfo.bathrooms", "bathrooms", h.Patch.Bathrooms, nil).
		Build()

	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// EvidenceHandler
// ─────────────────────────────────────────────────────────────────────────────

type EvidenceHandler struct {
	Patch *dto.UpdateEvidenceDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyEvidenceRepository
}

func (h *EvidenceHandler) Name() string     { return "evidence" }
func (h *EvidenceHandler) FKColumn() string { return "edvidence_id" }
func (h *EvidenceHandler) Present() bool    { return h.Patch != nil }

func (h *EvidenceHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyEdvidence{
		Title:       derefStr(h.Patch.Title),
		Description: derefStr(h.Patch.Description),
		FileID:      derefUint64Ptr(h.Patch.FileID),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *EvidenceHandler) Update(ctx context.Context, id uint64) error {
	fields := patch.WithMask(h.Mask).
		String("evidence.title", "title", h.Patch.Title, "").
		String("evidence.description", "description", h.Patch.Description, "").
		ForeignKey("evidence.file_id", "file_id", h.Patch.FileID, nil).
		Build()

	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

type LocationHandler struct {
	Patch *dto.UpdateLocationDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyLocationRepository
}

func (h *LocationHandler) Name() string     { return "location" }
func (h *LocationHandler) FKColumn() string { return "location_id" }
func (h *LocationHandler) Present() bool    { return h.Patch != nil }

// Insert - Create new location
func (h *LocationHandler) Insert(ctx context.Context) (uint64, error) {
	if h.Patch.Latitude == nil || h.Patch.Longitude == nil {
		return 0, fmt.Errorf("latitude and longitude are required for location")
	}

	location := &domain.PropertyLocation{
		AddressDetail: derefStr(h.Patch.AddressDetail),
		MapURL:        derefStr(h.Patch.MapURL),
		Latitude:      h.Patch.Latitude,
		Longitude:     h.Patch.Longitude,
		PositionType:  h.getPositionType(),
		RegionID:      derefUint64Ptr(h.Patch.RegionID),
		ProvinceID:    h.Patch.ProvinceID,
		WardID:        h.Patch.WardID,
	}

	if err := h.Repo.Create(ctx, location); err != nil {
		return 0, fmt.Errorf("create location: %w", err)
	}

	log.Printf("✓ Created location ID=%d with province_id=%v, ward_id=%v",
		location.ID, safeStr(location.ProvinceID), safeStr(location.WardID))

	return location.ID, nil
}

// Update - Update existing location
func (h *LocationHandler) Update(ctx context.Context, id uint64) error {
	existing, err := h.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("location not found")
	}

	fields := make(map[string]interface{})

	// Chỉ update các field được phép trong mask
	if h.Mask.Allows("location.address_detail") && h.Patch.AddressDetail != nil {
		if *h.Patch.AddressDetail != existing.AddressDetail {
			fields["address_detail"] = *h.Patch.AddressDetail
		}
	}

	if h.Mask.Allows("location.latitude") && h.Patch.Latitude != nil {
		if *h.Patch.Latitude != *existing.Latitude {
			fields["latitude"] = *h.Patch.Latitude
		}
	}

	if h.Mask.Allows("location.longitude") && h.Patch.Longitude != nil {
		if *h.Patch.Longitude != *existing.Longitude {
			fields["longitude"] = *h.Patch.Longitude
		}
	}

	if h.Mask.Allows("location.province_id") && h.Patch.ProvinceID != nil {
		if existing.ProvinceID == nil || *h.Patch.ProvinceID != *existing.ProvinceID {
			fields["province_id"] = *h.Patch.ProvinceID
			log.Printf("  - province_id: %v -> %d",
				safeStr(existing.ProvinceID), *h.Patch.ProvinceID)
		}
	}

	// WardID - tương tự
	if h.Mask.Allows("location.ward_id") && h.Patch.WardID != nil {
		if existing.WardID == nil || *h.Patch.WardID != *existing.WardID {
			fields["ward_id"] = *h.Patch.WardID
			log.Printf("  - ward_id: %v -> %d",
				safeStr(existing.WardID), *h.Patch.WardID)
		}
	}

	if len(fields) == 0 {
		return nil
	}

	return h.Repo.UpdateFields(ctx, id, fields)
}

func (h *LocationHandler) getPositionType() enums.EPositionType {
	if h.Patch.PositionType != nil {
		return *h.Patch.PositionType
	}
	return enums.EPositionType(10)
}

func safeStr(p *uint64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *p)
}
