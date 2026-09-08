package domain

// LocationSearchResult represents a joined result between District and Province
type LocationSearchResult struct {
	DistrictID   uint64 `json:"districtId"`
	DistrictName string `json:"districtName"`
	DistrictType string `json:"districtType"`
	ProvinceID   uint64 `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
	ProvinceType string `json:"provinceType"`
	FullAddress  string `json:"fullAddress"`
}
