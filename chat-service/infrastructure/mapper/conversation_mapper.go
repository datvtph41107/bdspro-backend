package mapper

import (
	"chat/models"
	chatpb "pb/types/chat"
	"time"
)

type ConversationMapper struct {
}

func NewConversationMapper() *ConversationMapper {
	return &ConversationMapper{}
}

func (m *ConversationMapper) ConversationToPb(conv *models.ConversationModel) *chatpb.Conversation {
	result := &chatpb.Conversation{
		ConversationId:  int32(conv.ID),
		Name:            *conv.Name,
		Type:            chatpb.ConversationTypeEnum(conv.Type),
		UpdatedAt:       conv.UpdatedAt.Format(time.DateTime),
		Avatar:          conv.Avatar,
		UnreadCount:     int32(conv.UnreadCount),
		BackgroundImage: conv.BackgroundImage,
		ParticipantIds:  conv.ParticipantIDs,
	}

	if conv.LatestMessage.ID != 0 {
		result.LastMessage = &chatpb.Message{
			MessageId:   int32(conv.LatestMessage.ID),
			SenderId:    int32(conv.LatestMessage.SenderID),
			Content:     conv.LatestMessage.Content,
			ContentType: string(conv.LatestMessage.ContentType),
			FileUrl:     conv.LatestMessage.FileURL,
			CreatedAt:   conv.LatestMessage.CreatedAt.Local().Format(time.RFC3339),
		}
	}

	return result
}

func (m *ConversationMapper) ConversationsToPb(convs []*models.ConversationModel) []*chatpb.Conversation {
	result := make([]*chatpb.Conversation, len(convs))
	for i, conv := range convs {
		result[i] = m.ConversationToPb(conv)
	}
	return result
}
