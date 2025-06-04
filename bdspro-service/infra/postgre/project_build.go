package postgre

import (
	"bdspro/internal/domain"
	"common/case/crud2"

	"github.com/gin-gonic/gin"
)

type PostgreProjectBuild struct {
	crud2.BaseRepo[domain.ProjectBuild]
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *PostgreProjectBuild) GetAll() ([]domain.ProjectBuild, error) {
	var entities []domain.ProjectBuild
	err := r.DB.Preload("Project").
		Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *PostgreProjectBuild) GetByID(c *gin.Context, id uint64) (*domain.ProjectBuild, error) {
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
