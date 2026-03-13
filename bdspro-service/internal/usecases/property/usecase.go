package property_usecases

import (
	"common/pkg/fieldmask"
	"common/pkg/patch"
	"context"

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
		RecordStatus:   enums.EPropertyStatus(derefUint32(h.Patch.RecordStatus)),
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
		Uint32("info.legal_status", "legal_status", h.Patch.LegalStatus, 0).
		Uint32("info.record_status", "record_status", h.Patch.RecordStatus, 0).
		ForeignKey("info.avatar_id", "avatar_id", h.Patch.AvatarID, nil).
		ForeignKey("info.property_type_id", "property_type_id", h.Patch.PropertyTypeID, nil).
		ForeignKey("info.project_id", "project_id", h.Patch.ProjectID, nil).
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

type LocationHandler struct {
	Patch *dto.UpdateLocationDTO
	Mask  fieldmask.FieldMask
	Repo  property_repo.PropertyLocationRepository
}

func (h *LocationHandler) Name() string     { return "location" }
func (h *LocationHandler) FKColumn() string { return "location_id" }
func (h *LocationHandler) Present() bool    { return h.Patch != nil }

func (h *LocationHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyLocation{
		AddressDetail: derefStr(h.Patch.AddressDetail),
		MapURL:        derefStr(h.Patch.MapURL),
		Latitude:      h.Patch.Latitude,
		Longitude:     h.Patch.Longitude,
		RegionID:      derefUint64Ptr(h.Patch.RegionID),
		ProvinceID:    derefUint64Ptr(h.Patch.ProvinceID),
		DistrictID:    derefUint64Ptr(h.Patch.DistrictID),
		WardID:        derefUint64Ptr(h.Patch.WardID),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *LocationHandler) Update(ctx context.Context, id uint64) error {
	fields := patch.WithMask(h.Mask).
		String("location.address_detail", "address_detail", h.Patch.AddressDetail, "").
		String("location.map_url", "map_url", h.Patch.MapURL, "").
		NullableFloat("location.latitude", "latitude", h.Patch.Latitude, nil).
		NullableFloat("location.longitude", "longitude", h.Patch.Longitude, nil).
		ForeignKey("location.region_id", "region_id", h.Patch.RegionID, nil).
		ForeignKey("location.province_id", "province_id", h.Patch.ProvinceID, nil).
		ForeignKey("location.district_id", "district_id", h.Patch.DistrictID, nil).
		ForeignKey("location.ward_id", "ward_id", h.Patch.WardID, nil).
		Build()

	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

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
	if h.Patch.Basic != nil {
		e.DocumentNo = derefStr(h.Patch.Basic.DocumentNo)
		e.IssuringAuth = derefStr(h.Patch.Basic.IssuringAuth)
		e.DocumentType = enums.EDocType(derefUint32(h.Patch.Basic.DocumentType))
		e.Plot = derefUint32(h.Patch.Basic.Plot)
		e.Sheet = derefUint32(h.Patch.Basic.Sheet)
	}
	if h.Patch.Area != nil {
		e.AreaTotal = h.Patch.Area.AreaTotal
		e.AreaLand = h.Patch.Area.AreaLand
		e.AreaPlant = h.Patch.Area.AreaPlant
	}
	if h.Patch.Dimension != nil {
		e.FrontWidth = h.Patch.Dimension.FrontWidth
		e.Depth = h.Patch.Dimension.Depth
		e.StreetWidth = h.Patch.Dimension.StreetWidth
	}
	if h.Patch.Expiry != nil {
		e.ExpiredLand = ParseDate(h.Patch.Expiry.ExpiredLand)
		e.ExpiredPlant = ParseDate(h.Patch.Expiry.ExpiredPlant)
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *LandInfoHandler) Update(ctx context.Context, id uint64) error {
	b := patch.WithMask(h.Mask).
		String("land_info.note", "note", h.Patch.Note, "").
		String("land_info.land_note", "land_note", h.Patch.LandNote, "").
		Int32("land_info.purpose_used", "purpose_used", h.Patch.PurposeUsed, 0)

	if h.Patch.Basic != nil {
		b.String("land_info.basic.document_no", "document_no", h.Patch.Basic.DocumentNo, "").
			String("land_info.basic.issuring_auth", "issuring_auth", h.Patch.Basic.IssuringAuth, "").
			Uint32("land_info.basic.document_type", "document_type", h.Patch.Basic.DocumentType, 0).
			Uint32("land_info.basic.plot", "plot", h.Patch.Basic.Plot, 0).
			Uint32("land_info.basic.sheet", "sheet", h.Patch.Basic.Sheet, 0)
	}
	if h.Patch.Area != nil {
		b.NullableFloat("land_info.area.area_total", "area_total", h.Patch.Area.AreaTotal, nil).
			NullableFloat("land_info.area.area_land", "area_land", h.Patch.Area.AreaLand, nil).
			NullableFloat("land_info.area.area_plant", "area_plant", h.Patch.Area.AreaPlant, nil)
	}
	if h.Patch.Dimension != nil {
		b.NullableFloat("land_info.dimension.front_width", "front_width", h.Patch.Dimension.FrontWidth, nil).
			NullableFloat("land_info.dimension.depth", "depth", h.Patch.Dimension.Depth, nil).
			NullableFloat("land_info.dimension.street_width", "street_width", h.Patch.Dimension.StreetWidth, nil)
	}
	if h.Patch.Expiry != nil {
		b.DateString("land_info.expiry.expired_land", "expired_land", h.Patch.Expiry.ExpiredLand, nil).
			DateString("land_info.expiry.expired_plant", "expired_plant", h.Patch.Expiry.ExpiredPlant, nil)
	}

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

func (h *BuildingHandler) Name() string     { return "property_building_info" }
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
		String("building_info.note", "note", h.Patch.Note, "").
		Uint32("building_info.build_status", "build_status", h.Patch.BuildStatus, 0).
		Uint32("building_info.building_type", "building_type", h.Patch.BuildingType, 0).
		Uint32("building_info.direction", "direction", h.Patch.Direction, 0).
		Uint32("building_info.balcony_direction", "balcony_direction", h.Patch.BalconyDir, 0).
		NullableFloat("building_info.area_actual", "area_actual", h.Patch.AreaActual, nil).
		NullableFloat("building_info.area_floor", "area_floor", h.Patch.AreaFloor, nil).
		NullableFloat("building_info.area_construction", "area_construction", h.Patch.AreaConstruction, nil).
		NullableUint32("building_info.floors", "floors", h.Patch.Floors, nil).
		NullableUint32("building_info.room_number", "room_number", h.Patch.RoomNumber, nil).
		NullableUint32("building_info.bedrooms", "bedrooms", h.Patch.Bedrooms, nil).
		NullableUint32("building_info.bathrooms", "bathrooms", h.Patch.Bathrooms, nil).
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
