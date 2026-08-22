package handler_grpc

import (
	_err "common/domain/err"
	"context"
	"errors"
	"log"
	tqdpb "pb/types/tqd"
	"strings"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	qh_usecase "tqd/internal/usecase/qh"
	tqdutils "tqd/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func planningProjectStatusError(err error) error {
	var dto *_err.ErrorDTO
	if errors.As(err, &dto) {
		switch dto.Code {
		case 400:
			return status.Error(codes.InvalidArgument, dto.Message)
		case 404:
			return status.Error(codes.NotFound, dto.Message)
		}
	}
	return status.Error(codes.Internal, err.Error())
}

type QHPlanningGrpcHandler struct {
	tqdpb.UnimplementedQHPlanningServiceServer
	projectUsecase  qh_usecase.IQHPlanningProjectUsecase
	documentUsecase qh_usecase.IQHPlanningDocumentUsecase
	eventUsecase    qh_usecase.IQHPlanningEventUsecase
	folderUsecase   qh_usecase.IQHPlanningFolderUsecase
	aiJobUsecase    qh_usecase.IProAIJobUsecase
	mapper          *mapper.QHPlanningMapper
	validator       *validator.QHPlanningValidator
}

func NewQHPlanningGrpcHandler(
	projectUsecase qh_usecase.IQHPlanningProjectUsecase,
	documentUsecase qh_usecase.IQHPlanningDocumentUsecase,
	eventUsecase qh_usecase.IQHPlanningEventUsecase,
	folderUsecase qh_usecase.IQHPlanningFolderUsecase,
	aiJobUsecase qh_usecase.IProAIJobUsecase,
	mapper *mapper.QHPlanningMapper,
	validator *validator.QHPlanningValidator,
) *QHPlanningGrpcHandler {
	return &QHPlanningGrpcHandler{
		projectUsecase:  projectUsecase,
		documentUsecase: documentUsecase,
		eventUsecase:    eventUsecase,
		folderUsecase:   folderUsecase,
		aiJobUsecase:    aiJobUsecase,
		mapper:          mapper,
		validator:       validator,
	}
}

// QHPlanningProject APIs

func (h *QHPlanningGrpcHandler) CreatePlanningProject(ctx context.Context, req *tqdpb.PlanningProject) (*tqdpb.PlanningProject, error) {
	project := h.mapper.ToProjectDomain(req)
	if err := h.validator.ValidateProject(project); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.projectUsecase.Create(ctx, project)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToProjectProto(result), nil
}

func (h *QHPlanningGrpcHandler) GetPlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProject, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.projectUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return h.mapper.ToProjectProto(result), nil
}

// CountPlanningProject @Summary Đếm timeline, tin, relation, tài liệu, layer của đồ án
// @Description Đếm số sự kiện (timeline), tin SEO, quan hệ đồ án, tài liệu và layer liên quan thuộc planning project. Số tin gọi CRM CountSeoRefs.
// @Tags Planning Project Admin
// @Produce json
// @Param id path int true "Planning project ID"
// @Success 200 {object} tqdpb.PlanningProjectCountResponse
// @Router /v2/tqd/admin/planning-projects/{id}/count [get]
func (h *QHPlanningGrpcHandler) CountPlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProjectCountResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.projectUsecase.Count(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return &tqdpb.PlanningProjectCountResponse{
		Timeline: result.Timeline,
		Tin:      result.Tin,
		Relation: result.Relation,
		Doc:      result.Doc,
		Layer:    result.Layer,
	}, nil
}

func (h *QHPlanningGrpcHandler) UpdatePlanningProject(ctx context.Context, req *tqdpb.PlanningProject) (*tqdpb.PlanningProject, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	project := h.mapper.ToProjectDomain(req)
	if req.ClientView == nil {
		existing, err := h.projectUsecase.GetByID(ctx, req.Id)
		if err != nil {
			return nil, planningProjectStatusError(err)
		}
		project.ClientView = existing.ClientView
	}
	if err := h.validator.ValidateProject(project); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.projectUsecase.Update(ctx, req.Id, project)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToProjectProto(result), nil
}

func (h *QHPlanningGrpcHandler) DeletePlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	success, err := h.projectUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{
		Success: success,
		Message: "Project deleted successfully",
	}, nil
}

func (h *QHPlanningGrpcHandler) ListPlanningProjects(ctx context.Context, req *tqdpb.ListPlanningProjectsRequest) (*tqdpb.ListPlanningProjectsResponse, error) {
	if req == nil {
		req = &tqdpb.ListPlanningProjectsRequest{}
	}
	listReq := &qh_dto.ListPlanningProjectsRequest{
		Search:         strings.TrimSpace(req.Search),
		PlanningType:   req.PlanningType,
		PlanningLevel:  req.PlanningLevel,
		ValidityStatus: req.ValidityStatus,
		ProcessStatus:  req.ProcessStatus,
		AdmID:          req.AdmId,
		LayerIDs:       tqdutils.MergePositiveUint64(req.GetLayerId(), tqdutils.ParseUintCSV(req.GetLayerIds())),
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.projectUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningProjectsResponse{
		Data:  h.mapper.ToProjectProtoList(results),
		Total: total,
	}, nil
}

// ListPlanningProjectLayers danh sách layer thuộc đồ án quy hoạch.
// @Summary Danh sách layer theo đồ án
// @Description Lấy danh sách layer theo planningId (qh_layers.planning_project_id hoặc gán qua layer_plan_doc_relation).
// @Tags Planning Project Admin
// @Produce json
// @Security BearerAuth
// @Param planningId path int true "Planning project ID"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} tqdpb.ListLayersResponse
// @Router /v2/tqd/admin/planning-projects/{planningId}/layers [get]
func (h *QHPlanningGrpcHandler) ListPlanningProjectLayers(ctx context.Context, req *tqdpb.ListPlanningProjectLayersRequest) (*tqdpb.ListLayersResponse, error) {
	if req == nil || req.PlanningId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningId is required")
	}

	listReq := &qh_dto.ListPlanningProjectLayersRequest{
		PlanningID: req.PlanningId,
	}
	listReq.Page = req.Page
	listReq.Size = req.Size

	layers, total, err := h.projectUsecase.ListLayers(ctx, listReq)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}

	return &tqdpb.ListLayersResponse{
		Data:     h.mapper.ToLayerProtoList(layers),
		Total:    total,
		Page:     req.Page,
		PageSize: req.Size,
	}, nil
}

// UnassignPlanningProjectLayer gỡ layer khỏi đồ án — chỉ xóa layer_plan_doc_relation.
// @Summary Gỡ layer khỏi đồ án
// @Description Xóa toàn bộ layer_plan_doc_relation của layer với tài liệu thuộc đồ án. Không xóa layer.
// @Tags Planning Project Admin
// @Produce json
// @Security BearerAuth
// @Param planningId path int true "Planning project ID"
// @Param layerId path int true "Layer ID"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/planning-projects/{planningId}/layers/{layerId} [delete]
func (h *QHPlanningGrpcHandler) UnassignPlanningProjectLayer(ctx context.Context, req *tqdpb.UnassignPlanningProjectLayerRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := h.validator.ValidateUnassignLayer(req.PlanningId, req.LayerId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := h.documentUsecase.UnassignLayer(ctx, &qh_dto.UnassignPlanningProjectLayerRequest{
		PlanningID: req.PlanningId,
		LayerID:    req.LayerId,
	})
	if err != nil {
		return nil, planningProjectStatusError(err)
	}

	return &tqdpb.SubmitResponse{
		Success: true,
		Message: "Layer unassigned successfully",
	}, nil
}

// CreatePlanningProjectFromFolder tạo đồ án + toàn bộ tài liệu từ 1 folder admin chọn (upload multipart qua gateway).
// Mỗi file trở thành 1 PlanningDocument trạng thái Pending, chờ job classify quét và gọi Gemini phân loại.
func (h *QHPlanningGrpcHandler) CreatePlanningProjectFromFolder(ctx context.Context, req *tqdpb.CreatePlanningProjectFromFolderRequest) (*tqdpb.CreatePlanningProjectFromFolderResponse, error) {
	if req.FolderName == "" {
		return nil, status.Error(codes.InvalidArgument, "folder_name is required")
	}
	if len(req.Files) == 0 {
		return nil, status.Error(codes.InvalidArgument, "files is required")
	}

	project := h.mapper.ToProjectDomainFromFolderRequest(req)
	files := h.mapper.ToFolderFileInputs(req.Files)

	resultProject, resultDocuments, aiJob, err := h.folderUsecase.CreateFromFolder(ctx, project, req.FolderName, files, req.JobType)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.CreatePlanningProjectFromFolderResponse{
		Project:   h.mapper.ToProjectProto(resultProject),
		Documents: h.mapper.ToDocumentProtoList(resultDocuments),
		AiJob:     h.mapper.ToProAIJobProto(aiJob),
	}, nil
}

// AdminAssignProvinces gán provinces cho planning project
func (h *QHPlanningGrpcHandler) AdminAssignProvinces(ctx context.Context, req *tqdpb.AdminAssignProvincesRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningId is required")
	}
	if err := h.projectUsecase.AssignProvinces(ctx, req.PlanningId, req.ProvinceIds); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{Success: true, Message: "provinces assigned"}, nil
}

// AdminUnassignProvinces gỡ provinces khỏi planning project (xóa bản ghi)
func (h *QHPlanningGrpcHandler) AdminUnassignProvinces(ctx context.Context, req *tqdpb.AdminUnassignProvincesRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningId is required")
	}
	if err := h.projectUsecase.UnassignProvinces(ctx, req.PlanningId, req.ProvinceIds); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{Success: true, Message: "provinces unassigned"}, nil
}

// AdminAssignWards gán wards cho planning project
func (h *QHPlanningGrpcHandler) AdminAssignWards(ctx context.Context, req *tqdpb.AdminAssignWardsRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningId is required")
	}
	if err := h.projectUsecase.AssignWards(ctx, req.PlanningId, req.WardIds); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{Success: true, Message: "wards assigned"}, nil
}

// AdminUnassignWards gỡ wards khỏi planning project (xóa bản ghi)
func (h *QHPlanningGrpcHandler) AdminUnassignWards(ctx context.Context, req *tqdpb.AdminUnassignWardsRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningId is required")
	}
	if err := h.projectUsecase.UnassignWards(ctx, req.PlanningId, req.WardIds); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{Success: true, Message: "wards unassigned"}, nil
}

func (h *QHPlanningGrpcHandler) ListProAIJobs(ctx context.Context, req *tqdpb.ListProAIJobsRequest) (*tqdpb.ListProAIJobsResponse, error) {
	listReq := &qh_dto.ListProAIJobsRequest{
		Search:        req.Search,
		JobType:       req.JobType,
		ProcessStatus: req.ProcessStatus,
	}
	listReq.Page = req.Page
	listReq.Size = req.Size

	results, total, err := h.aiJobUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.ListProAIJobsResponse{
		Data:  h.mapper.ToProAIJobProtoList(results),
		Total: total,
	}, nil
}

func (h *QHPlanningGrpcHandler) GetProAIJob(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.ProAIJob, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.aiJobUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return h.mapper.ToProAIJobProto(result), nil
}

// ApprovePlanningProject phê duyệt đồ án (chỉ khi đã Classified), đồng thời phê duyệt các tài liệu đã Classified.
func (h *QHPlanningGrpcHandler) ApprovePlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProject, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.projectUsecase.Approve(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToProjectProto(result), nil
}

// QHPlanningDocument APIs

func (h *QHPlanningGrpcHandler) CreatePlanningDocument(ctx context.Context, req *tqdpb.PlanningDocument) (*tqdpb.PlanningDocument, error) {
	document := h.mapper.ToDocumentDomain(req)
	if err := h.validator.ValidateDocument(document); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.documentUsecase.Create(ctx, document)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToDocumentProto(result), nil
}

func (h *QHPlanningGrpcHandler) GetPlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocument, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.documentUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return h.mapper.ToDocumentProto(result), nil
}

// UpdatePlanningDocument cập nhật toàn bộ thông tin tài liệu quy hoạch.
// @Summary Cập nhật tài liệu quy hoạch
// @Tags Planning Document Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID tài liệu"
// @Param body body tqdpb.PlanningDocument true "Thông tin tài liệu"
// @Success 200 {object} tqdpb.PlanningDocument
// @Router /v2/tqd/admin/planning-documents/{id} [put]
func (h *QHPlanningGrpcHandler) UpdatePlanningDocument(ctx context.Context, req *tqdpb.PlanningDocument) (*tqdpb.PlanningDocument, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	document := h.mapper.ToDocumentDomain(req)
	if err := h.validator.ValidateDocument(document); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.documentUsecase.Update(ctx, req.Id, document)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToDocumentProto(result), nil
}

func (h *QHPlanningGrpcHandler) DeletePlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	success, err := h.documentUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &tqdpb.SubmitResponse{
		Success: success,
		Message: "Document deleted successfully",
	}, nil
}

// CreatePlanningDocumentHistory tạo lịch sử pháp lý cho 1 tài liệu quy hoạch.
// @Summary Tạo lịch sử tài liệu quy hoạch
// @Tags Planning Document Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param planningDocumentId path int true "Planning document ID"
// @Param request body tqdpb.PlanningDocumentHistory true "Lịch sử tài liệu"
// @Success 200 {object} tqdpb.PlanningDocumentHistory
// @Router /v2/tqd/admin/planning-documents/{planningDocumentId}/histories [post]
func (h *QHPlanningGrpcHandler) CreatePlanningDocumentHistory(ctx context.Context, req *tqdpb.PlanningDocumentHistory) (*tqdpb.PlanningDocumentHistory, error) {
	history := h.mapper.ToDocumentHistoryDomain(req)
	if err := h.validator.ValidateDocumentHistory(history); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.documentUsecase.CreateHistory(ctx, history)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToDocumentHistoryProto(result), nil
}

// DeletePlanningDocumentHistory xóa cứng lịch sử tài liệu.
// @Summary Xóa lịch sử tài liệu quy hoạch
// @Tags Planning Document Admin
// @Produce json
// @Security BearerAuth
// @Param planningDocumentId path int true "Planning document ID"
// @Param id path int true "History ID"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/planning-documents/{planningDocumentId}/histories/{id} [delete]
func (h *QHPlanningGrpcHandler) DeletePlanningDocumentHistory(ctx context.Context, req *tqdpb.DeletePlanningDocumentHistoryRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.Id == 0 || req.PlanningDocumentId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningDocumentId and id are required")
	}
	success, err := h.documentUsecase.DeleteHistory(ctx, req.PlanningDocumentId, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return &tqdpb.SubmitResponse{
		Success: success,
		Message: "Document history deleted successfully",
	}, nil
}

// ListPlanningDocumentHistories danh sách lịch sử tài liệu quy hoạch.
// @Summary Danh sách lịch sử tài liệu quy hoạch
// @Description Danh sách lịch sử theo docId, sắp xếp theo eventDate giảm dần.
// @Tags Planning Document Admin
// @Produce json
// @Security BearerAuth
// @Param planningDocumentId path int true "Planning document ID"
// @Param page query int false "Trang"
// @Param size query int false "Số bản ghi"
// @Success 200 {object} tqdpb.ListPlanningDocumentHistoriesResponse
// @Router /v2/tqd/admin/planning-documents/{planningDocumentId}/histories [get]
func (h *QHPlanningGrpcHandler) ListPlanningDocumentHistories(ctx context.Context, req *tqdpb.ListPlanningDocumentHistoriesRequest) (*tqdpb.ListPlanningDocumentHistoriesResponse, error) {
	if req == nil || req.PlanningDocumentId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningDocumentId is required")
	}
	listReq := &qh_dto.ListPlanningDocumentHistoriesRequest{
		PlanningDocumentID: req.PlanningDocumentId,
	}
	listReq.Page = req.Page
	listReq.Size = req.Size
	results, total, err := h.documentUsecase.ListHistories(ctx, listReq)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return &tqdpb.ListPlanningDocumentHistoriesResponse{
		Data:  h.mapper.ToDocumentHistoryProtoList(results),
		Total: total,
	}, nil
}

// ListPlanningDocuments danh sách tài liệu quy hoạch.
// @Summary Danh sách tài liệu quy hoạch
// @Tags Planning Document Admin
// @Produce json
// @Security BearerAuth
// @Param planningProjectId query int false "ID đồ án"
// @Param keyword query string false "Tìm theo title/code/description"
// @Param layerId query int false "Lọc theo 1 layer"
// @Param layerIds query string false "CSV layer id, lọc tài liệu gắn layer qua layer_plan_doc_relation (vd. 1279,1278)"
// @Param parcelId query int false "Lọc theo thửa (parcel → regions → layers)"
// @Param documentType query string false "Loại tài liệu"
// @Param status query string false "Trạng thái hiệu lực"
// @Param page query int false "Trang"
// @Param size query int false "Số bản ghi"
// @Success 200 {object} tqdpb.ListPlanningDocumentsResponse
// @Router /v2/tqd/admin/planning-documents [get]
func (h *QHPlanningGrpcHandler) ListPlanningDocuments(ctx context.Context, req *tqdpb.ListPlanningDocumentsRequest) (*tqdpb.ListPlanningDocumentsResponse, error) {
	if req == nil {
		req = &tqdpb.ListPlanningDocumentsRequest{}
	}
	planningProjectID := req.PlanningProjectId
	if planningProjectID == 0 {
		planningProjectID = req.ProjectId
	}

	listReq := &qh_dto.ListPlanningDocumentsRequest{
		PlanningProjectID: planningProjectID,
		DocumentType:      req.DocumentType,
		ValidityStatus:    req.Status,
		Keyword:           strings.TrimSpace(req.GetKeyword()),
		LayerID:           req.LayerId,
		LayerIDs:          tqdutils.MergePositiveUint64(req.LayerId, tqdutils.ParseUintCSV(req.GetLayerIds())),
		ParcelID:          req.GetParcelId(),
	}
	log.Printf("[ListPlanningDocuments] planningProjectID=%v, parcelId=%v, layerIds=%v, keyword=%v, documentType=%v, validityStatus=%v", planningProjectID, req.GetParcelId(), req.GetLayerIds(), listReq.Keyword, req.DocumentType, req.Status)
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.documentUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningDocumentsResponse{
		Data:  h.mapper.ToDocumentProtoList(results),
		Total: total,
	}, nil
}

// ApprovePlanningDocument phê duyệt 1 tài liệu đã được AI phân loại (chỉ khi đang ở trạng thái Classified).
func (h *QHPlanningGrpcHandler) ApprovePlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocument, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.documentUsecase.Approve(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToDocumentProto(result), nil
}

// RetryPlanningDocument đưa tài liệu Failed về Pending để job classify chạy lại.
func (h *QHPlanningGrpcHandler) RetryPlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocument, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.documentUsecase.Retry(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToDocumentProto(result), nil
}

// AssignPlanningDocumentsToLayer gán tài liệu vào layer.
// @Summary Gán tài liệu vào layer
// @Description Tạo bản ghi layer_plan_doc_relation từ layerIds và docIds (cartesian product).
// @Tags Layer Doc Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AssignPlanningDocumentsToLayerRequest true "layerIds và docIds"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/layer-docs [post]
func (h *QHPlanningGrpcHandler) AssignPlanningDocumentsToLayer(ctx context.Context, req *tqdpb.AssignPlanningDocumentsToLayerRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := h.validator.ValidateAssignLayer(req.LayerId, req.LayerIds, req.DocIds); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := h.documentUsecase.AssignToLayer(ctx, &qh_dto.AssignPlanningDocumentsToLayerRequest{
		LayerID:  req.LayerId,
		LayerIDs: req.LayerIds,
		DocIDs:   req.DocIds,
	})
	if err != nil {
		return nil, planningProjectStatusError(err)
	}

	return &tqdpb.SubmitResponse{
		Success: true,
		Message: "Documents assigned to layer successfully",
	}, nil
}

// UnassignPlanningDocumentsFromLayer gỡ tài liệu khỏi layer.
// @Summary Gỡ tài liệu khỏi layer
// @Description Xóa cứng bản ghi layer_plan_doc_relation theo layerIds và docIds.
// @Tags Layer Doc Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AssignPlanningDocumentsToLayerRequest true "layerIds và docIds"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/layer-docs [delete]
func (h *QHPlanningGrpcHandler) UnassignPlanningDocumentsFromLayer(ctx context.Context, req *tqdpb.AssignPlanningDocumentsToLayerRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := h.validator.ValidateAssignLayer(req.LayerId, req.LayerIds, req.DocIds); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := h.documentUsecase.UnassignFromLayer(ctx, &qh_dto.AssignPlanningDocumentsToLayerRequest{
		LayerID:  req.LayerId,
		LayerIDs: req.LayerIds,
		DocIDs:   req.DocIds,
	})
	if err != nil {
		return nil, planningProjectStatusError(err)
	}

	return &tqdpb.SubmitResponse{
		Success: true,
		Message: "Documents unassigned from layer successfully",
	}, nil
}

// RetryFailedPlanningDocuments đưa tất cả tài liệu Failed của đồ án về Pending.
func (h *QHPlanningGrpcHandler) RetryFailedPlanningDocuments(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProject, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.projectUsecase.RetryFailedDocuments(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return h.mapper.ToProjectProto(result), nil
}

func (h *QHPlanningGrpcHandler) GetPlanningProjectSeoSource(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProjectSeoSourceResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	source, err := h.projectUsecase.GetPlanningProjectSeoSource(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return planningProjectSeoSourceToProto(source), nil
}

func (h *QHPlanningGrpcHandler) UpdatePlanningProjectSeoID(ctx context.Context, req *tqdpb.UpdatePlanningProjectSeoIDRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningProjectId == 0 || req.SeoId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningProjectId và seoId là bắt buộc")
	}

	if err := h.projectUsecase.UpdatePlanningProjectSeoID(ctx, req.PlanningProjectId, req.SeoId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.SubmitResponse{
		Success: true,
		Message: "planning project seo_id updated successfully",
	}, nil
}

func (h *QHPlanningGrpcHandler) GetPlanningDocumentSeoSource(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocumentSeoSourceResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	source, err := h.documentUsecase.GetPlanningDocumentSeoSource(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return planningDocumentSeoSourceToProto(source), nil
}

func (h *QHPlanningGrpcHandler) UpdatePlanningDocumentSeoID(ctx context.Context, req *tqdpb.UpdatePlanningDocumentSeoIDRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.PlanningDocumentId == 0 || req.SeoId == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningDocumentId và seoId là bắt buộc")
	}

	if err := h.documentUsecase.UpdatePlanningDocumentSeoID(ctx, req.PlanningDocumentId, req.SeoId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.SubmitResponse{
		Success: true,
		Message: "planning document seo_id updated successfully",
	}, nil
}

// AdminGenPlanningDocumentSeo tạo seo_domain trên CRM cho danh sách tài liệu quy hoạch.
// @Summary Admin tạo SEO domain từ planning document
// @Description Nhận danh sách docIds, gọi CRM sinh bản ghi seo_domain tương ứng.
// @Tags Planning Document Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AdminGenPlanningDocumentSeoRequest true "Danh sách docIds"
// @Router /v2/tqd/admin/planning-documents/gen-seo [post]
func (h *QHPlanningGrpcHandler) AdminGenPlanningDocumentSeo(ctx context.Context, req *tqdpb.AdminGenPlanningDocumentSeoRequest) (*tqdpb.AdminGenPlanningDocumentSeoResponse, error) {
	if req == nil || len(req.DocIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "docIds is required")
	}

	log.Printf("[AdminGenPlanningDocumentSeo] docIds=%v", req.DocIds)

	result, err := h.documentUsecase.AdminGenPlanningDocumentSeo(ctx, req.DocIds)
	if err != nil {
		log.Printf("[AdminGenPlanningDocumentSeo] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result == nil {
		return &tqdpb.AdminGenPlanningDocumentSeoResponse{}, nil
	}

	items := make([]*tqdpb.AdminGenPlanningDocumentSeoItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, &tqdpb.AdminGenPlanningDocumentSeoItem{
			PlanningDocumentId: item.PlanningDocumentID,
			SeoId:              item.SeoID,
			CanonicalUrl:       item.CanonicalURL,
			Slug:               item.Slug,
			Created:            item.Created,
			Message:            item.Message,
		})
	}

	return &tqdpb.AdminGenPlanningDocumentSeoResponse{
		Data:    items,
		Total:   result.Total,
		Created: result.Created,
		Skipped: result.Skipped,
		Failed:  result.Failed,
		Errors:  result.Errors,
	}, nil
}

// AdminGenPlanningProjectSeo tạo seo_domain trên CRM cho danh sách đồ án quy hoạch.
// @Summary Admin tạo SEO domain từ planning project
// @Description Nhận danh sách planningProjectIds, gọi CRM sinh bản ghi seo_domain tương ứng.
// @Tags Planning Project Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AdminGenPlanningProjectSeoRequest true "Danh sách planningProjectIds"
// @Router /v2/tqd/admin/planning-projects/gen-seo [post]
func (h *QHPlanningGrpcHandler) AdminGenPlanningProjectSeo(ctx context.Context, req *tqdpb.AdminGenPlanningProjectSeoRequest) (*tqdpb.AdminGenPlanningProjectSeoResponse, error) {
	if req == nil || len(req.PlanningProjectIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "planningProjectIds is required")
	}

	log.Printf("[AdminGenPlanningProjectSeo] planningProjectIds=%v", req.PlanningProjectIds)

	result, err := h.projectUsecase.AdminGenPlanningProjectSeo(ctx, req.PlanningProjectIds)
	if err != nil {
		log.Printf("[AdminGenPlanningProjectSeo] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result == nil {
		return &tqdpb.AdminGenPlanningProjectSeoResponse{}, nil
	}

	items := make([]*tqdpb.AdminGenPlanningProjectSeoItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, &tqdpb.AdminGenPlanningProjectSeoItem{
			PlanningProjectId: item.PlanningProjectID,
			SeoId:             item.SeoID,
			CanonicalUrl:      item.CanonicalURL,
			Slug:              item.Slug,
			Created:           item.Created,
			Message:           item.Message,
		})
	}

	return &tqdpb.AdminGenPlanningProjectSeoResponse{
		Data:    items,
		Total:   result.Total,
		Created: result.Created,
		Skipped: result.Skipped,
		Failed:  result.Failed,
		Errors:  result.Errors,
	}, nil
}

func planningProjectSeoSourceToProto(source *qh_dto.PlanningProjectSeoSource) *tqdpb.PlanningProjectSeoSourceResponse {
	if source == nil || source.PlanningProjectID == 0 {
		return nil
	}
	resp := &tqdpb.PlanningProjectSeoSourceResponse{
		PlanningProjectId: source.PlanningProjectID,
		Name:              source.Name,
		Slug:              source.Slug,
	}
	if source.SeoID != nil && *source.SeoID > 0 {
		resp.SeoId = *source.SeoID
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		resp.ProvinceId = *source.ProvinceID
	}
	return resp
}

func planningDocumentSeoSourceToProto(source *qh_dto.PlanningDocumentSeoSource) *tqdpb.PlanningDocumentSeoSourceResponse {
	if source == nil || source.PlanningDocumentID == 0 {
		return nil
	}
	resp := &tqdpb.PlanningDocumentSeoSourceResponse{
		PlanningDocumentId: source.PlanningDocumentID,
		Title:              source.Title,
		Slug:               source.Slug,
		Filepath:           source.Filepath,
		Thumbnail:          source.Thumbnail,
		Ext:                source.Ext,
	}
	if source.SeoID != nil && *source.SeoID > 0 {
		resp.SeoId = *source.SeoID
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		resp.ProvinceId = *source.ProvinceID
	}
	return resp
}

// QHPlanningEvent APIs

// CreatePlanningEvent tạo mốc lịch sử pháp lý cho đồ án.
// @Summary Tạo sự kiện timeline đồ án
// @Description Tạo sự kiện lịch sử pháp lý gắn với planningProjectId.
// @Tags Planning Event Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.PlanningEvent true "Sự kiện timeline"
// @Success 200 {object} tqdpb.PlanningEvent
// @Router /v2/tqd/admin/planning-events [post]
func (h *QHPlanningGrpcHandler) CreatePlanningEvent(ctx context.Context, req *tqdpb.PlanningEvent) (*tqdpb.PlanningEvent, error) {
	event := h.mapper.ToEventDomain(req)
	if err := h.validator.ValidateEvent(event); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.eventUsecase.Create(ctx, event)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToEventProto(result), nil
}

// GetPlanningEvent lấy chi tiết 1 sự kiện timeline.
// @Summary Chi tiết sự kiện timeline
// @Tags Planning Event Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} tqdpb.PlanningEvent
// @Router /v2/tqd/admin/planning-events/{id} [get]
func (h *QHPlanningGrpcHandler) GetPlanningEvent(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningEvent, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.eventUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToEventProto(result), nil
}

// UpdatePlanningEvent cập nhật sự kiện timeline.
// @Summary Cập nhật sự kiện timeline
// @Tags Planning Event Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param request body tqdpb.PlanningEvent true "Sự kiện timeline"
// @Success 200 {object} tqdpb.PlanningEvent
// @Router /v2/tqd/admin/planning-events/{id} [put]
func (h *QHPlanningGrpcHandler) UpdatePlanningEvent(ctx context.Context, req *tqdpb.PlanningEvent) (*tqdpb.PlanningEvent, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	event := h.mapper.ToEventDomain(req)
	if err := h.validator.ValidateEvent(event); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.eventUsecase.Update(ctx, req.Id, event)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToEventProto(result), nil
}

// DeletePlanningEvent xóa cứng sự kiện timeline (DELETE SQL, không soft delete).
// @Summary Xóa sự kiện timeline
// @Description Xóa bản ghi qh_planning_events bằng DELETE, không đánh dấu deleted_at.
// @Tags Planning Event Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/planning-events/{id} [delete]
func (h *QHPlanningGrpcHandler) DeletePlanningEvent(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	success, err := h.eventUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return &tqdpb.SubmitResponse{
		Success: success,
		Message: "Event deleted successfully",
	}, nil
}

// ListPlanningEvents danh sách sự kiện timeline cho admin.
// @Summary Danh sách sự kiện timeline
// @Tags Planning Event Admin
// @Produce json
// @Security BearerAuth
// @Param planningProjectId query int false "Planning project ID"
// @Param projectId query int false "Alias project ID"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} tqdpb.ListPlanningEventsResponse
// @Router /v2/tqd/admin/planning-events [get]
func (h *QHPlanningGrpcHandler) ListPlanningEvents(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*tqdpb.ListPlanningEventsResponse, error) {
	if req == nil {
		req = &tqdpb.ListPlanningEventsRequest{}
	}
	planningProjectID := req.PlanningProjectId
	if planningProjectID == 0 {
		planningProjectID = req.ProjectId
	}

	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: planningProjectID,
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.eventUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningEventsResponse{
		Data:  h.mapper.ToEventProtoList(results),
		Total: total,
	}, nil
}

// QHPlanningRelation APIs

// CreatePlanningRelation tạo quan hệ parent/child/peer giữa hai đồ án.
// @Summary Tạo quan hệ đồ án
// @Description fromId là đồ án nguồn (parent nếu không phải peer), toId là đồ án đích (child hoặc peer).
// @Tags Planning Relation Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.PlanningRelation true "Quan hệ đồ án"
// @Success 200 {object} tqdpb.PlanningRelation
// @Router /v2/tqd/admin/planning-relations [post]
func (h *QHPlanningGrpcHandler) CreatePlanningRelation(ctx context.Context, req *tqdpb.PlanningRelation) (*tqdpb.PlanningRelation, error) {
	relation := h.mapper.ToRelationDomain(req)
	if err := h.validator.ValidateRelation(relation); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.eventUsecase.CreateRelation(ctx, relation)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToRelationAdminProto(result), nil
}

// GetPlanningRelation lấy chi tiết 1 quan hệ đồ án.
// @Summary Chi tiết quan hệ đồ án
// @Tags Planning Relation Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Relation ID"
// @Success 200 {object} tqdpb.PlanningRelation
// @Router /v2/tqd/admin/planning-relations/{id} [get]
func (h *QHPlanningGrpcHandler) GetPlanningRelation(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningRelation, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.eventUsecase.GetRelationByID(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToRelationAdminProto(result), nil
}

// UpdatePlanningRelation cập nhật quan hệ đồ án.
// @Summary Cập nhật quan hệ đồ án
// @Tags Planning Relation Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Relation ID"
// @Param request body tqdpb.PlanningRelation true "Quan hệ đồ án"
// @Success 200 {object} tqdpb.PlanningRelation
// @Router /v2/tqd/admin/planning-relations/{id} [put]
func (h *QHPlanningGrpcHandler) UpdatePlanningRelation(ctx context.Context, req *tqdpb.PlanningRelation) (*tqdpb.PlanningRelation, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	relation := h.mapper.ToRelationDomain(req)
	if err := h.validator.ValidateRelation(relation); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := h.eventUsecase.UpdateRelation(ctx, req.Id, relation)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return h.mapper.ToRelationAdminProto(result), nil
}

// DeletePlanningRelation xóa cứng quan hệ đồ án (DELETE SQL, không soft delete).
// @Summary Xóa quan hệ đồ án
// @Description Xóa bản ghi qh_planning_relations bằng DELETE, không đánh dấu deleted_at.
// @Tags Planning Relation Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Relation ID"
// @Success 200 {object} tqdpb.SubmitResponse
// @Router /v2/tqd/admin/planning-relations/{id} [delete]
func (h *QHPlanningGrpcHandler) DeletePlanningRelation(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	success, err := h.eventUsecase.DeleteRelation(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	return &tqdpb.SubmitResponse{
		Success: success,
		Message: "Relation deleted successfully",
	}, nil
}

// ListPlanningRelations danh sách quan hệ đồ án cho admin.
// @Summary Danh sách quan hệ đồ án
// @Tags Planning Relation Admin
// @Produce json
// @Security BearerAuth
// @Param planningProjectId query int false "Planning project ID"
// @Param projectId query int false "Alias project ID"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} tqdpb.ListPlanningRelationsResponse
// @Router /v2/tqd/admin/planning-relations [get]
func (h *QHPlanningGrpcHandler) ListPlanningRelations(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*tqdpb.ListPlanningRelationsResponse, error) {
	if req == nil {
		req = &tqdpb.ListPlanningEventsRequest{}
	}
	planningProjectID := req.PlanningProjectId
	if planningProjectID == 0 {
		planningProjectID = req.ProjectId
	}

	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: planningProjectID,
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.eventUsecase.GetRelations(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningRelationsResponse{
		Data:  h.mapper.ToRelationAdminProtoList(results),
		Total: total,
	}, nil
}

// Client APIs

func (h *QHPlanningGrpcHandler) ClientGetPlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProject, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.projectUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, planningProjectStatusError(err)
	}
	if !result.ClientView {
		return nil, status.Error(codes.NotFound, "Không tìm thấy dữ liệu")
	}
	return h.mapper.ToProjectProto(result), nil
}

func (h *QHPlanningGrpcHandler) ClientListPlanningProjects(ctx context.Context, req *tqdpb.ClientListPlanningProjectsRequest) (*tqdpb.ListPlanningProjectsResponse, error) {
	if req == nil {
		req = &tqdpb.ClientListPlanningProjectsRequest{}
	}
	listReq := &qh_dto.ListPlanningProjectsRequest{
		Search:         strings.TrimSpace(req.GetSearch()),
		ClientViewOnly: true,
		AdmID:          req.GetAdmId(),
	}
	listReq.Page = uint32(req.GetPage())
	listReq.Size = uint32(req.GetSize())

	results, total, err := h.projectUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningProjectsResponse{
		Data:  h.mapper.ToProjectProtoList(results),
		Total: total,
	}, nil
}

func (h *QHPlanningGrpcHandler) ClientGetPlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocument, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.documentUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := h.requireClientVisibleProject(ctx, result.PlanningProjectID); err != nil {
		return nil, err
	}
	return h.mapper.ToDocumentProto(result), nil
}

func (h *QHPlanningGrpcHandler) ClientListPlanningDocuments(ctx context.Context, req *tqdpb.ListPlanningDocumentsRequest) (*tqdpb.ListPlanningDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	planningProjectID := req.PlanningProjectId
	if planningProjectID == 0 {
		planningProjectID = req.ProjectId
	}
	if err := h.requireClientVisibleProject(ctx, planningProjectID); err != nil {
		return nil, err
	}
	return h.ListPlanningDocuments(ctx, req)
}

// ClientListPlanningEvents trả lịch sử pháp lý của đồ án theo timeline.
// @Summary Lịch sử đồ án quy hoạch
// @Description Danh sách sự kiện pháp lý của đồ án, sắp xếp theo ngày sự kiện giảm dần. Dùng để hiển thị timeline.
// @Tags Planning Project Client
// @Accept json
// @Produce json
// @Param project_id path int true "Planning project ID"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} tqdpb.ListPlanningEventsResponse
// @Router /v2/tqd/client/planning-projects/{project_id}/timeline [get]
func (h *QHPlanningGrpcHandler) ClientListPlanningEvents(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*tqdpb.ListPlanningEventsResponse, error) {
	if req == nil || req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}
	if err := h.requireClientVisibleProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}

	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: req.ProjectId,
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.eventUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningEventsResponse{
		Data:  h.mapper.ToEventProtoList(results),
		Total: total,
	}, nil
}

func (h *QHPlanningGrpcHandler) ClientListPlanningRelations(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*tqdpb.ListPlanningRelationsResponse, error) {
	if req == nil || req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}
	if err := h.requireClientVisibleProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}

	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: req.ProjectId,
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, _, err := h.eventUsecase.GetRelations(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	visible := make([]qh_domain.QHPlanningRelationItem, 0, len(results))
	for i := range results {
		if clientVisibleRelation(&results[i], req.ProjectId) {
			visible = append(visible, results[i])
		}
	}

	return &tqdpb.ListPlanningRelationsResponse{
		Data:  h.mapper.ToRelationProtoList(visible, req.ProjectId),
		Total: int64(len(visible)),
	}, nil
}

// ClientCountPlanningProject @Summary Đếm timeline, tin, relation, tài liệu, layer của đồ án (client)
// @Description Đếm số sự kiện (timeline), tin SEO, quan hệ đồ án, tài liệu và layer liên quan thuộc planning project.
// @Tags Planning Project Client
// @Produce json
// @Param project_id path int true "Planning project ID"
// @Success 200 {object} tqdpb.PlanningProjectCountResponse
// @Router /v2/tqd/client/planning-projects/{project_id}/count [get]
func (h *QHPlanningGrpcHandler) ClientCountPlanningProject(ctx context.Context, req *tqdpb.PlanningProjectIdRequest) (*tqdpb.PlanningProjectCountResponse, error) {
	if req == nil || req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}
	if err := h.requireClientVisibleProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}
	return h.CountPlanningProject(ctx, &tqdpb.IdRequest{Id: req.ProjectId})
}

func (h *QHPlanningGrpcHandler) requireClientVisibleProject(ctx context.Context, id uint64) error {
	project, err := h.projectUsecase.GetByID(ctx, id)
	if err != nil {
		return planningProjectStatusError(err)
	}
	if !project.ClientView {
		return status.Error(codes.NotFound, "Không tìm thấy dữ liệu")
	}
	return nil
}

func clientVisibleRelation(item *qh_domain.QHPlanningRelationItem, projectID uint64) bool {
	if item == nil {
		return false
	}
	related := []*qh_domain.QHPlanningProject{item.Parent, item.Child, item.Peer}
	for _, project := range related {
		if project == nil || project.ID == projectID {
			continue
		}
		if !project.ClientView {
			return false
		}
	}
	return true
}
