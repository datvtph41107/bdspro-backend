package mapper

import (
	_utils "common/utils"
	"hub/helpers"
	"hub/internal/domain"
	"hub/internal/enums"
	hubpb "pb/types/hub"
)

// ── domain → proto ────────────────────────────────────────────────────────────

// ApplinkToProto baseURL từ config — link = baseURL + code
func ApplinkToProto(a *domain.Applink, baseURL string) *hubpb.ApplinkInfo {
	code := helpers.ApplinkCodeToBase64(a.Code)
	link := baseURL + code
	return &hubpb.ApplinkInfo{
		Code:       code,
		RefId:      a.RefID,
		Action:     int32(a.Action),
		ActionName: a.Action.String(),
		CreatedAt:  _utils.FormatTimeToString(&a.CreatedAt),
		Link:       link,
	}
}

// ── proto → domain ────────────────────────────────────────────────────────────

// ProtoActionToDomain cast int32 từ request → enums.ApplinkAction
// Validation IsValid() thực hiện ở usecase layer
func ProtoActionToDomain(action int32) enums.ApplinkAction {
	return enums.ApplinkAction(action)
}
