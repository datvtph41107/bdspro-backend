package u_product_access

import (
	"bdspro/internal/enums"
	product_accecss "bdspro/v1/mod/global/product_access"
	user_dto "bdspro/v1/mod/user/dto/request"
	u_product "bdspro/v1/mod/user/product"
	"bdspro/v1/models"
	"bdspro/v1/providers"
	"common/case/crud2"
	_db "common/db"
	pb_noti "pb/types/notification"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductAccessService struct {
	crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]
	UserProductService *u_product.UserProductService
	NotiService        *providers.TransmitProvider
}

func NewUserProductAccessService(
	repo *product_accecss.ProductAccessRepo,
	UserProductService *u_product.UserProductService,
	NotiService *providers.TransmitProvider,
) *ProductAccessService {
	return &ProductAccessService{
		BaseService:        crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]{Repo: repo},
		UserProductService: UserProductService,
		NotiService:        NotiService,
	}
}

func (s *ProductAccessService) CreateBatch(c *gin.Context, id uint64, dtos []AccessRequest) (*[]models.ProductAccess, error) {
	var entities []models.ProductAccess
	product, err := s.UserProductService.RequiredOwnerProduct(c, id)
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
		_db.SaveWithAudit(c, product)
	}
	// err := _db.SaveWithAudit(c, &entities)
	err = _db.DB.WithContext(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}, {Name: "target_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"fields", "commission", "target_type"}),
		}).
		Save(&entities).Error

	data := []*pb_noti.NotiNewRequest{}
	for i := 0; i < len(entities); i++ {
		data = append(data, &pb_noti.NotiNewRequest{
			Title:   "Chia sẻ sản phẩm",
			Message: "Bạn nhận được chia sẻ sản phẩm ",
			Type:    1,
			UserID:  entities[i].TargetId,
		})
		// entities[i].ProductID = id
	}

	s.NotiService.NotiClient.SendBatch(c, &pb_noti.NotiBatchRequest{
		Datas: data,
	})

	return &entities, err
}

func (s *ProductAccessService) BulkSave(c *gin.Context, id uint64, dto BulkRequest) (*[]models.ProductAccess, error) {
	var entities []models.ProductAccess
	product, err := s.UserProductService.RequiredOwnerProduct(c, id)
	if err != nil {
		return nil, err
	}

	// đổi trạng thái sang nội bộ nếu riêng tư
	if product.SaleVisibility == enums.EVisiblePrivate {
		product.SaleVisibility = enums.EVisibleInternal
		_db.SaveWithAudit(c, product)
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
		query := _db.DB.Model(&models.ProductAccess{})
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

	data := []*pb_noti.NotiNewRequest{}
	for i := 0; i < len(entities); i++ {
		data = append(data, &pb_noti.NotiNewRequest{
			Title:   "Chia sẻ sản phẩm",
			Message: "Bạn nhận được chia sẻ sản phẩm ",
			Type:    1,
			UserID:  entities[i].TargetId,
		})
		// entities[i].ProductID = id
	}

	s.NotiService.NotiClient.SendBatch(c, &pb_noti.NotiBatchRequest{
		Datas: data,
	})

	return &entities, err
}

func (s *ProductAccessService) Search(c *gin.Context, dto user_dto.ProductAccessSearch) (*[]models.ProductAccessQuery, error) {
	// var entities []models.ProductAccess
	// productId := *dto.ProductID
	// if dto.ProductID
	_, err := s.UserProductService.RequiredOwnerProduct(c, dto.ProductID)
	if err != nil {
		return nil, err
	}
	results, err := s.Repo.OwnerGetByProduct(c, dto.ProductID, dto.TargetType)

	return results, err
}
