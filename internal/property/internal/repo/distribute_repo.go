package repo

import (
	"context"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type DistributeRepository interface {
	Save(ctx context.Context, entity *domain.DistributionEntity) error
	FindByID(ctx context.Context, id uint64) (*domain.DistributionEntity, error)
	AssignPartner(ctx context.Context, distributeID uint64, productID uint64, partnerIDs []uint64) error
	FindByIdWithPartner(ctx context.Context, disId uint64) (*dto.DistributionWithPartners, error)
	GetListDistrByProduct(ctx context.Context, productID uint64, dto *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error)
	GetListPartnersOfDistribute(ctx context.Context, disId uint64) ([]*domain.ProductUser, error)
	Revoke(ctx context.Context, id uint64, revokedAt time.Time, revoke dto.RevokeDTO) error
	RevokeUser(ctx context.Context, distributeID uint64, originProfileID uint64) (int64, error)
}
