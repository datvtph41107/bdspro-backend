package repository

import (
	"context"

	_dto "common/domain/dto"
	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// ProjectBuildRepository persists project_build rows and related list queries.
type ProjectBuildRepository interface {
	Search(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error)
	SearchByCreatedBy(ctx context.Context, createdBy uint64, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.ProjectBuild, error)
	Create(ctx context.Context, build *entity.ProjectBuild) (*entity.ProjectBuild, error)
	Update(ctx context.Context, build *entity.ProjectBuild) (*entity.ProjectBuild, error)
	Delete(ctx context.Context, id uint64) error
}

type projectBuildPostgres struct {
	db *gorm.DB
}

func NewProjectBuildPostgres(db *gorm.DB) ProjectBuildRepository {
	return &projectBuildPostgres{db: db}
}

func (r *projectBuildPostgres) Search(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.ProjectBuild{}).Where("deleted_at IS NULL")
	if text != "" {
		q = q.Where("name ILIKE ?", "%"+text+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.ProjectBuild
	err := q.Order("name ASC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&items).Error
	return items, total, err
}

func (r *projectBuildPostgres) SearchByCreatedBy(ctx context.Context, createdBy uint64, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error) {
	if r.db == nil || createdBy == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.ProjectBuild{}).
		Where("deleted_at IS NULL AND created_by = ?", createdBy)
	if text != "" {
		q = q.Where("name ILIKE ?", "%"+text+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.ProjectBuild
	err := q.Order("id DESC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&items).Error
	return items, total, err
}

func (r *projectBuildPostgres) GetByID(ctx context.Context, id uint64) (*entity.ProjectBuild, error) {
	if r.db == nil || id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var build entity.ProjectBuild
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Apartments", "floor IS NOT NULL AND floor <> 0").
		Preload("Attributes").
		First(&build).Error
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func (r *projectBuildPostgres) Create(ctx context.Context, build *entity.ProjectBuild) (*entity.ProjectBuild, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := r.db.WithContext(ctx).Create(build).Error; err != nil {
		return nil, err
	}
	return build, nil
}

func (r *projectBuildPostgres) Update(ctx context.Context, build *entity.ProjectBuild) (*entity.ProjectBuild, error) {
	if r.db == nil || build == nil || build.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err := r.db.WithContext(ctx).Save(build).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, build.ID)
}

func (r *projectBuildPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Delete(&entity.ProjectBuild{}, id).Error
}
