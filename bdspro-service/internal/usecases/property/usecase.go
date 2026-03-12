package property_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	property_repo "bdspro/internal/repo/property"
	"common/pkg/optional"
	"context"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// InfoHandler
// ─────────────────────────────────────────────────────────────────────────────
type InfoHandler struct {
	Patch *dto.InfoUpdate
	Repo  property_repo.PropertyInfoRepository
}

func (h *InfoHandler) Name() string     { return "info" }
func (h *InfoHandler) FKColumn() string { return "property_info_id" }
func (h *InfoHandler) Present() bool    { return h.Patch != nil }

func (h *InfoHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyInfo{
		Title:          h.Patch.Title.OrElse(""),
		Note:           h.Patch.Note.OrElse(""),
		UnitCode:       h.Patch.UnitCode.OrElse(""),
		Identifier:     h.Patch.Identifier.OrElse(""),
		Level:          h.Patch.Level.OrElse(""),
		Scope:          enums.EPropertyScope(h.Patch.Scope.OrElse(0)),
		SourceType:     enums.EPropertySourceType(h.Patch.SourceType.OrElse(0)),
		LegalStatus:    enums.EHouseCertificate(h.Patch.LegalStatus.OrElse(0)),
		RecordStatus:   enums.EPropertyStatus(h.Patch.RecordStatus.OrElse(0)),
		AvatarID:       h.Patch.AvatarID.Ptr(),
		PropertyTypeID: h.Patch.PropertyTypeID.Ptr(),
		ProjectID:      h.Patch.ProjectID.Ptr(),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *InfoHandler) Update(ctx context.Context, id uint64) error {
	fields := make(map[string]any)
	h.Patch.Title.ApplyTo(fields, "title")
	h.Patch.Note.ApplyTo(fields, "note")
	h.Patch.UnitCode.ApplyTo(fields, "unit_code")
	h.Patch.Identifier.ApplyTo(fields, "identifier")
	h.Patch.Level.ApplyTo(fields, "level")
	h.Patch.Scope.ApplyTo(fields, "scope")
	h.Patch.SourceType.ApplyTo(fields, "source_type")
	h.Patch.LegalStatus.ApplyTo(fields, "legal_status")
	h.Patch.RecordStatus.ApplyTo(fields, "record_status")
	h.Patch.AvatarID.ApplyTo(fields, "avatar_id")
	h.Patch.PropertyTypeID.ApplyTo(fields, "property_type_id")
	h.Patch.ProjectID.ApplyTo(fields, "project_id")
	h.Patch.OriginProfileID.ApplyTo(fields, "origin_profile_id")
	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// LocationHandler
// ─────────────────────────────────────────────────────────────────────────────
type LocationHandler struct {
	Patch *dto.LocationUpdate
	Repo  property_repo.PropertyLocationRepository
}

func (h *LocationHandler) Name() string     { return "location" }
func (h *LocationHandler) FKColumn() string { return "location_id" }
func (h *LocationHandler) Present() bool    { return h.Patch != nil }

func (h *LocationHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyLocation{
		AddressDetail: h.Patch.AddressDetail.OrElse(""),
		MapURL:        h.Patch.MapURL.OrElse(""),
		Latitude:      h.Patch.Latitude.Ptr(),
		Longitude:     h.Patch.Longitude.Ptr(),
		RegionID:      h.Patch.RegionID.Ptr(),
		ProvinceID:    h.Patch.ProvinceID.Ptr(),
		DistrictID:    h.Patch.DistrictID.Ptr(),
		WardID:        h.Patch.WardID.Ptr(),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *LocationHandler) Update(ctx context.Context, id uint64) error {
	fields := make(map[string]any)
	h.Patch.AddressDetail.ApplyTo(fields, "address_detail")
	h.Patch.MapURL.ApplyTo(fields, "map_url")
	h.Patch.Latitude.ApplyTo(fields, "latitude")
	h.Patch.Longitude.ApplyTo(fields, "longitude")
	h.Patch.RegionID.ApplyTo(fields, "region_id")
	h.Patch.ProvinceID.ApplyTo(fields, "province_id")
	h.Patch.DistrictID.ApplyTo(fields, "district_id")
	h.Patch.WardID.ApplyTo(fields, "ward_id")
	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// LandInfoHandler (có sub‑groups)
// ─────────────────────────────────────────────────────────────────────────────
type LandInfoHandler struct {
	Patch *dto.LandInfoUpdate
	Repo  property_repo.PropertyLandInfoRepository
}

func (h *LandInfoHandler) Name() string     { return "land_info" }
func (h *LandInfoHandler) FKColumn() string { return "land_info_id" }
func (h *LandInfoHandler) Present() bool    { return h.Patch != nil }

func (h *LandInfoHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyLandInfo{
		Note:        h.Patch.Note.OrElse(""),
		LandNote:    h.Patch.LandNote.OrElse(""),
		PurposeUsed: enums.EPurposeUsed(h.Patch.PurposeUsed.OrElse(0)),
	}
	if h.Patch.Basic != nil {
		e.DocumentNo = h.Patch.Basic.DocumentNo.OrElse("")
		e.IssuringAuth = h.Patch.Basic.IssuringAuth.OrElse("")
		e.DocumentType = enums.EDocType(h.Patch.Basic.DocumentType.OrElse(0))
		e.Plot = h.Patch.Basic.Plot.OrElse(0)
		e.Sheet = h.Patch.Basic.Sheet.OrElse(0)
	}
	if h.Patch.Area != nil {
		e.AreaTotal = h.Patch.Area.AreaTotal.Ptr()
		e.AreaLand = h.Patch.Area.AreaLand.Ptr()
		e.AreaPlant = h.Patch.Area.AreaPlant.Ptr()
	}
	if h.Patch.Dimension != nil {
		e.FrontWidth = h.Patch.Dimension.FrontWidth.Ptr()
		e.Depth = h.Patch.Dimension.Depth.Ptr()
		e.StreetWidth = h.Patch.Dimension.StreetWidth.Ptr()
	}
	if h.Patch.Expiry != nil {
		e.ExpiredLand = parseTimePtr(h.Patch.Expiry.ExpiredLand)
		e.ExpiredPlant = parseTimePtr(h.Patch.Expiry.ExpiredPlant)
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *LandInfoHandler) Update(ctx context.Context, id uint64) error {
	fields := make(map[string]any)
	h.Patch.Note.ApplyTo(fields, "note")
	h.Patch.LandNote.ApplyTo(fields, "land_note")
	h.Patch.PurposeUsed.ApplyTo(fields, "purpose_used")

	if h.Patch.Basic != nil {
		h.Patch.Basic.DocumentNo.ApplyTo(fields, "document_no")
		h.Patch.Basic.IssuringAuth.ApplyTo(fields, "issuring_auth")
		h.Patch.Basic.DocumentType.ApplyTo(fields, "document_type")
		h.Patch.Basic.Plot.ApplyTo(fields, "plot")
		h.Patch.Basic.Sheet.ApplyTo(fields, "sheet")
	}
	if h.Patch.Area != nil {
		h.Patch.Area.AreaTotal.ApplyTo(fields, "area_total")
		h.Patch.Area.AreaLand.ApplyTo(fields, "area_land")
		h.Patch.Area.AreaPlant.ApplyTo(fields, "area_plant")
	}
	if h.Patch.Dimension != nil {
		h.Patch.Dimension.FrontWidth.ApplyTo(fields, "front_width")
		h.Patch.Dimension.Depth.ApplyTo(fields, "depth")
		h.Patch.Dimension.StreetWidth.ApplyTo(fields, "street_width")
	}
	if h.Patch.Expiry != nil {
		applyTimeField(fields, "expired_land", h.Patch.Expiry.ExpiredLand)
		applyTimeField(fields, "expired_plant", h.Patch.Expiry.ExpiredPlant)
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
	Patch *dto.BuildingUpdate
	Repo  property_repo.PropertyBuildingInfoRepository
}

func (h *BuildingHandler) Name() string     { return "building_info" }
func (h *BuildingHandler) FKColumn() string { return "building_info_id" }
func (h *BuildingHandler) Present() bool    { return h.Patch != nil }

func (h *BuildingHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyBuildingInfo{
		Note:             h.Patch.Note.OrElse(""),
		BuildStatus:      enums.EBuildStatus(h.Patch.BuildStatus.OrElse(0)),
		BuildingType:     enums.EBuildingType(h.Patch.BuildingType.OrElse(0)),
		Direction:        enums.EHouseOrient(h.Patch.Direction.OrElse(0)),
		BalconyDirection: enums.EHouseOrient(h.Patch.BalconyDir.OrElse(0)),
		AreaActual:       h.Patch.AreaActual.Ptr(),
		AreaFloor:        h.Patch.AreaFloor.Ptr(),
		AreaConstruction: h.Patch.AreaConstruction.Ptr(),
		Floors:           h.Patch.Floors.Ptr(),
		RoomNumber:       h.Patch.RoomNumber.Ptr(),
		Bedrooms:         h.Patch.Bedrooms.Ptr(),
		Bathrooms:        h.Patch.Bathrooms.Ptr(),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *BuildingHandler) Update(ctx context.Context, id uint64) error {
	fields := make(map[string]any)
	h.Patch.Note.ApplyTo(fields, "note")
	h.Patch.BuildStatus.ApplyTo(fields, "build_status")
	h.Patch.BuildingType.ApplyTo(fields, "building_type")
	h.Patch.Direction.ApplyTo(fields, "direction")
	h.Patch.BalconyDir.ApplyTo(fields, "balcony_direction")
	h.Patch.AreaActual.ApplyTo(fields, "area_actual")
	h.Patch.AreaFloor.ApplyTo(fields, "area_floor")
	h.Patch.AreaConstruction.ApplyTo(fields, "area_construction")
	h.Patch.Floors.ApplyTo(fields, "floors")
	h.Patch.RoomNumber.ApplyTo(fields, "room_number")
	h.Patch.Bedrooms.ApplyTo(fields, "bedrooms")
	h.Patch.Bathrooms.ApplyTo(fields, "bathrooms")
	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─────────────────────────────────────────────────────────────────────────────
// EvidenceHandler
// ─────────────────────────────────────────────────────────────────────────────
type EvidenceHandler struct {
	Patch *dto.EvidenceUpdate
	Repo  property_repo.PropertyEvidenceRepository
}

func (h *EvidenceHandler) Name() string     { return "evidence" }
func (h *EvidenceHandler) FKColumn() string { return "edvidence_id" }
func (h *EvidenceHandler) Present() bool    { return h.Patch != nil }

func (h *EvidenceHandler) Insert(ctx context.Context) (uint64, error) {
	e := &domain.PropertyEdvidence{
		Title:       h.Patch.Title.OrElse(""),
		Description: h.Patch.Description.OrElse(""),
		FileID:      h.Patch.FileID.Ptr(),
	}
	if err := h.Repo.Create(ctx, e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (h *EvidenceHandler) Update(ctx context.Context, id uint64) error {
	fields := make(map[string]any)
	h.Patch.Title.ApplyTo(fields, "title")
	h.Patch.Description.ApplyTo(fields, "description")
	h.Patch.FileID.ApplyTo(fields, "file_id")
	if len(fields) == 0 {
		return nil
	}
	return h.Repo.UpdateFields(ctx, id, fields)
}

// ─── helpers ─────────────────────────────────────────────────────────────────
func parseTimePtr(t optional.String) *time.Time {
	if !t.IsPresent() {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, t.Value())
	if err != nil {
		return nil
	}
	return &parsed
}

func applyTimeField(fields map[string]any, col string, t optional.String) {
	if t.IsPresent() {
		if parsed, err := time.Parse(time.RFC3339, t.Value()); err == nil {
			fields[col] = parsed
		}
	} else if t.IsNull() {
		fields[col] = nil
	}
}
