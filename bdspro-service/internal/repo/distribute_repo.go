package repo

import (
	"context"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type DistributeRepository interface {
	Save(ctx context.Context, entities []*domain.DistributionEntity) error
	AssignPartner(ctx context.Context, partners ...*domain.ProductUser) error
	FindByIdWithPartner(ctx context.Context, disId uint64) (*domain.DistributionEntity, error)         // đây sẽ là hàm handle lấy ra distribute và thông tin của các partner
	GetListDistrByProduct(ctx context.Context, productID uint64) ([]*domain.DistributionEntity, error) // lấy ra list distribute với thông tin các user partner
	GetPartnersOfDistribute(ctx context.Context, disId uint64) ([]*domain.ProductUser, error)
	Revoke(ctx context.Context, id uint64, revokedAt time.Time, revoke dto.RevokeDTO) error
	CountActiveByProduct(ctx context.Context, productID uint64) (int64, error)
	RemovePartnersFromDistribute(ctx context.Context, distributeID uint64, partnerIDs []uint64) error
	FindByID(ctx context.Context, id uint64) (*domain.DistributionEntity, error)
	FindActiveByProductAndPartner(ctx context.Context, productID, partnerID uint64) (*domain.DistributionEntity, error)
}
