package client

import (
	"chat/env"
	"chat/internal/dto"
	"context"

	"pb/clients"
	chatpb "pb/types/chat"
)

// @bind: chat/internal/interface.IUserClient
type UserClient struct {
	Client *clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{Client: clients.MakeUserGrpcClient(env.ENV_RUNTIME)}
}

func (c *UserClient) GetUserById(ctx context.Context, id uint64) (*dto.UserProfile, error) {
	users, err := c.Client.GetProfileByIds(ctx, []uint64{id})
	if err != nil {
		return nil, err
	}

	if len(users.Profiles) == 0 {
		return nil, nil
	}

	profile := &dto.UserProfile{
		ID:           users.Profiles[0].Id,
		FullName:     users.Profiles[0].FullName,
		Avatar:       users.Profiles[0].Avatar,
		TickVerified: users.Profiles[0].TickVerified,
	}

	return profile, nil
}

func (c *UserClient) GetUserByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error) {
	user, err := c.Client.GetProfileByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	profiles := make([]*dto.UserProfile, len(user.Profiles))
	for i, profile := range user.Profiles {
		profiles[i] = &dto.UserProfile{
			ID:           profile.Id,
			FullName:     profile.FullName,
			Avatar:       profile.Avatar,
			TickVerified: profile.TickVerified,
		}
	}

	return profiles, nil
}

func (c *UserClient) MapUserProfile(ctx context.Context, pbConversations []*chatpb.Conversation) {
	userIds := make(map[uint64]struct{}, len(pbConversations))
	for _, conv := range pbConversations {
		for _, id := range conv.ParticipantIds {
			userIds[id] = struct{}{}
		}
	}

	users, err := c.Client.GetMapByIDs(ctx, userIds)
	if err != nil {
		return
	}

	for _, conv := range pbConversations {
		avatars := make([]string, 0)
		for _, id := range conv.ParticipantIds {
			if user, ok := users[id]; ok {
				avatar := user.Avatar
				if avatar != "" {
					avatars = append(avatars, avatar)
				}
			}
			if len(avatars) >= 4 {
				break
			}
		}
		conv.Avatars = avatars
	}
}

func (c *UserClient) GetAvatars(ctx context.Context, ids []uint64) ([]string, error) {
	users, err := c.Client.GetProfileByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	avatars := make([]string, 0)
	for _, user := range users.Profiles {
		avatars = append(avatars, user.Avatar)
	}

	return avatars, nil
}
