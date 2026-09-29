package domain

type DocumentEntity struct {
	ID       uint64 `gorm:"primaryKey"`
	FileName string `gorm:"not null"`
	FileUrl  string `gorm:"not null"`
	FileType string `gorm:"not null"`
	LeadID   uint64 `gorm:"not null"`
}

func (DocumentEntity) TableName() string {
	return "tb_document"
}