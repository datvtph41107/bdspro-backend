package handler

import (
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type TagHandler struct {
	userpb.UnimplementedTagServiceServer
	TagUsecase *usecases.TagUsecase
	TagMapper  *mapper.TagMapper
}

func NewTagHandler(tagUsecase *usecases.TagUsecase, tagMapper *mapper.TagMapper) *TagHandler {
	return &TagHandler{
		TagUsecase: tagUsecase,
		TagMapper:  tagMapper,
	}
}

func (h *TagHandler) AdminCreateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error) {
	dtoReq := h.TagMapper.AdminTagRequestToDTO(req)
	result, err := h.TagUsecase.AdminCreateTag(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return h.TagMapper.TagEntityToPb(result), nil
}

func (h *TagHandler) AdminUpdateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error) {
	dtoReq := h.TagMapper.AdminTagRequestToDTO(req)
	result, err := h.TagUsecase.AdminUpdateTag(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return h.TagMapper.TagEntityToPb(result), nil
}

func (h *TagHandler) AdminDeleteTag(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	if err := h.TagUsecase.AdminDeleteTag(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (h *TagHandler) ListTagsByType(ctx context.Context, req *userpb.TagListRequest) (*userpb.TagListResponse, error) {
	dtoReq := h.TagMapper.TagListRequestToDTO(req)
	result, total, err := h.TagUsecase.ListTagsByType(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return &userpb.TagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: total,
	}, nil
}

func (h *TagHandler) AdminListTags(ctx context.Context, req *userpb.AdminTagListRequest) (*userpb.AdminTagListResponse, error) {
	dtoReq := h.TagMapper.AdminTagListRequestToDTO(req)
	result, total, err := h.TagUsecase.AdminListTags(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return &userpb.AdminTagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: total,
	}, nil
}

// @Summary Lấy danh sách tags của profile theo type (Public API)
// @Description API public không cần token để lấy danh sách tags có type là mainArea của profile
// @Tags Tag
// @Accept json
// @Produce json
// @Param profileId path int true "ID của Profile"
// @Param tagType query int false "Loại tag (mặc định là 20 - mainArea)"
// @Success 200 {object} userpb.TagListResponse
// @Router /v2/pub/user/profile/tags/{profileId} [get]
func (h *TagHandler) GetProfileTagsByType(ctx context.Context, req *userpb.GetProfileTagsByTypeRequest) (*userpb.TagListResponse, error) {
	profileID, tagType := h.TagMapper.GetProfileTagsByTypeRequestToDTO(req)
	result, total, err := h.TagUsecase.GetProfileTagsByType(ctx, profileID, tagType)
	if err != nil {
		return nil, err
	}
	return &userpb.TagListResponse{
		Data:  h.TagMapper.TagsToPb(result),
		Total: total,
	}, nil
}
