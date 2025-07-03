package dto

type CreateConversationRequest struct {
	Name              string  `json:"name"`
	Type              int32   `json:"type"`
	CreatedBy         int32   `json:"created_by"`
	Members           []int32 `json:"members"`
	IsBroadcast       bool    `json:"is_broadcast"`
	ForbidForward     bool    `json:"forbid_forward"`
	Avatar            string  `json:"avatar"`
	BackgroundImageId *uint64 `json:"background_image_id"`
}

type RemoveUserRequest struct {
	ConversationId int64 `json:"conversation_id"`
	TargetUserId   int64 `json:"target_user_id"`
}

type AddUserRequest struct {
	ConversationId int64 `json:"conversation_id"`
	UserId         int64 `json:"user_id"`
}
