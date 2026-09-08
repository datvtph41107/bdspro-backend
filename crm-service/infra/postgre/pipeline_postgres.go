package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.PipelineRepo
type PostgrePipeline struct {
	db *gorm.DB
}

func NewPostgrePipeline(db *gorm.DB) *PostgrePipeline {
	return &PostgrePipeline{db: db}
}

func (r *PostgrePipeline) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerOf, dto dto.PipelineSearchDTO, withStages bool) ([]domain.PipelineEntity, int64, error) {
	var entities []domain.PipelineEntity
	query := r.db.WithContext(c).Model(&domain.PipelineEntity{})

	query = query.Where("(owner_id = ? AND owner_type = ?) OR is_custom = ?", ownerId, ownerType, false)

	querySearch := query
	if withStages {
		querySearch = querySearch.Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC").Where("deleted_at is null")
		})
	}

	if err := querySearch.
		Preload("DefaultStage").
		Find(&entities).
		Order("order_number ASC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *PostgrePipeline) Create(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	if err := r.db.WithContext(c).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgrePipeline) Update(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	if err := r.db.WithContext(c).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgrePipeline) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).
		Model(&domain.PipelineEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *PostgrePipeline) GetByID(c context.Context, id uint64) (*domain.PipelineEntity, error) {
	var entity domain.PipelineEntity
	if err := r.db.WithContext(c).
		Model(&domain.PipelineEntity{}).
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("DefaultStage").
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgrePipeline) GetDefault(c context.Context, organizationId uint64) (*domain.PipelineEntity, error) {
	var entity domain.PipelineEntity
	if err := r.db.WithContext(c).
		Model(&domain.PipelineEntity{}).
		Where("(owner_id = ? AND owner_type = ?) OR is_custom = ?",
			organizationId, enums.EOwnerOfOrgnization, false).
		Order("is_default DESC").
		Preload("DefaultStage").
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}