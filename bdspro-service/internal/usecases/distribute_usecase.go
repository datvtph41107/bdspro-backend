package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	"context"
	"fmt"
	"time"
)

type DistributionUsecase struct {
	distributionRepo repo.DistributeRepository
	transaction      provider.TransactionProvider
	userProvider     provider.IUserProvider
}

func NewDistributionUsecase(
	r repo.DistributeRepository,
	transaction provider.TransactionProvider,
	userProvider provider.IUserProvider,
) *DistributionUsecase {
	return &DistributionUsecase{
		distributionRepo: r,
		transaction:      transaction,
		userProvider:     userProvider,
	}
}

func (uc *DistributionUsecase) CreateDistribute(ctx context.Context, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &dto.DistributeResultDTO{}
	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()

		partnerCommission := calculatePartnerCommission(req.Price, req.CommissionType, req.CommissionValue)

		dist := &domain.DistributionEntity{
			Price:             req.Price,
			CanDeal:           req.CanDeal,
			ChannelPrice:      req.ChannelPrice,
			CommissionType:    req.CommissionType,
			CommissionValue:   req.CommissionValue,
			Note:              req.Note,
			DistributedAt:     &now,
			FromDate:          req.FromDate,
			ToDate:            req.ToDate,
			Status:            enums.DISTRIBUTION_STATUS_ACTIVE,
			PartnerCommission: partnerCommission,
			ProductID:         req.ProductID,
		}

		if err := uc.distributionRepo.Save(txCtx, []*domain.DistributionEntity{dist}); err != nil {
			return fmt.Errorf("failed to create distribution entity: %w", err)
		}

		createdCount, err := uc.assignPartners(txCtx, req.PartnerIDs, req.ProductID, dist.ID)
		if err != nil {
			return fmt.Errorf("failed to assign partners: %w", err)
		}

		resp.Price = dist.Price
		resp.CanDeal = dist.CanDeal
		resp.ChannelPrice = dist.ChannelPrice
		resp.Status = dist.Status
		resp.CommissionType = dist.CommissionType
		resp.CommissionValue = dist.CommissionValue
		resp.Note = dist.Note
		resp.PartnerCommission = dist.PartnerCommission
		resp.CreatedCount = createdCount
		resp.FromDate = dist.FromDate
		resp.ToDate = dist.ToDate

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (uc *DistributionUsecase) UpdateDistribute(ctx context.Context, distID uint64, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &dto.DistributeResultDTO{}
	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		dist, err := uc.distributionRepo.FindByID(txCtx, distID)
		if err != nil {
			return fmt.Errorf("distribution not found: %w", err)
		}

		// Cập nhật các trường
		dist.Price = req.Price
		dist.CanDeal = req.CanDeal
		dist.ChannelPrice = req.ChannelPrice
		dist.CommissionType = req.CommissionType
		dist.CommissionValue = req.CommissionValue
		dist.Note = req.Note
		dist.FromDate = req.FromDate
		dist.ToDate = req.ToDate
		dist.PartnerCommission = calculatePartnerCommission(req.Price, req.CommissionType, req.CommissionValue)

		if dist.Status != enums.DISTRIBUTION_STATUS_REVOKED && time.Now().After(*dist.ToDate) {
			dist.Status = enums.DISTRIBUTION_STATUS_EXPRIED
		}

		if err := uc.distributionRepo.Save(txCtx, []*domain.DistributionEntity{dist}); err != nil {
			return fmt.Errorf("failed to update distribution entity: %w", err)
		}

		existingPartners, _ := uc.distributionRepo.GetPartnersOfDistribute(txCtx, dist.ID)
		toRemove := []uint64{}
		for _, p := range existingPartners {
			if !contains(req.PartnerIDs, p.ProfileID) {
				toRemove = append(toRemove, p.ProfileID)
			}
		}
		if len(toRemove) > 0 {
			_ = uc.distributionRepo.RemovePartnersFromDistribute(txCtx, dist.ID, toRemove)
		}

		createdCount, err := uc.assignPartners(txCtx, req.PartnerIDs, req.ProductID, dist.ID)
		if err != nil {
			return fmt.Errorf("failed to assign partners: %w", err)
		}

		// Gán resultDTO
		resp.Price = dist.Price
		resp.CanDeal = dist.CanDeal
		resp.ChannelPrice = dist.ChannelPrice
		resp.Status = dist.Status
		resp.CommissionType = dist.CommissionType
		resp.CommissionValue = dist.CommissionValue
		resp.Note = dist.Note
		resp.PartnerCommission = dist.PartnerCommission
		resp.CreatedCount = createdCount
		resp.FromDate = dist.FromDate
		resp.ToDate = dist.ToDate

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (uc *DistributionUsecase) assignPartners(ctx context.Context, partnerIDs []uint64, productID, distID uint64) (int, error) {
	if len(partnerIDs) == 0 {
		return 0, nil
	}

	partners := make([]*domain.ProductUser, 0, len(partnerIDs))
	for _, pid := range partnerIDs {
		partners = append(partners, &domain.ProductUser{
			ProfileID:    pid,
			ProductID:    productID,
			DistributeID: distID,
			IsOwner:      false,
		})
	}
	if err := uc.distributionRepo.AssignPartner(ctx, partners...); err != nil {
		return 0, err
	}

	return len(partners), nil
}

func (uc *DistributionUsecase) RevokeDistribute(ctx context.Context, distributeID uint64, revoke dto.RevokeDTO) error {
	now := time.Now()
	return uc.distributionRepo.Revoke(ctx, distributeID, now, revoke)
}

func (uc *DistributionUsecase) GetDistributeWithPartners(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error) {
	dist, err := uc.distributionRepo.FindByIdWithPartner(ctx, disID)
	if err != nil {
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}
	if dist == nil {
		return nil, fmt.Errorf("distribution %d not found", disID)
	}

	return dist, nil
}

func (uc *DistributionUsecase) GetListDistributes(ctx context.Context, productID uint64, dto *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error) {
	distributes, err := uc.distributionRepo.GetListDistrByProduct(ctx, productID, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get distribute list: %w", err)
	}

	return distributes, nil
}

func (uc *DistributionUsecase) ListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*domain.ProductUser, error) {
	partners, err := uc.distributionRepo.GetPartnersOfDistribute(ctx, disID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partners: %w", err)
	}
	return partners, nil
}

func calculatePartnerCommission(price int64, commissionType enums.CommissionType, commissionValue float64) float64 {
	if commissionValue <= 0 {
		return 0
	}

	switch commissionType {
	case enums.CommissionTypePercent:
		return float64(price) * commissionValue / 100
	case enums.CommissionTypeVND:
		return commissionValue
	default:
		return 0
	}
}

func contains(slice []uint64, val uint64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// func (uc *DistributionUsecase) CreateDistribute(ctx context.Context, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error) {
// 	resp := &dto.DistributeResultDTO{}
// 	now := time.Now()
// 	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
// 		var toCreate []*domain.DistributionEntity
// 		for _, partnerID := range req.PartnerIDs {
// 			toCreate = append(toCreate, &domain.DistributionEntity{
// 				ProductID:     req.ProductID,
// 				PartnerID:     partnerID,
// 				Status:        enums.DistributionStatusActive,
// 				Policy:        req.Policy,
// 				DistributedAt: now,
// 			})
// 		}
// 		if len(toCreate) > 0 {
// 			if err := uc.distributionRepo.Create(txCtx, toCreate); err != nil {
// 				return err
// 			}
// 			resp.CreatedCount = len(toCreate)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }
