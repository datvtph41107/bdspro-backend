package admin_postgres

import (
	"bdspro/internal/domain"
	"common/case/crud"
	_db "common/db"
)

// @bind: bdspro/internal/repo/admin.AdminProjectRepo
type AdminProjectPostgres struct {
	crud.BaseRepo[domain.Project]
}

func NewAdminProjectPostgres(db *_db.DBTransaction) *AdminProjectPostgres {
	return crud.MakeCrudRepo(&AdminProjectPostgres{}, db).(*AdminProjectPostgres)
}
