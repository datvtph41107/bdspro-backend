package mapper

import (
	"chat/models"
	_utils "common/utils"
	chatpb "pb/types/chat"
	"strconv"
	"time"
)

type MessageMapper struct{}

func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

func (m *MessageMapper) MessageToPb(msg *models.MessageModel) *chatpb.Message {
	if msg == nil {
		return nil
	}

	if msg.Recall {
		return &chatpb.Message{
			MessageId: int32(msg.ID),
			SenderId:  int32(msg.SenderID),
			IsRecall:  true,
			CreatedAt: _utils.FormatTimeToString(&msg.CreatedAt),
		}
	}

	pb := &chatpb.Message{
		MessageId:   int32(msg.ID),
		SenderId:    int32(msg.SenderID),
		Content:     msg.Content,
		ContentType: string(msg.ContentType),
		FileUrl:     msg.FileURL,
		RoomId:      "",
		CreatedAt:   _utils.FormatTimeToString(&msg.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(&msg.UpdatedAt),
		IsRecall:    msg.Recall,
		ExtraId:     msg.ExtraId,
	}
	if msg.ConversationID != 0 {
		pb.RoomId = strconv.FormatUint(msg.ConversationID, 10)
	}

	if msg.ReplyMessage != nil {
		pb.ReplyMessage = &chatpb.Message{
			MessageId:   int32(msg.ReplyMessage.ID),
			SenderId:    int32(msg.ReplyMessage.SenderID),
			Content:     msg.ReplyMessage.Content,
			ContentType: string(msg.ReplyMessage.ContentType),
			FileUrl:     msg.ReplyMessage.FileURL,
			RoomId:      strconv.FormatUint(msg.ReplyMessage.ConversationID, 10),
			CreatedAt:   _utils.FormatTimeToString(&msg.ReplyMessage.CreatedAt),
			UpdatedAt:   _utils.FormatTimeToString(&msg.ReplyMessage.UpdatedAt), // THÊM
			IsRecall:    msg.ReplyMessage.Recall,
			ExtraId:     msg.ReplyMessage.ExtraId,
		}

		if msg.ReplyMessage.Reactions != nil {
			pb.ReplyMessage.Reactions = make([]*chatpb.MessageReaction, len(msg.ReplyMessage.Reactions))
			for i, r := range msg.ReplyMessage.Reactions {
				pb.ReplyMessage.Reactions[i] = &chatpb.MessageReaction{
					UserId:    int32(r.UserID),
					MessageId: int32(r.MessageID),
					CreatedAt: r.CreatedAt.Local().Format(time.RFC3339),
					Reaction:  r.Reaction,
					Id:        uint64(r.ID),
				}
			}
		}
	}

	if msg.Reactions != nil {
		pb.Reactions = make([]*chatpb.MessageReaction, len(msg.Reactions))
		for i, r := range msg.Reactions {
			pb.Reactions[i] = &chatpb.MessageReaction{
				UserId:    int32(r.UserID),
				MessageId: int32(r.MessageID),
				CreatedAt: r.CreatedAt.Local().Format(time.RFC3339),
				Reaction:  r.Reaction,
				Id:        uint64(r.ID),
				RoomId:    strconv.FormatUint(r.ConversationId, 10), // THÊM
			}
		}
	}
	if msg.PinnedAt != nil {
		timeStr := _utils.FormatTimeToString(msg.PinnedAt)
		pb.PinnedAt = &timeStr
	}
	if msg.ReadAt != nil {
		readAt := _utils.FormatTimeToString(msg.ReadAt)
		pb.ReadAt = &readAt
	}
	return pb
}
