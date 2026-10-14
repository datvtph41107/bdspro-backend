package services

import (
	"context"

	_errors "common/errors"
	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/mapper"
	"userv4/modules/bdspro/repository"
)

// AdminProjectService owns /v2/bdspro/admin/project/* (migrated from bdspro-service).
type AdminProjectService interface {
	List(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ProjectListResponse, error)
	Create(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error)
	Update(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error)
	Delete(ctx context.Context, id uint64) error
}

type LocalAdminProjectService struct {
	repo repository.ProjectAdminRepository
}

func NewLocalAdminProjectService(repo repository.ProjectAdminRepository) *LocalAdminProjectService {
	return &LocalAdminProjectService{repo: repo}
}

func (s *LocalAdminProjectService) List(ctx context.Context, req dto.CatalogSearchRequest) (*bdspropb.ProjectListResponse, error) {
	items, total, err := s.repo.List(ctx, req.Pagable(), req.Text)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ProjectListResponse{
		Data:          mapper.ProjectsToPbList(items),
		TotalElements: total,
	}, nil
}

func (s *LocalAdminProjectService) Create(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	if req.GetName() == "" {
		return nil, _errors.ReturnError(400, "Tên dự án không được để trống")
	}
	if req.GetDeveloperId() == 0 {
		return nil, _errors.ReturnError(400, "Developer ID không được để trống")
	}
	entity := mapper.ProjectFromPb(req)
	created, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return mapper.ProjectToPb(created), nil
}

func (s *LocalAdminProjectService) Update(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	if req.GetId() == 0 {
		return nil, _errors.ReturnError(400, "ID không hợp lệ")
	}
	if req.GetName() == "" {
		return nil, _errors.ReturnError(400, "Tên dự án không được để trống")
	}
	if req.GetDeveloperId() == 0 {
		return nil, _errors.ReturnError(400, "Developer ID không được để trống")
	}
	entity := mapper.ProjectFromPb(req)
	updated, err := s.repo.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}
	return mapper.ProjectToPb(updated), nil
}

func (s *LocalAdminProjectService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.ReturnError(400, "ID không hợp lệ")
	}
	return s.repo.Delete(ctx, id)
}
