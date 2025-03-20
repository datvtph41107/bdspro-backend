package data

type FollowInfoRes struct {
	NumFollower  int64 `json:"numFollower"`
	NumFollowing int64 `json:"numFollowing"`
	NumFriend    int64 `json:"numFriend"`
}
