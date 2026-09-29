package domain

import _models "common/models"

// Project struct represents the projects table
// TableName sets the table name explicitly
type BuildRoom struct {
	_models.BaseEntity
	BlockID      uint64 `gorm:"not null" json:"blockId"`
	Name         string `gorm:"size:255;not null" json:"name"`
	Floor        uint   `gorm:"type:int" json:"floor"`
	Order        uint   `gorm:"column:ord;type:int" json:"order"`
	Note         string `gorm:"type:text" json:"note"`
	UnitCode     string `gorm:"type:varchar(20)" json:"unitCode"`
	Area         uint   `gorm:"type:int" json:"area"`   // Diện tích
	NumBedroom   *int   `json:"numBedroom,omitempty"`   // Số phòng ngủ (optional)
	NumBathroom  *int   `json:"numBathroom,omitempty"`  // Số WC (optional)
	Furniture    string `json:"furniture,omitempty"`    // Nội thất
	BlueprintUrl string `json:"blueprintUrl,omitempty"` // Link bản vẽ
}

func (BuildRoom) TableName() string {
	return "build_room"
}

type BuildRoomItem struct {
	ID           uint64 `gorm:"primaryKey" json:"-"`
	BlockID      uint64 `gorm:"type:not null" json:"-"`
	Name         string `gorm:"size:255;not null" json:"name"`
	Floor        uint   `gorm:"type:int" json:"floor"`
	Order        uint   `gorm:"column:ord;type:int" json:"order"`
	Note         string `gorm:"type:text" json:"note"`
	UnitCode     string `gorm:"type:varchar(20)" json:"unitCode"`
	Area         uint   `gorm:"type:int" json:"area"`   // Diện tích
	NumBedroom   *int   `json:"numBedroom,omitempty"`   // Số phòng ngủ (optional)
	NumBathroom  *int   `json:"numBathroom,omitempty"`  // Số WC (optional)
	Furniture    string `json:"furniture,omitempty"`    // Nội thất
	BlueprintUrl string `json:"blueprintUrl,omitempty"` // Link bản vẽ
}

func (BuildRoomItem) TableName() string {
	return "build_room"
}
