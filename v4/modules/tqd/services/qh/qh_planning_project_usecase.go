package qh_usecase

import (
	_err "common/domain/err"
	"context"
	"errors"
	"fmt"
	"strings"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/provider"
	repo "userv4/modules/tqd/repository"

	"gorm.io/gorm"
)

type IQHPlanningProjectUsecase interface {
	Create(ctx context.Context, entity *qh.QHPlanningProject) (*qh.QHPlanningProject, error)
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningProject) (*qh.QHPlanningProject, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningProject, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh.QHPlanningProject, int64, error)
	ListLayers(ctx context.Context, req *qh_dto.ListPlanningProjectLayersRequest) ([]qh.QHLayer, int64, error)
	// Approve phê duyệt đồ án (chỉ khi đã Classified) và phê duyệt luôn các tài liệu đã Classified thuộc đồ án.
	Approve(ctx context.Context, id uint64) (*qh.QHPlanningProject, error)
	// RetryFailedDocuments đưa tất cả tài liệu Failed của đồ án về Pending.
	RetryFailedDocuments(ctx context.Context, id uint64) (*qh.QHPlanningProject, error)

	GetPlanningProjectSeoSource(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectSeoSource, error)
	UpdatePlanningProjectSeoID(ctx context.Context, planningProjectID uint64, seoID uint64) error
	AdminGenPlanningProjectSeo(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error)
	Count(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectCount, error)
	// Assign/unassign provinces and wards
	AssignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error
	UnassignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error
	AssignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error
	UnassignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error
}

type qhPlanningProjectUsecase struct {
	repo           repo.IQHPlanningProjectRepo
	documentRepo   repo.IQHPlanningDocumentRepo
	eventRepo      repo.IQHPlanningEventRepo
	aiJobRepo      repo.IProAIJobRepo
	crmSeoProvider provider.CrmSeoProvider
}

// NewQHPlanningProjectUsecase @bind: internal/usecase/qh.IQHPlanningProjectUsecase
func NewQHPlanningProjectUsecase(
	repo repo.IQHPlanningProjectRepo,
	documentRepo repo.IQHPlanningDocumentRepo,
	eventRepo repo.IQHPlanningEventRepo,
	aiJobRepo repo.IProAIJobRepo,
	crmSeoProvider provider.CrmSeoProvider,
) IQHPlanningProjectUsecase {
	return &qhPlanningProjectUsecase{
		repo:           repo,
		documentRepo:   documentRepo,
		eventRepo:      eventRepo,
		aiJobRepo:      aiJobRepo,
		crmSeoProvider: crmSeoProvider,
	}
}

func (u *qhPlanningProjectUsecase) Create(ctx context.Context, entity *qh.QHPlanningProject) (*qh.QHPlanningProject, error) {
	// Tạo qua API QH: bắt đầu ở Pending (chờ xử lý / phân loại).
	entity.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusPending)
	if strings.TrimSpace(entity.Metadata) == "" {
		entity.Metadata = "{}"
	}
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Update(ctx context.Context, id uint64, entity *qh.QHPlanningProject) (*qh.QHPlanningProject, error) {
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningProjectUsecase) GetByID(ctx context.Context, id uint64) (*qh.QHPlanningProject, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy dữ liệu",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh.QHPlanningProject, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningProjectUsecase) ListLayers(ctx context.Context, req *qh_dto.ListPlanningProjectLayersRequest) ([]qh.QHLayer, int64, error) {
	if req == nil || req.PlanningID == 0 {
		return nil, 0, &_err.ErrorDTO{
			Code:    400,
			Message: "planningId is required",
		}
	}
	layers, total, err := u.repo.ListLayers(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return layers, total, nil
}

func (u *qhPlanningProjectUsecase) Approve(ctx context.Context, id uint64) (*qh.QHPlanningProject, error) {
	project, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy đồ án",
		}
	}
	if project.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Đồ án chưa phân loại xong (chưa Classified), không thể phê duyệt",
		}
	}

	if err := u.documentRepo.ApproveClassifiedByProject(ctx, id); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt các tài liệu của đồ án: " + err.Error(),
		}
	}

	if err := u.repo.UpdateProcessStatus(ctx, id, enums.PlanningProcessStatusApproved); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt đồ án",
		}
	}
	_ = u.aiJobRepo.UpdateProcessStatusByProjectID(ctx, id, enums.PlanningProcessStatusApproved)

	project.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusApproved)
	return project, nil
}

func (u *qhPlanningProjectUsecase) RetryFailedDocuments(ctx context.Context, id uint64) (*qh.QHPlanningProject, error) {
	project, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy đồ án",
		}
	}

	affected, err := u.documentRepo.ResetFailedClassifyByProject(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể retry các tài liệu lỗi: " + err.Error(),
		}
	}
	if affected == 0 {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Không có tài liệu nào ở trạng thái Lỗi để retry",
		}
	}

	if err := u.repo.UpdateProcessStatus(ctx, id, enums.PlanningProcessStatusPending); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Đã reset tài liệu nhưng không cập nhật được trạng thái đồ án",
		}
	}
	_ = u.aiJobRepo.UpdateProcessStatusByProjectID(ctx, id, enums.PlanningProcessStatusPending)

	project.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusPending)
	return project, nil
}

func (u *qhPlanningProjectUsecase) GetPlanningProjectSeoSource(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectSeoSource, error) {
	if planningProjectID == 0 {
		return nil, fmt.Errorf("planning_project_id is required")
	}
	return u.repo.GetPlanningProjectSeoSource(ctx, planningProjectID)
}

func (u *qhPlanningProjectUsecase) UpdatePlanningProjectSeoID(ctx context.Context, planningProjectID uint64, seoID uint64) error {
	if planningProjectID == 0 || seoID == 0 {
		return fmt.Errorf("planning_project_id và seo_id là bắt buộc")
	}
	return u.repo.UpdatePlanningProjectSeoID(ctx, planningProjectID, seoID)
}

func (u *qhPlanningProjectUsecase) AssignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error {
	if planningProjectID == 0 {
		return fmt.Errorf("planning_project_id is required")
	}
	return u.repo.AssignProvinces(ctx, planningProjectID, provinceIDs)
}

func (u *qhPlanningProjectUsecase) UnassignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error {
	if planningProjectID == 0 {
		return fmt.Errorf("planning_project_id is required")
	}
	return u.repo.UnassignProvinces(ctx, planningProjectID, provinceIDs)
}

func (u *qhPlanningProjectUsecase) AssignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error {
	if planningProjectID == 0 {
		return fmt.Errorf("planning_project_id is required")
	}
	return u.repo.AssignWards(ctx, planningProjectID, wardIDs)
}

func (u *qhPlanningProjectUsecase) UnassignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error {
	if planningProjectID == 0 {
		return fmt.Errorf("planning_project_id is required")
	}
	return u.repo.UnassignWards(ctx, planningProjectID, wardIDs)
}

func (u *qhPlanningProjectUsecase) AdminGenPlanningProjectSeo(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error) {
	if len(planningProjectIDs) == 0 {
		return nil, fmt.Errorf("planningProjectIds không được để trống")
	}
	if u.crmSeoProvider == nil {
		return nil, fmt.Errorf("crm seo provider chưa được cấu hình")
	}

	uniqueIDs := make([]uint64, 0, len(planningProjectIDs))
	seen := make(map[uint64]struct{}, len(planningProjectIDs))
	for _, id := range planningProjectIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return nil, fmt.Errorf("planningProjectIds không hợp lệ")
	}

	sources, err := u.repo.GetPlanningProjectSeoSourcesByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, fmt.Errorf("get planning project seo sources failed: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("không tìm thấy đồ án hợp lệ")
	}

	foundIDs := make([]uint64, 0, len(sources))
	prevSeoIDByProject := make(map[uint64]uint64, len(sources))
	for _, source := range sources {
		foundIDs = append(foundIDs, source.PlanningProjectID)
		if source.SeoID != nil && *source.SeoID > 0 {
			prevSeoIDByProject[source.PlanningProjectID] = *source.SeoID
		}
	}

	result, err := u.crmSeoProvider.GenerateSeoFromPlanningProjectIDs(ctx, foundIDs)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &qh_dto.AdminGenPlanningProjectSeoResult{}
	}

	for i := range result.Items {
		item := &result.Items[i]
		if item.PlanningProjectID == 0 || item.SeoID == 0 {
			continue
		}

		prevSeoID := prevSeoIDByProject[item.PlanningProjectID]
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

		if err := u.repo.UpdatePlanningProjectSeoID(ctx, item.PlanningProjectID, item.SeoID); err != nil {
			item.Message = fmt.Sprintf("cập nhật seo_id thất bại: %v", err)
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("planningProjectId=%d: %v", item.PlanningProjectID, err))
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
			result.Errors = append(result.Errors, fmt.Sprintf("planningProjectId=%d: đồ án không tồn tại", id))
		}
		if result.Total < int64(len(uniqueIDs)) {
			result.Total = int64(len(uniqueIDs))
		}
	}

	return result, nil
}

func (u *qhPlanningProjectUsecase) Count(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectCount, error) {
	if planningProjectID == 0 {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "id is required",
		}
	}
	if _, err := u.repo.GetByID(ctx, planningProjectID); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	type tinResult struct {
		total int64
		err   error
	}
	tinCh := make(chan tinResult, 1)
	go func() {
		if u.crmSeoProvider == nil {
			tinCh <- tinResult{err: fmt.Errorf("crm seo provider chưa được cấu hình")}
			return
		}
		total, err := u.crmSeoProvider.CountSeoByPlanningProject(ctx, planningProjectID)
		tinCh <- tinResult{total: total, err: err}
	}()

	timeline, err := u.eventRepo.CountByProject(ctx, planningProjectID)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}
	relation, err := u.eventRepo.CountRelationsByProject(ctx, planningProjectID)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}
	doc, err := u.documentRepo.CountByProject(ctx, planningProjectID)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}
	layer, err := u.repo.CountLayersByProject(ctx, planningProjectID)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}

	tin := <-tinCh
	if tin.err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: tin.err.Error()}
	}

	return &qh_dto.PlanningProjectCount{
		Timeline: timeline,
		Tin:      tin.total,
		Relation: relation,
		Doc:      doc,
		Layer:    layer,
	}, nil
}
