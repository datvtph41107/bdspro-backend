package dto

type RegionRequest struct {
	Text     string `form:"text"`
	Level    *int   `form:"level"`
	ParentID *uint  `form:"parentId"`
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}
