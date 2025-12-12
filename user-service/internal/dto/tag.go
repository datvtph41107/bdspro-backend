package dto

import "user/enums"

type AdminTagRequest struct {
	ID        uint64
	Name      string
	TagType   enums.TagType
	IsActive  bool
	IsDefault bool
}

type TagListRequest struct {
	TagType enums.TagType
	Page    uint32
	Size    uint32
	Text    string
	Sort    string
}

type AdminTagListRequest struct {
	TagTypes []enums.TagType
	Page     uint32
	Size     uint32
	Text     string
	Sort     string

	IsDefault *bool
}
