package services

import (
	"context"
	"fmt"
	"time"

	_enum "common/domain/enum"
	"mdn4/internal/interface/repo"
	"mdn4/modules/auth/dto"
	"mdn4/modules/user/domain/user"
)

type AuthUserProfileUsecase struct {
	UserInfoRepo repo.ProfileAccountRepository
}

func NewUserInfoUsecase(userInfoRepo repo.ProfileAccountRepository) *AuthUserProfileUsecase {
	return &AuthUserProfileUsecase{
		UserInfoRepo: userInfoRepo,
	}
}

// LockAccount khóa tài khoản (tạm thời hoặc vĩnh viễn)
func (u *AuthUserProfileUsecase) LockAccount(ctx context.Context, req *dto.LockAccountRequest) (*dto.LockAccountResponse, error) {
	userInfo, err := u.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Lỗi khi kiểm tra trạng thái user: " + err.Error(),
		}, err
	}

	if userInfo.IsLocked() {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "User đã bị khóa",
		}, fmt.Errorf("user already locked")
	}

	err = u.UserInfoRepo.LockUser(ctx, req.ProfileID, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Lỗi khi khóa user: " + err.Error(),
		}, err
	}

	response := &dto.LockAccountResponse{
		Success:   true,
		Message:   "Khóa user thành công",
		ProfileID: req.ProfileID,
		LockType:  req.LockType,
		LockedAt:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	if req.LockType == _enum.EUserStatusTemporaryLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// UnlockAccount mở khóa tài khoản
func (u *AuthUserProfileUsecase) UnlockAccount(ctx context.Context, req *dto.UnlockAccountRequest) (*dto.UnlockAccountResponse, error) {
	userInfo, err := u.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "Lỗi khi kiểm tra trạng thái user: " + err.Error(),
		}, err
	}

	if !userInfo.IsLocked() {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "User không bị khóa",
		}, fmt.Errorf("user is not locked")
	}

	err = u.UserInfoRepo.UnlockUser(ctx, req.ProfileID, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "Lỗi khi mở khóa user: " + err.Error(),
		}, err
	}

	return &dto.UnlockAccountResponse{
		Success:    true,
		Message:    "Mở khóa user thành công",
		ProfileID:  req.ProfileID,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetAccountStatus lấy trạng thái tài khoản
func (u *AuthUserProfileUsecase) GetAccountStatus(ctx context.Context, req *dto.GetAccountStatusRequest) (*dto.GetAccountStatusResponse, error) {
	userInfo, err := u.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.GetAccountStatusResponse{}, err
	}

	var lockType _enum.EUserStatus
	if userInfo.IsTemporarilyLocked() {
		lockType = _enum.EUserStatusTemporaryLocked
	} else if userInfo.IsPermanentlyLocked() {
		lockType = _enum.EUserStatusPermanentlyLocked
	}

	var isExpired bool
	if userInfo.IsTemporarilyLocked() {
		isExpired = userInfo.IsLockExpired()
	}

	response := &dto.GetAccountStatusResponse{
		ProfileID:  userInfo.ProfileID,
		Status:     uint32(userInfo.Status),
		StatusText: userInfo.Status.String(),
		IsLocked:   userInfo.IsLocked(),
		LockType:   lockType,
		CanLogin:   userInfo.CanLogin(),
		IsExpired:  isExpired,
	}

	if userInfo.IsLocked() {
		if userInfo.LockedAt != nil {
			response.LockedAt = userInfo.LockedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		response.LockReason = userInfo.LockReason
		response.LockedBy = userInfo.LockedBy

		if userInfo.LockedUntil != nil {
			response.LockedUntil = userInfo.LockedUntil.Format("2006-01-02T15:04:05Z07:00")
		}
	}

	return response, nil
}

// CheckAndUnlockExpiredAccounts kiểm tra và tự động mở khóa các tài khoản hết hạn
func (u *AuthUserProfileUsecase) CheckAndUnlockExpiredAccounts(ctx context.Context) error {
	return u.UserInfoRepo.CheckLockExpiration(ctx)
}

// LockMultipleUsers khóa nhiều user cùng lúc
func (u *AuthUserProfileUsecase) LockMultipleUsers(ctx context.Context, req *dto.LockMultipleUsersRequest) (*dto.LockMultipleUsersResponse, error) {
	if req.LockType == _enum.EUserStatusTemporaryLocked && req.Duration <= 0 {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Thời gian khóa tạm thời phải lớn hơn 0",
		}, fmt.Errorf("invalid duration for temporary lock")
	}

	err := u.UserInfoRepo.LockMultipleUsers(ctx, req.ProfileIDs, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi khóa nhiều user: " + err.Error(),
		}, err
	}

	response := &dto.LockMultipleUsersResponse{
		Success:    true,
		Message:    "Khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		LockType:   req.LockType,
		LockedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	if req.LockType == _enum.EUserStatusTemporaryLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// UnlockMultipleUsers mở khóa nhiều user cùng lúc
func (u *AuthUserProfileUsecase) UnlockMultipleUsers(ctx context.Context, req *dto.UnlockMultipleUsersRequest) (*dto.UnlockMultipleUsersResponse, error) {
	err := u.UserInfoRepo.UnlockMultipleUsers(ctx, req.ProfileIDs, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi mở khóa nhiều user: " + err.Error(),
		}, err
	}

	return &dto.UnlockMultipleUsersResponse{
		Success:    true,
		Message:    "Mở khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetAuthUsersByProfileIDs lấy thông tin account status theo danh sách profile IDs
func (u *AuthUserProfileUsecase) GetAuthUsersByProfileIDs(ctx context.Context, req []uint64) ([]user.UserProfileEntity, error) {
	return u.UserInfoRepo.GetUsersByProfileIDs(ctx, req)
}
