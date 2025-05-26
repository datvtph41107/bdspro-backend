package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IRestRepo[T IRestModel] interface {
	Create(c *gin.Context, entity *RestModel[T]) error
	CreateBatch(c *gin.Context, entity *[]RestModel[T]) error
	Update(c *gin.Context, id uint64, entity *RestModel[T]) error
	Delete(c *gin.Context, id uint64) error
	DeleteBatch(c *gin.Context, ids []uint64) error
	Search(c *gin.Context, profileId uint64, offset, limit int) ([]RestModel[T], error)
	GetByID(c *gin.Context, id uint64) (*RestModel[T], error)
}

type RestRepo[M IRestModel] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *RestRepo[T]) Create(c *gin.Context, entity *RestModel[T]) error {
	return r.DB.WithContext(c).
		// Model(entity).
		Create(entity).
		Error
}

func (r *RestRepo[T]) CreateBatch(c *gin.Context, entity *[]RestModel[T]) error {
	return r.DB.WithContext(c).
		// Model(entity).
		Create(entity).
		Error
}

// Cập nhật bản ghi
func (r *RestRepo[T]) Update(c *gin.Context, id uint64, entity *RestModel[T]) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *RestRepo[T]) Delete(c *gin.Context, id uint64) error {
	return r.DB.Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// Xóa nhiều bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *RestRepo[T]) DeleteBatch(c *gin.Context, ids []uint64) error {
	return r.DB.Model(new(T)).
		Where("id in (?)", ids).
		Update("deleted_at", time.Now()).Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *RestRepo[T]) Search(c *gin.Context, profileId uint64, offset, limit int) ([]RestModel[T], error) {
	var entities []RestModel[T]
	err := r.DB.Where("deleted_at IS NULL and Rest_id = ?", profileId).
		Offset(offset).
		Limit(limit).
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *RestRepo[T]) GetByID(c *gin.Context, id uint64) (*RestModel[T], error) {
	var entity *RestModel[T]
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
