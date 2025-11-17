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
