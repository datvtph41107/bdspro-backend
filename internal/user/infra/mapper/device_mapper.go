package mapper

import (
	_utils "common/utils"
	authpb "pb/types/auth"
	"user/internal/domain/auth"
	"user/internal/dto"
)

// DeviceMapper chuyển đổi giữa device domain, dto và protobuf
type DeviceMapper struct {
}

// NewDeviceMapper khởi tạo DeviceMapper
func NewDeviceMapper() *DeviceMapper {
	return &DeviceMapper{}
}

// MapCreateDeviceRequestPbToDTO chuyển proto request sang DTO request
func (m *DeviceMapper) MapCreateDeviceRequestPbToDTO(req *authpb.CreateDeviceRequest) *dto.CreateDeviceRequest {
	if req == nil {
		return nil
	}

	var lastSeenPtr *string
	if req.LastSeenAt != "" {
		lastSeen := req.LastSeenAt
		lastSeenPtr = &lastSeen
	}

	return &dto.CreateDeviceRequest{
		DeviceID:     req.DeviceId,
		DeviceName:   req.DeviceName,
		DeviceType:   req.DeviceType,
		Platform:     req.Platform,
		OSVersion:    req.OsVersion,
		AppVersion:   req.AppVersion,
		BuildNumber:  req.BuildNumber,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
		Locale:       req.Locale,
		Timezone:     req.Timezone,
		PushToken:    req.PushToken,
		IPAddress:    req.IpAddress,
		UserAgent:    req.UserAgent,
		LastSeenAt:   lastSeenPtr,
	}
}

// MapDeviceToDTO chuyển Domain sang DTO response
func (m *DeviceMapper) MapDeviceToDTO(device *auth.DeviceEntity) *dto.DeviceResponse {
	if device == nil {
		return nil
	}

	return &dto.DeviceResponse{
		ID:             device.ID,
		DeviceID:       device.DeviceID,
		DeviceName:     device.DeviceName,
		DeviceType:     device.DeviceType,
		Platform:       device.Platform,
		OSVersion:      device.OSVersion,
		AppVersion:     device.AppVersion,
		BuildNumber:    device.BuildNumber,
		Manufacturer:   device.Manufacturer,
		Model:          device.Model,
		Locale:         device.Locale,
		Timezone:       device.Timezone,
		PushToken:      device.PushToken,
		IPAddress:      device.IPAddress,
		UserAgent:      device.UserAgent,
		ProfileID:      device.ProfileID,
		OrganizationID: device.OrganizationID,
		AuthID:         device.AuthID,
		LastSeenAt:     device.LastSeenAt,
		CreatedAt:      device.CreatedAt,
		UpdatedAt:      device.UpdatedAt,
	}
}

// MapDeviceToPb chuyển domain sang protobuf
func (m *DeviceMapper) MapDeviceToPb(device *auth.DeviceEntity) *authpb.Device {
	if device == nil {
		return nil
	}

	var lastSeen string
	if device.LastSeenAt != nil {
		lastSeen = _utils.FormatTimeToString(device.LastSeenAt)
	}

	var createdAt string
	if device.CreatedAt != nil {
		createdAt = _utils.FormatTimeToString(device.CreatedAt)
	}

	var updatedAt string
	if device.UpdatedAt != nil {
		updatedAt = _utils.FormatTimeToString(device.UpdatedAt)
	}

	var profileID uint64
	if device.ProfileID != nil {
		profileID = *device.ProfileID
	}

	var organizationID uint64
	if device.OrganizationID != nil {
		organizationID = *device.OrganizationID
	}

	var authID uint64
	if device.AuthID != nil {
		authID = *device.AuthID
	}

	return &authpb.Device{
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
		IpAddress:      device.IPAddress,
		UserAgent:      device.UserAgent,
		ProfileId:      profileID,
		OrganizationId: organizationID,
		AuthId:         authID,
		LastSeenAt:     lastSeen,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
