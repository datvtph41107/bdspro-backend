package models

type BlockEntity struct {
	ProfileId uint64 `gorm:"column:profile_id;primaryKey"`
	BlockedId uint64 `gorm:"column:blocked_id;primaryKey"`
}

func (BlockEntity) TableName() string {
	return "block"
}
