package handler

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"crm/internal"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
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
// @Router /stage/list/{pipelineId} [get]
func (s *StageService) Search(ctx context.Context, req *crmpb.StageSearchRequest) (*crmpb.StageListDTO, error) {
	stages, total, err := s.UC.Search(ctx, req.PipelineId, dto.StageSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		StageName: req.StageName,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*crmpb.StageDTO, len(stages))
	for i, stage := range stages {
		result[i] = mapper.StageDomainToPb(&stage)
	}
	return &crmpb.StageListDTO{
		Data:  result,
		Total: uint32(total),
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
	domain := mapper.StagePbToDomain(req)

	if req.PipelineId == 0 || req.StageName == "" {
		return nil, _errors.ReturnError(service.StageCreateFieldsRequired)
	}

	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}

	return mapper.StageDomainToPb(result), nil
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
	domain := mapper.StagePbToDomain(req)

	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}

	return mapper.StageDomainToPb(result), nil
}

// @Summary Xóa giai đoạn
// @Description Xóa giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /stage/{id} [delete]
func (s *StageService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(service.IDRequired)
	}

	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Delete stage success",
	}, nil
}
