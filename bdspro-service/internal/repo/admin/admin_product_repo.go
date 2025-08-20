package admin_repo

import (
	"bdspro/internal/domain"
	"common/case/crud"
)

type AdminProductRepo interface {
	crud.IBaseRepo[domain.Product]
}
