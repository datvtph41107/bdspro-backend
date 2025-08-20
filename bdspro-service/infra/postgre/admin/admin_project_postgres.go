package admin_postgres

import (
	"bdspro/internal/domain"
	"common/case/crud"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/admin.AdminProjectRepo
type AdminProjectPostgres struct {
	crud.BaseRepo[domain.Project]
}

func NewAdminProjectPostgres(db *gorm.DB) *AdminProjectPostgres {
	return &AdminProjectPostgres{
		BaseRepo: crud.BaseRepo[domain.Project]{
			DB: db,
		},
	}
}
