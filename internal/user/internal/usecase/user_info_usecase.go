package usecase

import (
	"context"
	"fmt"
	"time"

	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/interface/repo"
)

type AuthUserProfileUsecase struct {
	UserInfoRepo repo.UserInfoRepository
}

func NewUserInfoUsecase(userInfoRepo repo.UserInfoRepository) *AuthUserProfileUsecase {
	return &AuthUserProfileUsecase{
		UserInfoRepo: userInfoRepo,
	}
}

// LockAccount khóa tài khoản (tạm thời hoặc vĩnh viễn)
func (u *AuthUserProfileUsecase) LockAccount(ctx context.Context, req *dto.LockAccountRequest) (*dto.LockAccountResponse, error) {
	// Validate duration cho khóa tạm thời
	// if req.LockType == uint32(auth.StatusTemporarilyLocked) && req.Duration <= 0 {
	// 	return &dto.LockAccountResponse{
	// 		Success: false,
	// 		Message: "Thời gian khóa tạm thời phải lớn hơn 0",
	// 	}, fmt.Errorf("invalid duration for temporary lock")
	// }

	// Kiểm tra user đã bị khóa chưa
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

	// Thực hiện khóa user
	err = u.UserInfoRepo.LockUser(ctx, req.ProfileID, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Lỗi khi khóa user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.LockAccountResponse{
		Success:   true,
		Message:   "Khóa user thành công",
		ProfileID: req.ProfileID,
		LockType:  req.LockType,
		LockedAt:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Thêm thông tin thời gian khóa đến cho khóa tạm thời
	if req.LockType == auth.StatusTemporarilyLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// UnlockAccount mở khóa tài khoản
func (u *AuthUserProfileUsecase) UnlockAccount(ctx context.Context, req *dto.UnlockAccountRequest) (*dto.UnlockAccountResponse, error) {
	// Kiểm tra user có bị khóa không
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

	// Thực hiện mở khóa user
	err = u.UserInfoRepo.UnlockUser(ctx, req.ProfileID, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "Lỗi khi mở khóa user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.UnlockAccountResponse{
		Success:    true,
		Message:    "Mở khóa user thành công",
		ProfileID:  req.ProfileID,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// GetAccountStatus lấy trạng thái tài khoản
func (u *AuthUserProfileUsecase) GetAccountStatus(ctx context.Context, req *dto.GetAccountStatusRequest) (*dto.GetAccountStatusResponse, error) {
	// Lấy thông tin user
	userInfo, err := u.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.GetAccountStatusResponse{}, err
	}

	// Xác định trạng thái text
	var statusText string
	switch userInfo.Status {
	case auth.StatusActive:
		statusText = "Hoạt động"
	case auth.StatusInactive:
		statusText = "Vô hiệu hóa"
	case auth.StatusTemporarilyLocked:
		statusText = "Khóa tạm thời"
	case auth.StatusPermanentlyLocked:
		statusText = "Khóa vĩnh viễn"
	default:
		statusText = "Không xác định"
	}

	// Xác định loại khóa
	var lockType auth.AccountStatus
	if userInfo.IsTemporarilyLocked() {
		lockType = auth.StatusTemporarilyLocked
	} else if userInfo.IsPermanentlyLocked() {
		lockType = auth.StatusPermanentlyLocked
	}

	// Kiểm tra khóa có hết hạn không
	var isExpired bool
	if userInfo.IsTemporarilyLocked() {
		isExpired = userInfo.IsLockExpired()
	}

	// Tạo response
	response := &dto.GetAccountStatusResponse{
		ProfileID:  userInfo.ProfileID,
		Status:     uint32(userInfo.Status),
		StatusText: statusText,
		IsLocked:   userInfo.IsLocked(),
		LockType:   lockType,
		CanLogin:   userInfo.CanLogin(),
		IsExpired:  isExpired,
	}

	// Thêm thông tin khóa nếu có
	if userInfo.IsLocked() {
		response.LockedAt = userInfo.LockedAt.Format("2006-01-02T15:04:05Z07:00")
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
	// Validate duration cho khóa tạm thời
	if req.LockType == auth.StatusTemporarilyLocked && req.Duration <= 0 {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Thời gian khóa tạm thời phải lớn hơn 0",
		}, fmt.Errorf("invalid duration for temporary lock")
	}

	// Thực hiện khóa nhiều user
	err := u.UserInfoRepo.LockMultipleUsers(ctx, req.ProfileIDs, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi khóa nhiều user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.LockMultipleUsersResponse{
		Success:    true,
		Message:    "Khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		LockType:   req.LockType,
		LockedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Thêm thông tin thời gian khóa đến cho khóa tạm thời
	if req.LockType == auth.StatusTemporarilyLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// UnlockMultipleUsers mở khóa nhiều user cùng lúc
func (u *AuthUserProfileUsecase) UnlockMultipleUsers(ctx context.Context, req *dto.UnlockMultipleUsersRequest) (*dto.UnlockMultipleUsersResponse, error) {
	// Thực hiện mở khóa nhiều user
	err := u.UserInfoRepo.UnlockMultipleUsers(ctx, req.ProfileIDs, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi mở khóa nhiều user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.UnlockMultipleUsersResponse{
		Success:    true,
		Message:    "Mở khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// GetAuthUsersByProfileIDs lấy thông tin AuthUser theo danh sách profile IDs
func (u *AuthUserProfileUsecase) GetAuthUsersByProfileIDs(ctx context.Context, req []uint64) ([]auth.AuthUser, error) {
	// Lấy thông tin users từ database
	users, err := u.UserInfoRepo.GetUsersByProfileIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return users, nil
}
