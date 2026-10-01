package dto

// TagIDPath is DELETE /v2/user/admin/tags/{id}.
type TagIDPath struct {
	ID uint64 `path:"id"`
}

// TagListByTypeRequest is GET /v2/user/profile/tags/{tagType}.
type TagListByTypeRequest struct {
	TagType uint32 `path:"tagType"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
	Text    string `form:"text"`
	Sort    string `form:"sort"`
}

// AdminTagListRequest is GET /v2/user/admin/tags.
type AdminTagListRequest struct {
	TagTypes  []uint32 `form:"tagTypes"`
	Page      uint32   `form:"page"`
	Size      uint32   `form:"size"`
	Text      string   `form:"text"`
	Sort      string   `form:"sort"`
	IsDefault *bool    `form:"isDefault"`
}

// AdminTagWriteRequest is POST/PUT admin tag body (+ path id on update).
type AdminTagWriteRequest struct {
	ID        uint64 `path:"id" json:"id"`
	Name      string `json:"name"`
	TagType   uint32 `json:"tagType"`
	IsActive  bool   `json:"isActive"`
	IsDefault bool   `json:"isDefault"`
}

// GetProfileTagsByTypeRequest is GET /v2/pub/user/profile/tags/{tagType}.
type GetProfileTagsByTypeRequest struct {
	TagType   uint32 `path:"tagType"`
	ProfileID uint64 `form:"profileId"`
}
