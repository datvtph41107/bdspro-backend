package postgres

import (
	_dto "common/dto"
	"context"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"

	"gorm.io/gorm"
)

// type GroupDealModel struct {
// 	Id        uint64    `gorm:"primaryKey;autoIncrement"`
// 	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_deals_created_at"`
// 	UpdatedAt time.Time `gorm:"autoUpdateTime"`
// 	CreatedBy uint32    `gorm:"index:idx_group_deals_created_by"`
// 	UpdatedBy uint32

// 	OwnerId   uint32            `gorm:"not null;index:idx_group_deals_owner_id"`
// 	OwnerType enums.OwnerType   `gorm:"not null;index:idx_group_deals_owner_type"`
// 	Name      string            `gorm:"not null;index:idx_group_deals_name"`
// 	Amount    float64           `gorm:"not null"`
// 	Status    entity.DealStatus `gorm:"not null;index:idx_group_deals_status"`
// 	DealType  enums.DealType    `gorm:"not null;index:idx_group_deals_deal_type"`
// 	DeletedAt *time.Time        `gorm:"index"`
// 	IsDeleted bool              `gorm:"default:false"`
// }

// func (GroupDealModel) TableName() string {
// 	return "group_deals"
// }

// func GroupDealModelToEntity(model *GroupDealModel) *entity.Deal {
// 	return &entity.Deal{
// 		// Id:        model.Id,
// 		// CreatedAt: model.CreatedAt,
// 		// UpdatedAt: model.UpdatedAt,
// 		// CreatedBy: model.CreatedBy,
// 		// UpdatedBy: model.UpdatedBy,
// 		BaseEntity: _models.BaseEntity{
// 			ID:        model.Id,
// 			CreatedAt: &model.CreatedAt,
// 			UpdatedAt: &model.UpdatedAt,
// 			CreatedBy: &model.CreatedBy,
// 			UpdatedBy: &model.UpdatedBy,
// 		},
// 		GroupId: model.GroupId,
// 		Name:    model.Name,
// 		Amount:  model.Amount,
// 		Status:  model.Status,
// 	}
// }

// func GroupDealEntityToModel(entity *entity.Deal) *GroupDealModel {
// 	return &GroupDealModel{
// 		Id:        entity.Id,
// 		CreatedAt: entity.CreatedAt,
// 		UpdatedAt: entity.UpdatedAt,
// 		CreatedBy: entity.CreatedBy,
// 		UpdatedBy: entity.UpdatedBy,
// 		GroupId:   entity.GroupId,
// 		Name:      entity.Name,
// 		Amount:    entity.Amount,
// 		Status:    entity.Status,
// 	}
// }

// @bind: organization/internal/domain/repository.GroupDealRepository
type GroupDealPostgresRepository struct {
	db *gorm.DB
}

func NewGroupDealPostgresRepository(db *gorm.DB) *GroupDealPostgresRepository {
	return &GroupDealPostgresRepository{db: db}
}

func (r *GroupDealPostgresRepository) CreateDealProduct(ctx context.Context, dealProduct []*entity.DealProduct) ([]*entity.DealProduct, error) {
	if err := r.db.WithContext(ctx).Create(dealProduct).Error; err != nil {
		return nil, err
	}
	return dealProduct, nil
}

func (r *GroupDealPostgresRepository) GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*entity.DealProduct, error) {
	var dealProducts []*entity.DealProduct
	if err := r.db.WithContext(ctx).
		Where("deal_id = ?", dealID).
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&dealProducts).Error; err != nil {
		return nil, err
	}
	return dealProducts, nil
}

func (r *GroupDealPostgresRepository) Create(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	// model := GroupDealEntityToModel(deal)
	deal.Products = make([]entity.DealProduct, len(deal.ProductIds))
	for i, productID := range deal.ProductIds {
		deal.Products[i] = entity.DealProduct{
			DealID:    deal.ID,
			ProductID: productID,
		}
	}
	if err := r.db.WithContext(ctx).Create(deal).Error; err != nil {
		return nil, err
	}
	return deal, nil
}

func (r *GroupDealPostgresRepository) Update(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).Where("id = ?", deal.ID).Updates(deal).Error; err != nil {
		return nil, err
	}
	return deal, nil
}

func (r *GroupDealPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupDealPostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	var model entity.Deal
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GroupDealPostgresRepository) DetailByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	var model entity.Deal
	if err := r.db.WithContext(ctx).
		Preload("Products").
		Preload("Members", "member_type = 30").
		Preload("Customers", "member_type = 10").
		Preload("Partners", "member_type = 20").
		Preload("BankAccount", "deleted_at is null").
		Preload("Documents", "deleted_at is null and doc_owner = 10").
		Where("id = ? AND deleted_at is null", id).
		First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GroupDealPostgresRepository) GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	var deals []*entity.Deal
	var total int64

	offset := (page - 1) * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", groupID, enums.OwnerTypeGroup, false).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", groupID, enums.OwnerTypeGroup, false).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, uint32(total), nil
}

func (r *GroupDealPostgresRepository) GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	var deals []*entity.Deal
	var total int64

	offset := (page - 1) * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", organizationID, enums.OwnerTypeOrganization, false).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", organizationID, enums.OwnerTypeOrganization, false).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, uint32(total), nil
}

func (r *GroupDealPostgresRepository) UpdateStatus(ctx context.Context, id uint64, status enums.DealStatus) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).Where("id = ? AND is_deleted = ?", id, false).Update("status", status).Error
}

func (r *GroupDealPostgresRepository) UpdateStatusAndCancelReason(ctx context.Context, id uint64, status enums.DealStatus, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"status":        status,
			"cancel_reason": cancelReason,
		}).Error
}

func (r *GroupDealPostgresRepository) UpdateSetting(ctx context.Context, id uint64, setting *dto.DealSetting) error {
	tx := r.db.WithContext(ctx).Begin()
	var model *entity.Deal
	err := tx.Model(&entity.Deal{}).Where("id = ?", id).First(&model).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	
	err = tx.Model(&entity.Deal{}).Where("id = ? AND is_deleted = ?", id, false).Updates(model).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	err = tx.Model(&entity.BankAccount{}).Where("id = ?", model.BankAccountID).Updates(map[string]interface{}{
		"bank_name":      setting.BankName,
		"account_number": setting.AccountNumber,
		"account_name":   setting.AccountName,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
