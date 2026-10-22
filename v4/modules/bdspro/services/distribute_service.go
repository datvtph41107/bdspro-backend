package services

import (
	"context"
	"log"
	"time"

	_util "common/domain/util"
	_errors "common/errors"
	_utils "common/utils"
	iusecase "mdn4/internal/interface"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
	"mdn4/modules/bdspro/repository"
)

type DistributionService interface {
	CreateDistribute(ctx context.Context, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error)
	UpdateDistribute(ctx context.Context, distID uint64, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error)
	RevokeDistribute(ctx context.Context, distributeID uint64, revoke dto.RevokeDTO) error
	GetDistributeWithPartners(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error)
	GetListDistributesWithPartners(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error)
	GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*entity.ProductUser, error)
	RevokeUserFromDistribute(ctx context.Context, distributeID, originProfileID uint64) error
}

// @bind: mdn4/modules/bdspro/services.DistributionService
type LocalDistributionService struct {
	distributionRepo repository.DistributeRepository
	productPriceRepo repository.ProductPriceRepository
	transaction      iusecase.ITransaction
}

func NewLocalDistributionService(
	distributionRepo repository.DistributeRepository,
	productPriceRepo repository.ProductPriceRepository,
	transaction iusecase.ITransaction,
) *LocalDistributionService {
	return &LocalDistributionService{
		distributionRepo: distributionRepo,
		productPriceRepo: productPriceRepo,
		transaction:      transaction,
	}
}

func (s *LocalDistributionService) CreateDistribute(ctx context.Context, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(false); err != nil {
		return nil, _errors.BadRequestException("Validate: ", err.Error())
	}
	var resp *dto.DistributeResultDTO
	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		dist := &entity.Distribution{
			ProductID:     req.ProductID,
			CanDeal:       req.CanDeal,
			Note:          req.Note,
			FromDate:      req.FromDate,
			ToDate:        req.ToDate,
			DistributedAt: &now,
		}
		if err := s.distributionRepo.Save(txCtx, dist); err != nil {
			log.Printf("create distribution failed: %v", err)
			return _errors.InternalServerException("create distribution failed")
		}
		productPrice := s.createProductPriceFromDistribute(
			req.ProductID, dist.ID, req.Price, req.CommissionType, req.CommissionValue, req.ChannelPrice,
			_utils.GetProfileIdWithContext(ctx),
		)
		if err := s.productPriceRepo.Create(txCtx, productPrice); err != nil {
			log.Printf("create product price failed: %v", err)
			return _errors.InternalServerException("create product price failed")
		}
		dist.PriceID = &productPrice.ID
		if err := s.distributionRepo.Save(txCtx, dist); err != nil {
			return _errors.InternalServerException("update distribution priceId failed")
		}
		if err := s.distributionRepo.AssignPartner(txCtx, dist.ID, dist.ProductID, req.PartnerIDs); err != nil {
			return err
		}
		resp = s.mapToResultDTO(dist, productPrice, len(req.PartnerIDs))
		log.Printf("distribution created product=%s partners=%d", _util.GetProductCode(dist.ProductID), len(req.PartnerIDs))
		return nil
	})
	return resp, err
}

func (s *LocalDistributionService) UpdateDistribute(ctx context.Context, distID uint64, req *dto.DistributeDTO) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(true); err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}
	var resp *dto.DistributeResultDTO
	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		dist, err := s.distributionRepo.FindByID(txCtx, distID)
		if err != nil {
			return _errors.InternalServerException("load distribution failed")
		}
		if dist.RevokedAt != nil || (dist.ToDate != nil && time.Now().After(*dist.ToDate)) {
			return _errors.BadRequestException("distribution has been revoked")
		}
		shouldCreateNewPrice := false
		var oldPrice *entity.ProductPrice
		if dist.PriceID != nil {
			oldPrice, err = s.productPriceRepo.GetByID(txCtx, *dist.PriceID)
			if err != nil {
				log.Printf("get old product price failed: %v", err)
			}
		}
		if oldPrice != nil {
			oldPriceValue := int64(0)
			if oldPrice.SalePrice != nil {
				oldPriceValue = int64(*oldPrice.SalePrice)
			}
			oldCommissionType := enums.CommissionType(0)
			if oldPrice.SaleCommissionType != nil {
				oldCommissionType = enums.CommissionType(*oldPrice.SaleCommissionType)
			}
			oldCommissionValue := float64(0)
			if oldPrice.SalePrice != nil && oldPrice.SaleCommission != nil {
				if oldCommissionType == enums.CommissionTypePercent {
					oldCommissionValue = (*oldPrice.SaleCommission / *oldPrice.SalePrice) * 100
				} else if oldCommissionType == enums.CommissionTypeVND {
					oldCommissionValue = *oldPrice.SaleCommission
				}
			}
			if oldPriceValue != req.Price || oldCommissionType != req.CommissionType || oldCommissionValue != req.CommissionValue {
				shouldCreateNewPrice = true
			}
		} else {
			shouldCreateNewPrice = true
		}
		dist.CanDeal = req.CanDeal
		dist.Note = req.Note
		dist.FromDate = req.FromDate
		dist.ToDate = req.ToDate
		var productPrice *entity.ProductPrice
		if shouldCreateNewPrice {
			productPrice = s.createProductPriceFromDistribute(
				dist.ProductID, dist.ID, req.Price, req.CommissionType, req.CommissionValue, req.ChannelPrice,
				_utils.GetProfileIdWithContext(ctx),
			)
			if err := s.productPriceRepo.Create(txCtx, productPrice); err != nil {
				return _errors.InternalServerException("create product price failed")
			}
			dist.PriceID = &productPrice.ID
		} else if dist.PriceID != nil {
			productPrice, err = s.productPriceRepo.GetByID(txCtx, *dist.PriceID)
			if err == nil && productPrice != nil {
				productPrice.ChannelPrice = req.ChannelPrice
				_ = s.productPriceRepo.UpdatePrice(txCtx, dist.PriceID, productPrice)
			}
		}
		if err := s.distributionRepo.Save(txCtx, dist); err != nil {
			return _errors.InternalServerException("update distribution failed")
		}
		resp = s.mapToResultDTO(dist, productPrice, len(req.PartnerIDs))
		if dist.ToDate != nil && time.Now().After(*dist.ToDate) {
			resp.Status = enums.DISTRIBUTION_STATUS_EXPRIED
		}
		return nil
	})
	return resp, err
}

func (s *LocalDistributionService) mapToResultDTO(dist *entity.Distribution, productPrice *entity.ProductPrice, createdCount int) *dto.DistributeResultDTO {
	result := &dto.DistributeResultDTO{
		CanDeal:      dist.CanDeal,
		Note:         dist.Note,
		Status:       enums.DISTRIBUTION_STATUS_ACTIVE,
		CreatedCount: createdCount,
		FromDate:     dist.FromDate,
		ToDate:       dist.ToDate,
	}
	if productPrice != nil {
		result.ChannelPrice = productPrice.ChannelPrice
		if productPrice.SalePrice != nil {
			result.Price = int64(*productPrice.SalePrice)
		}
		if productPrice.SaleCommission != nil {
			result.PartnerCommission = *productPrice.SaleCommission
		}
		if productPrice.SaleCommissionType != nil {
			result.CommissionType = enums.CommissionType(*productPrice.SaleCommissionType)
		}
		if productPrice.SalePrice != nil && productPrice.SaleCommission != nil && productPrice.SaleCommissionType != nil {
			commissionType := enums.CommissionType(*productPrice.SaleCommissionType)
			if commissionType == enums.CommissionTypePercent {
				result.CommissionValue = (*productPrice.SaleCommission / *productPrice.SalePrice) * 100
			} else if commissionType == enums.CommissionTypeVND {
				result.CommissionValue = *productPrice.SaleCommission
			}
		}
	}
	return result
}

func (s *LocalDistributionService) RevokeDistribute(ctx context.Context, distributeID uint64, revoke dto.RevokeDTO) error {
	if err := s.distributionRepo.Revoke(ctx, distributeID, time.Now(), revoke); err != nil {
		return _errors.InternalServerException("Failed to revoke distribution")
	}
	return nil
}

func (s *LocalDistributionService) GetDistributeWithPartners(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error) {
	dist, err := s.distributionRepo.FindByIDWithPartner(ctx, disID)
	if err != nil {
		return nil, _errors.InternalServerException("Failed to fetch distribution")
	}
	dist.Status = displayDistributionStatus(dist)
	return dist, nil
}

func (s *LocalDistributionService) GetListDistributesWithPartners(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error) {
	list, err := s.distributionRepo.GetListByProduct(ctx, productID, filter)
	if err != nil {
		return nil, _errors.InternalServerException("Failed to get distribution list")
	}
	for _, dist := range list {
		dist.Status = displayDistributionStatus(dist)
	}
	return list, nil
}

func (s *LocalDistributionService) GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*entity.ProductUser, error) {
	return s.distributionRepo.GetListPartnersOfDistribute(ctx, disID)
}

func (s *LocalDistributionService) RevokeUserFromDistribute(ctx context.Context, distributeID, originProfileID uint64) error {
	return s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		remainingCount, err := s.distributionRepo.RevokeUser(txCtx, distributeID, originProfileID)
		if err != nil {
			return _errors.InternalServerException(err.Error())
		}
		if remainingCount == 0 {
			revoke := dto.RevokeDTO{
				Reason:      uint32(enums.DISTRIBUTION_REASON_OTHER),
				ReasonOther: "Xóa đối tác",
			}
			if err := s.distributionRepo.Revoke(txCtx, distributeID, time.Now(), revoke); err != nil {
				return _errors.InternalServerException("Failed to revoke distribute")
			}
		}
		return nil
	})
}

func displayDistributionStatus(dist *dto.DistributionWithPartners) enums.DistributionStatus {
	now := time.Now()
	if dist.RevokedAt != nil {
		return enums.DISTRIBUTION_STATUS_REVOKED
	}
	if dist.ToDate != nil && now.After(*dist.ToDate) {
		return enums.DISTRIBUTION_STATUS_EXPRIED
	}
	return enums.DISTRIBUTION_STATUS_ACTIVE
}

func (s *LocalDistributionService) createProductPriceFromDistribute(
	productID, distributeID uint64,
	price int64,
	commissionType enums.CommissionType,
	commissionValue float64,
	channelPrice bool,
	changedBy uint64,
) *entity.ProductPrice {
	priceFloat := float64(price)
	commissionTypeUint32 := uint32(commissionType)
	return &entity.ProductPrice{
		ProductID:          &productID,
		DistributeID:       &distributeID,
		Currency:           "VND",
		SalePrice:          &priceFloat,
		SaleCommission:     &commissionValue,
		SaleCommissionType: &commissionTypeUint32,
		ChannelPrice:       channelPrice,
		ChangeNote:         "Tạo từ kênh phân phối",
		ChangedBy:          changedBy,
	}
}
