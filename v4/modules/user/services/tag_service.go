package services

import (
	"context"

	userpb "pb/types/user"
	mapper "userv4/modules/user/mapper"
	usecases "userv4/modules/user/services/user"
)

// TagService owns tag HTTP surface migrated into modules/user.
type TagService interface {
	AdminCreateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error)
	AdminUpdateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error)
	AdminDeleteTag(ctx context.Context, id uint64) error
	AdminListTags(ctx context.Context, req *userpb.AdminTagListRequest) (*userpb.AdminTagListResponse, error)
	ListTagsByType(ctx context.Context, req *userpb.TagListRequest) (*userpb.TagListResponse, error)
	GetProfileTagsByType(ctx context.Context, req *userpb.GetProfileTagsByTypeRequest) (*userpb.TagListResponse, error)
}

// LocalTagService adapts TagUsecase into modules/user.
type LocalTagService struct {
	usecase *usecases.TagUsecase
	mapper  *mapper.TagMapper
}

func NewLocalTagService(usecase *usecases.TagUsecase, mapper *mapper.TagMapper) *LocalTagService {
	return &LocalTagService{usecase: usecase, mapper: mapper}
}

var _ TagService = (*LocalTagService)(nil)

func (s *LocalTagService) AdminCreateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error) {
	result, err := s.usecase.AdminCreateTag(ctx, s.mapper.AdminTagRequestToDTO(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.TagEntityToPb(result), nil
}

func (s *LocalTagService) AdminUpdateTag(ctx context.Context, req *userpb.AdminTagRequest) (*userpb.TagItem, error) {
	result, err := s.usecase.AdminUpdateTag(ctx, s.mapper.AdminTagRequestToDTO(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.TagEntityToPb(result), nil
}

func (s *LocalTagService) AdminDeleteTag(ctx context.Context, id uint64) error {
	return s.usecase.AdminDeleteTag(ctx, id)
}

func (s *LocalTagService) AdminListTags(ctx context.Context, req *userpb.AdminTagListRequest) (*userpb.AdminTagListResponse, error) {
	result, total, err := s.usecase.AdminListTags(ctx, s.mapper.AdminTagListRequestToDTO(req))
	if err != nil {
		return nil, err
	}
	return &userpb.AdminTagListResponse{Data: s.mapper.TagsToPb(result), Total: total}, nil
}

func (s *LocalTagService) ListTagsByType(ctx context.Context, req *userpb.TagListRequest) (*userpb.TagListResponse, error) {
	result, total, err := s.usecase.ListTagsByType(ctx, s.mapper.TagListRequestToDTO(req))
	if err != nil {
		return nil, err
	}
	return &userpb.TagListResponse{Data: s.mapper.TagsToPb(result), Total: total}, nil
}

func (s *LocalTagService) GetProfileTagsByType(ctx context.Context, req *userpb.GetProfileTagsByTypeRequest) (*userpb.TagListResponse, error) {
	profileID, tagType := s.mapper.GetProfileTagsByTypeRequestToDTO(req)
	result, total, err := s.usecase.GetProfileTagsByType(ctx, profileID, tagType)
	if err != nil {
		return nil, err
	}
	return &userpb.TagListResponse{Data: s.mapper.TagsToPb(result), Total: total}, nil
}
