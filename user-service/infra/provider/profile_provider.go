package provider

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	"time"
	userenums "user/enums"
	"user/internal/dto"
	internalenums "user/internal/enums"
	"user/internal/interface/repo"
	"user/internal/models"
	"user/internal/utils"
)

// @bind: user/internal/interface/providers.ProfileProvider
type ProfileProvider struct {
	profileRepo repo.IProfileRepo
}

func NewProfileProvider(profileRepo repo.IProfileRepo) *ProfileProvider {
	return &ProfileProvider{profileRepo: profileRepo}
}

func (p *ProfileProvider) GetByProfileID(ctx context.Context, profileID uint64) (*dto.ProfileDTO, error) {
	profile, err := p.profileRepo.GetUserDetailByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return entityToProfileDTO(profile), nil
}

func (p *ProfileProvider) GetByProfileIDV3(ctx context.Context, profileID uint64) (*_dto.UserV3DTO, error) {
	profile, err := p.profileRepo.GetUserDetailByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return &_dto.UserV3DTO{
		ProfileID:         profile.ProfileID,
		FullName:          profile.FullName,
		Avatar:            profile.Avatar,
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),
	}, nil
}

func (p *ProfileProvider) SaveProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	entity := profileDTOToEntity(profile)
	models.SetAllVisibilityPublic(entity)
	_, err := p.profileRepo.CreateUserProfile(ctx, entity)
	return err
}

func (p *ProfileProvider) CreateProfile(ctx context.Context, authID uint64, planID *uint64, name, email, phone, avatar string) (*dto.ProfileDTO, error) {
	entity := models.MakeUserProfileEntity(authID, planID, name, email, phone, avatar)
	if authID > 0 {
		entity.ProfileID = authID
	}
	created, err := p.profileRepo.CreateUserProfile(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entityToProfileDTO(created), nil
}

func (p *ProfileProvider) UpdateProfile(ctx context.Context, profile *dto.ProfileDTO) error {
	existing, err := p.profileRepo.GetUserDetailByID(ctx, profile.ProfileID)
	if err != nil {
		return err
	}

	if profile.FullName != "" {
		existing.FullName = profile.FullName
	}
	if profile.Email != "" {
		existing.Email = profile.Email
	}
	if profile.Phone != "" {
		existing.Phone = profile.Phone
	}
	if profile.Avatar != "" {
		existing.Avatar = profile.Avatar
	}
	if profile.PlanID != nil {
		existing.PlanID = profile.PlanID
	}
	if profile.RoleType != 0 {
		existing.RoleType = userenums.ERole(profile.RoleType)
	}

	_, err = p.profileRepo.UpdateUserProfile(ctx, existing)
	return err
}

func (p *ProfileProvider) SoftDeleteProfile(ctx context.Context, profileID uint64) error {
	return p.profileRepo.DeleteUserProfile(ctx, profileID, nil, "soft delete admin profile")
}

func (p *ProfileProvider) HardDeleteProfile(ctx context.Context, profileID uint64) error {
	return p.profileRepo.DeleteUserProfile(ctx, profileID, nil, "hard delete profile")
}

func (p *ProfileProvider) CheckAccountStatus(ctx context.Context, profileID uint64) (bool, bool, error) {
	profile, err := p.profileRepo.GetUserDetailByID(ctx, profileID)
	if err != nil {
		return false, true, err
	}
	return profile.Status.IsLocked(), profile.Status.CanLogin(), nil
}

func (p *ProfileProvider) LockAccount(ctx context.Context, profileID uint64) error {
	return p.profileRepo.UpdateUserProfileStatus(ctx, profileID, _enum.EUserStatusTemporaryLocked)
}

func (p *ProfileProvider) MakeUserProfileEntity(authID uint64, planID *uint64, name, email, phone, avatar string) *dto.ProfileDTO {
	return &dto.ProfileDTO{
		ReferralCode: utils.GenerateReferral(authID),
		ProfileID:    authID,
		PlanID:       planID,
		FullName:     name,
		Email:        email,
		Phone:        phone,
		Avatar:       avatar,
	}
}

func (p *ProfileProvider) GetPlanInfo(profileID uint64) (*uint64, *time.Time) {
	profile, err := p.profileRepo.GetByProfileID(profileID)
	if err != nil || profile == nil {
		return nil, nil
	}
	return profile.PlanID, profile.PlanAt
}

func (p *ProfileProvider) GetMapProfileByIds(ctx context.Context, ids []uint64) (map[uint64]*dto.ProfileDTO, error) {
	profiles, err := p.profileRepo.GetProfileByIds(ids)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]*dto.ProfileDTO, len(profiles))
	for _, profile := range profiles {
		if profile == nil {
			continue
		}
		result[profile.ProfileID] = entityToProfileDTO(profile)
	}
	return result, nil
}

func entityToProfileDTO(profile *models.UserProfileEntity) *dto.ProfileDTO {
	if profile == nil {
		return nil
	}
	return &dto.ProfileDTO{
		ProfileID:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
		PlanID:       profile.PlanID,
		RoleType:     internalenums.ERole(profile.RoleType),
	}
}

func profileDTOToEntity(profile *dto.ProfileDTO) *models.UserProfileEntity {
	if profile == nil {
		return nil
	}
	entity := &models.UserProfileEntity{
		ProfileID:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
		PlanID:       profile.PlanID,
	}
	if profile.RoleType != 0 {
		entity.RoleType = userenums.ERole(profile.RoleType)
	}
	return entity
}
