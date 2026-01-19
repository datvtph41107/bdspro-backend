package domain

import _models "common/models"

type ProductNoteFile struct {
	_models.BaseEntity
	NoteID   uint64 `gorm:"not null;index"`
	FileName string `gorm:"type:varchar(255);"`
	FileType string `gorm:"type:varchar(50);"`
	URL      string `gorm:"type:text;not null"`
	Size     int64  `gorm:""`
}

func (ProductNoteFile) TableName() string {
	return "product_note_files"
}
