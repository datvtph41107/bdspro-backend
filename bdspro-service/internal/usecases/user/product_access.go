package user_usecase

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

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserProductAccessUC struct {
	// crud2.BaseService[domain.ProductAccess, *repo.ProductAccessRepo]
	Repo       repo.ProductAccessRepo
	UProductUC *UserProductUC
	// NotiService *providers.TransmitProvider
	Client iusecase.IClientUC
}

func NewUserProductAccessUsecase(
	repo repo.ProductAccessRepo,
	UProductUC *UserProductUC,
	Client iusecase.IClientUC,
) *UserProductAccessUC {
	return &UserProductAccessUC{
		Repo:       repo,
		UProductUC: UProductUC,
		Client:     Client,
	}
}

func (s *UserProductAccessUC) CreateBatch(c context.Context, id uint64, dtos []dto.ProductAccessRequest) (*[]domain.ProductAccess, error) {
	var entities []domain.ProductAccess
	product, err := s.UProductUC.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entities, &dtos)
	for i := 0; i < len(entities); i++ {
		entities[i].ProductID = id
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
			UserID:  entities[i].TargetId,
		})
		// entities[i].ProductID = id
	}

	s.Client.SendBatch(c, &delivery_dto.NotiBatchDTO{
		Datas: data,
	})

	return &entities, err
}

func (s *UserProductAccessUC) BulkSave(c context.Context, id uint64, dto dto.ProductAccessBulkRequest) (*[]domain.ProductAccess, error) {
	var entities []domain.ProductAccess
	product, err := s.UProductUC.RequiredOwner(c, id)
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
			entities[i].ProductID = id
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
					{Name: "product_id"},
					{Name: "target_id"},
					{Name: "target_type"},
				},
				DoUpdates: clause.Set(assignments),
			}).
			Save(&entities).Error
		if err != nil {
			return nil, err
		}
	}

	if len(dto.Deletes) > 0 {
		query := _db.DB.Model(&domain.ProductAccess{})
		for i, cond := range dto.Deletes {
			if i == 0 {
				query = query.Where("target_id = ? AND target_type = ?", cond.TargetID, cond.TargetType)
			} else {
				query = query.Or("target_id = ? AND target_type = ?", cond.TargetID, cond.TargetType)
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
			UserID:  entities[i].TargetId,
		})
		// entities[i].ProductID = id
	}

	s.Client.SendBatch(c, &delivery_dto.NotiBatchDTO{
		Datas: data,
	})

	return &entities, err
}

func (s *UserProductAccessUC) Search(c *gin.Context, dto dto.ProductAccessSearch) (*[]domain.ProductAccessQuery, error) {
	// var entities []domain.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.UProductUC.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}
	results, err := s.Repo.OwnerSearchProduct(c, dto.ProductID, &dto)

	return results, err
}

func (s *UserProductAccessUC) UpdateCommission(c context.Context, dto dto.ProductAccessUpdate) error {
	// var entities []domain.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.UProductUC.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return err
	}

	profileId := _utils.GetProfileIdWithContext(c)

	return s.Repo.UpdateCommission(c, profileId, dto)
}

func (s *UserProductAccessUC) RequiredOwner(c context.Context, id uint64) (*domain.ProductAccess, error) {
	asset, err := s.Repo.GetByID(id)
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

func (s *UserProductAccessUC) Delete(c context.Context, id uint64) error {
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
