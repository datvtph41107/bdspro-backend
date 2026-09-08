package domain

// LocationSearchResultV2 represents a joined result between Ward and Province
type LocationSearchResultV2 struct {
	WardID       uint64 `json:"wardId"`
	WardName     string `json:"wardName"`
	WardType     string `json:"wardType"`
	ProvinceID   uint64 `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
	ProvinceType string `json:"provinceType"`
	FullAddress  string `json:"fullAddress"`
}
