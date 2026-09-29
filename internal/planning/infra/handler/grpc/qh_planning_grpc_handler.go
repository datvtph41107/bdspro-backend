package handler_grpc

import (
	_err "common/domain/err"
	"context"
	"errors"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	qh_dto "tqd/internal/dto/qh"
	qh_usecase "tqd/internal/usecase/qh"

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

func (h *QHPlanningGrpcHandler) UpdatePlanningProject(ctx context.Context, req *tqdpb.PlanningProject) (*tqdpb.PlanningProject, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	project := h.mapper.ToProjectDomain(req)
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
	listReq := &qh_dto.ListPlanningProjectsRequest{
		Search:         req.Search,
		PlanningType:   req.PlanningType,
		PlanningLevel:  req.PlanningLevel,
		ValidityStatus: req.ValidityStatus,
		ProcessStatus:  req.ProcessStatus,
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
		return nil, status.Error(codes.Internal, err.Error())
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

func (h *QHPlanningGrpcHandler) ListPlanningDocuments(ctx context.Context, req *tqdpb.ListPlanningDocumentsRequest) (*tqdpb.ListPlanningDocumentsResponse, error) {
	listReq := &qh_dto.ListPlanningDocumentsRequest{
		PlanningProjectID: req.ProjectId,
		DocumentType:      req.DocumentType,
	}
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

// Client APIs

func (h *QHPlanningGrpcHandler) ClientGetPlanningProject(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningProject, error) {
	return h.GetPlanningProject(ctx, req)
}

func (h *QHPlanningGrpcHandler) ClientListPlanningProjects(ctx context.Context, req *tqdpb.ListPlanningProjectsRequest) (*tqdpb.ListPlanningProjectsResponse, error) {
	return h.ListPlanningProjects(ctx, req)
}

func (h *QHPlanningGrpcHandler) ClientGetPlanningDocument(ctx context.Context, req *tqdpb.IdRequest) (*tqdpb.PlanningDocument, error) {
	return h.GetPlanningDocument(ctx, req)
}

func (h *QHPlanningGrpcHandler) ClientListPlanningDocuments(ctx context.Context, req *tqdpb.ListPlanningDocumentsRequest) (*tqdpb.ListPlanningDocumentsResponse, error) {
	return h.ListPlanningDocuments(ctx, req)
}

func (h *QHPlanningGrpcHandler) ClientListPlanningEvents(ctx context.Context, req *tqdpb.ListPlanningEventsRequest) (*tqdpb.ListPlanningEventsResponse, error) {
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
	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: req.ProjectId,
	}
	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	results, total, err := h.eventUsecase.GetRelations(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningRelationsResponse{
		Data:  h.mapper.ToRelationProtoList(results, req.ProjectId),
		Total: total,
	}, nil
}
