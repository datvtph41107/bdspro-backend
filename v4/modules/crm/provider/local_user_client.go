package provider

import (
	"context"
	"strings"

	repo "mdn4/internal/interface/repo"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/user/domain/user"
	sharepb "pb/types/shared"
)

// LocalUserClient resolves profiles from the v4 user DB (in-process, no gRPC).
// @bind: mdn4/modules/crm/provider.UserClient
type LocalUserClient struct {
	profiles repo.IProfileRepo
}

func NewLocalUserClient(profiles repo.IProfileRepo) *LocalUserClient {
	return &LocalUserClient{profiles: profiles}
}

func (c *LocalUserClient) GetProfileByPhones(ctx context.Context, phones []string) (map[string]*domain.Profile, error) {
	_ = ctx
	out := make(map[string]*domain.Profile)
	if c.profiles == nil || len(phones) == 0 {
		return out, nil
	}

	trimmed := make([]string, 0, len(phones))
	for _, p := range phones {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		trimmed = append(trimmed, t)
	}
	if len(trimmed) == 0 {
		return out, nil
	}

	entities, err := c.profiles.GetProfileByPhones(trimmed)
	if err != nil {
		return nil, err
	}
	for _, u := range entities {
		if u == nil || u.ProfileID == 0 {
			continue
		}
		prof := userProfileToCRM(u)
		out[u.Phone] = prof
		if t := strings.TrimSpace(u.Phone); t != "" && t != u.Phone {
			out[t] = prof
		}
	}

	for _, p := range phones {
		key := strings.TrimSpace(p)
		if key == "" || out[key] != nil {
			continue
		}
		u, err := c.profiles.GetByPhone(key)
		if err != nil || u == nil || u.ProfileID == 0 {
			continue
		}
		prof := userProfileToCRM(u)
		out[key] = prof
		if p != key {
			out[p] = prof
		}
	}
	return out, nil
}

func (c *LocalUserClient) GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error) {
	_ = ctx
	if c.profiles == nil || profileId == 0 {
		return &domain.Profile{}, nil
	}
	u, err := c.profiles.GetByProfileID(profileId)
	if err != nil || u == nil {
		return nil, err
	}
	p := userProfileToCRM(u)
	p.Phone = ""
	return p, nil
}

func (c *LocalUserClient) GetProfileWithPhoneById(ctx context.Context, profileId uint64) (*domain.Profile, error) {
	_ = ctx
	if c.profiles == nil || profileId == 0 {
		return &domain.Profile{}, nil
	}
	u, err := c.profiles.GetByProfileID(profileId)
	if err != nil || u == nil {
		return nil, err
	}
	return userProfileToCRM(u), nil
}

func (c *LocalUserClient) ValidateProfileIds(ctx context.Context, profileIds []uint64) (bool, error) {
	_ = ctx
	if len(profileIds) == 0 {
		return true, nil
	}
	if c.profiles == nil {
		return false, nil
	}
	ids := make([]uint64, 0, len(profileIds))
	for _, id := range profileIds {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return true, nil
	}
	list, err := c.profiles.GetProfileByIds(ids)
	if err != nil {
		return false, err
	}
	return len(list) == len(ids), nil
}

func (c *LocalUserClient) GetProfileByIds(ctx context.Context, ids []uint64) ([]*sharepb.ProfileItem, error) {
	_ = ctx
	if c.profiles == nil || len(ids) == 0 {
		return nil, nil
	}
	list, err := c.profiles.GetProfileByIds(ids)
	if err != nil {
		return nil, err
	}
	out := make([]*sharepb.ProfileItem, 0, len(list))
	for _, u := range list {
		if u == nil {
			continue
		}
		out = append(out, &sharepb.ProfileItem{
			Id:           u.ProfileID,
			FullName:     u.FullName,
			Avatar:       u.Avatar,
			Phone:        u.Phone,
			TickVerified: u.TickVerified,
		})
	}
	return out, nil
}

func userProfileToCRM(u *user.UserProfileEntity) *domain.Profile {
	if u == nil {
		return &domain.Profile{}
	}
	return &domain.Profile{
		ProfileId:    u.ProfileID,
		Phone:        u.Phone,
		FullName:     u.FullName,
		Avatar:       u.Avatar,
		TickVerified: u.TickVerified,
	}
}
