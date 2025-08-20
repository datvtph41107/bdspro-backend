package admin_repo

import (
	"bdspro/internal/domain"
	"common/case/crud"
)

type AdminProjectRepo interface {
	crud.IBaseRepo[domain.Project]
}
