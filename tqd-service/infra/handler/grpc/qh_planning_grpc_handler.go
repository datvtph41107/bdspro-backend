package handler_grpc

import (
	"context"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	qh_dto "tqd/internal/dto/qh"
	qh_usecase "tqd/internal/usecase/qh"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QHPlanningGrpcHandler struct {
	tqdpb.UnimplementedQHPlanningServiceServer
	projectUsecase  qh_usecase.IQHPlanningProjectUsecase
	documentUsecase qh_usecase.IQHPlanningDocumentUsecase
	mapper          *mapper.QHPlanningMapper
	validator       *validator.QHPlanningValidator
}

func NewQHPlanningGrpcHandler(
	projectUsecase qh_usecase.IQHPlanningProjectUsecase,
	documentUsecase qh_usecase.IQHPlanningDocumentUsecase,
	mapper *mapper.QHPlanningMapper,
	validator *validator.QHPlanningValidator,
) *QHPlanningGrpcHandler {
	return &QHPlanningGrpcHandler{
		projectUsecase:  projectUsecase,
		documentUsecase: documentUsecase,
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
	}
	listReq.Page = req.Page
	listReq.Size = req.Size

	results, total, err := h.projectUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningProjectsResponse{
		Data:  h.mapper.ToProjectProtoList(results),
		Total: total,
	}, nil
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
		PlanningProjectID: req.PlanningProjectId,
		DocumentType:      req.DocumentType,
		ValidityStatus:    req.ValidityStatus,
	}
	listReq.Page = req.Page
	listReq.Size = req.Size

	results, total, err := h.documentUsecase.GetList(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListPlanningDocumentsResponse{
		Data:  h.mapper.ToDocumentProtoList(results),
		Total: total,
	}, nil
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
