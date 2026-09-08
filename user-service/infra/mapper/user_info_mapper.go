package mapper

import (
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/domain/auth"
)

type UserInfoMapper struct {
}

func NewUserInfoMapper() *UserInfoMapper {
	return &UserInfoMapper{}
}

func (m *UserInfoMapper) MapAuthUserStatusInfoToPb(user *auth.AuthUser) *userpb.AuthUserStatusInfo {
	result := &userpb.AuthUserStatusInfo{
		ProfileId:   user.ProfileID,
		Status:      uint32(user.Status),
		StatusText:  user.Status.Str(),
		IsLocked:    user.IsLocked(),
		LockType:    uint32(user.Status),
		LockedAt:    _utils.FormatTimeToString(user.LockedAt),
		LockedUntil: _utils.FormatTimeToString(user.LockedUntil),
		LockReason:  user.LockReason,
		LockedBy:    user.LockedBy,
		CanLogin:    user.CanLogin(),
		IsExpired:   user.IsLockExpired(),
	}

	if user.Role != nil {
		result.Role = &sharepb.Role{
			Id:       user.Role.ID,
			RoleName: user.Role.RoleName,
			RoleKey:  user.Role.RoleKey,
		}
		if user.Role.Color != nil {
			result.Role.Color = &sharepb.Color{
				Id:              user.Role.Color.ID,
				Name:            user.Role.Color.Name,
				ContentColor:    user.Role.Color.ContentColor,
				BackgroundColor: user.Role.Color.BackgroundColor,
				ColorKey:        user.Role.Color.ColorKey,
				HexCode:         user.Role.Color.HexCode,
				Description:     &user.Role.Color.Description,
				IsActive:        user.Role.Color.IsActive,
			}
		}
	}

	return result
}
