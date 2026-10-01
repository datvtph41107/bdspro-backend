package repository

import (
	"context"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
)

type IQHPlanningDocumentRepo interface {
	Create(ctx context.Context, entity *qh.QHPlanningDocument) error
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningDocument) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh.QHPlanningDocument, int64, error)
	CountByProject(ctx context.Context, projectID uint64) (int64, error)
	GetAll(ctx context.Context) ([]qh.QHPlanningDocument, error)
	GetDetail(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error)

	// ListByProcessStatus lấy các tài liệu đang ở 1 trạng thái xử lý (dùng cho job classify), order theo id tăng dần.
	ListByProcessStatus(ctx context.Context, status uint32, limit int) ([]qh.QHPlanningDocument, error)
	// CountByProjectAndStatuses đếm số tài liệu của 1 đồ án theo danh sách trạng thái (dùng để quyết định trạng thái đồ án).
	CountByProjectAndStatuses(ctx context.Context, projectID uint64, statuses []uint32) (int64, error)
	// ApproveClassifiedByProject chuyển toàn bộ tài liệu đã Classified của 1 đồ án sang Approved (khi admin duyệt cả đồ án).
	ApproveClassifiedByProject(ctx context.Context, projectID uint64) error
	// ResetClassifyForRetry đưa 1 tài liệu Failed về Pending và xóa classify_error.
	ResetClassifyForRetry(ctx context.Context, id uint64) error
	// ResetFailedClassifyByProject đưa tất cả tài liệu Failed của đồ án về Pending; trả về số bản ghi cập nhật.
	ResetFailedClassifyByProject(ctx context.Context, projectID uint64) (int64, error)
	// AssignToLayers tạo/khôi phục bản ghi layer_plan_doc_relation cho từng cặp layerIds × docIds.
	AssignToLayers(ctx context.Context, layerIDs, docIDs []uint64) error
	// UnassignFromLayers xóa cứng bản ghi layer_plan_doc_relation theo layerIds × docIds.
	UnassignFromLayers(ctx context.Context, layerIDs, docIDs []uint64) error
	// UnassignLayer xóa toàn bộ layer_plan_doc_relation của layer với tài liệu thuộc đồ án.
	UnassignLayer(ctx context.Context, planningID, layerID uint64) error
	GetPlanningDocumentSeoSource(ctx context.Context, documentID uint64) (*qh_dto.PlanningDocumentSeoSource, error)
	GetPlanningDocumentSeoSourcesByIDs(ctx context.Context, documentIDs []uint64) ([]qh_dto.PlanningDocumentSeoSource, error)
	UpdatePlanningDocumentSeoID(ctx context.Context, documentID uint64, seoID uint64) error
	CreateHistory(ctx context.Context, entity *qh.QHPlanningDocumentEvent) error
	DeleteHistory(ctx context.Context, id uint64) error
	GetHistoryByID(ctx context.Context, id uint64) (*qh.QHPlanningDocumentEvent, error)
	ListHistories(ctx context.Context, req *qh_dto.ListPlanningDocumentHistoriesRequest) ([]qh.QHPlanningDocumentEvent, int64, error)
}
