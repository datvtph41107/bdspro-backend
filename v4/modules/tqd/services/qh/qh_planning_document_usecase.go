package qh_usecase

import (
	_err "common/domain/err"
	"context"
	"errors"
	"fmt"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/provider"
	repo "userv4/modules/tqd/repository"

	"gorm.io/gorm"
)

type IQHPlanningDocumentUsecase interface {
	Create(ctx context.Context, entity *qh.QHPlanningDocument) (*qh.QHPlanningDocument, error)
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningDocument) (*qh.QHPlanningDocument, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh.QHPlanningDocument, int64, error)
	// Approve phê duyệt 1 tài liệu đã được AI phân loại (chỉ khi đang ở trạng thái Classified).
	Approve(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error)
	// Retry đưa tài liệu Failed về Pending để job classify chạy lại.
	Retry(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error)
	AssignToLayer(ctx context.Context, req *qh_dto.AssignPlanningDocumentsToLayerRequest) error
	UnassignFromLayer(ctx context.Context, req *qh_dto.AssignPlanningDocumentsToLayerRequest) error
	UnassignLayer(ctx context.Context, req *qh_dto.UnassignPlanningProjectLayerRequest) error
	CreateHistory(ctx context.Context, entity *qh.QHPlanningDocumentEvent) (*qh.QHPlanningDocumentEvent, error)
	DeleteHistory(ctx context.Context, documentID uint64, historyID uint64) (bool, error)
	ListHistories(ctx context.Context, req *qh_dto.ListPlanningDocumentHistoriesRequest) ([]qh.QHPlanningDocumentEvent, int64, error)
	GetPlanningDocumentSeoSource(ctx context.Context, documentID uint64) (*qh_dto.PlanningDocumentSeoSource, error)
	UpdatePlanningDocumentSeoID(ctx context.Context, documentID uint64, seoID uint64) error
	AdminGenPlanningDocumentSeo(ctx context.Context, documentIDs []uint64) (*qh_dto.AdminGenPlanningDocumentSeoResult, error)
}

type qhPlanningDocumentUsecase struct {
	repo           repo.IQHPlanningDocumentRepo
	projectRepo    repo.IQHPlanningProjectRepo
	aiJobRepo      repo.IProAIJobRepo
	crmSeoProvider provider.CrmSeoProvider
}

// NewQHPlanningDocumentUsecase @bind: internal/usecase/qh.IQHPlanningDocumentUsecase
func NewQHPlanningDocumentUsecase(
	repo repo.IQHPlanningDocumentRepo,
	projectRepo repo.IQHPlanningProjectRepo,
	aiJobRepo repo.IProAIJobRepo,
	crmSeoProvider provider.CrmSeoProvider,
) IQHPlanningDocumentUsecase {
	return &qhPlanningDocumentUsecase{
		repo:           repo,
		projectRepo:    projectRepo,
		aiJobRepo:      aiJobRepo,
		crmSeoProvider: crmSeoProvider,
	}
}

func (u *qhPlanningDocumentUsecase) Create(ctx context.Context, entity *qh.QHPlanningDocument) (*qh.QHPlanningDocument, error) {
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningDocumentUsecase) Update(ctx context.Context, id uint64, entity *qh.QHPlanningDocument) (*qh.QHPlanningDocument, error) {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu: " + err.Error(),
		}
	}
	return u.GetByID(ctx, id)
}

func (u *qhPlanningDocumentUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningDocumentUsecase) GetByID(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningDocumentUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh.QHPlanningDocument, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningDocumentUsecase) Approve(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error) {
	document, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy tài liệu",
		}
	}
	if document.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Tài liệu chưa được AI phân loại xong (chưa ở trạng thái Classified), không thể phê duyệt",
		}
	}

	document.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusApproved)
	if err := u.repo.Update(ctx, id, document); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt tài liệu",
		}
	}
	return document, nil
}

func (u *qhPlanningDocumentUsecase) Retry(ctx context.Context, id uint64) (*qh.QHPlanningDocument, error) {
	document, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy tài liệu",
		}
	}
	if document.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Chỉ có thể retry tài liệu đang ở trạng thái Lỗi (Failed)",
		}
	}

	if err := u.repo.ResetClassifyForRetry(ctx, id); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể đưa tài liệu về chờ xử lý: " + err.Error(),
		}
	}

	// Đưa đồ án về Pending nếu đang Failed để worker/admin theo dõi lại.
	if project, err := u.projectRepo.GetByID(ctx, document.PlanningProjectID); err == nil && project != nil {
		if project.ProcessStatus == enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed) {
			_ = u.projectRepo.UpdateProcessStatus(ctx, project.ID, enums.PlanningProcessStatusPending)
			_ = u.aiJobRepo.UpdateProcessStatusByProjectID(ctx, project.ID, enums.PlanningProcessStatusPending)
		}
	}

	return u.repo.GetByID(ctx, id)
}

func (u *qhPlanningDocumentUsecase) AssignToLayer(ctx context.Context, req *qh_dto.AssignPlanningDocumentsToLayerRequest) error {
	layerIDs, docIDs, err := normalizeLayerDocIDs(req)
	if err != nil {
		return err
	}
	if err := u.repo.AssignToLayers(ctx, layerIDs, docIDs); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể gán tài liệu vào layer: " + err.Error(),
		}
	}
	return nil
}

func (u *qhPlanningDocumentUsecase) UnassignFromLayer(ctx context.Context, req *qh_dto.AssignPlanningDocumentsToLayerRequest) error {
	layerIDs, docIDs, err := normalizeLayerDocIDs(req)
	if err != nil {
		return err
	}
	if err := u.repo.UnassignFromLayers(ctx, layerIDs, docIDs); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể gỡ tài liệu khỏi layer: " + err.Error(),
		}
	}
	return nil
}

func (u *qhPlanningDocumentUsecase) UnassignLayer(ctx context.Context, req *qh_dto.UnassignPlanningProjectLayerRequest) error {
	if req == nil || req.PlanningID == 0 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "planningId is required",
		}
	}
	if req.LayerID == 0 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "layerId is required",
		}
	}
	if err := u.repo.UnassignLayer(ctx, req.PlanningID, req.LayerID); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể gỡ layer: " + err.Error(),
		}
	}
	return nil
}

func (u *qhPlanningDocumentUsecase) CreateHistory(ctx context.Context, entity *qh.QHPlanningDocumentEvent) (*qh.QHPlanningDocumentEvent, error) {
	if entity == nil || entity.PlanningDocumentID == 0 {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "planningDocumentId is required",
		}
	}
	if _, err := u.repo.GetByID(ctx, entity.PlanningDocumentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy tài liệu",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if err := u.repo.CreateHistory(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	created, err := u.repo.GetHistoryByID(ctx, entity.ID)
	if err != nil {
		return entity, nil
	}
	return created, nil
}

func (u *qhPlanningDocumentUsecase) DeleteHistory(ctx context.Context, documentID uint64, historyID uint64) (bool, error) {
	if historyID == 0 {
		return false, &_err.ErrorDTO{
			Code:    400,
			Message: "id is required",
		}
	}
	history, err := u.repo.GetHistoryByID(ctx, historyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy lịch sử tài liệu",
			}
		}
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if documentID > 0 && history.PlanningDocumentID != documentID {
		return false, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy lịch sử tài liệu",
		}
	}
	if err := u.repo.DeleteHistory(ctx, historyID); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningDocumentUsecase) ListHistories(ctx context.Context, req *qh_dto.ListPlanningDocumentHistoriesRequest) ([]qh.QHPlanningDocumentEvent, int64, error) {
	if req == nil || req.PlanningDocumentID == 0 {
		return nil, 0, &_err.ErrorDTO{
			Code:    400,
			Message: "planningDocumentId is required",
		}
	}
	if _, err := u.repo.GetByID(ctx, req.PlanningDocumentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy tài liệu",
			}
		}
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	entities, total, err := u.repo.ListHistories(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningDocumentUsecase) GetPlanningDocumentSeoSource(ctx context.Context, documentID uint64) (*qh_dto.PlanningDocumentSeoSource, error) {
	if documentID == 0 {
		return nil, fmt.Errorf("planning_document_id is required")
	}
	return u.repo.GetPlanningDocumentSeoSource(ctx, documentID)
}

func (u *qhPlanningDocumentUsecase) UpdatePlanningDocumentSeoID(ctx context.Context, documentID uint64, seoID uint64) error {
	if documentID == 0 || seoID == 0 {
		return fmt.Errorf("planning_document_id và seo_id là bắt buộc")
	}
	return u.repo.UpdatePlanningDocumentSeoID(ctx, documentID, seoID)
}

func (u *qhPlanningDocumentUsecase) AdminGenPlanningDocumentSeo(ctx context.Context, documentIDs []uint64) (*qh_dto.AdminGenPlanningDocumentSeoResult, error) {
	if len(documentIDs) == 0 {
		return nil, fmt.Errorf("docIds không được để trống")
	}
	if u.crmSeoProvider == nil {
		return nil, fmt.Errorf("crm seo provider chưa được cấu hình")
	}

	uniqueIDs := uniquePositiveIDs(documentIDs)
	if len(uniqueIDs) == 0 {
		return nil, fmt.Errorf("docIds không hợp lệ")
	}

	sources, err := u.repo.GetPlanningDocumentSeoSourcesByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, fmt.Errorf("get planning document seo sources failed: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("không tìm thấy tài liệu hợp lệ")
	}

	foundIDs := make([]uint64, 0, len(sources))
	prevSeoIDByDoc := make(map[uint64]uint64, len(sources))
	for _, source := range sources {
		foundIDs = append(foundIDs, source.PlanningDocumentID)
		if source.SeoID != nil && *source.SeoID > 0 {
			prevSeoIDByDoc[source.PlanningDocumentID] = *source.SeoID
		}
	}

	result, err := u.crmSeoProvider.GenerateSeoFromPlanningDocumentIDs(ctx, foundIDs)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &qh_dto.AdminGenPlanningDocumentSeoResult{}
	}

	for i := range result.Items {
		item := &result.Items[i]
		if item.PlanningDocumentID == 0 || item.SeoID == 0 {
			continue
		}

		prevSeoID := prevSeoIDByDoc[item.PlanningDocumentID]
		switch {
		case prevSeoID == 0:
			item.Created = true
			item.Message = "created"
		case prevSeoID == item.SeoID:
			item.Created = false
			item.Message = "skipped"
		default:
			item.Created = true
			item.Message = "linked"
		}

		if err := u.repo.UpdatePlanningDocumentSeoID(ctx, item.PlanningDocumentID, item.SeoID); err != nil {
			item.Message = fmt.Sprintf("cập nhật seo_id thất bại: %v", err)
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("docId=%d: %v", item.PlanningDocumentID, err))
		}
	}

	if len(foundIDs) < len(uniqueIDs) {
		foundSet := make(map[uint64]struct{}, len(foundIDs))
		for _, id := range foundIDs {
			foundSet[id] = struct{}{}
		}
		for _, id := range uniqueIDs {
			if _, ok := foundSet[id]; ok {
				continue
			}
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("docId=%d: tài liệu không tồn tại", id))
		}
		if result.Total < int64(len(uniqueIDs)) {
			result.Total = int64(len(uniqueIDs))
		}
	}

	return result, nil
}

func uniquePositiveIDs(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func mergeLayerIDs(layerID uint64, layerIDs []uint64) []uint64 {
	ids := make([]uint64, 0, len(layerIDs)+1)
	if layerID > 0 {
		ids = append(ids, layerID)
	}
	ids = append(ids, layerIDs...)
	return uniquePositiveIDs(ids)
}

func normalizeLayerDocIDs(req *qh_dto.AssignPlanningDocumentsToLayerRequest) ([]uint64, []uint64, error) {
	if req == nil {
		return nil, nil, &_err.ErrorDTO{Code: 400, Message: "request is required"}
	}
	layerIDs := mergeLayerIDs(req.LayerID, req.LayerIDs)
	if len(layerIDs) == 0 {
		return nil, nil, &_err.ErrorDTO{Code: 400, Message: "layerIds is required"}
	}
	docIDs := uniquePositiveIDs(req.DocIDs)
	if len(docIDs) == 0 {
		return nil, nil, &_err.ErrorDTO{Code: 400, Message: "docIds is required"}
	}
	return layerIDs, docIDs, nil
}
