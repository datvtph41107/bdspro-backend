package repo

import (
	"context"
	"hub/internal/domain"
)

type IApplinkRepo interface {
	// ExistsByCode kiểm tra code đã tồn tại — dùng trước khi insert
	ExistsByCode(ctx context.Context, code uint64) (bool, error)

	// Create lưu applink mới, populate ID + timestamps từ DB
	Create(ctx context.Context, a *domain.Applink) error

	// GetByCode tìm applink theo code (uint64) — hit index idx_applink_code
	GetByCode(ctx context.Context, code uint64) (*domain.Applink, error)

	// GetByRefIDAndAction tìm applink theo ref_id + action — nếu có thì return, dùng để tránh tạo trùng
	GetByRefIDAndAction(ctx context.Context, refID uint64, action int) (*domain.Applink, error)
}
