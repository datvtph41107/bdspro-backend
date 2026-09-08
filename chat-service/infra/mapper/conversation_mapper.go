package mapper

import (
	"chat/internal/domain"
	_utils "common/utils"
	chatpb "pb/types/chat"
	"time"
)

type ConversationMapper struct {
	backgroundMapper *BackgroundMapper
}

func NewConversationMapper(
	backgroundMapper *BackgroundMapper,
) *ConversationMapper {
	return &ConversationMapper{
		backgroundMapper: backgroundMapper,
	}
}

func (m *ConversationMapper) ConversationToPb(conv *domain.Conversation) *chatpb.Conversation {
	updatedAtStr := ""
	if conv.UpdatedAt != nil {
		updatedAtStr = _utils.FormatTimeToString(conv.UpdatedAt)
	}

	result := &chatpb.Conversation{
		ConversationId: int32(conv.ID),
		Name:           *conv.Name,
		Type:           chatpb.ConversationTypeEnum(conv.Type),
		UpdatedAt:      updatedAtStr,
		Avatar:         conv.Avatar,
		UnreadCount:    int32(conv.UnreadCount),
		ParticipantIds: conv.ParticipantIDs,
		ReceiverId:     conv.ReceiverId,
		CreatedBy:      conv.CreatedBy,
	}

	if conv.BackgroundImage != nil {
		result.BackgroundImage = m.backgroundMapper.BackgroundImageToPb(conv.BackgroundImage)
	}

	if conv.LatestMessage.ID != 0 {
		latestCreatedAtStr := ""
		if conv.LatestMessage.CreatedAt != nil {
			latestCreatedAtStr = conv.LatestMessage.CreatedAt.Local().Format(time.RFC3339)
		}
		result.LastMessage = &chatpb.Message{
			MessageId:   int32(conv.LatestMessage.ID),
			SenderId:    int32(conv.LatestMessage.SenderID),
			Content:     conv.LatestMessage.Content,
			ContentType: string(conv.LatestMessage.ContentType),
			FileUrl:     conv.LatestMessage.FileURL,
			CreatedAt:   latestCreatedAtStr,
		}
		if conv.LatestMessage.CreatedAt != nil {
			result.UpdatedAt = _utils.FormatTimeToString(conv.LatestMessage.CreatedAt)
		}
	}

	return result
}

func (m *ConversationMapper) ConversationsToPb(convs []*domain.Conversation) []*chatpb.Conversation {
	result := make([]*chatpb.Conversation, len(convs))
	for i, conv := range convs {
		result[i] = m.ConversationToPb(conv)
	}
	return result
}

func (m *ConversationMapper) ConversationToPbGetConversationById(conversation *domain.Conversation, includeMembers bool) *chatpb.GetConversationByIdResponse {
	var participants []*chatpb.Participant
	if includeMembers {
		for _, participant := range conversation.Participants {
			participants = append(participants, &chatpb.Participant{
				UserId:    participant.UserID,
				Role:      chatpb.ParticipantRoleEnum(participant.Role),
				JoinedAt:  participant.JoinedAt.Format(time.RFC3339),
				MuteNotif: participant.MuteNotif,
			})
		}
	}
	createdAtStr := ""
	if conversation.CreatedAt != nil {
		createdAtStr = conversation.CreatedAt.Local().Format(time.RFC3339)
	}
	updatedAtStr := ""
	if conversation.UpdatedAt != nil {
		updatedAtStr = conversation.UpdatedAt.Local().Format(time.RFC3339)
	}

	result := &chatpb.GetConversationByIdResponse{
		ConversationId:  int32(conversation.ID),
		Name:            *conversation.Name,
		Type:            chatpb.ConversationTypeEnum(conversation.Type),
		CreatedBy:       int32(conversation.CreatedBy),
		IsBroadcast:     conversation.IsBroadcast,
		ForbidForward:   conversation.ForbidForward,
		CreatedAt:       createdAtStr,
		UpdatedAt:       updatedAtStr,
		Participants:    participants,
		BackgroundImage: m.backgroundMapper.BackgroundImageToPb(conversation.BackgroundImage),
		Avatar:          conversation.Avatar,
		ReceiverId:      conversation.ReceiverId,
	}

	if conversation.BackgroundImage != nil {
		result.BackgroundImage = m.backgroundMapper.BackgroundImageToPb(conversation.BackgroundImage)
	}

	return result
}
