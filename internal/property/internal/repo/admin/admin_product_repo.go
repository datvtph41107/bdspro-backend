package admin_repo

import (
	"bdspro/internal/domain"
	case_archived "common/case/archived"
	"common/case/crud"
)

type AdminProductRepo interface {
	crud.ICrudRepo[domain.Product]
	case_archived.IArchivedRepo[domain.Product]
}
