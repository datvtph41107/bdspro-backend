package admin_postgres

import (
	"bdspro/internal/domain"
	"common/case/crud"
	_db "common/db"
)

// @bind: bdspro/internal/repo/admin.AdminProjectRepo
type AdminProjectPostgres struct {
	*crud.CrudRepo[domain.Project]
}

func NewAdminProjectPostgres(db *_db.TransactionRepo) *AdminProjectPostgres {
	x := &AdminProjectPostgres{
		CrudRepo: &crud.CrudRepo[domain.Project]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}
