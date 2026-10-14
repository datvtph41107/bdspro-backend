package services

import (
	"context"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/mapper"
	"userv4/modules/bdspro/repository"
)

// CatalogService serves public list catalog routes migrated from bdspro-service ListUsecase.
type CatalogService interface {
	ListProjects(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
	ListProjectBuildItems(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
}

type LocalCatalogService struct {
	repo repository.CatalogRepository
}

func NewLocalCatalogService(repo repository.CatalogRepository) *LocalCatalogService {
	return &LocalCatalogService{repo: repo}
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
