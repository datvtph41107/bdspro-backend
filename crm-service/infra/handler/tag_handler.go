package handler

import (
	"context"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/mapper"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TagService struct {
	crmpb.UnimplementedTagServiceServer
	TagUsecase *usecase.TagUsecase
	TagMapper  *mapper.TagMapper
}

func NewTagService(tagUsecase *usecase.TagUsecase, tagMapper *mapper.TagMapper) *TagService {
	return &TagService{
		TagUsecase: tagUsecase,
		TagMapper:  tagMapper,
	}
}

// @Summary Tạo tag mới (Admin)
// @Description Tạo tag mới cho hệ thống
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body crmpb.AdminTagRequest true "Thông tin tag"
// @Success 200 {object} crmpb.TagItem "Tag đã tạo"
// @Router /v2/crm/tag/admin [post]
func (h *TagService) AdminCreateTag(ctx context.Context, req *crmpb.AdminTagRequest) (*crmpb.TagItem, error) {
	dtoReq := h.TagMapper.AdminTagRequestToDTO(req)
	result, err := h.TagUsecase.AdminCreateTag(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return h.TagMapper.TagEntityToPb(result), nil
}

// @Summary Cập nhật tag (Admin)
// @Description Cập nhật thông tin tag
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID tag"
// @Param request body crmpb.AdminTagRequest true "Thông tin tag"
// @Success 200 {object} crmpb.TagItem "Tag đã cập nhật"
// @Router /v2/crm/tag/admin/{id} [put]
func (h *TagService) AdminUpdateTag(ctx context.Context, req *crmpb.AdminTagRequest) (*crmpb.TagItem, error) {
	dtoReq := h.TagMapper.AdminTagRequestToDTO(req)
	result, err := h.TagUsecase.AdminUpdateTag(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return h.TagMapper.TagEntityToPb(result), nil
}

// @Summary Xóa tag (Admin)
// @Description Xóa tag khỏi hệ thống
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID tag"
// @Success 200 {object} sharepb.Empty "Xóa thành công"
// @Router /v2/crm/tag/admin/{id} [delete]
func (h *TagService) AdminDeleteTag(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	if err := h.TagUsecase.AdminDeleteTag(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// @Summary Lấy danh sách tag cho contact
// @Description Lấy danh sách tag có thể gán cho contact
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query uint32 false "Trang"
// @Param size query uint32 false "Kích thước trang"
// @Param text query string false "Tìm kiếm"
// @Success 200 {object} crmpb.TagListResponse "Danh sách tag"
// @Router /v2/crm/tag/list [get]
func (h *TagService) ListTagsForContact(ctx context.Context, req *crmpb.TagListRequest) (*crmpb.TagListResponse, error) {
	dtoReq := h.TagMapper.TagListRequestToDTO(req)
	result, total, err := h.TagUsecase.ListTagsForContact(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.TagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: total,
	}, nil
}

// @Summary Lấy danh sách tag (Admin)
// @Description Lấy danh sách tất cả tag trong hệ thống
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query uint32 false "Trang"
// @Param size query uint32 false "Kích thước trang"
// @Param text query string false "Tìm kiếm"
// @Success 200 {object} crmpb.TagListResponse "Danh sách tag"
// @Router /v2/crm/tag/admin/list [get]
func (h *TagService) AdminListTags(ctx context.Context, req *crmpb.AdminTagListRequest) (*crmpb.TagListResponse, error) {
	dtoReq := h.TagMapper.AdminTagListRequestToDTO(req)
	result, total, err := h.TagUsecase.AdminListTags(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.TagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: total,
	}, nil
}

// @Summary Lấy danh sách tag của contact
// @Description Lấy danh sách tag đã gán cho contact
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID contact"
// @Success 200 {object} crmpb.TagListResponse "Danh sách tag"
// @Router /v2/crm/tag/contact/{id} [get]
func (h *TagService) ListTagsByContactID(ctx context.Context, req *sharepb.IdRequest) (*crmpb.TagListResponse, error) {
	result, err := h.TagUsecase.ListTagsByContactID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get tags: %v", err)
	}
	return &crmpb.TagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: uint32(len(result)),
	}, nil
}