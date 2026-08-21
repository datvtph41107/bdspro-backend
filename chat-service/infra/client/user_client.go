package client

import (
	"chat/internal/dto"
	_utils "common/utils"
	"context"
	"fmt"
	"pb/clients"
	chatpb "pb/types/chat"
)

// @bind: chat/internal/interface.IUserClient
type UserClient struct {
	Client *clients.UserGrpcClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{Client: rpcClient}
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
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	userIds := make(map[uint64]struct{}, len(pbConversations))
	for _, conv := range pbConversations {
		if conv.ReceiverId != nil {
			fmt.Println("conv.ReceiverId", *conv.ReceiverId, currentUserId, conv.CreatedBy)
			if *conv.ReceiverId == currentUserId {
				userIds[conv.CreatedBy] = struct{}{}
			} else {
				userIds[*conv.ReceiverId] = struct{}{}
			}
		}
		for _, id := range conv.ParticipantIds {
			if id == currentUserId {
				continue
			}
			userIds[id] = struct{}{}
		}
	}

	users, err := c.Client.GetMapByIDs(ctx, userIds)
	if err != nil {
		return
	}

	for _, conv := range pbConversations {
		if conv.ReceiverId != nil {
			if *conv.ReceiverId == currentUserId {
				if user, ok := users[conv.CreatedBy]; ok {
					conv.Avatar = user.Avatar
					conv.Name = user.FullName
				}
			} else {
				if user, ok := users[*conv.ReceiverId]; ok {
					conv.Avatar = user.Avatar
					conv.Name = user.FullName
				}
			}
		}
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

func (c *UserClient) MapAvatarAndNameToPb(ctx context.Context, conversation *chatpb.GetConversationByIdResponse) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	if conversation.ReceiverId != nil {
		if *conversation.ReceiverId != currentUserId {
			user, _ := c.GetUserById(ctx, *conversation.ReceiverId)
			if user != nil {
				conversation.Avatar = user.Avatar
				conversation.Name = user.FullName
			}
		} else {
			user, _ := c.GetUserById(ctx, uint64(conversation.CreatedBy))
			if user != nil {
				conversation.Avatar = user.Avatar
				conversation.Name = user.FullName
			}
		}
	} else if len(conversation.Participants) > 0 {
		userIds := make([]uint64, 0)
		for _, participant := range conversation.Participants {
			userIds = append(userIds, participant.UserId)
		}
		avatars, _ := c.GetAvatars(ctx, userIds)
		conversation.Avatars = avatars
	}
}

func (c *UserClient) MapAvatarAndNameToMessagePbList(ctx context.Context, messages []*chatpb.Message) {
	userIds := make(map[uint64]struct{}, len(messages))
	for _, message := range messages {
		userIds[uint64(message.SenderId)] = struct{}{}
	}
	users, _ := c.Client.GetMapByIDs(ctx, userIds)
	for _, message := range messages {
		if user, ok := users[uint64(message.SenderId)]; ok {
			message.Avatar = &user.Avatar
			message.FullName = &user.FullName
		}
	}
}
