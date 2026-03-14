package validator

import (
	"common/pkg/validate"
	"common/pkg/validate/grpcerror"
	"fmt"
	bdspropb "pb/types/bdspro"
)

// ─────────────────────────────────────────────────────────────────────────────
// Validation layer
//
// Quy tắc:
//   - Mỗi sub-domain có validator function riêng → tái sử dụng ở Create
//   - Sub-domain nil → skip (không validate field không được gửi)
//   - Tất cả lỗi tích lũy qua validate.Error.Merge → client biết hết 1 lần
//   - Chỉ format validation ở đây — business validation thuộc usecase
// ─────────────────────────────────────────────────────────────────────────────

func ValidateUpdateRequest(req *bdspropb.UpdatePropertyRequest) error {
	e := validate.New()

	if req.Id == 0 {
		e.Add("id", "required")
	}

	e.Merge(validateInfoPayload(req.Info))
	e.Merge(validateLocationPayload(req.Location))
	e.Merge(validateLandInfoPayload(req.LandInfo))
	e.Merge(validateBuildingPayload(req.BuildingInfo))
	e.Merge(validateEvidencePayload(req.Evidence))
	e.Merge(validateMediaPayload(req.Media))
	e.Merge(validateExternalRefPayload(req.ExternalRefs))
	// lineage_meta: không validate format đặc biệt ngoài date
	e.Merge(validateLineageMetaPayload(req.LineageMeta))

	return grpcerror.FromValidation(e)
}

func validateInfoPayload(p *bdspropb.UpdateInfoPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "info.title", p.Title, 255)
	validate.StringMaxLen(e, "info.note", p.Note, 1000)
	validate.StringMaxLen(e, "info.unit_code", p.UnitCode, 50)
	validate.StringMaxLen(e, "info.identifier", p.Identifier, 100)
	validate.StringMaxLen(e, "info.level", p.Level, 100)
	return e
}

func validateLocationPayload(p *bdspropb.UpdateLocationPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "location.address_detail", p.AddressDetail, 100)
	validate.URLOpt(e, "location.map_url", p.MapUrl)
	validate.Latitude(e, "location.latitude", p.Latitude)
	validate.Longitude(e, "location.longitude", p.Longitude)
	return e
}

func validateLandInfoPayload(p *bdspropb.UpdateLandInfoPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "landInfo.note", p.Note, 1000)
	validate.StringMaxLen(e, "landInfo.landNote", p.LandNote, 1000)

	validate.DateStringOpt(e, "landInfo.expiredLand", p.ExpiredLand)
	validate.DateStringOpt(e, "landInfo.expiredPlant", p.ExpiredPlant)
	// Cross-field: nếu cả 2 có giá trị → land phải trước plant
	validate.DateOrder(e,
		"landInfo.expiredLand",
		"landInfo.expiredPlant",
		p.ExpiredLand,
		p.ExpiredPlant,
	)
	validate.Float64NonNegative(e, "landInfo.areaTotal", p.AreaTotal)
	validate.Float64NonNegative(e, "landInfo.areaLand", p.AreaLand)
	validate.Float64NonNegative(e, "landInfo.areaPlant", p.AreaPlant)
	validate.Float64NonNegative(e, "landInfo.frontWidth", p.FrontWidth)
	validate.Float64NonNegative(e, "landInfo.depth", p.Depth)
	validate.Float64NonNegative(e, "landInfo.streetWidth", p.StreetWidth)
	return e
}

func validateBuildingPayload(p *bdspropb.UpdateBuildingPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "buildingInfo.note", p.Note, 1000)
	validate.Float64NonNegative(e, "buildingInfo.area_actual", p.AreaActual)
	validate.Float64NonNegative(e, "buildingInfo.area_floor", p.AreaFloor)
	validate.Float64NonNegative(e, "buildingInfo.area_construction", p.AreaConstruction)
	return e
}

func validateEvidencePayload(p *bdspropb.UpdateEvidencePayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "evidence.title", p.Title, 255)
	validate.StringMaxLen(e, "evidence.description", p.Description, 2000)
	return e
}

func validateMediaPayload(p *bdspropb.UpdateMediaPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	for i, item := range p.Items {
		f := func(field string) string { return fmt.Sprintf("media[%d].%s", i, field) }
		// INSERT (id=0): media_url và media_type bắt buộc
		if item.Id == 0 {
			url := item.MediaUrl
			validate.Required(e, f("media_url"), &url)
		}
		mt := item.MediaType
		validate.MediaType(e, f("media_type"), &mt)
	}
	return e
}

func validateExternalRefPayload(p *bdspropb.UpdateExternalRefPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	for i, item := range p.Items {
		f := func(field string) string { return fmt.Sprintf("external_refs[%d].%s", i, field) }
		// INSERT (id=0): external_ref_id bắt buộc
		if item.Id == 0 {
			v := item.ExternalRefId
			validate.RequiredID(e, f("external_ref_id"), &v)
		}
	}
	return e
}

func validateLineageMetaPayload(p *bdspropb.UpdateLineageMetaPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "lineage_meta.national_id", p.NationalId, 100)
	validate.DateStringOpt(e, "lineage_meta.verified_national_at", p.VerifiedNationalAt)
	return e
}
