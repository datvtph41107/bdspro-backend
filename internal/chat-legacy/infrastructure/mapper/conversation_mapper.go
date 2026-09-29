package mapper

import (
	"chat/models"
	_utils "common/utils"
	chatpb "pb/types/chat"
	"time"
)

type ConversationMapper struct {
	backgroundMapper *BackgroundMapper
	messageMapper    *MessageMapper
}

func NewConversationMapper(
	backgroundMapper *BackgroundMapper,
	messageMapper *MessageMapper,
) *ConversationMapper {
	return &ConversationMapper{
		backgroundMapper: backgroundMapper,
		messageMapper:    messageMapper,
	}
}

func (m *ConversationMapper) ConversationToPb(conv *models.ConversationModel) *chatpb.Conversation {
	result := &chatpb.Conversation{
		ConversationId: int32(conv.ID),
		Name:           *conv.Name,
		Type:           chatpb.ConversationTypeEnum(conv.Type),
		UpdatedAt:      _utils.FormatTimeToString(&conv.UpdatedAt),
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
		result.LastMessage = m.messageMapper.MessageToPb(&conv.LatestMessage)
		result.UpdatedAt = _utils.FormatTimeToString(&conv.LatestMessage.CreatedAt)
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

func (m *ConversationMapper) ConversationToPbGetConversationById(conversation *models.ConversationModel, includeMembers bool) *chatpb.GetConversationByIdResponse {
	var participants []*chatpb.Participant
	if includeMembers {
		for _, participant := range conversation.Participants {
			participants = append(participants, &chatpb.Participant{
				UserId:    (participant.UserID),
				Role:      chatpb.ParticipantRoleEnum(participant.Role),
				JoinedAt:  participant.JoinedAt.Local().Format(time.RFC3339),
				MuteNotif: participant.MuteNotif,
			})
		}
	}
	result := &chatpb.GetConversationByIdResponse{
		ConversationId:  int32(conversation.ID),
		Name:            *conversation.Name,
		Type:            chatpb.ConversationTypeEnum(conversation.Type),
		CreatedBy:       int32(conversation.CreatedBy),
		IsBroadcast:     conversation.IsBroadcast,
		ForbidForward:   conversation.ForbidForward,
		CreatedAt:       _utils.FormatTimeToString(&conversation.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(&conversation.UpdatedAt),
		Participants:    participants,
		BackgroundImage: m.backgroundMapper.BackgroundImageToPb(conversation.BackgroundImage),
		Avatar:          conversation.Avatar,
		ReceiverId:      conversation.ReceiverId,
	}

	if conversation.BackgroundImage != nil {
		result.BackgroundImage = m.backgroundMapper.BackgroundImageToPb(conversation.BackgroundImage)
	}

	if conversation.LatestMessage.ID != 0 {
		result.LastMessage = m.messageMapper.MessageToPb(&conversation.LatestMessage)
		result.UpdatedAt = _utils.FormatTimeToString(&conversation.LatestMessage.UpdatedAt)
	}

	return result
}
