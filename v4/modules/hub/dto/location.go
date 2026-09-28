package dto

// GetProvincesRequest mirrors hubpb.GetProvincesRequest query params.
type GetProvincesRequest struct {
	Keyword string `form:"keyword"`
}

// GetDistrictsRequest mirrors hubpb.GetDistrictsRequest query params.
type GetDistrictsRequest struct {
	ProvinceID string `form:"provinceId"`
	Keyword    string `form:"keyword"`
}

// GetWardsRequest mirrors hubpb.GetWardsRequest query params.
type GetWardsRequest struct {
	DistrictID string `form:"districtId"`
	Keyword    string `form:"keyword"`
}

// SearchLocationRequest mirrors hubpb.SearchLocationRequest query params.
type SearchLocationRequest struct {
	Keyword string `form:"keyword"`
}
