package models

type Profile struct {
	ProfileId uint64
	Phone     string
	FullName  string
}

func (Profile) TableName() string {
	return "profile_transfer"
}
