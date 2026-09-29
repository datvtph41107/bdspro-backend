package case_archived

import (
	_db "common/db"
	"context"
)

type IArchivedRepo[T any] interface {
	// _db.IRepo
	Archive(c context.Context, id uint64, archived bool) error
	Approve(c context.Context, id uint64) error
	Reject(c context.Context, id uint64) error
}

type ArchivedRepo[T any] struct {
	*_db.TransactionRepo
	implements IArchivedRepo[T]
}

func (r *ArchivedRepo[T]) Init(repo IArchivedRepo[T], db *_db.TransactionRepo) {
	r.implements = repo
	r.TransactionRepo = db
}

func (r *ArchivedRepo[T]) Archive(c context.Context, id uint64, archived bool) error {
	return r.TransactionRepo.GetDB(c).
		Model(new(T)).
		Where("id = ?", id).
		Update("archived", archived).
		Error
}

func (r *ArchivedRepo[T]) Approve(c context.Context, id uint64) error {
	return r.TransactionRepo.GetDB(c).
		Model(new(T)).
		Where("id = ?", id).
		Update("approved", true).
		Error
}

func (r *ArchivedRepo[T]) Reject(c context.Context, id uint64) error {
	return r.TransactionRepo.GetDB(c).
		Model(new(T)).
		Where("id = ?", id).
		Update("approved", false).
		Error
}
