package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IQHPlanningDocumentRepo interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) error
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error)
	GetAll(ctx context.Context) ([]qh_domain.QHPlanningDocument, error)
	GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)

	// ListByProcessStatus lấy các tài liệu đang ở 1 trạng thái xử lý (dùng cho job classify), order theo id tăng dần.
	ListByProcessStatus(ctx context.Context, status uint32, limit int) ([]qh_domain.QHPlanningDocument, error)
	// CountByProjectAndStatuses đếm số tài liệu của 1 đồ án theo danh sách trạng thái (dùng để quyết định trạng thái đồ án).
	CountByProjectAndStatuses(ctx context.Context, projectID uint64, statuses []uint32) (int64, error)
	// ApproveClassifiedByProject chuyển toàn bộ tài liệu đã Classified của 1 đồ án sang Approved (khi admin duyệt cả đồ án).
	ApproveClassifiedByProject(ctx context.Context, projectID uint64) error
	// ResetClassifyForRetry đưa 1 tài liệu Failed về Pending và xóa classify_error.
	ResetClassifyForRetry(ctx context.Context, id uint64) error
	// ResetFailedClassifyByProject đưa tất cả tài liệu Failed của đồ án về Pending; trả về số bản ghi cập nhật.
	ResetFailedClassifyByProject(ctx context.Context, projectID uint64) (int64, error)
}
