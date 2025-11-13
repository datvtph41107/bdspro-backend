package handler

import (
	_enum "common/domain/enum"
	_utils "common/utils"
	"context"
	"log"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/dto"
	"user/internal/models"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalHandler struct {
	userpb.UnimplementedInternalUserServiceServer
	profileUsecase *usecases.ProfileUsecase
	adminUsecase   *usecases.AdminUsecase
}

func NewInternalHandler(profileUsecase *usecases.ProfileUsecase, adminUsecase *usecases.AdminUsecase) *InternalHandler {
	return &InternalHandler{
		profileUsecase: profileUsecase,
		adminUsecase:   adminUsecase,
	}
}

// CreateProfile tạo profile mới (dùng bởi auth-service)
func (h *InternalHandler) CreateProfile(ctx context.Context, req *userpb.CreateProfileRequest) (*userpb.CreateProfileResponse, error) {
	// Tạo UserProfileEntity từ request
	profile := &models.UserProfileEntity{
		// ProfileID:    req.ProfileId,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Avatar:       req.Avatar,
		ReferralCode: req.ReferralCode,
	}

	// Gọi repo để tạo profile
	createdProfile, err := h.profileUsecase.ProfileRepo.CreateUserProfile(ctx, profile)
	if err != nil {
		log.Printf("Failed to create profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to create profile: %v", err)
	}

	// Trả về response
	return &userpb.CreateProfileResponse{
		ProfileId:    createdProfile.ProfileID,
		FullName:     createdProfile.FullName,
		Email:        createdProfile.Email,
		Phone:        createdProfile.Phone,
		Avatar:       createdProfile.Avatar,
		ReferralCode: createdProfile.ReferralCode,
		PlanId:       req.PlanId,
	}, nil
}

// UpdateProfile cập nhật profile (dùng bởi auth-service)
func (h *InternalHandler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*sharepb.Empty, error) {
	// Lấy profile hiện tại
	existingProfile, err := h.profileUsecase.ProfileRepo.GetUserDetailByID(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Failed to get profile: %v", err)
		return nil, status.Errorf(codes.NotFound, "Profile not found")
	}

	// Cập nhật các field nếu có
	if req.FullName != nil {
		existingProfile.FullName = *req.FullName
	}
	if req.Email != nil {
		existingProfile.Email = *req.Email
	}
	if req.Phone != nil {
		existingProfile.Phone = *req.Phone
	}
	if req.Avatar != nil {
		existingProfile.Avatar = *req.Avatar
	}

	// Gọi repo để cập nhật
	_, err = h.profileUsecase.ProfileRepo.UpdateUserProfile(ctx, existingProfile)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to update profile: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// SoftDeleteProfile xóa mềm profile (dùng bởi auth-service)
func (h *InternalHandler) SoftDeleteProfile(ctx context.Context, req *userpb.SoftDeleteProfileRequest) (*sharepb.Empty, error) {
	// Lấy thông tin người xóa từ context (nếu có)
	var deletedBy *uint64
	if profileID := _utils.GetProfileIdWithContext(ctx); profileID > 0 {
		deletedBy = &profileID
	}

	reason := "Xóa từ auth service"

	// Xóa profile
	err := h.profileUsecase.ProfileRepo.DeleteUserProfile(ctx, req.ProfileId, deletedBy, reason)
	if err != nil {
		log.Printf("Failed to soft delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to soft delete profile: %v", err)
	}

	// Xóa TẤT CẢ auth method của user
	err = h.profileUsecase.AuthProvider.DeleteAllAuthMethodsByUserId(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Warning: Failed to delete all auth methods for user %d: %v", req.ProfileId, err)
		// Không return error vì profile đã được xóa
	}

	return &sharepb.Empty{}, nil
}

// HardDeleteProfile xóa vĩnh viễn profile khỏi database (dùng bởi auth-service)
func (h *InternalHandler) HardDeleteProfile(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	// Lấy thông tin người xóa từ context (nếu có)
	var deletedBy *uint64
	if profileID := _utils.GetProfileIdWithContext(ctx); profileID > 0 {
		deletedBy = &profileID
	}

	reason := "Hard delete từ auth service"

	// Xóa profile
	err := h.profileUsecase.ProfileRepo.DeleteUserProfile(ctx, req.Id, deletedBy, reason)
	if err != nil {
		log.Printf("Failed to hard delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hard delete profile: %v", err)
	}

	// Xóa TẤT CẢ auth method của user
	err = h.profileUsecase.AuthProvider.DeleteAllAuthMethodsByUserId(ctx, req.Id)
	if err != nil {
		log.Printf("Warning: Failed to delete all auth methods for user %d: %v", req.Id, err)
		// Không return error vì profile đã được xóa
	}

	return &sharepb.Empty{}, nil
}

// GetProfileByID lấy profile theo ID
func (h *InternalHandler) GetProfileByID(ctx context.Context, req *userpb.GetProfileByIDRequest) (*userpb.ProfileResponse, error) {
	profile, err := h.profileUsecase.ProfileRepo.GetUserDetailByID(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Failed to get profile by ID: %v", err)
		return nil, status.Errorf(codes.NotFound, "Profile not found")
	}

	return &userpb.ProfileResponse{
		ProfileId:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
		CreatedAt:    _utils.FormatTimeToString(profile.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(profile.UpdatedAt),
	}, nil
}

func (h *InternalHandler) LockUser(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.Empty, error) {
	err := h.adminUsecase.LockUser(ctx, &dto.LockUserRequest{
		ProfileID: req.Id,
		Reason:    req.Text,
		LockType:  _enum.EUserStatusTemporaryLocked,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to lock account: %v", err)
	}

	return &sharepb.Empty{}, nil
}
