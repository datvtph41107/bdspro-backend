package handler

import (
	"context"
	authpb "pb/types/auth"

	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUserProfileHandler struct {
	authpb.UnimplementedUserInfoServiceServer
	UserInfoUsecase *usecase.AuthUserProfileUsecase
}

func NewUserInfoHandler(userInfoUsecase *usecase.AuthUserProfileUsecase) *AuthUserProfileHandler {
	return &AuthUserProfileHandler{
		UserInfoUsecase: userInfoUsecase,
	}
}

// @Summary Khóa tài khoản
// @Description Khóa tài khoản tạm thời hoặc vĩnh viễn
// @Tags Account Management
// @Accept json
// @Produce json
// @Param body body authpb.LockAccountRequest true "Thông tin khóa tài khoản"
// @Success 200 {object} authpb.LockAccountResponse
// @Failure 400 {object} authpb.LockAccountResponse
// @Failure 500 {object} authpb.LockAccountResponse
// @Router /account/lock [post]
func (h *AuthUserProfileHandler) LockAccount(ctx context.Context, req *authpb.LockAccountRequest) (*authpb.LockAccountResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.LockAccountRequest{
		ProfileID: req.ProfileId,
		LockType:  auth.AccountStatus(req.LockType),
		Duration:  int(req.Duration),
		Reason:    req.Reason,
		LockedBy:  req.LockedBy,
	}

	// Call usecase
	result, err := h.UserInfoUsecase.LockAccount(ctx, dtoReq)
	if err != nil {
		return &authpb.LockAccountResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to lock account: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.LockAccountResponse{
		Success:     result.Success,
		Message:     result.Message,
		ProfileId:   result.ProfileID,
		LockType:    uint32(result.LockType),
		LockedAt:    result.LockedAt,
		LockedUntil: result.LockedUntil,
	}, nil
}

// @Summary Mở khóa tài khoản
// @Description Mở khóa tài khoản đã bị khóa
// @Tags Account Management
// @Accept json
// @Produce json
// @Param body body authpb.UnlockAccountRequest true "Thông tin mở khóa tài khoản"
// @Success 200 {object} authpb.UnlockAccountResponse
// @Failure 400 {object} authpb.UnlockAccountResponse
// @Failure 500 {object} authpb.UnlockAccountResponse
// @Router /account/unlock [post]
func (h *AuthUserProfileHandler) UnlockAccount(ctx context.Context, req *authpb.UnlockAccountRequest) (*authpb.UnlockAccountResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.UnlockAccountRequest{
		ProfileID:  req.ProfileId,
		Reason:     req.Reason,
		UnlockedBy: req.UnlockedBy,
	}

	// Call usecase
	result, err := h.UserInfoUsecase.UnlockAccount(ctx, dtoReq)
	if err != nil {
		return &authpb.UnlockAccountResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to unlock account: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.UnlockAccountResponse{
		Success:    result.Success,
		Message:    result.Message,
		ProfileId:  result.ProfileID,
		UnlockedAt: result.UnlockedAt,
	}, nil
}

// @Summary Lấy trạng thái tài khoản
// @Description Lấy thông tin trạng thái tài khoản
// @Tags Account Management
// @Accept json
// @Produce json
// @Param profileId path uint64 true "ID profile"
// @Success 200 {object} authpb.GetAccountStatusResponse
// @Failure 400 {object} authpb.GetAccountStatusResponse
// @Failure 500 {object} authpb.GetAccountStatusResponse
// @Router /account/status/{profileId} [get]
func (h *AuthUserProfileHandler) GetAccountStatus(ctx context.Context, req *authpb.GetAccountStatusRequest) (*authpb.GetAccountStatusResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.GetAccountStatusRequest{
		ProfileID: req.ProfileId,
	}

	// Call usecase
	result, err := h.UserInfoUsecase.GetAccountStatus(ctx, dtoReq)
	if err != nil {
		return &authpb.GetAccountStatusResponse{}, status.Errorf(codes.Internal, "failed to get account status: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.GetAccountStatusResponse{
		ProfileId:   result.ProfileID,
		Status:      result.Status,
		StatusText:  result.StatusText,
		IsLocked:    result.IsLocked,
		LockType:    uint32(result.LockType),
		LockedAt:    result.LockedAt,
		LockedUntil: result.LockedUntil,
		LockReason:  result.LockReason,
		LockedBy:    result.LockedBy,
		CanLogin:    result.CanLogin,
		IsExpired:   result.IsExpired,
	}, nil
}

// @Summary Khóa nhiều user
// @Description Khóa nhiều user cùng lúc
// @Tags Account Management
// @Accept json
// @Produce json
// @Param body body authpb.LockMultipleUsersRequest true "Thông tin khóa nhiều user"
// @Success 200 {object} authpb.LockMultipleUsersResponse
// @Failure 400 {object} authpb.LockMultipleUsersResponse
// @Failure 500 {object} authpb.LockMultipleUsersResponse
// @Router /account/lock-multiple [post]
func (h *AuthUserProfileHandler) LockMultipleUsers(ctx context.Context, req *authpb.LockMultipleUsersRequest) (*authpb.LockMultipleUsersResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.LockMultipleUsersRequest{
		ProfileIDs: req.ProfileIds,
		LockType:   auth.AccountStatus(req.LockType),
		Duration:   int(req.Duration),
		Reason:     req.Reason,
		LockedBy:   req.LockedBy,
	}

	// Call usecase
	result, err := h.UserInfoUsecase.LockMultipleUsers(ctx, dtoReq)
	if err != nil {
		return &authpb.LockMultipleUsersResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to lock multiple users: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.LockMultipleUsersResponse{
		Success:     result.Success,
		Message:     result.Message,
		ProfileIds:  result.ProfileIDs,
		LockType:    uint32(result.LockType),
		LockedAt:    result.LockedAt,
		LockedUntil: result.LockedUntil,
	}, nil
}

// @Summary Mở khóa nhiều user
// @Description Mở khóa nhiều user cùng lúc
// @Tags Account Management
// @Accept json
// @Produce json
// @Param body body authpb.UnlockMultipleUsersRequest true "Thông tin mở khóa nhiều user"
// @Success 200 {object} authpb.UnlockMultipleUsersResponse
// @Failure 400 {object} authpb.UnlockMultipleUsersResponse
// @Failure 500 {object} authpb.UnlockMultipleUsersResponse
// @Router /account/unlock-multiple [post]
func (h *AuthUserProfileHandler) UnlockMultipleUsers(ctx context.Context, req *authpb.UnlockMultipleUsersRequest) (*authpb.UnlockMultipleUsersResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.UnlockMultipleUsersRequest{
		ProfileIDs: req.ProfileIds,
		Reason:     req.Reason,
		UnlockedBy: req.UnlockedBy,
	}

	// Call usecase
	result, err := h.UserInfoUsecase.UnlockMultipleUsers(ctx, dtoReq)
	if err != nil {
		return &authpb.UnlockMultipleUsersResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to unlock multiple users: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.UnlockMultipleUsersResponse{
		Success:    result.Success,
		Message:    result.Message,
		ProfileIds: result.ProfileIDs,
		UnlockedAt: result.UnlockedAt,
	}, nil
}
