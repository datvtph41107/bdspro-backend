package dto

import _dto "common/domain/dto"

type ProductHistorySearch struct {
	_dto.Pagable
}

type InfoAreaResponse struct {
	Total    float64 `gorm:"column:total_area" json:"totalArea"`
	Avaiable float64 `gorm:"column:avaiable_area" json:"avaiableArea"`
	Area     float64 `gorm:"column:split_area" json:"splitArea"`
}
