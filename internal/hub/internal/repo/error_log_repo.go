package repo

import (
	"context"
	"hub/internal/domain"
)

type IErrorLogRepo interface {
	Insert(ctx context.Context, log *domain.ErrorLog) error
}
