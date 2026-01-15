package dto

import (
	_dto "common/domain/dto"
	sharepb "pb/types/shared"

	"github.com/lib/pq"
)

type CreateProductNoteRequest struct {
	ProductID  uint64                `json:"productId" validate:"required"`
	Content    string                `json:"content" validate:"required"`
	MentionIDs []uint64              `json:"mentionIds"`
	Files      []ProductNoteFileItem `json:"fileIds"`
}

type UpdateProductNoteRequest struct {
	NoteID     uint64                `json:"noteId" validate:"required"`
	Content    string                `json:"content" validate:"required"`
	MentionIDs []uint64              `json:"mentionIds"`
	Files      []ProductNoteFileItem `json:"fileIds"`
}

type GetProductNotesFilter struct {
	_dto.Pagable
	ProductID  uint64  `json:"productId" validate:"required"`
	OnlyPinned *bool   `json:"onlyPinned"`
	AuthorID   *uint64 `json:"authorId"`
	Query      *string `json:"query"`
}

type ProductNoteItemRespone struct {
	ID        uint64 `gorm:"column:id"`
	Content   string `gorm:"column:content"`
	IsPinned  bool   `gorm:"column:is_pinned"`
	AuthorId  uint64 `gorm:"column:author_id"`
	CreatedAt string `gorm:"column:created_at"`

	MentionIDs pq.Int64Array `gorm:"column:mention_ids;type:bigint[]"`
	FileIDs    pq.Int64Array `gorm:"column:file_ids;type:bigint[]"`

	// enrich sau
	Author   *sharepb.ProfileItem     `gorm:"-"`
	Mentions []ProductNoteMentionItem `gorm:"-"`
	Files    []ProductNoteFileItem    `gorm:"-"`
}

type ProductNoteMentionItem struct {
	UserId uint64 `json:"id"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Role   string `json:"role,omitempty"`
}

type ProductNoteFileItem struct {
	ID       uint64 `json:"id"`
	FileName string `json:"fileName"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
}
