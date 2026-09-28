package dto

// GetFAQSimpleRequest mirrors hubpb.GetFAQSimpleRequest query params.
type GetFAQSimpleRequest struct {
	GroupKey string `form:"groupKey"`
	Page     int32  `form:"page"`
	Size     int32  `form:"size"`
	Text     string `form:"text"`
}
