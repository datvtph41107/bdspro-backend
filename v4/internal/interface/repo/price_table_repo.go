package repo

import (
	_crud "common/domain/crud"
	"userv4/internal/models"
)

type IPriceTableRepo interface {
	_crud.ICrudRepo[models.PriceTableDomain]
}
