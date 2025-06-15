package postgre

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_db "common/db"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProjectRepo
type PostgreProject struct {
	// crud2.BaseRepo[domain.Project]
	DB *gorm.DB
}

func NewPostgreProject(db *gorm.DB) *PostgreProject {
	return &PostgreProject{
		DB: db,
	}
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *PostgreProject) GetByID(c *gin.Context, id uint64) (*domain.Project, error) {
	log.Printf("LOAD_DETAIL", id)
	var entity *domain.Project
	err := _db.DB.
		Preload("Builds").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PostgreProject) GetAll() ([]domain.Project, error) {
	var entities []domain.Project
	err := r.DB.Where("deleted_at IS NULL").
		Preload("Developer").
		Find(&entities).Error
	return entities, err
}

func (r *PostgreProject) SearchItem(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, int64, error) {
	var entities []domain.ProjectItem
	err := r.DB.Where("deleted_at IS NULL").
		// Preload("Developer").
		Where("LOWER(name) LIKE '%' || LOWER(?) || '%'", dto.Text).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).Error
	return entities, 0, err
}

func (r *PostgreProject) GetAllItem(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, error) {
	var items []domain.ProjectItem
	err := r.DB.
		Model(&domain.Project{}).
		Where("deleted_at IS NULL").
		Preload("Developer").
		// Preload("Builds").
		Find(&items).Error

	return items, err
}
