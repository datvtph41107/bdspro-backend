package postgre

import (
	base_enum "base/enum"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.RuleRepo
type PostgreRule struct {
	db *gorm.DB
}

func NewPostgreRule(db *gorm.DB) *PostgreRule {
	return &PostgreRule{db: db}
}

func (r *PostgreRule) Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error) {
	var entities []domain.RuleEntity
	query := r.db.WithContext(c).Model(&domain.RuleEntity{})

	// query = query.Where("owner_id = ? AND owner_type = ?", ownerId, ownerType)

	if err := query.
		Debug().
		Order("updated_at DESC").
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

func (r *PostgreRule) Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	if err := r.db.WithContext(c).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreRule) Update(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error) {
	if err := r.db.WithContext(c).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PostgreRule) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).
		Model(&domain.RuleEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *PostgreRule) GetByID(c context.Context, id uint64) (*domain.RuleEntity, error) {
	var entity domain.RuleEntity
	if err := r.db.WithContext(c).
		Model(&domain.RuleEntity{}).
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreRule) Active(c context.Context, id uint64, active bool) error {
	return r.db.WithContext(c).
		Model(&domain.RuleEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("active", active).
		Error
}