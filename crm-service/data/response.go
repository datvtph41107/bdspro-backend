package data

import sharepb "pb/types/shared"

type FollowInfoRes struct {
	NumFollower      int64                  `json:"numFollower"`
	NumFollowing     int64                  `json:"numFollowing"`
	NumFriend        int64                  `json:"numFriend"`
	ContactId        *uint64                `json:"contactId,omitempty"`
	BlockId          *uint64                `json:"blockId,omitempty"`
	FollowId         *uint64                `json:"followId,omitempty"`
	FriendCommons    []*sharepb.UserV3Proto `json:"friendCommons,omitempty"`
	NumFriendCommons uint64                 `json:"numFriendCommons"`
}