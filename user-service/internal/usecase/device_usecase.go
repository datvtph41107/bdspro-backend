package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"strconv"
	"strings"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/interface/repo"
)

// DeviceUsecase xử lý logic nghiệp vụ liên quan đến thiết bị
type DeviceUsecase struct {
	deviceRepo repo.DeviceRepository
}

// NewDeviceUsecase khởi tạo DeviceUsecase
func NewDeviceUsecase(deviceRepo repo.DeviceRepository) *DeviceUsecase {
	return &DeviceUsecase{
		deviceRepo: deviceRepo,
	}
}

// RegisterDevice tạo mới hoặc cập nhật thông tin thiết bị
func (u *DeviceUsecase) RegisterDevice(ctx context.Context, req *dto.CreateDeviceRequest) (*auth.DeviceEntity, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "Yêu cầu không hợp lệ")
	}

	if strings.TrimSpace(req.DeviceID) == "" {
		return nil, _errors.ReturnError(400, "deviceId là bắt buộc")
	}

	lastSeen := _utils.TimeNowPtr()
	if req.LastSeenAt != nil && strings.TrimSpace(*req.LastSeenAt) != "" {
		parsed := _utils.ParseStringToTime(strings.TrimSpace(*req.LastSeenAt))
		if parsed == nil {
			return nil, _errors.ReturnError(400, "lastSeenAt không hợp lệ, định dạng phải là RFC3339")
		}
		lastSeen = parsed
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	orgID := _utils.GetOrganizationIdFromContext(ctx)
	authID := _utils.GetAuthIdFromContext(ctx)

	existing, err := u.deviceRepo.FindByDeviceIDAndModelAndManufacturer(ctx, req.DeviceID, req.Model, req.Manufacturer)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể truy vấn thông tin thiết bị")
	}

	if existing == nil {
		device := &auth.DeviceEntity{
			DeviceID:       req.DeviceID,
			DeviceName:     req.DeviceName,
			DeviceType:     req.DeviceType,
			Platform:       req.Platform,
			OSVersion:      req.OSVersion,
			AppVersion:     req.AppVersion,
			BuildNumber:    req.BuildNumber,
			Manufacturer:   req.Manufacturer,
			Model:          req.Model,
			Locale:         req.Locale,
			Timezone:       req.Timezone,
			IPAddress:      req.IPAddress,
			UserAgent:      req.UserAgent,
			LastSeenAt:     lastSeen,
			ProfileID:      toUint64Pointer(profileID),
			OrganizationID: toUint64Pointer(orgID),
			AuthID:         toUint64Pointer(authID),
		}

		// Chỉ set PushToken nếu có giá trị (không rỗng)
		if strings.TrimSpace(req.PushToken) != "" {
			device.PushToken = strings.TrimSpace(req.PushToken)
		}

		created, createErr := u.deviceRepo.Create(ctx, device)
		if createErr != nil {
			return nil, _errors.ReturnError(500, "Không thể lưu thông tin thiết bị")
		}
		return created, nil
	}

	// Cập nhật nếu đã tồn tại
	existing.DeviceName = req.DeviceName
	existing.DeviceType = req.DeviceType
	existing.Platform = req.Platform
	existing.OSVersion = req.OSVersion
	existing.AppVersion = req.AppVersion
	existing.BuildNumber = req.BuildNumber
	existing.Manufacturer = req.Manufacturer
	existing.Model = req.Model
	existing.Locale = req.Locale
	existing.Timezone = req.Timezone
	existing.IPAddress = req.IPAddress
	existing.UserAgent = req.UserAgent
	existing.LastSeenAt = lastSeen

	// Chỉ update PushToken nếu có giá trị (không rỗng), nếu rỗng thì giữ nguyên giá trị cũ
	if strings.TrimSpace(req.PushToken) != "" {
		existing.PushToken = strings.TrimSpace(req.PushToken)
	}

	// Chỉ cập nhật profileId, organizationId, authId nếu có giá trị từ context
	// Không ghi đè nếu thiếu (giữ nguyên giá trị cũ trong DB)
	if profileID > 0 {
		existing.ProfileID = toUint64Pointer(profileID)
	}
	if orgID > 0 {
		existing.OrganizationID = toUint64Pointer(orgID)
	}
	if authID > 0 {
		existing.AuthID = toUint64Pointer(authID)
	}

	updated, updateErr := u.deviceRepo.Update(ctx, existing)
	if updateErr != nil {
		return nil, _errors.ReturnError(500, "Không thể cập nhật thông tin thiết bị")
	}

	return updated, nil
}

// UpdateDeviceWithAuthInfo cập nhật authId và profileId vào device đã tồn tại
func (u *DeviceUsecase) UpdateDeviceWithAuthInfo(ctx context.Context, deviceID string, authID, profileID uint64) error {
	if strings.TrimSpace(deviceID) == "" {
		return nil // Không có deviceID thì bỏ qua
	}

	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		return err
	}
	existing, err := u.deviceRepo.FindByDeviceID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		// Device chưa tồn tại, không cần update
		return nil
	}

	// Cập nhật authId và profileId
	if authID > 0 {
		existing.AuthID = toUint64Pointer(authID)
	}
	if profileID > 0 {
		existing.ProfileID = toUint64Pointer(profileID)
	}

	_, err = u.deviceRepo.Update(ctx, existing)
	return err
}

func toUint64Pointer(value uint64) *uint64 {
	if value == 0 {
		return nil
	}
	v := value
	return &v
}
