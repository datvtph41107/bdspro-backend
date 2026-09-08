package handler

import (
	_dto "common/domain/dto"
	_routes "common/routes"
	"context"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
)

type PipelineService struct {
	crmpb.UnimplementedPipelineServiceServer
	UC *usecase.PipelineUsecase
}

func NewPipelineService(pipelineUsecase *usecase.PipelineUsecase) *PipelineService {
	return &PipelineService{UC: pipelineUsecase}
}

// @Summary Tìm kiếm quy trình
// @Description Tìm kiếm quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param dto query dto.PipelineSearchDTO true "Pipeline search dto"
// @Security BearerAuth
// @Router /pipeline/list [get]
func (s *PipelineService) Search(ctx context.Context, req *crmpb.PipelineSearchRequest) (*crmpb.PipelineListDTO, error) {
	pipelines, total, err := s.UC.Search(ctx,
		dto.PipelineSearchDTO{
			Pagable: _dto.Pagable{
				Page: uint32(req.Page),
				Size: uint32(req.Size),
			},
			PipelineName: req.PipelineName,
			Active:       req.Active,
			Colors:       req.Colors,
		}, false)
	if err != nil {
		return nil, err
	}

	result := make([]*crmpb.PipelineDTO, len(pipelines))
	for i, pipeline := range pipelines {
		result[i] = mapper.PipelineDomainToPb(&pipeline)
	}

	return &crmpb.PipelineListDTO{
		Total: uint32(total),
		Data:  result,
	}, nil
}

// @Summary Lấy chi tiết quy trình
// @Description Lấy chi tiết quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /pipeline/detail/{id} [get]
func (s *PipelineService) Detail(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	pipeline, err := s.UC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return mapper.PipelineDomainToPb(pipeline), nil
}

// @Summary Tạo quy trình
// @Description Tạo quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param entity body crmpb.PipelineDTO true "Pipeline entity"
// @Security BearerAuth
// @Router /pipeline [post]
func (s *PipelineService) Create(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	domain := mapper.PipelinePbToDomain(req)

	// if req.OwnerId == 0 || req.OwnerType == 0 {
	// 	return nil, &_routes.Except{
	// 		Code:    int(codes.InvalidArgument),
	// 		Message: "ownerId and ownerType are required",
	// 	}
	// }

	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}

	return mapper.PipelineDomainToPb(result), nil
}

// @Summary Cập nhật quy trình
// @Description Cập nhật quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param entity body crmpb.PipelineSaveRequest true "Pipeline entity"
// @Security BearerAuth
// @Param id path int true "ID"
// @Router /pipeline/{id} [put]
func (s *PipelineService) Update(ctx context.Context, req *crmpb.PipelineSaveRequest) (*crmpb.PipelineDTO, error) {
	dto := mapper.PipelineSavePbToDomain(req)

	result, err := s.UC.Update(ctx, req.Id, dto)
	if err != nil {
		return nil, err
	}

	return mapper.PipelineDomainToPb(result), nil
}

// @Summary Xóa quy trình
// @Description Xóa quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /pipeline/{id} [delete]
func (s *PipelineService) Delete(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    int(codes.InvalidArgument),
			Message: "id is required",
		}
	}

	err := s.UC.Delete(ctx, req.Id)
	return nil, err
}

// @Summary Lấy quy trình mặc định
// @Description Lấy quy trình mặc định
// @Tags Quy trình
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /pipeline/default [get]
func (s *PipelineService) GetDefault(ctx context.Context, req *sharepb.Empty) (*crmpb.PipelineDTO, error) {
	pipeline, stage, err := s.UC.GetDefault(ctx)
	if err != nil {
		return nil, err
	}

	result := mapper.PipelineDomainToPb(pipeline)
	if stage != nil {
		result.DefaultStage = mapper.StageDomainToPb(stage)
	}
	return result, nil
}

// @Summary Lấy quy trình và stage
// @Description Lấy quy trình và stage
// @Tags Quy trình
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /pipeline/with-stages [get]
func (s *PipelineService) GetWithStages(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineListDTO, error) {
	dto := dto.PipelineSearchDTO{
		Pagable: _dto.Pagable{
			Page: 1,
			Size: 1,
		},
	}
	pipelines, total, err := s.UC.Search(ctx, dto, true)
	if err != nil {
		return nil, err
	}

	result := make([]*crmpb.PipelineDTO, len(pipelines))
	for i, pipeline := range pipelines {
		result[i] = mapper.PipelineDomainToPb(&pipeline)
	}

	return &crmpb.PipelineListDTO{
		Total: uint32(total),
		Data:  result,
	}, nil
}