package domain

type DealOfGroup struct {
	ID      uint64 `gorm:"primaryKey;column:id" json:"id"`
	DealID  uint64 `gorm:"column:deal_id;not null" json:"dealId"`
	GroupID uint64 `gorm:"column:group_id;not null" json:"groupId"`
}

func (DealOfGroup) TableName() string {
	return "deal_of_group"
}
