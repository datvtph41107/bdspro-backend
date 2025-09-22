package admin_handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminProjectHandler struct {
	bdspropb.UnimplementedAdminProjectServiceServer
	AdminProjectUsecase *admin_usecases.AdminProjectUsecase
	ProjectMapper       *mapper.ProjectMapper
}

func NewAdminProjectHandler(
	uc *admin_usecases.AdminProjectUsecase,
	projectMapper *mapper.ProjectMapper,
) *AdminProjectHandler {
	return &AdminProjectHandler{
		AdminProjectUsecase: uc,
		ProjectMapper:       projectMapper,
	}
}

func (h *AdminProjectHandler) GetList(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ProjectListResponse, error) {
	result, total, err := h.AdminProjectUsecase.GetList(ctx,
		&dto.TextSearchRequest{
			Pagable: _dto.Pagable{
				Page: req.Page,
				Size: req.Size,
			},
			Text: req.Text,
		})
	if err != nil {
		return nil, err
	}

	projects := h.ProjectMapper.MapProjectPbList(result)
	return &bdspropb.ProjectListResponse{
		Data:          projects,
		TotalElements: total,
	}, nil
}

func (h *AdminProjectHandler) Create(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	project := h.ProjectMapper.ProjectPbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "Tên dự án không được để trống")
	}
	if req.DeveloperId == 0 {
		return nil, _errors.ReturnError(400, "Developer ID không được để trống")
	}

	project, err := h.AdminProjectUsecase.Create(ctx, project)
	if err != nil {
		return nil, err
	}
	return h.ProjectMapper.MapProjectPb(project), nil
}

func (h *AdminProjectHandler) Update(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	project := h.ProjectMapper.ProjectPbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "Tên dự án không được để trống")
	}
	if req.DeveloperId == 0 {
		return nil, _errors.ReturnError(400, "Developer ID không được để trống")
	}

	project, err := h.AdminProjectUsecase.Update(ctx, req.Id, project)
	if err != nil {
		return nil, err
	}
	return h.ProjectMapper.MapProjectPb(project), nil
}

func (h *AdminProjectHandler) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	_, err := h.AdminProjectUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "success",
	}, nil
}
