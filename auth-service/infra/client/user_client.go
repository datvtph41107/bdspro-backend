package client

import (
	"auth/env"
	"auth/internal/dto"
	"auth/internal/utils"
	"context"
	"log"
	"pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"
)

// @bind: auth/internal/interface/provider.ProfileProvider
type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *UserClient) GetByProfileID(ctx context.Context, profileID uint64) (*dto.ProfileDTO, error) {
	if c.UserGrpcClient == nil {
		log.Printf("Internal user service client is nil")
		return nil, nil
	}

	resp, err := c.UserGrpcClient.InternalClient.GetProfileByID(ctx, &userpb.GetProfileByIDRequest{
		ProfileId: profileID,
	})
	if err != nil {
		log.Printf("Failed to get profile by ID: %v", err)
		return nil, err
	}

	return &dto.ProfileDTO{
		ProfileID:    resp.ProfileId,
		FullName:     resp.FullName,
		Email:        resp.Email,
		Phone:        resp.Phone,
		Avatar:       resp.Avatar,
		ReferralCode: resp.ReferralCode,
		// PlanID:       resp.PlanId,
	}, nil
}

func (c *UserClient) SaveProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	if c.UserGrpcClient == nil || c.UserGrpcClient.InternalClient == nil {
		log.Printf("Internal user service client is nil")
		return nil
	}

	// Gọi sang user-service để tạo profile
	_, err := c.UserGrpcClient.InternalClient.CreateProfile(ctx, &userpb.CreateProfileRequest{
		ProfileId:    profile.ProfileID,
		PlanId:       profile.PlanID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
	})

	if err != nil {
		log.Printf("Failed to create profile: %v", err)
		return err
	}

	return nil
}

func (c *UserClient) CreateProfile(ctx context.Context, authID uint64, planID *uint64, name, email, phone, avatar string) (*dto.ProfileDTO, error) {
	if c.UserGrpcClient == nil || c.UserGrpcClient.InternalClient == nil {
		log.Printf("Internal user service client is nil")
		return nil, nil
	}

	referralCode := utils.GenerateReferral(authID)

	resp, err := c.UserGrpcClient.InternalClient.CreateProfile(ctx, &userpb.CreateProfileRequest{
		// ProfileId:    authID,
		PlanId:       planID,
		FullName:     name,
		Email:        email,
		Phone:        phone,
		Avatar:       avatar,
		ReferralCode: referralCode,
	})

	if err != nil {
		log.Printf("Failed to create profile: %v", err)
		return nil, err
	}

	return &dto.ProfileDTO{
		ProfileID:    resp.ProfileId,
		FullName:     resp.FullName,
		Email:        resp.Email,
		Phone:        resp.Phone,
		Avatar:       resp.Avatar,
		ReferralCode: resp.ReferralCode,
		PlanID:       resp.PlanId,
		DeletedAt:    nil,
	}, nil
}

func (c *UserClient) UpdateProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	if c.UserGrpcClient == nil || c.UserGrpcClient.InternalClient == nil {
		log.Printf("Internal user service client is nil")
		return nil
	}

	_, err := c.UserGrpcClient.InternalClient.UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		ProfileId: profile.ProfileID,
		FullName:  &profile.FullName,
		Email:     &profile.Email,
		Phone:     &profile.Phone,
		Avatar:    &profile.Avatar,
		PlanId:    profile.PlanID,
	})

	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		return err
	}

	return nil
}

func (c *UserClient) SoftDeleteProfile(ctx context.Context, profileID uint64) error {
	if c.UserGrpcClient == nil || c.UserGrpcClient.InternalClient == nil {
		log.Printf("Internal user service client is nil")
		return nil
	}

	_, err := c.UserGrpcClient.InternalClient.SoftDeleteProfile(ctx, &userpb.SoftDeleteProfileRequest{
		ProfileId: profileID,
	})

	if err != nil {
		log.Printf("Failed to soft delete profile: %v", err)
		return err
	}

	return nil
}

func (c *UserClient) HardDeleteProfile(ctx context.Context, profileID uint64) error {
	if c.UserGrpcClient == nil || c.UserGrpcClient.InternalClient == nil {
		log.Printf("Internal user service client is nil")
		return nil
	}

	_, err := c.UserGrpcClient.InternalClient.HardDeleteProfile(ctx, &sharepb.IdRequest{
		Id: profileID,
	})

	if err != nil {
		log.Printf("Failed to hard delete profile: %v", err)
		return err
	}

	return nil
}

func (c *UserClient) MakeUserProfileEntity(authID uint64, planID *uint64, name, email, phone, avatar string) *dto.ProfileDTO {
	return &dto.ProfileDTO{
		ReferralCode: utils.GenerateReferral(authID),
		ProfileID:    authID,
		PlanID:       planID,
		FullName:     name,
		Email:        email,
		Phone:        phone,
		Avatar:       avatar,
		DeletedAt:    nil,
	}
}

func (c *UserClient) GetPlanInfo(profileID uint64) (*uint64, *time.Time) {
	return nil, nil
}

func (c *UserClient) CheckAccountStatus(ctx context.Context, profileID uint64) (isLocked bool, canLogin bool, err error) {
	if c.UserGrpcClient == nil {
		log.Printf("User service client is nil")
		return false, true, nil
	}

	resp, err := c.UserGrpcClient.Client.CheckAccountStatus(ctx, &userpb.CheckAccountStatusRequest{
		ProfileId: profileID,
	})
	if err != nil {
		log.Printf("Failed to check account status: %v", err)
		return false, true, err
	}

	return resp.IsLocked, resp.CanLogin, nil
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, ids []uint64) (map[uint64]*dto.ProfileDTO, error) {
	profileIdSet := make(map[uint64]struct{})
	for _, id := range ids {
		profileIdSet[id] = struct{}{}
	}
	result, err := c.UserGrpcClient.GetMapByIDs(ctx, profileIdSet)
	if err != nil {
		return nil, err
	}

	profileMap := make(map[uint64]*dto.ProfileDTO, len(result))
	for _, profile := range result {
		profileMap[profile.Id] = &dto.ProfileDTO{
			ProfileID: profile.Id,
			FullName:  profile.FullName,
			// Email:        profile.Email,
			Phone:  profile.Phone,
			Avatar: profile.Avatar,
			// ReferralCode: profile.ReferralCode,
			// PlanID:       profile.PlanId,
			// DeletedAt:    profile.DeletedAt,
		}
	}

	return profileMap, nil
}
