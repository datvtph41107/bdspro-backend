package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
)

type DistributionUsecase struct {
	distributionRepo repo.DistributeRepository
	transaction      provider.TransactionProvider
	notifyProvider   provider.NotificationProvider
}

func NewDistributionUsecase(
	r repo.DistributeRepository,
	transaction provider.TransactionProvider,
	notificationProvider provider.NotificationProvider,

) *DistributionUsecase {
	return &DistributionUsecase{
		distributionRepo: r,
		transaction:      transaction,
		notifyProvider:   notificationProvider,
	}
}

func (uc *DistributionUsecase) CreateDistribute(
	ctx context.Context,
	req *dto.DistributeDTO,
) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(); err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}

	var resp *dto.DistributeResultDTO
	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()

		dist := &domain.DistributionEntity{
			ProductID:       req.ProductID,
			Price:           req.Price,
			CanDeal:         req.CanDeal,
			ChannelPrice:    req.ChannelPrice,
			CommissionType:  req.CommissionType,
			CommissionValue: req.CommissionValue,
			PartnerCommission: calculatePartnerCommission(
				req.Price,
				req.CommissionType,
				req.CommissionValue,
			),
			Note:          req.Note,
			FromDate:      req.FromDate,
			ToDate:        req.ToDate,
			DistributedAt: &now,
		}

		if err := uc.distributionRepo.Save(txCtx, dist); err != nil {
			log.Printf("create distribution failed: %v", err)
			return _errors.InternalServerException("create distribution failed")
		}

		if err := uc.distributionRepo.AssignPartner(
			txCtx,
			dist.ID,
			dist.ProductID,
			req.PartnerIDs,
		); err != nil {
			return err
		}

		resp = mapToResultDTO(dist, len(req.PartnerIDs))

		uc.notifyProvider.RegistedEventProperty(ctx,
			&dto.CreatePropertyActivityDTO{
				SubjectID:   dist.ID,
				SubjectType: _enum.SUBJECT_TYPE_DISTRIBUTE,
				Action:      _enum.ACTION_DISTRIBUTE_CREATE,
				ActorID:     _utils.GetProfileIdWithContext(ctx),
				Description: fmt.Sprintf(
					"Sản phẩm %d được tạo kênh phân phối với %d đối tác",
					dist.ProductID,
					len(req.PartnerIDs),
				),
			})
		return nil
	})

	return resp, err
}

func (uc *DistributionUsecase) UpdateDistribute(
	ctx context.Context,
	distID uint64,
	req *dto.DistributeDTO,
) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(); err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}

	var resp *dto.DistributeResultDTO
	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		dist, err := uc.distributionRepo.FindByID(txCtx, distID)
		if err != nil {
			return _errors.InternalServerException("load distribution failed")
		}
		if dist == nil {
			return _errors.NotFoundException("distribution not found")
		}

		if dist.RevokedAt != nil || time.Now().After(*dist.ToDate) {
			return _errors.BadRequestException("distribution has been revoked")
		}

		dist.Price = req.Price
		dist.CanDeal = req.CanDeal
		dist.ChannelPrice = req.ChannelPrice
		dist.CommissionType = req.CommissionType
		dist.CommissionValue = req.CommissionValue
		dist.PartnerCommission = calculatePartnerCommission(
			req.Price,
			req.CommissionType,
			req.CommissionValue,
		)
		dist.Note = req.Note
		dist.FromDate = req.FromDate
		dist.ToDate = req.ToDate

		if time.Now().After(*dist.ToDate) {
			resp.Status = enums.DISTRIBUTION_STATUS_EXPRIED
		}

		if err := uc.distributionRepo.Save(txCtx, dist); err != nil {
			return _errors.InternalServerException("update distribution failed")
		}

		if err := uc.distributionRepo.AssignPartner(
			txCtx,
			dist.ID,
			dist.ProductID,
			req.PartnerIDs,
		); err != nil {
			return err
		}

		resp = mapToResultDTO(dist, len(req.PartnerIDs))
		return nil
	})

	return resp, err
}

func mapToResultDTO(dist *domain.DistributionEntity, createdCount int) *dto.DistributeResultDTO {
	return &dto.DistributeResultDTO{
		Price:             dist.Price,
		CanDeal:           dist.CanDeal,
		ChannelPrice:      dist.ChannelPrice,
		CommissionType:    dist.CommissionType,
		CommissionValue:   dist.CommissionValue,
		Note:              dist.Note,
		PartnerCommission: dist.PartnerCommission,
		Status:            enums.DISTRIBUTION_STATUS_ACTIVE,
		CreatedCount:      createdCount,
		FromDate:          dist.FromDate,
		ToDate:            dist.ToDate,
	}
}

func (uc *DistributionUsecase) RevokeDistribute(ctx context.Context, distributeID uint64, revoke dto.RevokeDTO) error {
	now := time.Now()
	if err := uc.distributionRepo.Revoke(ctx, distributeID, now, revoke); err != nil {
		return _errors.InternalServerException("Failed to revoke distribution")
	}
	return nil
}

func (uc *DistributionUsecase) GetDistributeWithPartners(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error) {
	dist, err := uc.distributionRepo.FindByIdWithPartner(ctx, disID)
	if err != nil {
		return nil, _errors.InternalServerException("Failed to fetch distribution")
	}
	if dist == nil {
		return nil, _errors.NotFoundException("Distribution not found")
	}
	dist.Status = displayStatus(dist)

	return dist, nil
}

func (uc *DistributionUsecase) GetListDistributesWithPartners(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error) {
	distributes, err := uc.distributionRepo.GetListDistrByProduct(ctx, productID, filter)
	if err != nil {
		return nil, _errors.InternalServerException("Failed to get distribution list")
	}

	for _, dist := range distributes {
		dist.Status = displayStatus(dist)
	}
	return distributes, nil
}

func displayStatus(dist *dto.DistributionWithPartners) enums.DistributionStatus {
	now := time.Now()
	if dist.RevokedAt != nil {
		return enums.DISTRIBUTION_STATUS_REVOKED
	}

	if dist.ToDate != nil && now.After(*dist.ToDate) {
		return enums.DISTRIBUTION_STATUS_EXPRIED
	}

	return enums.DISTRIBUTION_STATUS_ACTIVE
}

func (uc *DistributionUsecase) GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*domain.ProductUser, error) {
	partners, err := uc.distributionRepo.GetListPartnersOfDistribute(ctx, disID)
	if err != nil {
		return nil, _errors.InternalServerException("Failed to get partners")
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
