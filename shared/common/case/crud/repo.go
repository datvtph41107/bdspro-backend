package crud

import (
	_db "common/db"
	_dto "common/domain/dto"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// là usecase cơ bản cho logic thêm sửa xóa 1 bản ghi thông thường
type ICrudRepo[T any] interface {
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	GetAll(c context.Context) ([]T, error)
	GetByID(c context.Context, id uint64) (*T, error)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
	GetListByIDs(c context.Context, ids []uint64) ([]T, error)
	GetDetail(c context.Context, id uint64) (*T, error)
	BeforeSave(c context.Context, id *uint64, entity *T) error
	AfterSave(c context.Context, id *uint64, entity *T) error
}

type CrudRepo[T any] struct {
	*_db.TransactionRepo
	implements ICrudRepo[T]
}

func (r *CrudRepo[T]) Init(repo ICrudRepo[T], db *_db.TransactionRepo) {
	r.implements = repo
	r.TransactionRepo = db
}

// func (r *BaseRepo[T]) GetDB(c context.Context) *gorm.DB {
// 	return r.DBTransaction.GetDB(c, r.DBTransaction.DB)
// }

func (r *CrudRepo[T]) AfterSave(c context.Context, id *uint64, entity *T) error {
	r.GetDB(c)
	return nil
}

// Tạo mới bản ghi
func (r *CrudRepo[T]) Create(c context.Context, entity *T) error {
	if err := r.implements.BeforeSave(c, nil, entity); err != nil {
		return err
	}

	err := r.TransactionRepo.WithTransaction(c, func(c context.Context) error {
		exec := r.GetDB(c).Create(entity)
		if err := r.implements.AfterSave(c, nil, entity); err != nil {
			return err
		}
		return exec.Error
	})

	return err
}

// Cập nhật bản ghi
func (r *CrudRepo[T]) Update(c context.Context, id uint64, entity *T) error {
	if err := r.implements.BeforeSave(c, &id, entity); err != nil {
		return err
	}
	err := r.TransactionRepo.WithTransaction(c, func(c context.Context) error {
		err := r.GetDB(c).
			Model(entity).
			Where("id = ? and deleted_at is null", id).
			Updates(entity).Error
		if err != nil {
			return err
		}
		return r.implements.AfterSave(c, &id, entity)
	})
	return err
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *CrudRepo[T]) Delete(c context.Context, id uint64) error {
	return r.GetDB(c).
		Model(new(T)).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", time.Now()).
		Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *CrudRepo[T]) GetAll(c context.Context) ([]T, error) {
	var entities []T
	err := r.GetDB(c).
		Where("deleted_at is null").
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *CrudRepo[T]) GetByID(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.GetDB(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CrudRepo[T]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
	var entities []T
	var total int64
	err := r.GetDB(c).
		Where("deleted_at is null").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.GetDB(c).
		Model(new(T)).
		Where("deleted_at is null").
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	return entities, total, err
}

func (r *CrudRepo[T]) GetListByIDs(c context.Context, ids []uint64) ([]T, error) {
	var entities []T
	err := r.GetDB(c).
		Where("id IN (?) and deleted_at is null", ids).
		Find(&entities).Error
	return entities, err
}

func (r *CrudRepo[T]) GetDetail(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.GetDB(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *CrudRepo[T]) BeforeSave(c context.Context, id *uint64, entity *T) error {
	return nil
}
