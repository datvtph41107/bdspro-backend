package dto

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"user/internal/domain/auth"
)

// AuthSecurityDataDTO gom nhóm dữ liệu status, pin và thiết bị cho 1 authId
type AuthSecurityDataDTO struct {
	AuthID        uint64
	Status        *auth.UserStatusEntity
	PIN           *auth.UserPINEntity
	Devices       []*auth.DeviceEntity
	User          *_dto.UserV3DTO
	PersonConfigs []*PersonConfigDTO
}

type PersonConfigDTO struct {
	ID        uint64
	UserID    uint64
	Key       string
	Checked   bool
	IsDefault bool
	Channel   _enum.EChannelNotification
}
