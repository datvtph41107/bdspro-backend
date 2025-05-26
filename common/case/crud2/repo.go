package crud2

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	Create(c *gin.Context, entity *T) error
	CreateBatch(c *gin.Context, entity []T) error
	Update(c *gin.Context, id uint64, entity *T) error
	Delete(c *gin.Context, id uint64) error
	GetAll() ([]T, error)
	GetByID(c *gin.Context, id uint64) (*T, error)
}

type BaseRepo[T any] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *BaseRepo[T]) Create(c *gin.Context, entity *T) error {
	return r.DB.WithContext(c).Create(entity).Error
}

// Tạo mới bản ghi
func (r *BaseRepo[T]) CreateBatch(c *gin.Context, entity []T) error {
	return r.DB.WithContext(c).Create(&entity).Error
}

// Cập nhật bản ghi
func (r *BaseRepo[T]) Update(c *gin.Context, id uint64, entity *T) error {
	return r.DB.WithContext(c).Model(entity).Where("id = ? and deleted_at is null", id).Updates(entity).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *BaseRepo[T]) Delete(c *gin.Context, id uint64) error {
	return r.DB.Model(new(T)).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *BaseRepo[T]) DeleteBatch(c *gin.Context, ids []uint64) error {
	return r.DB.Model(new(T)).Where("id in (?)", ids).Update("deleted_at", time.Now()).Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetAll() ([]T, error) {
	var entities []T
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Limit(200).Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetByID(c *gin.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
