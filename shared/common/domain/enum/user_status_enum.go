package _enum

type EUserStatus int

const (
	EUserStatusInactive          EUserStatus = 0  // Chưa hoạt động
	EUserStatusActive            EUserStatus = 10 // Hoạt động
	EUserStatusTemporaryLocked   EUserStatus = 20 // Khóa tạm thời
	EUserStatusPermanentlyLocked EUserStatus = 30 // Khóa vĩnh viễn
	EUserStatusSuspended         EUserStatus = 40 // Ngừng hoạt động
)

var UserStatusNames = map[EUserStatus]string{
	EUserStatusInactive:          "Chưa hoạt động",
	EUserStatusActive:            "Hoạt động",
	EUserStatusTemporaryLocked:   "Khóa tạm thời",
	EUserStatusPermanentlyLocked: "Khóa vĩnh viễn",
	EUserStatusSuspended:         "Ngừng hoạt động",
}

func (e EUserStatus) IsValid() bool {
	switch e {
	case EUserStatusInactive, EUserStatusActive, EUserStatusTemporaryLocked, EUserStatusPermanentlyLocked, EUserStatusSuspended:
		return true
	default:
		return false
	}
}

func (e EUserStatus) String() string {
	if name, ok := UserStatusNames[e]; ok {
		return name
	}
	return "Unknown"
}

func (e EUserStatus) IsLocked() bool {
	return e == EUserStatusTemporaryLocked || e == EUserStatusPermanentlyLocked
}

// CanLogin mirrors user_profile account gate without lock expiry.
// Expiry-aware checks belong on UserProfileEntity.CanLogin (uses locked_until).
func (e EUserStatus) CanLogin() bool {
	switch e {
	case EUserStatusActive:
		return true
	default:
		return false
	}
}
