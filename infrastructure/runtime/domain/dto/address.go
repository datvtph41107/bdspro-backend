package _dto

type AddressV3DTO struct {
	ID           uint64  `json:"id"`
	Detail       string  `json:"detail"`
	ProvinceID   *uint64 `json:"provinceId"`
	ProvinceName string  `json:"provinceName"`
	DistrictID   *uint64 `json:"districtId"`
	DistrictName string  `json:"districtName"`
	WardID       *uint64 `json:"wardId"`
	WardName     string  `json:"wardName"`
}
