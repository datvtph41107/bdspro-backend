package u_asset_cost_type

import _dto "common/dto"

type SearchDTO struct {
	_dto.Page
	Text  string  `form:"text"`
	Types *uint64 `form:"types"` // 10, 20
}
