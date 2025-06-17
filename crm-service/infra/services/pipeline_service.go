package grpc_service

import (
	_dto "common/dto"
	_routes "common/routes"
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

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
// @Param ownerId path int true "Owner ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Param dto query dto.PipelineSearchDTO true "Pipeline search dto"
// @Security BearerAuth
// @Router /pipeline/list [get]
func (s *PipelineService) Search(ctx context.Context, req *crmpb.PipelineSearchRequest) (*crmpb.PipelineListDTO, error) {
	pipelines, total, err := s.UC.Search(ctx,
		dto.PipelineSearchDTO{
			Pagable: _dto.Pagable{
				Page: int(req.Page),
				Size: int(req.Size),
			},
			PipelineName: req.PipelineName,
			Active:       req.Active,
			Colors:       req.Colors,
		})
	if err != nil {
		return nil, err
	}

	result := make([]*crmpb.PipelineDTO, len(pipelines))
	for i, pipeline := range pipelines {
		result[i] = mapper.PipelineDomainToPb(&pipeline)
	}

	return &crmpb.PipelineListDTO{
		TotalElements: uint32(total),
		Data:          result,
	}, nil
}
func (s *PipelineService) Detail(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	return nil, nil
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

	result, err := s.UC.Create(ctx, 0, domain)
	if err != nil {
		return nil, err
	}

	return mapper.PipelineDomainToPb(result), nil
}
func (s *PipelineService) Update(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	domain := mapper.PipelinePbToDomain(req)

	result, err := s.UC.Update(ctx, req.Id, domain)
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
