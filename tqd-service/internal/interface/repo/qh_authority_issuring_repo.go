package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
)

// QHAuthorityIssuringRepository truy cập bảng qh_authority_issuring.
type QHAuthorityIssuringRepository interface {
	Create(ctx context.Context, row *qh_domain.QHAuthorityIssuring) error
	Update(ctx context.Context, row *qh_domain.QHAuthorityIssuring) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHAuthorityIssuring, error)
	// GetByCode — bản ghi active (chưa xóa mềm) theo code; dùng kiểm tra trùng.
	GetByCode(ctx context.Context, code string) (*qh_domain.QHAuthorityIssuring, error)
	List(ctx context.Context, offset, limit int) ([]qh_domain.QHAuthorityIssuring, int64, error)
}
