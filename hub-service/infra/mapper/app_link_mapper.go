package mapper

import (
	_utils "common/utils"
	"hub/helpers"
	"hub/internal/domain"
	"hub/internal/enums"
	hubpb "pb/types/hub"
)

// ── domain → proto ────────────────────────────────────────────────────────────

func ApplinkToProto(a *domain.Applink) *hubpb.ApplinkInfo {
	return &hubpb.ApplinkInfo{
		Code:       helpers.ApplinkCodeToBase64(a.Code), // uint64 → base64 khi expose
		RefId:      a.RefID,
		Action:     int32(a.Action),
		ActionName: a.Action.String(), // "user_profile", "chat_room"...
		CreatedAt:  _utils.FormatTimeToString(&a.CreatedAt),
	}
}

// ── proto → domain ────────────────────────────────────────────────────────────

// ProtoActionToDomain cast int32 từ request → enums.ApplinkAction
// Validation IsValid() thực hiện ở usecase layer
func ProtoActionToDomain(action int32) enums.ApplinkAction {
	return enums.ApplinkAction(action)
}
