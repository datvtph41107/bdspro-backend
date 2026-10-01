package entity

// ClickURL records a link-click event.
type ClickURL struct {
	ID         uint64  `gorm:"primaryKey;autoIncrement"`
	EventTime  int64   `gorm:"not null;index;column:event_time"`
	DomainID   uint64  `gorm:"not null;index;column:domain_id"`
	DomainType uint32  `gorm:"not null;index;column:domain_type"`
	ProfileID  *uint64 `gorm:"index;column:profile_id"`
}

func (ClickURL) TableName() string { return "click_urls" }
