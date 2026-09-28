package dto

// GetProvincesV2Request mirrors hubpb.GetProvincesV2Request query params.
type GetProvincesV2Request struct {
	Keyword string `form:"keyword"`
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
}

// GetProvinceV2ByCodeRequest mirrors hubpb.GetProvinceV2ByCodeRequest path and query params.
type GetProvinceV2ByCodeRequest struct {
	Code         int32 `path:"code"`
	IncludeWards bool  `form:"includeWards"`
}

// GetWardV2ByCodeRequest is the path payload for ward lookup by administrative code.
type GetWardV2ByCodeRequest struct {
	Code int32 `path:"code"`
}

// GetWardsV2Request mirrors hubpb.GetWardsV2Request query params.
type GetWardsV2Request struct {
	ProvinceID uint64 `form:"provinceId"`
	Keyword    string `form:"keyword"`
	Page       int32  `form:"page"`
	Size       int32  `form:"size"`
}

// SearchLocationV2Request mirrors hubpb.SearchLocationV2Request query params.
type SearchLocationV2Request struct {
	Keyword string `form:"keyword"`
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
}
