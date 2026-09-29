package models

type Profile struct {
	ProfileId uint64 `gorm:"column:profile_id;primaryKey"`
	Phone     string `gorm:"column:phone"`
	FullName  string `gorm:"column:full_name"`
}

func (Profile) TableName() string {
	return "profile_transfer"
}
