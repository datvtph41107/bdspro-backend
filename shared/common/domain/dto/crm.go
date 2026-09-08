package _dto

type FriendItemDTO struct {
	ID         uint64 `json:"id"`
	ReceiverID uint64 `json:"receiverId"`
	CreatedBy  uint64 `json:"createdBy"`
	Status     uint32 `json:"status"` // 10: Pending, 20: Accepted, 30: Declined
}

type RelationShipDTO struct {
	FriendStatus *FriendItemDTO
	Following    bool
	Blocked      bool
	NumFollower  int64
	NumFollowing int64
	NumFriend    int64
}

type ConversationMemberDTO struct {
	userId    uint64
	joinedAt  string
	leftAt    *string
	muteNotif bool
	role      uint32 // 0: Member, 1: Admin, 2: Banned
}

type ContactInfoV3DTO struct {
	NumFriend        int64        `json:"numFriend,omitempty"`
	NumFollower      int64        `json:"numFollower,omitempty"`
	NumFollowing     int64        `json:"numFollowing,omitempty"`
	ContactId        *uint64      `json:"contactId,omitempty"`
	BlockId          *uint64      `json:"blockId,omitempty"`
	ReportId         *uint64      `json:"reportId,omitempty"`
	FollowId         *uint64      `json:"followId,omitempty"`
	FriendCommons    []*UserV3DTO `json:"friendCommons"`
	NumFriendCommons uint64       `json:"numFriendCommons,omitempty"`

	Participant *ConversationMemberDTO `json:"participant,omitempty"`
}
