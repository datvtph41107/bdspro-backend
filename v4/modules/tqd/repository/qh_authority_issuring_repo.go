package repository

import (
	"context"

	qh "userv4/modules/tqd/entity/qh"
)

// QHAuthorityIssuringRepository truy cập bảng qh_authority_issuring.
type QHAuthorityIssuringRepository interface {
	Create(ctx context.Context, row *qh.QHAuthorityIssuring) error
	Update(ctx context.Context, row *qh.QHAuthorityIssuring) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHAuthorityIssuring, error)
	// GetByCode — bản ghi active (chưa xóa mềm) theo code; dùng kiểm tra trùng.
	GetByCode(ctx context.Context, code string) (*qh.QHAuthorityIssuring, error)
	List(ctx context.Context, offset, limit int) ([]qh.QHAuthorityIssuring, int64, error)
}
