package repository

import (
	"context"
	"time"

	_dto "common/domain/dto"
	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// ProjectAdminRepository is admin CRUD on projects.
type ProjectAdminRepository interface {
	List(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.Project, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.Project, error)
	Create(ctx context.Context, project *entity.Project) (*entity.Project, error)
	Update(ctx context.Context, id uint64, project *entity.Project) (*entity.Project, error)
	Delete(ctx context.Context, id uint64) error
}

type projectAdminPostgres struct {
	db *gorm.DB
}

func NewProjectAdminPostgres(db *gorm.DB) ProjectAdminRepository {
	return &projectAdminPostgres{db: db}
}

func (r *projectAdminPostgres) List(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.Project, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Project{}).Where("deleted_at IS NULL")
	if text != "" {
		q = q.Where("LOWER(name) LIKE '%' || LOWER(?) || '%'", text)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.Project
	err := q.Order("id DESC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&items).Error
	return items, total, err
}

func (r *projectAdminPostgres) GetByID(ctx context.Context, id uint64) (*entity.Project, error) {
	if r.db == nil || id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var project entity.Project
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectAdminPostgres) Create(ctx context.Context, project *entity.Project) (*entity.Project, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectAdminPostgres) Update(ctx context.Context, id uint64, project *entity.Project) (*entity.Project, error) {
	if r.db == nil || id == 0 || project == nil {
		return nil, gorm.ErrRecordNotFound
	}
	err := r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"name":         project.Name,
			"developer_id": project.DeveloperID,
			"description":  project.Description,
		}).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *projectAdminPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}
