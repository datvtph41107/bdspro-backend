package repo

import (
	_crud "common/domain/crud"
	"user/internal/models"
)

type IPriceTableRepo interface {
	_crud.ICrudRepo[models.PriceTableDomain]
}
