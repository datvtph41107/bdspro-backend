package dto

import (
	"bdspro/internal/domain"
	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	"time"
)

type CreateProductNoteRequest struct {
	ProductID uint64               `json:"productId" validate:"required"`
	Content   string               `json:"content" validate:"required"`
	Mentions  []MentionNoteDTO     `json:"mentionIds"`
	Files     []ProductNoteFileDTO `json:"files"`
}

type MentionNoteDTO struct {
	UserID   uint64 `json:"userId" validate:"required"`
	StartPos int32  `json:"startPos" validate:"required,min=0"`
	Length   int32  `json:"length" validate:"required,min=1"`
}

type ProductNoteFileDTO struct {
	URL      string `json:"url" validate:"required,url"`
	FileName string `json:"fileName" validate:"required"`
	FileType string `json:"fileType" validate:"required"`
	Size     int64  `json:"size" validate:"required,min=1"`
}

type UpdateProductNoteRequest struct {
	NoteID   uint64               `json:"noteId" validate:"required"`
	Content  string               `json:"content" validate:"required"`
	Mentions []MentionNoteDTO     `json:"mentionIds"`
	Files    []ProductNoteFileDTO `json:"files"`
}

type GetProductNotesFilter struct {
	_dto.Pagable
	ProductID  uint64  `json:"productId" validate:"required"`
	OnlyPinned *bool   `json:"onlyPinned"`
	AuthorID   *uint64 `json:"authorId"`
	Query      *string `json:"query"`
}

type ProductNoteMentionItem struct {
	UserId   uint64 `gorm:"column:user_id" json:"userId"`
	StartPos int32  `gorm:"column:start_pos" json:"startPos"`
	Length   int32  `gorm:"column:length" json:"length"`
	NoteID   uint64 `gorm:"column:note_id"`

	Name   string `gorm:"-" json:"name,omitempty"`
	Avatar string `gorm:"-" json:"avatar,omitempty"`
	Role   string `gorm:"-" json:"role,omitempty"`
}

type ProductNoteRow struct {
	ID                   uint64               `gorm:"column:id"`
	ProductID            uint64               `gorm:"column:product_id"`
	Content              string               `gorm:"column:content"`
	IsPinned             bool                 `gorm:"column:is_pinned"`
	AuthorId             uint64               `gorm:"column:author_id"`
	CreatedAt            string               `gorm:"column:created_at"`
	ProductNoteUpdatedAt *time.Time           `gorm:"column:product_note_updated_at"`
	Author               *sharepb.ProfileItem `gorm:"-"`
}

type ProductNote struct {
	Note     *ProductNoteRow
	Mentions []*ProductNoteMentionItem
	Files    []*domain.ProductNoteFile
}
