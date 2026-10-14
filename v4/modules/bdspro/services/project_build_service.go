package services

import (
	"context"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/mapper"
	"userv4/modules/bdspro/repository"
)

// ProjectBuildService owns /v2/bdspro/v2/project-build/* (migrated from bdspro-service).
type ProjectBuildService interface {
	ListItems(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
	ListMine(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error)
	GetByID(ctx context.Context, id uint64) (*bdspropb.ProjectBuild, error)
	Create(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error)
	Update(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error)
	Delete(ctx context.Context, id uint64) error
}

type LocalProjectBuildService struct {
	repo repository.ProjectBuildRepository
}

func NewLocalProjectBuildService(repo repository.ProjectBuildRepository) *LocalProjectBuildService {
	return &LocalProjectBuildService{repo: repo}
}

func (s *LocalProjectBuildService) ListItems(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error) {
	items, total, err := s.repo.Search(ctx, req.Pagable(), req.Text)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListItemResponse{
		Data:  mapper.ProjectBuildItemsToPbList(items),
		Total: uint32(total),
	}, nil
}

func (s *LocalProjectBuildService) ListMine(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ListItemResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	items, total, err := s.repo.SearchByCreatedBy(ctx, profileID, req.Pagable(), req.Text)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListItemResponse{
		Data:  mapper.ProjectBuildItemsToPbList(items),
		Total: uint32(total),
	}, nil
}

func (s *LocalProjectBuildService) GetByID(ctx context.Context, id uint64) (*bdspropb.ProjectBuild, error) {
	build, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ProjectBuildToPb(build), nil
}

func (s *LocalProjectBuildService) Create(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error) {
	if req.GetName() == "" {
		return nil, _errors.ReturnError(400, "Tên tòa nhà không được để trống")
	}
	entity := mapper.ProjectBuildFromPb(req)
	created, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapper.ProjectBuildToPb(created), nil
}

func (s *LocalProjectBuildService) Update(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error) {
	if req.GetId() == 0 {
		return nil, _errors.ReturnError(400, "ID không hợp lệ")
	}
	if req.GetName() == "" {
		return nil, _errors.ReturnError(400, "Tên tòa nhà không được để trống")
	}
	entity := mapper.ProjectBuildFromPb(req)
	updated, err := s.repo.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapper.ProjectBuildToPb(updated), nil
}

func (s *LocalProjectBuildService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.ReturnError(400, "ID không hợp lệ")
	}
	return s.repo.Delete(ctx, id)
}
