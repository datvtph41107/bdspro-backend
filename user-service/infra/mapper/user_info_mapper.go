package mapper

import (
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/domain/user"
)

type UserInfoMapper struct {
}

func NewUserInfoMapper() *UserInfoMapper {
	return &UserInfoMapper{}
}

func (m *UserInfoMapper) MapAuthUserStatusInfoToPb(profile *user.UserProfileEntity) *userpb.AuthUserStatusInfo {
	result := &userpb.AuthUserStatusInfo{
		ProfileId:   profile.ProfileID,
		Status:      uint32(profile.Status),
		StatusText:  profile.Status.String(),
		IsLocked:    profile.IsLocked(),
		LockType:    uint32(profile.Status),
		LockedAt:    _utils.FormatTimeToString(profile.LockedAt),
		LockedUntil: _utils.FormatTimeToString(profile.LockedUntil),
		LockReason:  profile.LockReason,
		LockedBy:    profile.LockedBy,
		CanLogin:    profile.CanLogin(),
		IsExpired:   profile.IsLockExpired(),
	}

	if profile.RoleID != nil && *profile.RoleID > 0 {
		result.Role = &sharepb.Role{
			Id:       *profile.RoleID,
			RoleName: profile.RoleName,
			RoleKey:  profile.RoleKeyData,
		}
		if profile.RoleColor != "" || profile.RoleBgColor != "" {
			result.Role.Color = &sharepb.Color{
				ContentColor:    profile.RoleColor,
				BackgroundColor: profile.RoleBgColor,
			}
		}
	}

	return result
}
