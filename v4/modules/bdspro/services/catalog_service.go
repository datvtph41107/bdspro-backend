package services

import (
	"context"

	_utils "common/utils"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/mapper"
	"mdn4/modules/bdspro/repository"
	bdspropb "pb/types/bdspro"
)

// CatalogService serves public list catalog routes migrated from bdspro-service ListUsecase.
type CatalogService interface {
	ListProjects(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
	ListProjectBuildItems(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
	// ListPropertyTypes(ctx context.Context, req dto.PropertyTypeListRequest) (*PropertyTypeListOutcome, error)
}

type LocalCatalogService struct {
	repo repository.CatalogRepository
	sync *_utils.SyncUtil
}

func NewLocalCatalogService(repo repository.CatalogRepository, sync *_utils.SyncUtil) CatalogService {
	return &LocalCatalogService{repo: repo, sync: sync}
}

func (s *LocalCatalogService) ListProjects(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error) {
	items, total, err := s.repo.SearchProjects(ctx, req.Pagable(), req.Text)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListItemResponse{
		Data:  mapper.ProjectItemsToPbList(items),
		Total: uint32(total),
	}, nil
}

func (s *LocalCatalogService) ListProjectBuildItems(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error) {
	items, total, err := s.repo.SearchProjectBuilds(ctx, req.Pagable(), req.Text)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListItemResponse{
		Data:  mapper.ProjectBuildItemsToPbList(items),
		Total: uint32(total),
	}, nil
}
