package client

import (
	"auth/internal/dto"
	"auth/internal/utils"
	"context"
	"time"
)

// @bind: auth/internal/interface/provider.ProfileProvider
type UserClient struct {
}

func NewUserClient() *UserClient {
	return &UserClient{}
}

func (c *UserClient) GetByProfileID(ctx context.Context, profileID uint64) (*dto.ProfileDTO, error) {
	return nil, nil
}

func (c *UserClient) SaveProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	return nil
}

func (c *UserClient) CreateProfile(ctx context.Context, authID uint64, planID *uint64, name, email, phone, avatar string) (*dto.ProfileDTO, error) {
	return nil, nil
}

func (c *UserClient) UpdateProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	return nil
}

func (c *UserClient) SoftDeleteProfile(ctx context.Context, profileID uint64) error {
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
