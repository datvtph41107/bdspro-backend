package usecase

import (
	_errors "common/errors"
	"context"
	"errors"
	"fmt"
	"user/internal"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"

	"gorm.io/gorm"
)

// AuthSecurityUsecase gom nhóm dữ liệu status, pin và thiết bị cho 1 authId
type AuthSecurityUsecase struct {
	statusRepo           repo.StatusRepository
	pinRepo              repo.PINRepository
	deviceRepo           repo.DeviceRepository
	authMethodRepo       repo.AuthMethodRepository
	profileProvider      providers.ProfileProvider
	notificationProvider providers.NotificationProvider
}

// NewAuthSecurityUsecase khởi tạo usecase
func NewAuthSecurityUsecase(
	statusRepo repo.StatusRepository,
	pinRepo repo.PINRepository,
	deviceRepo repo.DeviceRepository,
	authMethodRepo repo.AuthMethodRepository,
	profileProvider providers.ProfileProvider,
	notificationProvider providers.NotificationProvider,
) *AuthSecurityUsecase {
	return &AuthSecurityUsecase{
		statusRepo:           statusRepo,
		pinRepo:              pinRepo,
		deviceRepo:           deviceRepo,
		authMethodRepo:       authMethodRepo,
		profileProvider:      profileProvider,
		notificationProvider: notificationProvider,
	}
}

// GetAuthSecurityData lấy thông tin status, pin và devices theo authId
func (u *AuthSecurityUsecase) GetAuthSecurityData(ctx context.Context, authID uint64) (*dto.AuthSecurityDataDTO, error) {
	if authID == 0 {
		return nil, _errors.ReturnError(service.AuthIDRequiredLower)
	}

	var statusErr error
	status, statusErr := u.statusRepo.GetByID(ctx, authID)
	if statusErr != nil {
		if !errors.Is(statusErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("load user status: %w", statusErr)
		}
		status = nil
	}

	var pinErr error
	pin, pinErr := u.pinRepo.GetByAuthID(ctx, authID)
	if pinErr != nil {
		if !errors.Is(pinErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("load PIN info: %w", pinErr)
		}
		pin = nil
	}

	devices, err := u.deviceRepo.ListByAuthID(ctx, authID)
	if err != nil {
		return nil, fmt.Errorf("load device list: %w", err)
	}

	return &dto.AuthSecurityDataDTO{
		AuthID:  authID,
		Status:  status,
		PIN:     pin,
		Devices: devices,
	}, nil
}

// GetAuthSecurityDataByProfileID lấy thông tin bảo mật thông qua profileId (userId)
func (u *AuthSecurityUsecase) GetAuthSecurityDataByProfileID(ctx context.Context, profileID uint64) (*dto.AuthSecurityDataDTO, error) {
	if profileID == 0 {
		return nil, _errors.ReturnError(service.ProfileIDRequired)
	}

	auth, err := u.authMethodRepo.GetFirstByUserID(ctx, profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _errors.ReturnError(service.AccountForProfileNotFound)
		}
		return nil, fmt.Errorf("load account info: %w", err)
	}

	authSecurityData, err := u.GetAuthSecurityData(ctx, auth.ID)
	if err != nil {
		return nil, err
	}

	if u.profileProvider != nil {
		profile, profileErr := u.profileProvider.GetByProfileIDV3(ctx, profileID)
		if profileErr != nil {
			return nil, fmt.Errorf("load user profile info: %w", profileErr)
		}
		authSecurityData.User = profile
	}

	if u.notificationProvider != nil {
		personConfigs, cfgErr := u.notificationProvider.GetPersonConfigs(ctx, profileID)
		if cfgErr != nil {
			return nil, fmt.Errorf("load personal notification config: %w", cfgErr)
		}
		authSecurityData.PersonConfigs = personConfigs
	}

	return authSecurityData, nil
}
