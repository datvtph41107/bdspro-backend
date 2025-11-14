package usecase

import (
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/interface/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"strings"
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
func (u *DeviceUsecase) RegisterDevice(ctx context.Context, req *dto.CreateDeviceRequest) (*domain.DeviceEntity, error) {
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

	existing, err := u.deviceRepo.FindByDeviceID(ctx, req.DeviceID)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể truy vấn thông tin thiết bị")
	}

	if existing == nil {
		device := &domain.DeviceEntity{
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
			PushToken:      req.PushToken,
			IPAddress:      req.IPAddress,
			UserAgent:      req.UserAgent,
			LastSeenAt:     lastSeen,
			ProfileID:      toUint64Pointer(profileID),
			OrganizationID: toUint64Pointer(orgID),
			AuthID:         toUint64Pointer(authID),
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
	existing.PushToken = req.PushToken
	existing.IPAddress = req.IPAddress
	existing.UserAgent = req.UserAgent
	existing.LastSeenAt = lastSeen

	if pointer := toUint64Pointer(profileID); pointer != nil {
		existing.ProfileID = pointer
	}
	if pointer := toUint64Pointer(orgID); pointer != nil {
		existing.OrganizationID = pointer
	}
	if pointer := toUint64Pointer(authID); pointer != nil {
		existing.AuthID = pointer
	}

	updated, updateErr := u.deviceRepo.Update(ctx, existing)
	if updateErr != nil {
		return nil, _errors.ReturnError(500, "Không thể cập nhật thông tin thiết bị")
	}

	return updated, nil
}

func toUint64Pointer(value uint64) *uint64 {
	if value == 0 {
		return nil
	}
	v := value
	return &v
}
