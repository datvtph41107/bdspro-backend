package dto

import _dto "common/dto"

type ProductHistorySearch struct {
	_dto.Page
}

type InfoAreaResponse struct {
	Total    float64 `gorm:"column:total_area" json:"totalArea"`
	Avaiable float64 `gorm:"column:avaiable_area" json:"avaiableArea"`
	Area     float64 `gorm:"column:split_area" json:"splitArea"`
}
