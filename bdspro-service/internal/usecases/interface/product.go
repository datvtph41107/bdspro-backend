package iusecase

import (
	"bdspro/internal/enums"
	"context"
)

type IProductUsecase interface {
	GetOwnerID(c context.Context) (*uint64, error)
	HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error)
}
