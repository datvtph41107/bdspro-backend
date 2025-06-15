package grpc_service

import (
	_dto "common/dto"
	_models "common/models"
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"

	"google.golang.org/grpc/codes"
)

type PipelineService struct {
	crmpb.UnimplementedPipelineServiceServer
	UC *usecase.PipelineUsecase
}

func NewPipelineService(pipelineUsecase *usecase.PipelineUsecase) *PipelineService {
	return &PipelineService{UC: pipelineUsecase}
}

func (s *PipelineService) Search(ctx context.Context, req *crmpb.PipelineSearchRequest) (*crmpb.PipelineListDTO, error) {
	pipelines, total, err := s.UC.Search(ctx,
		req.OwnerId,
		enums.EOwnerType(req.OwnerType),
		dto.PipelineSearchDTO{
			Page: _dto.Pagable{
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
		result[i] = s.DomainToPb(&pipeline)
	}

	return &crmpb.PipelineListDTO{
		TotalElements: uint32(total),
		Data:          result,
	}, nil
}
func (s *PipelineService) Detail(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	return nil, nil
}
func (s *PipelineService) Create(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	domain := s.PbToDomain(req)

	if req.OwnerId == 0 || req.OwnerType == 0 {
		return nil, &_routes.Except{
			Code:    int(codes.InvalidArgument),
			Message: "ownerId and ownerType are required",
		}
	}

	result, err := s.UC.Create(ctx, 0, domain)
	if err != nil {
		return nil, err
	}

	return s.DomainToPb(result), nil
}
func (s *PipelineService) Update(ctx context.Context, req *crmpb.PipelineDTO) (*crmpb.PipelineDTO, error) {
	domain := s.PbToDomain(req)

	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}

	return s.DomainToPb(result), nil
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

// --map--

func (s *PipelineService) PbToDTO(pb *crmpb.PipelineDTO) *dto.PipelineDTO {
	return &dto.PipelineDTO{
		ID:           pb.Id,
		PipelineName: pb.PipelineName,
		Color:        pb.Color,
		Active:       pb.Active,
		OwnerID:      pb.OwnerId,
		OwnerType:    enums.EOwnerType(pb.OwnerType),
	}
}

func (s *PipelineService) PbToDomain(pb *crmpb.PipelineDTO) *domain.PipelineEntity {
	return &domain.PipelineEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		PipelineName: pb.PipelineName,
		Color:        pb.Color,
		Active:       pb.Active,
		OwnerID:      pb.OwnerId,
		OwnerType:    enums.EOwnerType(pb.OwnerType),
	}
}

func (s *PipelineService) DomainToPb(domain *domain.PipelineEntity) *crmpb.PipelineDTO {
	return &crmpb.PipelineDTO{
		Id:           domain.ID,
		PipelineName: domain.PipelineName,
		Color:        domain.Color,
		Active:       domain.Active,
		OwnerId:      domain.OwnerID,
		OwnerType:    uint32(domain.OwnerType),
	}
}
