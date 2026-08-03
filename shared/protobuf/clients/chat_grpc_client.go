package clients

import (
	"common"
	_dto "common/domain/dto"
	_provider "common/domain/provider"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	chatpb "pb/types/chat"
	"time"

	"google.golang.org/grpc"
)

type ChatGrpcClient struct {
	Client chatpb.ChatServiceClient
}

func NewChatGrpcClient(env_runtime string) _provider.ChatProvider {
	chatClient := &ChatGrpcClient{}
	go func() {
		name := "chat"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8208", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := chatpb.NewChatServiceClient(conn)
		chatClient.Client = client
	}()
	return chatClient
}

// func (c *ChatGrpcClient) SendMessage(ctx context.Context, message *chatpb.Message) error {
// 	_, err := c.client.SendMessage(ctx, &chatpb.SendMessageRequest{
// 		Message: message,
// 	})
// 	return err
// }

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
