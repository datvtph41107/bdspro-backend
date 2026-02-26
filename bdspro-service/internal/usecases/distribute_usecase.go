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
	_util "common/domain/util"
	_errors "common/errors"
	_utils "common/utils"
)

type DistributionUsecase struct {
	distributionRepo repo.DistributeRepository
	productPriceRepo repo.ProductPriceRepo
	productRepo      repo.ProductRepo
	productUserRepo  repo.ProductUserRepo
	transaction      provider.TransactionProvider
	notifyProvider   provider.NotificationProvider
	hubProvider      provider.HubProvider
}

func NewDistributionUsecase(
	r repo.DistributeRepository,
	productPriceRepo repo.ProductPriceRepo,
	productRepo repo.ProductRepo,
	productUserRepo repo.ProductUserRepo,
	transaction provider.TransactionProvider,
	notificationProvider provider.NotificationProvider,
	hubProvider provider.HubProvider,
) *DistributionUsecase {
	return &DistributionUsecase{
		distributionRepo: r,
		productPriceRepo: productPriceRepo,
		productRepo:      productRepo,
		transaction:      transaction,
		notifyProvider:   notificationProvider,
		hubProvider:      hubProvider,
	}
}

func (uc *DistributionUsecase) CreateDistribute(
	ctx context.Context,
	req *dto.DistributeDTO,
) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(false); err != nil {
		return nil, _errors.BadRequestException("Validate: ", err.Error())
	}

	var resp *dto.DistributeResultDTO
	err := uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()

		dist := &domain.DistributionEntity{
			ProductID:     req.ProductID,
			CanDeal:       req.CanDeal,
			Note:          req.Note,
			FromDate:      req.FromDate,
			ToDate:        req.ToDate,
			DistributedAt: &now,
		}

		if err := uc.distributionRepo.Save(txCtx, dist); err != nil {
			log.Printf("create distribution failed: %v", err)
			return _errors.InternalServerException("create distribution failed")
		}

		// Tạo ProductPrice mới
		productPrice := uc.createProductPriceFromDistribute(
			req.ProductID,
			dist.ID,
			req.Price,
			req.CommissionType,
			req.CommissionValue,
			req.ChannelPrice,
			_utils.GetProfileIdWithContext(ctx),
		)
		if err := uc.productPriceRepo.Create(txCtx, productPrice); err != nil {
			log.Printf("create product price failed: %v", err)
			return _errors.InternalServerException("create product price failed")
		}

		// Cập nhật PriceID vào distribution
		dist.PriceID = &productPrice.ID
		if err := uc.distributionRepo.Save(txCtx, dist); err != nil {
			log.Printf("update distribution priceId failed: %v", err)
			return _errors.InternalServerException("update distribution priceId failed")
		}

		if err := uc.distributionRepo.AssignPartner(
			txCtx,
			dist.ID,
			dist.ProductID,
			req.PartnerIDs,
		); err != nil {
			return err
		}

		resp = uc.mapToResultDTO(txCtx, dist, productPrice, len(req.PartnerIDs))

		// Lấy PropertyID từ Product
		var propertyID uint64
		product, err := uc.productRepo.GetProductByID(txCtx, &dist.ProductID)
		if err != nil || product == nil {
			log.Printf("get product failed: %v", err)
		} else if product.PropertyID != nil {
			propertyID = *product.PropertyID
		}

		uc.notifyProvider.RegistedEventProperty(ctx,
			&dto.CreatePropertyActivityDTO{
				SubjectID:   propertyID,
				SubjectType: _enum.SUBJECT_TYPE_DISTRIBUTE,
				Action:      _enum.ACTION_DISTRIBUTE_CREATE,
				ActorID:     _utils.GetProfileIdWithContext(ctx),
				Description: fmt.Sprintf(
					"Sản phẩm %s được tạo kênh phân phối với %d đối tác",
					_util.GetProductCode(dist.ProductID),
					len(req.PartnerIDs),
				),
			})
		go func() {
			cloneCtx := _utils.CloneContext(ctx)
			// actorID := _utils.GetProfileIdWithContext(cloneCtx)
			productUsers, err := uc.productUserRepo.GetByDistributeID(cloneCtx, dist.ID)
			if err != nil {
				return
			}
			ids := make([]uint64, len(productUsers))
			for idx, productUser := range productUsers {
				if productUser.OriginProfileID != nil {
					ids[idx] = *productUser.OriginProfileID
				}
			}
			uc.hubProvider.PutUpdate(cloneCtx, "distribute", dist.ID, ids, time.Now().Unix())
		}()
		return nil
	})

	return resp, err
}

func (uc *DistributionUsecase) UpdateDistribute(
	ctx context.Context,
	distID uint64,
	req *dto.DistributeDTO,
) (*dto.DistributeResultDTO, error) {
	if err := req.Validate(true); err != nil {
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

		// Kiểm tra xem Price hoặc CommissionValue có thay đổi không
		shouldCreateNewPrice := false
		var oldPrice *domain.ProductPrice
		if dist.PriceID != nil {
			oldPrice, err = uc.productPriceRepo.GetByID(txCtx, *dist.PriceID)
			if err != nil {
				log.Printf("get old product price failed: %v", err)
			}
		}

		// So sánh Price và CommissionValue
		if oldPrice != nil {
			oldPriceValue := int64(0)
			if oldPrice.SalePrice != nil {
				oldPriceValue = int64(*oldPrice.SalePrice)
			}
			oldCommissionType := enums.CommissionType(0)
			if oldPrice.SaleCommissionType != nil {
				oldCommissionType = enums.CommissionType(*oldPrice.SaleCommissionType)
			}

			// Tính CommissionValue cũ từ SaleCommission và SalePrice
			oldCommissionValue := float64(0)
			if oldPrice.SalePrice != nil && oldPrice.SaleCommission != nil {
				if oldCommissionType == enums.CommissionTypePercent {
					// CommissionValue = (SaleCommission / SalePrice) * 100
					oldCommissionValue = (*oldPrice.SaleCommission / *oldPrice.SalePrice) * 100
				} else if oldCommissionType == enums.CommissionTypeVND {
					oldCommissionValue = *oldPrice.SaleCommission
				}
			}

			// Kiểm tra xem Price hoặc CommissionValue có thay đổi không
			if oldPriceValue != req.Price || oldCommissionType != req.CommissionType || oldCommissionValue != req.CommissionValue {
				shouldCreateNewPrice = true
			}
		} else {
			// Nếu chưa có PriceID, luôn tạo mới
			shouldCreateNewPrice = true
		}

		dist.CanDeal = req.CanDeal
		dist.Note = req.Note
		dist.FromDate = req.FromDate
		dist.ToDate = req.ToDate

		if time.Now().After(*dist.ToDate) {
			resp.Status = enums.DISTRIBUTION_STATUS_EXPRIED
		}

		// Tạo ProductPrice mới nếu Price, CommissionValue hoặc ChannelPrice thay đổi
		var productPrice *domain.ProductPrice
		if shouldCreateNewPrice {
			productPrice = uc.createProductPriceFromDistribute(
				dist.ProductID,
				dist.ID,
				req.Price,
				req.CommissionType,
				req.CommissionValue,
				req.ChannelPrice,
				_utils.GetProfileIdWithContext(ctx),
			)
			if err := uc.productPriceRepo.Create(txCtx, productPrice); err != nil {
				log.Printf("create product price failed: %v", err)
				return _errors.InternalServerException("create product price failed")
			}

			// Cập nhật PriceID vào distribution
			dist.PriceID = &productPrice.ID
		} else {
			// Lấy ProductPrice hiện tại nếu không tạo mới
			if dist.PriceID != nil {
				productPrice, err = uc.productPriceRepo.GetByID(txCtx, *dist.PriceID)
				if err != nil {
					log.Printf("get product price failed: %v", err)
				} else if productPrice != nil {
					// Cập nhật ChannelPrice vào ProductPrice hiện tại
					productPrice.ChannelPrice = req.ChannelPrice
					if err := uc.productPriceRepo.UpdatePrice(txCtx, dist.PriceID, productPrice); err != nil {
						log.Printf("update product price channel price failed: %v", err)
					}
				}
			}
		}

		if err := uc.distributionRepo.Save(txCtx, dist); err != nil {
			return _errors.InternalServerException("update distribution failed")
		}

		// if err := uc.distributionRepo.AssignPartner(
		// 	txCtx,
		// 	dist.ID,
		// 	dist.ProductID,
		// 	req.PartnerIDs,
		// ); err != nil {
		// 	return err
		// }

		resp = uc.mapToResultDTO(txCtx, dist, productPrice, len(req.PartnerIDs))

		// Lấy PropertyID từ Product
		var propertyID uint64
		product, err := uc.productRepo.GetProductByID(txCtx, &dist.ProductID)
		if err != nil || product == nil {
			log.Printf("get product failed: %v", err)
		} else if product.PropertyID != nil {
			propertyID = *product.PropertyID
		}

		uc.notifyProvider.RegistedEventProperty(ctx,
			&dto.CreatePropertyActivityDTO{
				SubjectID:   propertyID,
				SubjectType: _enum.SUBJECT_TYPE_DISTRIBUTE,
				Action:      _enum.ACTION_DISTRIBUTE_UPDATE,
				ActorID:     _utils.GetProfileIdWithContext(ctx),
				Description: fmt.Sprintf(
					"Sản phẩm %s được cập nhật kênh phân phối",
					_util.GetProductCode(dist.ProductID),
				),
			})

		go func() {
			cloneCtx := _utils.CloneContext(ctx)
			// actorID := _utils.GetProfileIdWithContext(cloneCtx)
			productUsers, err := uc.productUserRepo.GetByDistributeID(cloneCtx, dist.ID)
			if err != nil {
				return
			}
			ids := make([]uint64, len(productUsers))
			for idx, productUser := range productUsers {
				if productUser.OriginProfileID != nil {
					ids[idx] = *productUser.OriginProfileID
				}
			}
			uc.hubProvider.PutUpdate(cloneCtx, "distribute", dist.ID, ids, time.Now().Unix())
		}()
		return nil
	})

	return resp, err
}

func (uc *DistributionUsecase) mapToResultDTO(
	ctx context.Context,
	dist *domain.DistributionEntity,
	productPrice *domain.ProductPrice,
	createdCount int,
) *dto.DistributeResultDTO {
	result := &dto.DistributeResultDTO{
		CanDeal:      dist.CanDeal,
		Note:         dist.Note,
		Status:       enums.DISTRIBUTION_STATUS_ACTIVE,
		CreatedCount: createdCount,
		FromDate:     dist.FromDate,
		ToDate:       dist.ToDate,
	}

	// Lấy giá trị từ ProductPrice nếu có
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
		// Tính CommissionValue từ SaleCommission và SalePrice
		if productPrice.SalePrice != nil && productPrice.SaleCommission != nil && productPrice.SaleCommissionType != nil {
			commissionType := enums.CommissionType(*productPrice.SaleCommissionType)
			if commissionType == enums.CommissionTypePercent {
				// CommissionValue = (SaleCommission / SalePrice) * 100
				result.CommissionValue = (*productPrice.SaleCommission / *productPrice.SalePrice) * 100
			} else if commissionType == enums.CommissionTypeVND {
				result.CommissionValue = *productPrice.SaleCommission
			}
		}
	}

	return result
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

		// if dist.Price != 0 {
		// 	dist.Price = dist.Price
		// }
		// if dist.CommissionType != 0 {
		// 	dist.CommissionType = dist.CommissionType
		// }
		// if dist.CommissionValue != 0 {
		// 	dist.CommissionValue = dist.CommissionValue
		// }
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

func (uc *DistributionUsecase) RevokeUserFromDistribute(ctx context.Context, distributeID uint64, originProfileID uint64) error {
	return uc.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// Xóa user khỏi product_user và lấy số lượng còn lại
		remainingCount, err := uc.distributionRepo.RevokeUser(txCtx, distributeID, originProfileID)
		if err != nil {
			return _errors.InternalServerException(err.Error())
		}

		// Nếu distribute chỉ còn 0 product_user, revoke distribute
		if remainingCount == 0 {
			now := time.Now()
			revoke := dto.RevokeDTO{
				Reason:      uint32(enums.DISTRIBUTION_REASON_OTHER),
				ReasonOther: "Xóa đối tác",
			}
			if err := uc.distributionRepo.Revoke(txCtx, distributeID, now, revoke); err != nil {
				return _errors.InternalServerException("Failed to revoke distribute")
			}
		}

		return nil
	})
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

// createProductPriceFromDistribute tạo ProductPrice từ thông tin distribution
func (uc *DistributionUsecase) createProductPriceFromDistribute(
	productID uint64,
	distributeID uint64,
	price int64,
	commissionType enums.CommissionType,
	commissionValue float64,
	channelPrice bool,
	changedBy uint64,
) *domain.ProductPrice {
	priceFloat := float64(price)
	// commissionFloat := calculatePartnerCommission(price, commissionType, commissionValue)
	commissionTypeUint32 := uint32(commissionType)

	return &domain.ProductPrice{
		ProductID:          &productID,
		DistributeID:       &distributeID,
		Currency:           "VND",
		SalePrice:          &priceFloat,
		SaleCommission:     &commissionValue,
		SaleCommissionType: &commissionTypeUint32,
		ChannelPrice:       channelPrice,
		PriceStatus:        enums.PriceStatusFixed,
		ChangeNote:         "Tạo từ kênh phân phối",
		ChangedBy:          changedBy,
	}
}
