package models

type Profile struct {
	ProfileId uint64
	Phone     string
}

func (Profile) TableName() string {
	return "sub_profile"
}
