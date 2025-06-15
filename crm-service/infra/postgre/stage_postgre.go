package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.StageRepo
type PostgreStage struct {
	db *gorm.DB
}

func NewPostgreStage(db *gorm.DB) *PostgreStage {
	return &PostgreStage{db: db}
}

func (r *PostgreStage) Search(c context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error) {
	var entities []domain.StageEntity
	query := r.db.WithContext(c).Model(&domain.StageEntity{})

	query = query.Where("pipeline_id = ?", pipelineId)

	if dto.StageName != "" {
		query = query.Where("stage_name LIKE ?", "%"+dto.StageName+"%")
	}

	if err := query.
		Order("order_number asc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).
		Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *PostgreStage) Create(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	if err := r.db.WithContext(c).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreStage) Update(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	if err := r.db.WithContext(c).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreStage) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *PostgreStage) GetByID(c context.Context, id uint64) (*domain.StageEntity, error) {
	var entity domain.StageEntity
	if err := r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreStage) GetOneStageDefault(c context.Context, pipelineId uint64) (*domain.StageEntity, error) {
	var entity domain.StageEntity
	if err := r.db.WithContext(c).
		Model(&domain.StageEntity{}).
		Where("pipeline_id = ? AND is_default = ", pipelineId, true).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
