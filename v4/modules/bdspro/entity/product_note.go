package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductNote struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy *uint64        `gorm:"column:created_by"`
	UpdatedBy *uint64        `gorm:"column:updated_by"`
	ProductID uint64         `gorm:"column:product_id;index;not null"`
	AuthorID  uint64         `gorm:"column:author_id;index;not null"`
	Content   string         `gorm:"column:content;type:text;not null"`
	IsPinned  *bool          `gorm:"column:is_pinned"`
	PinnedAt  *time.Time     `gorm:"column:pinned_at"`
}

func (ProductNote) TableName() string { return "product_notes" }

type ProductNoteMention struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	NoteID    uint64         `gorm:"column:note_id;index;not null"`
	UserID    uint64         `gorm:"column:user_id;not null"`
	StartPos  int32          `gorm:"column:start_pos"`
	Length    int32          `gorm:"column:length"`
}

func (ProductNoteMention) TableName() string { return "product_note_mentions" }

type ProductNoteFile struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	NoteID    uint64         `gorm:"column:note_id;index;not null"`
	FileName  string         `gorm:"column:file_name;size:255"`
	URL       string         `gorm:"column:url;type:text;not null"`
	Size      int64          `gorm:"column:size"`
	FileType  string         `gorm:"column:file_type;size:50"`
}

func (ProductNoteFile) TableName() string { return "product_note_files" }
