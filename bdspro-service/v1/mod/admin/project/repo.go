package a_project

import (
	"bdspro/v1/data"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	crud2.BaseRepo[models.Project]
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{
		BaseRepo: crud2.BaseRepo[models.Project]{DB: db},
	}
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *ProjectRepo) GetByID(c *gin.Context, id uint64) (*models.Project, error) {
	log.Printf("LOAD_DETAIL", id)
	var entity *models.Project
	err := _db.DB.
		Preload("Builds").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ProjectRepo) GetAll() ([]models.Project, error) {
	var entities []models.Project
	err := r.DB.Where("deleted_at IS NULL").
		Preload("Developer").
		Find(&entities).Error
	return entities, err
}

func (r *ProjectRepo) Search(dto *data.SearchName) ([]models.Project, error) {
	var entities []models.Project
	err := r.DB.Where("deleted_at IS NULL").
		Preload("Developer").
		Where("LOWER(name) LIKE '%' || LOWER(?) || '%'", dto.Name).
		Limit(dto.PageSize).
		Find(&entities).Error
	return entities, err
}
