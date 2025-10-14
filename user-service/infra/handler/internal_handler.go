package handler

import (
	_utils "common/utils"
	"context"
	"log"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/models"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalHandler struct {
	userpb.UnimplementedInternalUserServiceServer
	profileUsecase *usecases.ProfileUsecase
}

func NewInternalHandler(profileUsecase *usecases.ProfileUsecase) *InternalHandler {
	return &InternalHandler{
		profileUsecase: profileUsecase,
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
	err := h.profileUsecase.ProfileRepo.DeleteUserProfile(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Failed to soft delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to soft delete profile: %v", err)
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
