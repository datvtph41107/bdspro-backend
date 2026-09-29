package enums

// ─────────────────────────────────────────────────────────────────────────────
// hub/internal/enums/applink_action.go
// ─────────────────────────────────────────────────────────────────────────────

// XX = service  (01=user, 02=bdspro, 08=chat ...)
// YY = event    (00=default, 01=detail, 02=product ...)
type ApplinkAction int32

const (
	// ── Unspecified ─────────────────────────────────────
	ApplinkActionUnspecified ApplinkAction = 0

	// ── User service (01xx) ─────────────────────────────
	ApplinkActionUserDefault   ApplinkAction = 100
	ApplinkActionUserProfile   ApplinkAction = 101
	ApplinkActionUserSettings  ApplinkAction = 102
	ApplinkActionUserFollowers ApplinkAction = 103

	// ── BDSPro service (02xx) ───────────────────────────
	ApplinkActionBdsproDefault  ApplinkAction = 200
	ApplinkActionBdsproProperty ApplinkAction = 201
	ApplinkActionBdsproProduct  ApplinkAction = 202
	ApplinkActionBdsproProject  ApplinkAction = 203
	ApplinkActionBdsproSearch   ApplinkAction = 204

	// ── Chat service (08xx) ─────────────────────────────
	ApplinkActionChatDefault ApplinkAction = 800
	ApplinkActionChatRoom    ApplinkAction = 801
	ApplinkActionChatMessage ApplinkAction = 802
	ApplinkActionChatCall    ApplinkAction = 803
)

// Service trả về service code (XX)
func (a ApplinkAction) Service() int {
	return int(a) / 100
}

// Event trả về event code (YY)
func (a ApplinkAction) Event() int {
	return int(a) % 100
}

var applinkActionNames = map[ApplinkAction]string{
	ApplinkActionUnspecified: "unspecified",

	ApplinkActionUserDefault:   "user_default",
	ApplinkActionUserProfile:   "user_profile",
	ApplinkActionUserSettings:  "user_settings",
	ApplinkActionUserFollowers: "user_followers",

	ApplinkActionBdsproDefault:  "bdspro_default",
	ApplinkActionBdsproProperty: "bdspro_property",
	ApplinkActionBdsproProduct:  "bdspro_product",
	ApplinkActionBdsproProject:  "bdspro_project",
	ApplinkActionBdsproSearch:   "bdspro_search",

	ApplinkActionChatDefault: "chat_default",
	ApplinkActionChatRoom:    "chat_room",
	ApplinkActionChatMessage: "chat_message",
	ApplinkActionChatCall:    "chat_call",
}

// String implement fmt.Stringer
func (a ApplinkAction) String() string {
	if name, ok := applinkActionNames[a]; ok {
		return name
	}
	return "unknown"
}

// IsValid kiểm tra action có nằm trong danh sách hợp lệ không
func (a ApplinkAction) IsValid() bool {
	_, ok := applinkActionNames[a]
	return ok && a != ApplinkActionUnspecified
}
