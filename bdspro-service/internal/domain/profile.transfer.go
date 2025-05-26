package domain

type ProfileTransfer struct {
	ProfileId uint64
	Phone     string
	FullName  string
}

func (ProfileTransfer) TableName() string {
	return "profile_transfer"
}

type ProfileInfo struct {
	ProfileId uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
}

func (ProfileInfo) TableName() string {
	return "profile_transfer"
}
