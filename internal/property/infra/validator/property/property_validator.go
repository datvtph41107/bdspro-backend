package property_validator

import (
	"common/pkg/validate"
	"common/pkg/validate/grpcerror"
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
		return grpcerror.FromValidation(e)
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
