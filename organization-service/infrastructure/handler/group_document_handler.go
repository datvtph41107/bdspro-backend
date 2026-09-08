package handler

import (
	"context"

	organizationpb "pb/types/organization"

	"organization/env"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"

	"organization/pkg/utils"
)

type GroupDocumentHandler struct {
	organizationpb.UnimplementedGroupDocumentServiceServer
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

// @Summary Tạo tài liệu
// @Description Tạo tài liệu
// @Tags Tài liệu
// @Accept json
// @Produce json
// @Param group_document body organizationpb.CreateGroupDocumentRequest true "Thông tin tài liệu"
// @Param id path int true "ID của tài liệu"
// @Security BearerAuth
// @Router /group/document [post]
func (h *GroupDocumentHandler) CreateGroupDocument(ctx context.Context, req *organizationpb.CreateGroupDocumentRequest) (*organizationpb.CreateGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateCreateGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document := h.GroupDocumentTransformer.CreateGroupDocumentRequestToEntity(req)
	document.CreatedBy = utils.GetUserID(ctx, env.USER_CONTEXT)

	document, err := h.GroupDocumentUsecase.CreateGroupDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToCreateGroupDocumentResponse(document), nil
}

// @Summary Cập nhật tài liệu
// @Description Cập nhật tài liệu
// @Tags Tài liệu
// @Accept json
// @Produce json
// @Param group_document body organizationpb.UpdateGroupDocumentRequest true "Thông tin tài liệu"
// @Param id path int true "ID của tài liệu"
// @Security BearerAuth
// @Router /group/document/{id} [put]
func (h *GroupDocumentHandler) UpdateGroupDocument(ctx context.Context, req *organizationpb.UpdateGroupDocumentRequest) (*organizationpb.UpdateGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateUpdateGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document := h.GroupDocumentTransformer.UpdateGroupDocumentRequestToEntity(req)
	document.UpdatedBy = utils.GetUserID(ctx, env.USER_CONTEXT)

	document, err := h.GroupDocumentUsecase.UpdateGroupDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToUpdateGroupDocumentResponse(document), nil
}

// @Summary Xóa tài liệu
// @Description Xóa tài liệu
// @Tags Tài liệu
// @Accept json
// @Produce json
// @Param group_document body organizationpb.DeleteGroupDocumentRequest true "Thông tin tài liệu"
// @Param id path int true "ID của tài liệu"
// @Security BearerAuth
// @Router /group/document/{id} [delete]
func (h *GroupDocumentHandler) DeleteGroupDocument(ctx context.Context, req *organizationpb.DeleteGroupDocumentRequest) (*organizationpb.DeleteGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateDeleteGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDocumentUsecase.DeleteGroupDocument(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupDocumentResponse{Id: req.Id}, nil
}

// @Summary Lấy tài liệu
// @Description Lấy tài liệu
// @Tags Tài liệu
// @Accept json
// @Produce json
// @Param id path int true "ID của tài liệu"
// @Security BearerAuth
// @Router /group/document/detail/{id} [get]
func (h *GroupDocumentHandler) GetGroupDocument(ctx context.Context, req *organizationpb.GetGroupDocumentRequest) (*organizationpb.GetGroupDocumentResponse, error) {
	if err := h.GroupDocumentValidator.ValidateGetGroupDocumentRequest(req); err != nil {
		return nil, err
	}

	document, err := h.GroupDocumentUsecase.GetGroupDocumentByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.GroupDocumentTransformer.EntityToGetGroupDocumentResponse(document), nil
}

// @Summary Lấy danh sách tài liệu
// @Description Lấy danh sách tài liệu
// @Tags Tài liệu
// @Accept json
// @Produce json
// @Param group_document query organizationpb.GetGroupDocumentsRequest true "Thông tin tài liệu"
// @Param groupId path int true "ID của nhóm"
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng mục trên mỗi trang"
// @Security BearerAuth
// @Router /group/document/list/{groupId} [get]
func (h *GroupDocumentHandler) GetGroupDocuments(ctx context.Context, req *organizationpb.GetGroupDocumentsRequest) (*organizationpb.GetGroupDocumentsResponse, error) {
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
