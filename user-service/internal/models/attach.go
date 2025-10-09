package models

type AttachEntity struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	ImageURL  string  `gorm:"type:varchar(255);not null" json:"imageUrl"`
	ImageType string  `gorm:"type:varchar(255);not null" json:"imageType"`
	RateID    *uint64 `gorm:"not null" json:"rateId"`
	// Rate      *Rate   `gorm:"foreignKey:RateID" json:"rate"`
}

func (AttachEntity) TableName() string {
	return "attachs"
}
