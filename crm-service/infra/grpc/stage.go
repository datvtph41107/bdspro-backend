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
)

type StageService struct {
	crmpb.UnimplementedStageServiceServer
	UC *usecase.StageUsecase
}

func NewStageService(uc *usecase.StageUsecase) *StageService {
	return &StageService{
		UC: uc,
	}
}

// @Summary Tìm kiếm giai đoạn theo pipeline
// @Description Tìm kiếm giai đoạn theo pipeline
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param pipelineId path int true "Pipeline ID"
// @Param dto query dto.StageSearchDTO true "Stage search dto"
// @Security BearerAuth
// @Router /stage/{pipelineId} [get]
func (s *StageService) Search(ctx context.Context, req *crmpb.StageSearchRequest) (*crmpb.StageListDTO, error) {
	stages, total, err := s.UC.Search(ctx, req.PipelineId, dto.StageSearchDTO{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		StageName: req.StageName,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*crmpb.StageDTO, len(stages))
	for i, stage := range stages {
		result[i] = s.DomainToPb(&stage)
	}
	return &crmpb.StageListDTO{
		Data:          result,
		TotalElements: uint32(total),
	}, nil
}

func (s *StageService) Detail(ctx context.Context, req *crmpb.StageDTO) (*crmpb.StageDTO, error) {
	return nil, nil
}

// @Summary Tạo giai đoạn
// @Description Tạo giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param entity body domain.StageEntity true "Stage entity"
// @Security BearerAuth
// @Router /stage [post]
func (s *StageService) Create(ctx context.Context, req *crmpb.StageDTO) (*crmpb.StageDTO, error) {
	domain := s.PbToDomain(req)

	if req.PipelineId == 0 || req.StageName == "" {
		return nil, &_routes.Except{
			Code:    400,
			Message: "pipelineId and stageName are required",
		}
	}

	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}

	return s.DomainToPb(result), nil
}

// @Summary Cập nhật giai đoạn
// @Description Cập nhật giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param entity body domain.StageEntity true "Stage entity"
// @Security BearerAuth
// @Router /stage/{id} [put]
func (s *StageService) Update(ctx context.Context, req *crmpb.StageDTO) (*crmpb.StageDTO, error) {
	domain := s.PbToDomain(req)

	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}

	return s.DomainToPb(result), nil
}

// @Summary Xóa giai đoạn
// @Description Xóa giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /stage/{id} [delete]
func (s *StageService) Delete(ctx context.Context, req *crmpb.StageDTO) (*crmpb.StageDTO, error) {
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "id is required",
		}
	}

	err := s.UC.Delete(ctx, req.Id)
	return nil, err
}

// --map--

func (s *StageService) PbToDomain(pb *crmpb.StageDTO) *domain.StageEntity {
	return &domain.StageEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		StageName:   pb.StageName,
		PipelineID:  pb.PipelineId,
		OrderNumber: int(pb.OrderNumber),
		RuleID:      &pb.RuleId,
		Step:        enums.EStep(pb.Step),
		Active:      pb.Active,
	}
}

func (s *StageService) DomainToPb(domain *domain.StageEntity) *crmpb.StageDTO {
	return &crmpb.StageDTO{
		Id:          domain.ID,
		StageName:   domain.StageName,
		OrderNumber: int32(domain.OrderNumber),
		PipelineId:  domain.PipelineID,
		RuleId:      *domain.RuleID,
		Step:        uint32(domain.Step),
		Active:      domain.Active,
	}
}
