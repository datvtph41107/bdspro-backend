package iusecase

import (
	"context"
	"tx/internal/dto"
)

type IOrganizationClient interface {
	AddProductToDeal(ctx context.Context, req *dto.AddProductToDealRequest) error
}
