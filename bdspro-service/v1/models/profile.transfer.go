package models

type ProfileTransfer struct {
	ProfileId uint64
	Phone     string
	FullName  string
}

func (ProfileTransfer) TableName() string {
	return "profile_transfer"
}

type ProfileInfo struct {
	ProfileId uint64 `gorm:"column:profile_id" json:"profileId"`
	FullName  string `gorm:"column:full_name" json:"fullName"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
}

func (ProfileInfo) TableName() string {
	return "profile_transfer"
}
