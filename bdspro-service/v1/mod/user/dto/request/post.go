package user_dto

import _dto "common/dto"

type PostLinkSearch struct {
	_dto.Page
	ProductID *uint64 `json:"productId"`
}
