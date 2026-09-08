package mapper

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	sharepb "pb/types/shared"
	"user/internal/domain/auth"
	"user/internal/dto"
)

// AuthSecurityMapper chuyển đổi dữ liệu auth security sang proto
type AuthSecurityMapper struct {
}

func NewAuthSecurityMapper() *AuthSecurityMapper {
	return &AuthSecurityMapper{}
}

// MapToProto chuyển DTO sang sharepb.AuthDataV3Proto
func (m *AuthSecurityMapper) MapToProto(data *dto.AuthSecurityDataDTO) *sharepb.AuthDataV3Proto {
	if data == nil {
		return nil
	}

	result := &sharepb.AuthDataV3Proto{
		AuthId: data.AuthID,
	}

	if data.Status != nil {
		result.Status = m.mapStatus(data.Status)
	}
	if data.PIN != nil {
		result.Pin = m.mapPin(data.PIN)
	}
	if len(data.Devices) > 0 {
		result.Devices = m.mapDevices(data.Devices)
	}
	if data.User != nil {
		result.User = m.mapUser(data.User)
	}
	if len(data.PersonConfigs) > 0 {
		result.PersonConfigs = m.mapPersonConfigs(data.PersonConfigs)
	}

	return result
}

func (m *AuthSecurityMapper) mapStatus(status *auth.UserStatusEntity) *sharepb.UserStatusV3Proto {
	var lockedUntil string
	if !status.LockedUntil.IsZero() {
		lock := status.LockedUntil
		lockedUntil = _utils.FormatTimeToString(&lock)
	}

	var updatedAt string
	if !status.UpdatedAt.IsZero() {
		update := status.UpdatedAt
		updatedAt = _utils.FormatTimeToString(&update)
	}

	return &sharepb.UserStatusV3Proto{
		AuthId:      status.AuthID,
		Active:      status.Active,
		Verified:    status.Verified,
		LockedUntil: lockedUntil,
		UpdatedAt:   updatedAt,
	}
}

func (m *AuthSecurityMapper) mapPin(pin *auth.UserPINEntity) *sharepb.PinV3Proto {
	return &sharepb.PinV3Proto{
		// Pin:          pin.PIN,
		AuthId:       pin.AuthID,
		PinCheckTime: uint32(pin.PINCheckTime),
		PinDate:      _utils.FormatTimeToString(pin.PINDate),
		IsActive:     pin.IsActive,
		LockedUntil:  _utils.FormatTimeToString(pin.LockedUntil),
	}
}

func (m *AuthSecurityMapper) mapDevices(devices []*auth.DeviceEntity) []*sharepb.DeviceV3Proto {
	result := make([]*sharepb.DeviceV3Proto, 0, len(devices))
	for _, device := range devices {
		if device == nil {
			continue
		}
		result = append(result, m.mapDevice(device))
	}
	return result
}

func (m *AuthSecurityMapper) mapDevice(device *auth.DeviceEntity) *sharepb.DeviceV3Proto {
	return &sharepb.DeviceV3Proto{
		Id:             device.ID,
		DeviceId:       device.DeviceID,
		DeviceName:     device.DeviceName,
		DeviceType:     device.DeviceType,
		Platform:       device.Platform,
		OsVersion:      device.OSVersion,
		AppVersion:     device.AppVersion,
		BuildNumber:    device.BuildNumber,
		Manufacturer:   device.Manufacturer,
		Model:          device.Model,
		Locale:         device.Locale,
		Timezone:       device.Timezone,
		PushToken:      device.PushToken,
		LastSeenAt:     _utils.FormatTimeToString(device.LastSeenAt),
		AuthId:         derefUint64(device.AuthID),
		ProfileId:      derefUint64(device.ProfileID),
		OrganizationId: derefUint64(device.OrganizationID),
		IpAddress:      device.IPAddress,
		UserAgent:      device.UserAgent,
	}
}

func (m *AuthSecurityMapper) mapUser(user *_dto.UserV3DTO) *sharepb.UserV3Proto {
	if user == nil {
		return nil
	}

	result := &sharepb.UserV3Proto{
		ProfileId:         user.ProfileID,
		FullName:          user.FullName,
		Avatar:            user.Avatar,
		RoleRealEstate:    uint32(user.RoleRealEstate),
		MainAreaIds:       user.MainAreaIds,
		Visibility:        uint32(user.Visibility),
		ProfileVisibility: uint32(user.ProfileVisibility),
		StatusOnline:      uint32(user.StatusOnline),
	}

	return result
}

func (m *AuthSecurityMapper) mapPersonConfigs(configs []*dto.PersonConfigDTO) []*sharepb.PersonConfigV3Proto {
	if len(configs) == 0 {
		return []*sharepb.PersonConfigV3Proto{}
	}

	result := make([]*sharepb.PersonConfigV3Proto, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		result = append(result, &sharepb.PersonConfigV3Proto{
			Id:        cfg.ID,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			UserId:    cfg.UserID,
			Channel:   uint32(cfg.Channel),
		})
	}

	return result
}

func derefUint64(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}
