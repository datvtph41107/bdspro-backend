package a_project_build

import (
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectBuildRepo struct {
	crud2.BaseRepo[models.ProjectBuild]
}

func NewProjectBuildRepo(db *gorm.DB) *ProjectBuildRepo {
	return &ProjectBuildRepo{
		BaseRepo: crud2.BaseRepo[models.ProjectBuild]{DB: db},
	}
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *ProjectBuildRepo) GetAll() ([]models.ProjectBuild, error) {
	var entities []models.ProjectBuild
	err := r.DB.Preload("Project").
		Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *ProjectBuildRepo) GetByID(c *gin.Context, id uint64) (*models.ProjectBuild, error) {
	// log.Printf("LOAD_DETAIL", id)
	var entity *models.ProjectBuild
	err := _db.DB.
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Apartments", "floor is not null and floor <> 0").
		Preload("Attributes").
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
