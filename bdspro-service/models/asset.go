package models

type Asset struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnershipDocId uint64        `json:"ownership_doc_id"` // giấy tờ sở hữu
	OwnershipDoc   *OwnershipDoc `gorm:"foreignKey:OwnershipDocId" json:"ownershipDoc,omitempty"`
	ProvinceID     uint64        `json:"province_id"`
	Province       *Province     `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	DistrictID     uint64        `json:"district_id"`
	District       *District     `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	WardID         uint64        `json:"ward_id"`
	Ward           *Ward         `gorm:"foreignKey:WardID" json:"ward,omitempty"`
	Address        string        `gorm:"size:255" json:"address"`
	GoogleMapLink  string        `gorm:"size:255" json:"google_map_link"`
}
