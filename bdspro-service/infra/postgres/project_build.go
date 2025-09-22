package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProjectBuildRepo
type PostgreProjectBuild struct {
	DB *gorm.DB
}

func NewPostgreProjectBuild(db *gorm.DB) *PostgreProjectBuild {
	return &PostgreProjectBuild{
		DB: db,
	}
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *PostgreProjectBuild) GetAll(c context.Context) ([]domain.ProjectBuild, error) {
	var entities []domain.ProjectBuild
	err := GetDB(c, r.DB).Preload("Project").
		Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *PostgreProjectBuild) GetByID(c context.Context, id uint64) (*domain.ProjectBuild, error) {
	// log.Printf("LOAD_DETAIL", id)
	var entity *domain.ProjectBuild
	err := r.DB.
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Apartments", "floor is not null and floor <> 0").
		Preload("Attributes").
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PostgreProjectBuild) Search(c context.Context, dto *dto.TextSearchRequest) ([]domain.ProjectBuild, int64, error) {
	var entities []domain.ProjectBuild
	query := GetDB(c, r.DB).
		Where("deleted_at IS NULL AND name ILIKE ?", "%"+dto.Text+"%")

	if dto.Text != "" {
		query = query.Where("name ILIKE ?", "%"+dto.Text+"%")
	}

	err := query.Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	total := int64(0)
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, err
}

func (r *PostgreProjectBuild) Create(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error) {
	err := GetDB(c, r.DB).Create(dto).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *PostgreProjectBuild) Update(c context.Context, dto *domain.ProjectBuild) (*domain.ProjectBuild, error) {
	err := GetDB(c, r.DB).Save(dto).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *PostgreProjectBuild) Delete(c context.Context, id uint64) error {
	return GetDB(c, r.DB).Delete(&domain.ProjectBuild{}, id).Error
}
