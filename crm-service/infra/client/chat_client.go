package client

import (
	"context"
	"crm/env"
	chatpb "pb/types/chat"
	sharepb "pb/types/shared"

	"pb/clients"
)

type ChatClient struct {
	*clients.ChatGrpcClient
}

func NewChatClient() *ChatClient {
	return &ChatClient{
		ChatGrpcClient: clients.NewChatGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *ChatClient) GetListParticipant(ctx context.Context, conversationId uint64) ([]*chatpb.Participant, error) {
	req := &chatpb.GetListParticipantRequest{
		ConversationId: conversationId,
		Page:           0,
		Size:           1000, // Lấy tất cả members
	}
	resp, err := c.Client.GetListParticipant(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *ChatClient) GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error) {
	req := &sharepb.IdRequest{
		Id: conversationId,
	}
	resp, err := c.Client.GetMembersOfRoom(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.ProfileIds, nil
}

// MapParticipantToContact maps participant info to contact if contact's profileId matches participant's userId
func (c *ChatClient) MapParticipantToContact(ctx context.Context, conversationId uint64, contacts []*sharepb.ContactDTO) error {
	if conversationId == 0 {
		return nil
	}

	participants, err := c.GetListParticipant(ctx, conversationId)
	if err != nil {
		// Nếu lỗi thì bỏ qua, không fail toàn bộ request
		return nil
	}

	// Tạo map từ userId -> Participant
	participantMap := make(map[uint64]*chatpb.Participant)
	for _, p := range participants {
		participantMap[p.UserId] = p
	}

	// Map participant vào contact nếu profileId match
	for _, contact := range contacts {
		if contact.ProfileId != nil {
			if participant, ok := participantMap[*contact.ProfileId]; ok {
				contact.ConversationMember = &sharepb.ConversationMember{
					UserId:    participant.UserId,
					JoinedAt:  participant.JoinedAt,
					LeftAt:    participant.LeftAt,
					MuteNotif: participant.MuteNotif,
					Role:      int32(participant.Role),
				}
			}
		}
	}

	return nil
}
