package admin_repo

import (
	"bdspro/internal/domain"
	"common/case/crud"
)

type AdminPropertyTypeRepo interface {
	crud.IBaseRepo[domain.PropertyType]
}
