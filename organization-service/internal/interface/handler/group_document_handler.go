package handler

import (
	"context"

	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"

	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupDocumentHandler struct {
	grouppb.UnimplementedGroupDocumentServiceServer
	GroupDocumentUsecase     usecase.GroupDocumentUsecase
	GroupDocumentTransformer transformer.GroupDocumentTransformer
	GroupDocumentValidator   validator.GroupDocumentValidator
}

func NewGroupDocumentHandler(
	groupDocumentUsecase usecase.GroupDocumentUsecase,
	groupDocumentTransformer transformer.GroupDocumentTransformer,
	groupDocumentValidator validator.GroupDocumentValidator,
) *GroupDocumentHandler {
	return &GroupDocumentHandler{
		GroupDocumentUsecase:     groupDocumentUsecase,
		GroupDocumentTransformer: groupDocumentTransformer,
		GroupDocumentValidator:   groupDocumentValidator,
	}
}

func (h *GroupDocumentHandler) CreateGroupDocument(ctx context.Context, req *grouppb.CreateGroupDocumentRequest) (*grouppb.CreateGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateCreateGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document := h.GroupDocumentTransformer.CreateGroupDocumentRequestToEntity(req)
	document.CreatedBy = utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	document, err := h.GroupDocumentUsecase.CreateGroupDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToCreateGroupDocumentResponse(document), nil
}

func (h *GroupDocumentHandler) UpdateGroupDocument(ctx context.Context, req *grouppb.UpdateGroupDocumentRequest) (*grouppb.UpdateGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateUpdateGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document := h.GroupDocumentTransformer.UpdateGroupDocumentRequestToEntity(req)
	document.UpdatedBy = utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	document, err := h.GroupDocumentUsecase.UpdateGroupDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToUpdateGroupDocumentResponse(document), nil
}

func (h *GroupDocumentHandler) DeleteGroupDocument(ctx context.Context, req *grouppb.DeleteGroupDocumentRequest) (*grouppb.DeleteGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateDeleteGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDocumentUsecase.DeleteGroupDocument(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &grouppb.DeleteGroupDocumentResponse{Id: req.Id}, nil
}

func (h *GroupDocumentHandler) GetGroupDocument(ctx context.Context, req *grouppb.GetGroupDocumentRequest) (*grouppb.GetGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateGetGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document, err := h.GroupDocumentUsecase.GetGroupDocumentByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToGetGroupDocumentResponse(document), nil
}

func (h *GroupDocumentHandler) GetGroupDocuments(ctx context.Context, req *grouppb.GetGroupDocumentsRequest) (*grouppb.GetGroupDocumentsResponse, error) {
	if err := h.GroupDocumentValidator.ValidateGetGroupDocumentsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	documents, total, err := h.GroupDocumentUsecase.GetGroupDocumentsByGroupID(ctx, req.GroupId, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToGetGroupDocumentsResponse(documents, total), nil
} 