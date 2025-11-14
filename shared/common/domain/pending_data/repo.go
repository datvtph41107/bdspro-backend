package pending_data

import (
	"context"

	_crud "common/domain/crud"
)

// IPendingDataRepo định nghĩa các thao tác với bảng pending_data.
type IPendingDataRepo interface {
	_crud.ICrudRepo[PendingData]

	// GetByServiceName trả về toàn bộ pending data của một service cụ thể.
	GetByServiceName(ctx context.Context, serviceName string, limit int) ([]PendingData, error)

	// DeleteByIDs xóa (hoặc đánh dấu đã xử lý) các bản ghi theo danh sách id.
	DeleteByIDs(ctx context.Context, ids []uint64) error
}
