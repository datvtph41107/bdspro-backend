package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	delivery_dto "bdspro/internal/dto/delivery"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	iusecase "bdspro/internal/usecases/interface"
	_db "common/db"
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
	Client         iusecase.IClientUC
}

func NewShareAccessUsecase(
	repo repo.SharingAccessRepo,
	ProductUc *ProductUsecase,
	Client iusecase.IClientUC,
) *SharingAccessUsecase {
	return &SharingAccessUsecase{
		Repo:           repo,
		ProductUsecase: ProductUc,
		Client:         Client,
	}
}

func (s *SharingAccessUsecase) CreateBatch(c context.Context, domainId uint64, ownerType enums.EOwnerType, dtos []dto.SharingAccessRequest) (*[]domain.SharingAccess, error) {
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

	data := []delivery_dto.NotiDTO{}
	for i := 0; i < len(entities); i++ {
		data = append(data, delivery_dto.NotiDTO{
			Title:   "Chia sẻ sản phẩm",
			Message: "Bạn nhận được chia sẻ sản phẩm ",
			Type:    1,
			UserID:  entities[i].ToId,
		})
		// entities[i].ProductID = id
	}

	s.Client.SendBatch(c, &delivery_dto.NotiBatchDTO{
		Datas: data,
	})

	return &entities, err
}

func (s *SharingAccessUsecase) BulkSave(c context.Context, id uint64, ownerType enums.EOwnerType, dto dto.SharingAccessBulk) (*[]domain.SharingAccess, error) {
	var entities []domain.SharingAccess
	product, err := s.ProductUsecase.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	// đổi trạng thái sang nội bộ nếu riêng tư
	if product.SaleVisibility == enums.EVisiblePrivate {
		product.SaleVisibility = enums.EVisibleInternal
		_db.SaveWithContext(c, product)
	}
	if len(dto.Datas) > 0 {
		copier.Copy(&entities, &dto.Datas)
		for i := 0; i < len(entities); i++ {
			entities[i].DomainID = id
			entities[i].Domain = enums.EDomainAccess(ownerType)
			entities[i].DeletedAt = nil
		}

		assignments := append(
			clause.AssignmentColumns([]string{"fields", "commission"}),
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
		query := _db.DB.Model(&domain.SharingAccess{})
		for i, cond := range dto.Deletes {
			if i == 0 {
				query = query.Where("domain_id = ? AND domain = ? AND to_id = ? AND to_type = ?",
					id, ownerType, cond.ToID, cond.ToType)
			} else {
				query = query.Or("domain_id = ? AND domain = ? AND to_id = ? AND to_type = ?",
					id, ownerType, cond.ToID, cond.ToType)
			}
		}
		err = query.Update("deleted_at", time.Now()).Error
		if err != nil {
			return nil, err
		}
	}

	data := []delivery_dto.NotiDTO{}
	for i := 0; i < len(entities); i++ {
		data = append(data, delivery_dto.NotiDTO{
			Title:   "Chia sẻ sản phẩm",
			Message: "Bạn nhận được chia sẻ sản phẩm ",
			Type:    1,
			UserID:  entities[i].ToId,
		})
		// entities[i].ProductID = id
	}

	s.Client.SendBatch(c, &delivery_dto.NotiBatchDTO{
		Datas: data,
	})

	return &entities, err
}

func (s *SharingAccessUsecase) SearchForProduct(c context.Context, ownerType enums.EOwnerType, dto dto.SharingAccessSearch) (*[]domain.ProductAccessQuery, error) {
	// var entities []domain.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.ProductUsecase.RequiredOwner(c, dto.DomainID)
	if err != nil {
		return nil, err
	}
	results, err := s.Repo.OwnerSearchProduct(c, dto.DomainID, ownerType, &dto)

	return results, err
}

func (s *SharingAccessUsecase) UpdateCommission(c context.Context, dto dto.ProductAccessUpdate) error {
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

func (s *SharingAccessUsecase) Delete(c context.Context, id uint64) error {
	// var entities []domain.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.RequiredOwner(c, id)
	if err != nil {
		return err
	}

	// profileId := _utils.GetProfileIdWithContext(c)

	return s.Repo.Delete(c, id)
}
