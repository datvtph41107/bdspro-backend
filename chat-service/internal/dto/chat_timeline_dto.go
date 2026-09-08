package dto

type LoadMessageRequestDTO struct {
	ConversationID uint64 `json:"conversation_id"`
	Timestamp      int64  `json:"ts"`
	Sequence       uint64 `json:"seq"`
	Limit          int    `json:"limit"`
	Action         string `json:"action"`
}
