package clients

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"log"

	chatpb "pb/types/chat"

	"google.golang.org/grpc"
)

type ChatGrpcClient struct {
	Client chatpb.ChatServiceClient
}

func BindChatGrpcClient(conn grpc.ClientConnInterface) *ChatGrpcClient {
	return &ChatGrpcClient{Client: chatpb.NewChatServiceClient(conn)}
}

func (c *ChatGrpcClient) CreateConversation(ctx context.Context, name string, createdBy uint64, members []uint64) (uint64, error) {
	if c.Client == nil {
		log.Printf("Chat service client is nil")
		return 0, nil
	}

	membersInt32 := make([]int32, len(members))
	for i, member := range members {
		membersInt32[i] = int32(member)
	}

	response, err := c.Client.CreateConversation(ctx, &chatpb.CreateConversationRequest{
		Name:      name,
		Type:      chatpb.ConversationTypeEnum_ConversationTypeGroup,
		CreatedBy: int32(createdBy),
		Members:   membersInt32,
	})
	if err != nil {
		log.Printf("Failed to create conversation: %v", err)
		return 0, err
	}

	return uint64(response.ConversationId), nil
}

func (c *ChatGrpcClient) MapParticipantToContact(ctx context.Context, conversationId uint64, contacts []_dto.ContactInfoV3DTO) error {
	if conversationId == 0 {
		return nil
	}

	participants, err := c.Client.GetListParticipant(ctx, &chatpb.GetListParticipantRequest{
		ConversationId: conversationId,
	})
	if err != nil {
		// Nếu lỗi thì bỏ qua, không fail toàn bộ request
		return nil
	}

	// Tạo map từ userId -> Participant
	participantMap := make(map[uint64]*chatpb.Participant)
	for _, p := range participants.Data {
		participantMap[p.UserId] = p
	}

	// Map participant vào contact nếu profileId match
	// for _, contact := range contacts {
	// 	if contact.ContactId != nil {
	// 		if participant, ok := participantMap[*contact.ContactId]; ok {
	// 			contact.ConversationMember = &sharepb.ConversationMember{
	// 				UserId:    participant.UserId,
	// 				JoinedAt:  participant.JoinedAt,
	// 				LeftAt:    participant.LeftAt,
	// 				MuteNotif: participant.MuteNotif,
	// 				Role:      int32(participant.Role),
	// 			}
	// 		}
	// 	}
	// }

	return nil
}

func (c *ChatGrpcClient) SendProductToReceiver(ctx context.Context, productID uint64, content string, receiverID *uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	_, err := c.Client.SendToReceiver(ctx, &chatpb.SendMessageRequest{
		ConversationId: 1,
		SenderId:       int32(profileId),
		Content:        content,
		ContentType:    "product",
		FileUrl:        nil,
		ReplyId:        nil,
		ExtraId:        &productID,
		ReceiverId:     receiverID,
	})
	return err
}
