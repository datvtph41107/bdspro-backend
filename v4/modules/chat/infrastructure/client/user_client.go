package client

import (
	_utils "common/utils"
	"context"
	"fmt"
	"mdn4/modules/chat/internal/dto"
	"mdn4/modules/chat/models"
	chatpb "pb/types/chat"
	sharepb "pb/types/shared"
)

// ProfileByIDsFunc resolves profiles through the in-process v4 User module.
type ProfileByIDsFunc func(context.Context, []uint64) (map[uint64]*dto.UserProfile, error)

type UserClient struct {
	localProfileLookup ProfileByIDsFunc
}

func NewUserClient(lookup ProfileByIDsFunc) *UserClient {
	return &UserClient{localProfileLookup: lookup}
}

func (c *UserClient) GetUserById(ctx context.Context, id uint64) (*dto.UserProfile, error) {
	profiles, err := c.GetUserByIds(ctx, []uint64{id})
	if err != nil || len(profiles) == 0 {
		return nil, err
	}
	return profiles[0], nil
}

func (c *UserClient) GetUserByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error) {
	profiles, err := c.lookup(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.UserProfile, 0, len(ids))
	for _, id := range ids {
		if profile := profiles[id]; profile != nil {
			result = append(result, profile)
		}
	}
	return result, nil
}

func (c *UserClient) ProfilesByIDs(ctx context.Context, ids []uint64) (map[uint64]*sharepb.ProfileItem, error) {
	profiles, err := c.lookup(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make(map[uint64]*sharepb.ProfileItem, len(profiles))
	for id, profile := range profiles {
		if profile != nil {
			result[id] = &sharepb.ProfileItem{Id: id, FullName: profile.FullName, Avatar: profile.Avatar, TickVerified: profile.TickVerified}
		}
	}
	return result, nil
}

func (c *UserClient) MapUserProfile(ctx context.Context, conversations []*chatpb.Conversation) error {
	ids := make([]uint64, 0)
	seen := make(map[uint64]struct{})
	for _, conv := range conversations {
		if conv == nil {
			continue
		}
		appendUniqueID(&ids, seen, conv.CreatedBy)
		if conv.ReceiverId != nil {
			appendUniqueID(&ids, seen, *conv.ReceiverId)
		}
		for _, id := range conv.ParticipantIds {
			appendUniqueID(&ids, seen, id)
		}
	}
	users, err := c.ProfilesByIDs(ctx, ids)
	if err != nil {
		return err
	}
	currentUserID := _utils.GetProfileIdWithContext(ctx)
	for _, conv := range conversations {
		if conv == nil {
			continue
		}
		peerID := uint64(0)
		if conv.ReceiverId != nil {
			peerID = *conv.ReceiverId
			if peerID == currentUserID {
				peerID = conv.CreatedBy
			}
		} else if conv.Type == chatpb.ConversationTypeEnum_ConversationTypePrivate {
			for _, id := range conv.ParticipantIds {
				if id != currentUserID {
					peerID = id
					break
				}
			}
		}
		if peer := users[peerID]; peer != nil && conv.Type == chatpb.ConversationTypeEnum_ConversationTypePrivate {
			conv.Avatar = peer.Avatar
			conv.Name = peer.FullName
		}
		conv.Avatars = make([]string, 0, 4)
		for _, id := range conv.ParticipantIds {
			if profile := users[id]; profile != nil && profile.Avatar != "" {
				conv.Avatars = append(conv.Avatars, profile.Avatar)
				if len(conv.Avatars) == 4 {
					break
				}
			}
		}
	}
	return nil
}

func (c *UserClient) MapAvatarAndNameToPb(ctx context.Context, conversation *chatpb.GetConversationByIdResponse) error {
	if conversation == nil {
		return nil
	}
	ids := make([]uint64, 0, len(conversation.Participants)+2)
	seen := make(map[uint64]struct{})
	appendUniqueID(&ids, seen, uint64(conversation.CreatedBy))
	if conversation.ReceiverId != nil {
		appendUniqueID(&ids, seen, *conversation.ReceiverId)
	}
	for _, participant := range conversation.Participants {
		if participant != nil {
			appendUniqueID(&ids, seen, participant.UserId)
		}
	}
	users, err := c.ProfilesByIDs(ctx, ids)
	if err != nil {
		return err
	}
	currentUserID := _utils.GetProfileIdWithContext(ctx)
	peerID := uint64(0)
	if conversation.ReceiverId != nil {
		peerID = *conversation.ReceiverId
		if peerID == currentUserID {
			peerID = uint64(conversation.CreatedBy)
		}
	} else if conversation.Type == chatpb.ConversationTypeEnum_ConversationTypePrivate {
		for _, participant := range conversation.Participants {
			if participant != nil && participant.UserId != currentUserID {
				peerID = participant.UserId
				break
			}
		}
	}
	if peer := users[peerID]; peer != nil && conversation.Type == chatpb.ConversationTypeEnum_ConversationTypePrivate {
		conversation.Avatar = peer.Avatar
		conversation.Name = peer.FullName
	}
	conversation.Avatars = make([]string, 0, 4)
	for _, participant := range conversation.Participants {
		if participant == nil {
			continue
		}
		profile := users[participant.UserId]
		if profile == nil {
			continue
		}
		participant.Profile = profile
		if profile.Avatar != "" && len(conversation.Avatars) < 4 {
			conversation.Avatars = append(conversation.Avatars, profile.Avatar)
		}
	}
	return nil
}

func (c *UserClient) MapAvatarAndNameToMessagePbList(ctx context.Context, messages []*chatpb.Message) error {
	ids := make([]uint64, 0, len(messages))
	seen := make(map[uint64]struct{})
	for _, message := range messages {
		if message != nil {
			appendUniqueID(&ids, seen, uint64(message.SenderId))
			if message.ContentType == string(models.ContentTypeContact) && message.ExtraId != nil {
				appendUniqueID(&ids, seen, *message.ExtraId)
			}
		}
	}
	users, err := c.ProfilesByIDs(ctx, ids)
	if err != nil {
		return err
	}
	for _, message := range messages {
		if message == nil {
			continue
		}
		if user := users[uint64(message.SenderId)]; user != nil {
			message.Avatar = &user.Avatar
			message.FullName = &user.FullName
			if message.ContentType == string(models.ContentTypeContact) && message.ExtraId != nil {
				message.ContactUser = users[*message.ExtraId]
			}
		}
	}
	return nil
}

func (c *UserClient) GetAvatars(ctx context.Context, ids []uint64) ([]string, error) {
	users, err := c.lookup(ctx, ids)
	if err != nil {
		return nil, err
	}
	avatars := make([]string, 0, len(ids))
	for _, id := range ids {
		if user := users[id]; user != nil {
			avatars = append(avatars, user.Avatar)
		}
	}
	return avatars, nil
}

func (c *UserClient) lookup(ctx context.Context, ids []uint64) (map[uint64]*dto.UserProfile, error) {
	if len(ids) == 0 {
		return map[uint64]*dto.UserProfile{}, nil
	}
	if c.localProfileLookup == nil {
		return nil, fmt.Errorf("v4 user profile lookup is not configured")
	}
	return c.localProfileLookup(ctx, ids)
}

func appendUniqueID(ids *[]uint64, seen map[uint64]struct{}, id uint64) {
	if id == 0 {
		return
	}
	if _, ok := seen[id]; ok {
		return
	}
	seen[id] = struct{}{}
	*ids = append(*ids, id)
}
