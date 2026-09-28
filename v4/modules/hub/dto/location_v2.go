package dto

// GetProvincesV2Request mirrors hubpb.GetProvincesV2Request query params.
type GetProvincesV2Request struct {
	Keyword string `form:"keyword"`
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
}

// GetProvinceV2ByCodeRequest mirrors hubpb.GetProvinceV2ByCodeRequest query params.
type GetProvinceV2ByCodeRequest struct {
	IncludeWards bool `form:"includeWards"`
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
