package repo

import (
	_crud "common/domain/crud"
	"mdn4/modules/user/entity"
)

type IPriceTableRepo interface {
	_crud.ICrudRepo[entity.PriceTableDomain]
}
