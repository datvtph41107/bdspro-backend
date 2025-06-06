package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupDealModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_deals_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"index:idx_group_deals_created_by"`
	UpdatedBy uint32

	GroupId uint32            `gorm:"not null;index:idx_group_deals_group_id"`
	Name    string            `gorm:"not null;index:idx_group_deals_name"`
	Amount  float64           `gorm:"not null"`
	Status  entity.DealStatus `gorm:"not null;index:idx_group_deals_status"`
}

func (GroupDealModel) TableName() string {
	return "group_deals"
}

func GroupDealModelToEntity(model *GroupDealModel) *entity.Deal {
	return &entity.Deal{
		Id:        model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		GroupId:   model.GroupId,
		Name:      model.Name,
		Amount:    model.Amount,
		Status:    model.Status,
	}
}

func GroupDealEntityToModel(entity *entity.Deal) *GroupDealModel {
	return &GroupDealModel{
		Id:        entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		GroupId:   entity.GroupId,
		Name:      entity.Name,
		Amount:    entity.Amount,
		Status:    entity.Status,
	}
}

type groupDealPostgresRepository struct {
	db *gorm.DB
}

func NewGroupDealPostgresRepository(db *gorm.DB) repository.GroupDealRepository {
	return &groupDealPostgresRepository{db: db}
}

func (r *groupDealPostgresRepository) Create(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	model := GroupDealEntityToModel(deal)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupDealModelToEntity(model), nil
}

func (r *groupDealPostgresRepository) Update(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	model := GroupDealEntityToModel(deal)
	if err := r.db.WithContext(ctx).Model(&GroupDealModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupDealModelToEntity(model), nil
}

func (r *groupDealPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&GroupDealModel{}, id).Error
}

func (r *groupDealPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Deal, error) {
	var model GroupDealModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupDealModelToEntity(&model), nil
}

func (r *groupDealPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.Deal, error) {	
	var models []GroupDealModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Deal, len(models))
	for i, model := range models {
		entities[i] = GroupDealModelToEntity(&model)
	}	
	return entities, nil
}

func (r *groupDealPostgresRepository) UpdateStatus(ctx context.Context, id uint32, status entity.DealStatus) error {
	return r.db.WithContext(ctx).Model(&GroupDealModel{}).Where("id = ?", id).Update("status", status).Error
}

