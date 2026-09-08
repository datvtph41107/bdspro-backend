package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_db "common/db"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SharingAccessUsecase struct {
	Repo           repo.SharingAccessRepo
	ProductUsecase *ProductUsecase
	Client         provider.NotificationProvider
	ChatProvider   _provider.ChatProvider
	AssetRepo      repo.AssetRepo
}

func NewShareAccessUsecase(
	repo repo.SharingAccessRepo,
	ProductUc *ProductUsecase,
	Client provider.NotificationProvider,
	chatClient _provider.ChatProvider,
	assetRepo repo.AssetRepo,
) *SharingAccessUsecase {
	return &SharingAccessUsecase{
		Repo:           repo,
		ProductUsecase: ProductUc,
		Client:         Client,
		ChatProvider:   chatClient,
		AssetRepo:      assetRepo,
	}
}

func (s *SharingAccessUsecase) CreateBatch(c context.Context, domainId uint64, ownerType enums.EOwnerOf, dtos []dto.SharingAccessRequest) (*[]domain.SharingAccess, error) {
	var entities []domain.SharingAccess
	product, err := s.ProductUsecase.RequiredOwner(c, domainId)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entities, &dtos)
	for i := 0; i < len(entities); i++ {
		entities[i].DomainID = domainId
	}

	// đổi trạng thái sang nội bộ nếu riêng tư
	if product.SaleVisibility == enums.EVisiblePrivate {
		product.SaleVisibility = enums.EVisibleInternal
		_db.SaveWithContext(c, product)
	}
	// err := _db.SaveWithAudit(c, &entities)
	err = _db.DB.WithContext(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}, {Name: "target_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"fields", "commission", "target_type"}),
		}).
		Save(&entities).Error

	data := []_dto.NotificationDTO{}
	for i := 0; i < len(entities); i++ {
		data = append(data, _dto.NotificationDTO{
			Title:   "Chia sẻ sản phẩm",
			Message: []string{"Bạn nhận được chia sẻ sản phẩm "},
			OwnerID: entities[i].ToId,
			OwnerOf: _enum.EOwnerOf(entities[i].ToType),
		})
	}

	go func() {
		ctxClone := _utils.CloneContext(c)
		s.Client.SendBatch(ctxClone, &_dto.NotificationBatchDTO{
			Datas: data,
		})

		for _, entity := range entities {
			s.ChatProvider.SendProductToReceiver(ctxClone,
				product.ID,
				"",
				&entity.ToId,
			)
		}
	}()

	return &entities, err
}

func (s *SharingAccessUsecase) CanAccessAsset(c context.Context, assetId uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	asset, err := s.AssetRepo.GetByID(assetId)
	if err != nil {
		return err
	}
	if asset.CreatedBy == nil || *asset.CreatedBy != profileId {
		return &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}

	return nil
}

func (s *SharingAccessUsecase) BulkSave(c context.Context,
	id uint64,
	ownerType enums.EOwnerOf,
	dto *dto.SharingAccessBulk) (*[]domain.SharingAccess, error) {
	var entities []domain.SharingAccess
	var err error
	if dto.Domain == enums.EDomainProduct {
		product, err := s.ProductUsecase.RequiredOwner(c, id)
		if err != nil {
			return nil, err
		}
		// đổi trạng thái sang nội bộ nếu riêng tư
		if product.SaleVisibility == enums.EVisiblePrivate {
			product.SaleVisibility = enums.EVisibleInternal
			_db.SaveWithContext(c, product)
		}
	} else {
		err = s.CanAccessAsset(c, id)
		if err != nil {
			return nil, err
		}
	}

	if len(dto.Datas) > 0 {
		copier.Copy(&entities, &dto.Datas)
		for i := 0; i < len(entities); i++ {
			entities[i].DomainID = id
			entities[i].Domain = enums.EDomainAccess(dto.Domain)
			entities[i].DeletedAt = nil
		}

		assignments := append(
			clause.AssignmentColumns([]string{"fields", "commission", "permissions"}),
			clause.Assignment{
				Column: clause.Column{Name: "deleted_at"},
				Value:  gorm.Expr("NULL"),
			},
		)

		err = _db.DB.WithContext(c).Debug().
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "domain_id"},
					{Name: "domain"},
					{Name: "to_id"},
					{Name: "to_type"},
				},
				DoUpdates: clause.Set(assignments),
			}).
			Save(&entities).Error
		if err != nil {
			return nil, err
		}
	}

	if len(dto.Deletes) > 0 {
		query := _db.DB.Debug().Model(&domain.SharingAccess{})
		for i, cond := range dto.Deletes {
			if i == 0 {
				query = query.Where("domain_id = ? AND domain = ? AND to_id = ? AND to_type = ?",
					id, dto.Domain, cond.ToId, cond.ToType)
			} else {
				query = query.Or("domain_id = ? AND domain = ? AND to_id = ? AND to_type = ?",
					id, dto.Domain, cond.ToId, cond.ToType)
			}
		}
		err = query.Update("deleted_at", time.Now()).Error
		if err != nil {
			return nil, err
		}
	}

	data := []_dto.NotificationDTO{}
	for i := 0; i < len(entities); i++ {
		data = append(data, _dto.NotificationDTO{
			Title:   "Chia sẻ sản phẩm",
			Message: []string{"Bạn nhận được chia sẻ sản phẩm "},
			OwnerID: entities[i].ToId,
			OwnerOf: _enum.EOwnerOf(entities[i].ToType),
		})
		// entities[i].ProductID = id
	}

	go func() {
		ctxClone := _utils.CloneContext(c)
		s.Client.SendBatch(ctxClone, &_dto.NotificationBatchDTO{
			Datas: data,
		})

		if dto.Domain == enums.EDomainProduct {
			for _, entity := range entities {
				s.ChatProvider.SendProductToReceiver(ctxClone,
					entity.DomainID,
					"",
					&entity.ToId,
				)
			}
		}
	}()

	// Đếm số sharing access của product và cập nhật ShareCount
	if dto.Domain == enums.EDomainProduct {
		shareCount, countErr := s.Repo.CountSharingAccessByProductID(c, id)
		if countErr != nil {
			// Log error nhưng không fail toàn bộ operation
			// return nil, countErr
		} else {
			updateErr := s.ProductUsecase.ProductRepo.UpdateShareCount(c, id, shareCount)
			if updateErr != nil {
				// Log error nhưng không fail toàn bộ operation
				// return nil, updateErr
			}
		}
	}

	return &entities, err
}

func (s *SharingAccessUsecase) SearchForProduct(c context.Context, ownerType enums.EOwnerOf, dto *dto.SharingAccessSearch) (*[]domain.SharingAccess, error) {
	_, err := s.ProductUsecase.RequiredOwner(c, dto.DomainID)
	if err != nil {
		return nil, err
	}
	dto.Domain = enums.EDomainProduct
	results, err := s.Repo.OwnerSearchProduct(c, dto.DomainID, ownerType, dto)

	return results, err
}

func (s *SharingAccessUsecase) SearchForAsset(ctx context.Context,
	ownerType enums.EOwnerOf,
	dto *dto.SharingAccessSearch,
) (*[]domain.SharingAccess, error) {

	err := s.CanAccessAsset(ctx, dto.DomainID)
	if err != nil {
		return nil, err
	}
	// results, err := s.Repo.OwnerSearchAsset(ctx, dto.DomainID, ownerType, &dto)

	// todo: kiếm tra xem có quyền chia sẻ không

	dto.Domain = enums.EDomainAsset
	results, err := s.Repo.OwnerSearchAsset(ctx, dto.DomainID, ownerType, dto)

	return results, err
}

func (s *SharingAccessUsecase) UpdateCommission(c context.Context, dto dto.CommissionUpdate) error {
	// var entities []domain.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.ProductUsecase.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return err
	}

	profileId := _utils.GetProfileIdWithContext(c)

	return s.Repo.UpdateCommission(c, profileId, dto)
}

func (s *SharingAccessUsecase) RequiredOwner(c context.Context, id uint64) (*domain.SharingAccess, error) {
	asset, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	// profileId := _jwt.GetProfileId(c)
	profileId := _utils.GetProfileIdWithContext(c)

	if asset.CreatedBy == nil || *asset.CreatedBy != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return asset, nil
}

func (s *SharingAccessUsecase) ProductDelete(c context.Context, id uint64) error {
	_, err := s.RequiredOwner(c, id)
	if err != nil {
		return err
	}

	// profileId := _utils.GetProfileIdWithContext(c)

	return s.Repo.Delete(c, id)
}
