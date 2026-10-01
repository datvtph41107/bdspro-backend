package entity

type Profile struct {
	ProfileId    uint64
	Phone        string
	FullName     string
	Avatar       string
	TickVerified bool
}

func (Profile) TableName() string {
	return "profile_transfer"
}