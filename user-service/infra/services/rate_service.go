package services

import (
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type RateService struct {
	userpb.UnimplementedRateServiceServer
	rateUsecase *usecases.RateUsecase
	rateMapper  *mapper.RateMapper
}

func NewGrpcRateService(rateUsecase *usecases.RateUsecase, rateMapper *mapper.RateMapper) *RateService {
	return &RateService{
		rateUsecase: rateUsecase,
		rateMapper:  rateMapper,
	}
}

// @Summary Tạo đánh giá
// @Description Tạo đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rate body userpb.Rate true "Rate"
// @Success 200 {object} userpb.Rate
// @Router /v2/user/rate [post]
func (s *RateService) CreateRate(ctx context.Context, req *userpb.Rate) (*userpb.Rate, error) {
	rate := s.rateMapper.PbToEntityRate(req)

	rate, err := s.rateUsecase.CreateRate(ctx, rate)
	if err != nil {
		return nil, err
	}
	result := s.rateMapper.EntityToPbRate(rate)

	return result, nil
}

// @Summary Lấy danh sách đánh giá
// @Description Lấy danh sách đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Param targetId path uint64 true "Target ID"
// @Param targetType path uint32 true "Target Type"
// @Param dto query dto.RateSearchDTO true "Parent ID"
// @Success 200 {object} userpb.RateList
// @Router /v2/user/rate/list/{targetId}/{targetType} [get]
func (s *RateService) GetRateList(ctx context.Context, req *userpb.GetRateListRequest) (*userpb.RateList, error) {
	dto := s.rateMapper.PbToEntityRateSearch(req)
	rates, total, err := s.rateUsecase.GetRateList(ctx, req.TargetId, req.TargetType, dto)
	rateList := make([]*userpb.Rate, 0)
	for _, rate := range rates {
		r := s.rateMapper.EntityToPbRateItem(&rate)
		rateList = append(rateList, r)
	}
	if err != nil {
		return nil, err
	}
	return &userpb.RateList{
		Data:  rateList,
		Total: int32(total),
	}, nil
}

// @Summary Cập nhật đánh giá
// @Description Cập nhật đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rate body userpb.Rate true "Rate"
// @Param id path uint64 true "ID"
// @Success 200 {object} userpb.Rate
// @Router /v2/user/rate/{id} [put]
func (s *RateService) UpdateRate(ctx context.Context, req *userpb.Rate) (*userpb.Rate, error) {
	rate := s.rateMapper.PbToEntityRate(req)

	rate, err := s.rateUsecase.UpdateRate(ctx, rate)
	if err != nil {
		return nil, err
	}
	result := s.rateMapper.EntityToPbRate(rate)

	return result, nil
}

// @Summary Xóa đánh giá
// @Description Xóa đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} sharepb.Empty
// @Router /v2/user/rate/{id} [delete]
func (s *RateService) DeleteRate(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := s.rateUsecase.DeleteRate(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// @Summary Lấy danh sách đánh giá
// @Description Lấy danh sách đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rateId path uint64 true "Rate ID"
// @Param page query uint32 false "Page"
// @Param size query uint32 false "Size"
// @Success 200 {object} userpb.RateList
// @Router /v2/user/rate/history/{rateId} [get]
func (s *RateService) GetHistoryRate(ctx context.Context, req *userpb.GetHistoryRateRequest) (*userpb.RateList, error) {
	rates, total, err := s.rateUsecase.GetHistoryRate(ctx, req.RateId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	rateList := make([]*userpb.Rate, 0)
	for _, rate := range rates {
		r := s.rateMapper.EntityToPbRateItem(&rate)
		rateList = append(rateList, r)
	}
	return &userpb.RateList{
		Data:  rateList,
		Total: int32(total),
	}, nil
}

// @Summary Cập nhật ẩn đánh giá
// @Description Cập nhật ẩn đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Param hidden body userpb.UpdateHiddenRequest true "Hidden"
// @Success 200 {object} sharepb.Empty
// @Router /v2/user/rate/hidden/{id} [put]
func (s *RateService) UpdateHidden(ctx context.Context, req *userpb.UpdateHiddenRequest) (*sharepb.Empty, error) {
	err := s.rateUsecase.UpdateHidden(ctx, req.Id, req.Hidden)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// @Summary Lấy thông tin đánh giá
// @Description Lấy thông tin đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} userpb.RateStats
// @Router /v2/user/rate/stats/{id} [get]
func (s *RateService) GetRateStats(ctx context.Context, req *sharepb.IdRequest) (*userpb.RateStats, error) {
	stats, err := s.rateUsecase.GetRateStats(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return s.rateMapper.EntityToPbRateStats(&stats), nil
}
