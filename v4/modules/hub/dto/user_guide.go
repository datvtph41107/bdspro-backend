package dto

// GetUserGuideSimpleRequest mirrors hubpb.GetUserGuideSimpleRequest query params.
type GetUserGuideSimpleRequest struct {
	GroupKey string `form:"groupKey"`
	Page     int32  `form:"page"`
	Size     int32  `form:"size"`
	Text     string `form:"text"`
}

// GetAllUserGuidesByKeyRequest mirrors hubpb.GetAllUserGuidesByKeyRequest query params.
type GetAllUserGuidesByKeyRequest struct {
	Timestamp int64 `form:"timestamp"`
}
