package handler

import (
	"context"
	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	_utils "common/utils"

	_dto "common/domain/dto"
	"crm/infra/mapper"
)

type RateHandler struct {
	crmpb.UnimplementedRateServiceServer
	rateUsecase *usecase.RateUsecase
	userClient  *client.UserClient
	rateMapper  *mapper.RateMapper
}

// @bind: crm/infra/handler.RateHandler
func NewRateHandler(rateUsecase *usecase.RateUsecase, userClient *client.UserClient, rateMapper *mapper.RateMapper) *RateHandler {
	return &RateHandler{
		rateUsecase: rateUsecase,
		userClient:  userClient,
		rateMapper:  rateMapper,
	}
}

// CreateRate tạo đánh giá mới
// @Summary Tạo đánh giá mới
// @Description Tạo đánh giá mới cho một đối tượng
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param registry path string true "Registry"
// @Param request body dto.RateCreateRequest true "Thông tin đánh giá"
// @Success 200 {object} crmpb.CreateRateResponse
// @Router /rate/{registry} [post]
func (h *RateHandler) CreateRate(ctx context.Context, req *crmpb.CreateRateRequest) (*crmpb.RateResponse, error) {
	// Convert protobuf request to DTO
	rateReq := &dto.RateCreateRequest{
		Anonymous: req.Anonymous,
		OwnerID:   req.OwnerId,
		OwnerOf:   req.OwnerOf,
		Score:     uint8(req.Score),
		Comment:   req.Comment,
		ParentID:  req.ParentId,
		Attachs:   req.Attachs,
	}

	// Convert to domain
	rate := convertRateCreateRequestToDomain(rateReq)

	// Create rate
	createdRate, err := h.rateUsecase.CreateRate(ctx, req.Registry, rate)
	if err != nil {
		return nil, err
	}

	return h.convertRateToResponse(ctx, createdRate), nil
}

// GetRateList lấy danh sách đánh giá
// @Summary Lấy danh sách đánh giá
// @Description Lấy danh sách đánh giá theo target ID và type
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ownerId path uint64 true "ID đối tượng"
// @Param registry path string true "Registry"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} crmpb.GetRateListResponse
// @Router /rate/{registry}/list/{ownerId} [get]
func (h *RateHandler) GetRateList(ctx context.Context, req *crmpb.GetRateListRequest) (*crmpb.GetRateListResponse, error) {
	// Convert request to DTO
	searchReq := &dto.RateSearchDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		OwnerId:  req.OwnerId,
		ParentId: req.ParentId,
		Registry: req.Registry,
	}

	// Get rate list
	rates, total, err := h.rateUsecase.GetRateList(ctx, req.OwnerId, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	rateResponses := h.rateMapper.MapRateToPbs(ctx, rates)

	// Map user profile
	h.rateMapper.MapUserToPb(ctx, rateResponses)

	return &crmpb.GetRateListResponse{
		Data:  rateResponses,
		Total: total,
	}, nil
}

// GetRateStats lấy thống kê đánh giá
// @Summary Lấy thống kê đánh giá
// @Description Lấy thống kê đánh giá theo target ID và type
// @Tags Rate
// @Accept json
// @Produce json
// @Param ownerId path uint64 true "ID đối tượng"
// @Param registry path string true "Registry"
// @Success 200 {object} crmpb.RateStatsResponse
// @Router /rate/{registry}/stats/{ownerId} [get]
func (h *RateHandler) GetRateStats(ctx context.Context, req *crmpb.GetRateStatsRequest) (*crmpb.RateStatsResponse, error) {
	// Get rate stats
	stats, err := h.rateUsecase.GetRateStats(ctx, req.OwnerId, req.Registry)
	if err != nil {
		return nil, err
	}

	return &crmpb.RateStatsResponse{
		TotalReviews: stats.TotalReviews,
		AverageScore: stats.AverageScore,
		Star1Count:   stats.Star1Count,
		Star2Count:   stats.Star2Count,
		Star3Count:   stats.Star3Count,
		Star4Count:   stats.Star4Count,
		Star5Count:   stats.Star5Count,
	}, nil
}

// UpdateRate cập nhật đánh giá
// @Summary Cập nhật đánh giá
// @Description Cập nhật đánh giá của người dùng
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID đánh giá"
// @Param request body dto.RateUpdateRequest true "Thông tin cập nhật"
// @Success 200 {object} crmpb.UpdateRateResponse
// @Router /rate/{id} [put]
func (h *RateHandler) UpdateRate(ctx context.Context, req *crmpb.UpdateRateRequest) (*crmpb.UpdateRateResponse, error) {
	// Convert protobuf request to DTO
	rateReq := &dto.RateUpdateRequest{
		Score:   uint8(req.Score),
		Comment: req.Comment,
		Attachs: req.Attachs,
	}

	// Convert to domain
	rate := convertRateUpdateRequestToDomain(rateReq)

	// Update rate
	err := h.rateUsecase.UpdateRate(ctx, req.Id, rate)
	if err != nil {
		return nil, err
	}

	return &crmpb.UpdateRateResponse{
		Success: true,
	}, nil
}

// DeleteRate xóa đánh giá
// @Summary Xóa đánh giá
// @Description Xóa đánh giá của người dùng
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID đánh giá"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /rate/{id} [delete]
func (h *RateHandler) DeleteRate(ctx context.Context, req *crmpb.RegistryIdRequest) (*sharepb.SubmitResponse, error) {
	// Delete rate
	err := h.rateUsecase.DeleteRate(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Rate deleted successfully",
	}, nil
}

// GetRate lấy thông tin đánh giá
// @Summary Lấy thông tin đánh giá
// @Description Lấy thông tin chi tiết của một đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID đánh giá"
// @Success 200 {object} crmpb.RateResponse
// @Router /rate/detail/{id} [get]
func (h *RateHandler) GetRate(ctx context.Context, req *crmpb.RegistryIdRequest) (*crmpb.RateResponse, error) {
	// Get rate
	rate, err := h.rateUsecase.GetRate(ctx, req.Id, req.Registry)
	if err != nil {
		return nil, err
	}

	response := h.rateMapper.MapRateToPb(ctx, rate)

	// Map user profile
	h.rateMapper.MapUserToPbSingle(ctx, response)

	return response, nil
}

// UpdateHidden cập nhật trạng thái ẩn đánh giá
// @Summary Cập nhật trạng thái ẩn đánh giá
// @Description Cập nhật trạng thái ẩn/hiện của đánh giá
// @Tags Rate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID đánh giá"
// @Param body body crmpb.UpdateHiddenRequest true "Body"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /rate/hidden/{id} [put]
func (h *RateHandler) UpdateHidden(ctx context.Context, req *crmpb.UpdateHiddenRequest) (*sharepb.SubmitResponse, error) {
	// Update hidden status
	err := h.rateUsecase.UpdateHidden(ctx, req.Id, req.Hidden, req.Registry)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Rate hidden status updated successfully",
	}, nil
}

// Helper functions to convert between protobuf and domain/DTO
func convertRateCreateRequestToDomain(req *dto.RateCreateRequest) *domain.Rate {
	if req == nil {
		return nil
	}

	rate := &domain.Rate{
		Anonymous: req.Anonymous,
		OwnerID:   req.OwnerID,
		OwnerOf:   req.OwnerOf,
		Score:     req.Score,
		Comment:   req.Comment,
		ParentID:  req.ParentID,
	}

	if len(req.Attachs) > 0 {
		for _, attachURL := range req.Attachs {
			rate.Attachs = append(rate.Attachs, domain.FeedbackAttachEntity{
				FileURL:  attachURL,
				FileType: "image",
				FileName: "attachment",
			})
		}
	}

	return rate
}

func convertRateUpdateRequestToDomain(req *dto.RateUpdateRequest) *domain.Rate {
	if req == nil {
		return nil
	}

	rate := &domain.Rate{
		Score:   req.Score,
		Comment: req.Comment,
	}

	// Chỉ thêm attachs nếu có trong request (có thể là empty array để xóa tất cả)
	if req.Attachs != nil {
		rate.Attachs = make([]domain.FeedbackAttachEntity, 0)
		for _, attachURL := range req.Attachs {
			rate.Attachs = append(rate.Attachs, domain.FeedbackAttachEntity{
				FileURL:  attachURL,
				FileType: "image",
				FileName: "attachment",
			})
		}
	}

	return rate
}

func (h *RateHandler) convertRateToResponse(ctx context.Context, rate *domain.Rate) *crmpb.RateResponse {
	if rate == nil {
		return nil
	}

	response := &crmpb.RateResponse{
		Id:        rate.ID,
		Anonymous: rate.Anonymous,
		OwnerId:   rate.OwnerID,
		OwnerOf:   rate.OwnerOf,
		Score:     uint32(rate.Score),
		Comment:   rate.Comment,
		CreatedAt: _utils.FormatTimeToString(rate.CreatedAt),
	}

	if rate.ParentID != nil {
		response.ParentId = rate.ParentID
	}

	if rate.CreatedBy != nil {
		response.CreatedBy = rate.CreatedBy

		// Gọi user service để lấy thông tin user
		if !rate.Anonymous {
			profileIDSet := map[uint64]struct{}{
				*rate.CreatedBy: {},
			}
			userProfile, err := h.userClient.GetMapByIDs(ctx, profileIDSet)
			if err == nil && len(userProfile) > 0 {
				response.CreatedUser = userProfile[*rate.CreatedBy]
			}
		}
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, &crmpb.AttachEntity{
				Id:       attach.ID,
				FileUrl:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}