package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_dto "common/dto"
	"context"

	pb_bdspro "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type ProjectBuildService struct {
	pb_bdspro.UnimplementedProjectBuildServiceServer
	ProjectBuildUsecase *usecases.ProjectBuildUsecase
	Mapper              *mapper.ProjectBuildMapper
}

func NewProjectBuildService(projectBuildUsecase *usecases.ProjectBuildUsecase,
	mapper *mapper.ProjectBuildMapper,
) *ProjectBuildService {
	return &ProjectBuildService{
		ProjectBuildUsecase: projectBuildUsecase,
		Mapper:              mapper,
	}
}

// @Summary
// @Description Lấy danh sách tòa nhà
// @Tags ProjectBuild
// @Accept json
// @Produce json
// @Param text query string false "Text"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} pb_bdspro.ListItemResponse
// @Router /v2/bdspro/v2/project-build/items [get]
func (s *ProjectBuildService) GetProjectBuild(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	dto := &dto.TextSearchRequest{
		Text: req.Text,
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}
	items, total, err := s.ProjectBuildUsecase.Search(ctx, dto)
	if err != nil {
		return nil, err
	}
	results := s.Mapper.ProjectBuildItemToPbList(items)
	return &pb_bdspro.ListItemResponse{
		Data:  results,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy thông tin chi tiết của tòa nhà theo ID
// @Description Lấy thông tin chi tiết của tòa nhà theo ID
// @Tags ProjectBuild
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} pb_bdspro.ProjectBuild
// @Router /v2/bdspro/v2/project-build/detail/{id} [get]
func (s *ProjectBuildService) GetProjectBuildById(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.ProjectBuild, error) {
	id := req.Id
	projectBuild, err := s.ProjectBuildUsecase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.Mapper.ProjectBuildToPb(projectBuild), nil
}

// @Summary Tạo mới tòa nhà
// @Description Tạo mới tòa nhà
// @Tags ProjectBuild
// @Accept json
// @Produce json
// @Param projectBuild body pb_bdspro.ProjectBuild true "Thông tin tòa nhà"
// @Success 200 {object} pb_bdspro.ProjectBuild
// @Router /v2/bdspro/v2/project-build [post]
func (s *ProjectBuildService) CreateProjectBuild(ctx context.Context, req *pb_bdspro.ProjectBuild) (*pb_bdspro.ProjectBuild, error) {
	domain := s.Mapper.ProjectBuildToDomain(req)
	projectBuild, err := s.ProjectBuildUsecase.Create(ctx, domain)
	if err != nil {
		return nil, err
	}
	return s.Mapper.ProjectBuildToPb(projectBuild), nil
}

// @Summary Cập nhật thông tin tòa nhà
// @Description Cập nhật thông tin tòa nhà
// @Tags ProjectBuild
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param projectBuild body pb_bdspro.ProjectBuild true "Thông tin tòa nhà"
// @Success 200 {object} pb_bdspro.ProjectBuild
// @Router /v2/bdspro/v2/project-build/{id} [put]
func (s *ProjectBuildService) UpdateProjectBuild(ctx context.Context, req *pb_bdspro.ProjectBuild) (*pb_bdspro.ProjectBuild, error) {
	domain := s.Mapper.ProjectBuildToDomain(req)
	projectBuild, err := s.ProjectBuildUsecase.Update(ctx, domain)
	if err != nil {
		return nil, err
	}
	return s.Mapper.ProjectBuildToPb(projectBuild), nil
}

// @Summary Xóa tòa nhà
// @Description Xóa tòa nhà
// @Tags ProjectBuild
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} sharepb.Empty
// @Router /v2/bdspro/v2/project-build/{id} [delete]
func (s *ProjectBuildService) DeleteProjectBuild(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	id := req.Id
	err := s.ProjectBuildUsecase.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}
