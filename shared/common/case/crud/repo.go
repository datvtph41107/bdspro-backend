package crud

import (
	_db "common/db"
	_dto "common/dto"
	"context"
	"time"

	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	GetAll(c context.Context) ([]T, error)
	GetByID(c context.Context, id uint64) (*T, error)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
	GetDetail(c context.Context, id uint64) (*T, error)
	BeforeSave(c context.Context, id *uint64, entity *T) error
	AfterSave(c context.Context, id *uint64, entity *T) error

	// __internal__
	Self() IBaseRepo[T]
	GetDB(c context.Context) *gorm.DB
	Init(repo IBaseRepo[T], db *_db.DBTransaction)
}

type BaseRepo[T any] struct {
	*_db.DBTransaction
	implements IBaseRepo[T]
	// DB         *gorm.DB
}

func MakeCrudRepo[T any](repo IBaseRepo[T], db *_db.DBTransaction) IBaseRepo[T] {
	repo.Init(repo, db)
	return repo
}

func (r *BaseRepo[T]) Init(repo IBaseRepo[T], db *_db.DBTransaction) {
	r.implements = repo
	r.DBTransaction = db
}

func (r *BaseRepo[T]) Self() IBaseRepo[T] {
	return r.implements
}

// func (r *BaseRepo[T]) GetDB(c context.Context) *gorm.DB {
// 	return r.DBTransaction.GetDB(c, r.DBTransaction.DB)
// }

func (r *BaseRepo[T]) AfterSave(c context.Context, id *uint64, entity *T) error {
	return nil
}

// Tạo mới bản ghi
func (r *BaseRepo[T]) Create(c context.Context, entity *T) error {
	if err := r.Self().BeforeSave(c, nil, entity); err != nil {
		return err
	}

	err := r.DBTransaction.WithTransaction(c, func(c context.Context) error {
		exec := r.Self().GetDB(c).Create(entity)
		if err := r.Self().AfterSave(c, nil, entity); err != nil {
			return err
		}
		return exec.Error
	})

	return err
}

// Cập nhật bản ghi
func (r *BaseRepo[T]) Update(c context.Context, id uint64, entity *T) error {
	if err := r.Self().BeforeSave(c, &id, entity); err != nil {
		return err
	}
	err := r.DBTransaction.WithTransaction(c, func(c context.Context) error {
		err := r.Self().GetDB(c).
			Model(entity).
			Where("id = ? and deleted_at is null", id).
			Updates(entity).Error
		if err != nil {
			return err
		}
		return r.Self().AfterSave(c, &id, entity)
	})
	return err
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *BaseRepo[T]) Delete(c context.Context, id uint64) error {
	return r.Self().GetDB(c).
		Model(new(T)).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", time.Now()).
		Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetAll(c context.Context) ([]T, error) {
	var entities []T
	err := r.Self().GetDB(c).
		Where("deleted_at is null").
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetByID(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.Self().GetDB(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
	var entities []T
	var total int64
	err := r.Self().GetDB(c).
		Where("deleted_at is null").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.Self().GetDB(c).
		Model(new(T)).
		Where("deleted_at is null").
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	return entities, total, err
}

func (r *BaseRepo[T]) GetDetail(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.Self().GetDB(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) BeforeSave(c context.Context, id *uint64, entity *T) error {
	return nil
}
