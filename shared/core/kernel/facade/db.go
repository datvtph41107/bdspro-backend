package facade

import (
	"context"
	"core/kernel/contract"
)

type DBFacade struct {
	db contract.DB
}

func WithDB(ctx context.Context, db contract.DB) context.Context {
	return with(ctx, "facade_db", db)
}

func DB(ctx context.Context) DBFacade {
	if v := ctx.Value("facade_db"); v != nil {
		return DBFacade{db: v.(contract.DB)}
	}
	panic("DB facade not bound to context")
}

func (d DBFacade) Transaction(fn func(ctx context.Context) error) error {
	return d.db.Transaction(fn)
}
