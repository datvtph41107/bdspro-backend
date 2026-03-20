package mapper

import (
	"chat/models"
	"encoding/json"
	chatpb "pb/types/chat"
	"strconv"
	"time"
)

type MessageMapper struct{}

func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

func (m *MessageMapper) MessageToPb(msg *models.MessageModel) *chatpb.Message {
	pb := &chatpb.Message{
		MessageId:   int32(msg.ID),
		SenderId:    int32(msg.SenderID),
		Content:     msg.Content,
		ContentType: string(msg.ContentType),
		FileUrl:     msg.FileURL,
		RoomId:      strconv.FormatUint(msg.ConversationID, 10),
		CreatedAt:   msg.CreatedAt.Format(time.RFC3339),
		IsRecall:    msg.Recall,
		ExtraId:     msg.ExtraId,
	}

	// Parse system message data if content type is system
	if msg.ContentType == models.ContentTypeSystem && msg.Content != "" {
		var systemData map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Content), &systemData); err == nil {
			// Convert to protobuf system message data
			pb.SystemData = &chatpb.SystemMessageData{
				Action:    getString(systemData, "action"),
				ActorId:   getUint64(systemData, "actor_id"),
				TargetIds: getUint64Slice(systemData, "target_ids"),
				Extra:     getString(systemData, "extra"),
			}
		}
	}

	if msg.ReplyMessage != nil {
		pb.ReplyMessage = m.MessageToPb(msg.ReplyMessage)
	}

	return pb
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getUint64(data map[string]interface{}, key string) uint64 {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return uint64(num)
		}
	}
	return 0
}

func getUint64Slice(data map[string]interface{}, key string) []uint64 {
	if val, ok := data[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]uint64, 0, len(slice))
			for _, item := range slice {
				if num, ok := item.(float64); ok {
					result = append(result, uint64(num))
				}
			}
			return result
		}
	}
	return nil
}
