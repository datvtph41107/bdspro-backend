package dto

import (
	_dto "common/domain/dto"
	sharepb "pb/types/shared"

	"github.com/lib/pq"
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

type ProductNoteItemResponse struct {
	ID        uint64 `gorm:"column:id"`
	ProductID uint64 `gorm:"column:productId"`
	Content   string `gorm:"column:content"`
	IsPinned  bool   `gorm:"column:is_pinned"`
	AuthorId  uint64 `gorm:"column:author_id"`
	CreatedAt string `gorm:"column:created_at"`

	MentionIDs pq.Int64Array `gorm:"column:mention_ids;type:bigint[]"`
	FileIDs    pq.Int64Array `gorm:"column:file_ids;type:bigint[]"`

	// enrich sau
	Author   *sharepb.ProfileItem     `gorm:"-"`
	Mentions []ProductNoteMentionItem `gorm:"-"`
	Files    []ProductNoteFileDTO     `json:"-"`
}

type ProductNoteMentionItem struct {
	UserId   uint64 `json:"userId"`
	StartPos int32  `json:"startPos"`
	Length   int32  `json:"length"`

	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Role   string `json:"role,omitempty"`
}
