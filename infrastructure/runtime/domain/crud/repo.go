package _crud

import (
	_dto "common/domain/dto"
	"context"
)

// là usecase cơ bản cho logic thêm sửa xóa 1 bản ghi thông thường
type ICrudRepo[T any] interface {
	Create(c context.Context, entity *T) error
	CreateBatch(c context.Context, entity []*T) error
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
