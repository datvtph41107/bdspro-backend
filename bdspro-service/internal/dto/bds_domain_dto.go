package dto

import _dto "common/domain/dto"

// BdsDomainListRequest DTO cho request lấy danh sách bất động sản
type BdsDomainListRequest struct {
	_dto.Pagable
	Text       *string `json:"text,omitempty"`
	ProductID  *uint64 `json:"productId,omitempty"`
	AssetID    *uint64 `json:"assetId,omitempty"`
	ProvinceID *uint64 `json:"provinceId,omitempty"`
	WardID     *uint64 `json:"wardId,omitempty"`
	DistrictID *uint64 `json:"districtId,omitempty"`
}

// BdsDomainListResponse DTO cho response danh sách bất động sản
type BdsDomainListResponse struct {
	Data  []BdsDomainResponse `json:"data"`
	Total int64               `json:"total"`
}

// BdsDomainResponse DTO cho response một bất động sản
type BdsDomainResponse struct {
	ID             uint64   `json:"id"`
	ProductID      *uint64  `json:"productId,omitempty"`
	AssetID        *uint64  `json:"assetId,omitempty"`
	Area           *float64 `json:"area,omitempty"`
	AreaUse        *float64 `json:"areaUse,omitempty"`
	AreaBuild      *float64 `json:"areaBuild,omitempty"`
	AreaLand       *float64 `json:"areaLand,omitempty"`
	AreaFloor      *float64 `json:"areaFloor,omitempty"`
	AreaGreen      *float64 `json:"areaGreen,omitempty"`
	AreaParking    *float64 `json:"areaParking,omitempty"`
	AreaBalcony    *float64 `json:"areaBalcony,omitempty"`
	AreaTerrace    *float64 `json:"areaTerrace,omitempty"`
	Length         *float64 `json:"length,omitempty"`
	Width          *float64 `json:"width,omitempty"`
	Height         *float64 `json:"height,omitempty"`
	FrontWidth     *float64 `json:"frontWidth,omitempty"`
	BackWidth      *float64 `json:"backWidth,omitempty"`
	LeftWidth      *float64 `json:"leftWidth,omitempty"`
	RightWidth     *float64 `json:"rightWidth,omitempty"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Address        string   `json:"address,omitempty"`
	Street         string   `json:"street,omitempty"`
	StreetNumber   string   `json:"streetNumber,omitempty"`
	WardID         *uint64  `json:"wardId,omitempty"`
	DistrictID     *uint64  `json:"districtId,omitempty"`
	ProvinceID     *uint64  `json:"provinceId,omitempty"`
	RoadWidth      *float64 `json:"roadWidth,omitempty"`
	RoadType       *string  `json:"roadType,omitempty"`
	RoadAccess     *string  `json:"roadAccess,omitempty"`
	DistanceToRoad *float64 `json:"distanceToRoad,omitempty"`
	Orientation    *string  `json:"orientation,omitempty"`
	MainDirection  *string  `json:"mainDirection,omitempty"`
	CornerLot      *bool    `json:"cornerLot,omitempty"`
	AlleyAccess    *bool    `json:"alleyAccess,omitempty"`
	AlleyWidth     *float64 `json:"alleyWidth,omitempty"`
	Shape          *string  `json:"shape,omitempty"`
	Topography     *string  `json:"topography,omitempty"`
	Elevation      *float64 `json:"elevation,omitempty"`
	Note           string   `json:"note,omitempty"`
	ProfileID      *uint64  `json:"profileId,omitempty"`
	CreatedAt      string   `json:"createdAt,omitempty"`
	UpdatedAt      string   `json:"updatedAt,omitempty"`
}
