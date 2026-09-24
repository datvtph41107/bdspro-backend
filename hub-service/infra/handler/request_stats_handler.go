package handler

import (
	"context"

	_errors "common/errors"
	"hub/infra/mapper"
	"hub/internal/dto"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"

	"google.golang.org/protobuf/types/known/emptypb"
)

type RequestStatsHandler struct {
	hubpb.UnimplementedRequestStatsServiceServer
	RequestStatsMapper  *mapper.RequestStatsMapper
	RequestStatsUsecase _usecase.IRequestStatsUsecase
}

func NewRequestStatsHandler(
	requestStatsUsecase _usecase.IRequestStatsUsecase,
	requestStatsMapper *mapper.RequestStatsMapper,
) *RequestStatsHandler {
	return &RequestStatsHandler{
		RequestStatsMapper:  requestStatsMapper,
		RequestStatsUsecase: requestStatsUsecase,
	}
}

// @Summary Track request stats
// @Description Gateway fire-and-forget mỗi HTTP request để đếm traffic (buffer ~10s rồi flush DB)
// @Tags RequestStats
// @Accept json
// @Produce json
// @Param body body hubpb.TrackRequestRequest true "Thông tin request cần đếm"
// @Success 200 {object} google.protobuf.Empty
// @Router /v2/hub/request-stats/track [post]
func (h *RequestStatsHandler) TrackRequest(
	ctx context.Context,
	req *hubpb.TrackRequestRequest,
) (*emptypb.Empty, error) {
	dtoReq := h.RequestStatsMapper.MapTrackRequest(req)
	if errDTO := h.RequestStatsUsecase.TrackRequest(ctx, dtoReq); errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	return &emptypb.Empty{}, nil
}

// @Summary Dashboard thống kê request hệ thống
// @Description Tổng request, top path và phân bổ theo service theo lastDates ngày gần nhất (mặc định 3, giới hạn 1-30)
// @Tags RequestStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3)"
// @Success 200 {object} hubpb.GetRequestStatsDashboardResponse
// @Router /v2/hub/request-stats/dashboard [get]
func (h *RequestStatsHandler) GetRequestStatsDashboard(
	ctx context.Context,
	req *hubpb.GetRequestStatsDashboardRequest,
) (*hubpb.GetRequestStatsDashboardResponse, error) {
	if req == nil {
		req = &hubpb.GetRequestStatsDashboardRequest{}
	}

	resp, errDTO := h.RequestStatsUsecase.GetDashboard(ctx, &dto.GetRequestStatsDashboardRequest{
		LastDates: req.GetLastDates(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetRequestStatsDashboardResponse{
			Data:      []*hubpb.RequestTopPathItem{},
			ByService: []*hubpb.RequestServiceCountItem{},
		}, nil
	}

	data := make([]*hubpb.RequestTopPathItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.RequestTopPathItem{
			Method:  item.Method,
			Path:    item.Path,
			Service: item.Service,
			Count:   item.Count,
		})
	}

	byService := make([]*hubpb.RequestServiceCountItem, 0, len(resp.ByService))
	for _, item := range resp.ByService {
		byService = append(byService, &hubpb.RequestServiceCountItem{
			Service: item.Service,
			Count:   item.Count,
		})
	}

	return &hubpb.GetRequestStatsDashboardResponse{
		TotalRequests: resp.TotalRequests,
		Data:          data,
		ByService:     byService,
		LastDates:     resp.LastDates,
	}, nil
}

// @Summary Thống kê request theo ngày/giờ/phút
// @Description Tổng request theo gap=dates (YYYY-MM-DD), hours (YYYY-MM-DD HH:00) hoặc minutes (YYYY-MM-DD HH:MM, slot 5 phút) trong lastDates ngày gần nhất. minutes giới hạn lastDates tối đa 2.
// @Tags RequestStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3; gap=minutes tối đa 2)"
// @Param gap query string false "Khoảng gom: dates (mặc định) | hours | minutes"
// @Success 200 {object} hubpb.GetRequestDailyStatsResponse
// @Router /v2/hub/request-stats/daily [get]
func (h *RequestStatsHandler) GetRequestDailyStats(
	ctx context.Context,
	req *hubpb.GetRequestDailyStatsRequest,
) (*hubpb.GetRequestDailyStatsResponse, error) {
	if req == nil {
		req = &hubpb.GetRequestDailyStatsRequest{}
	}

	resp, errDTO := h.RequestStatsUsecase.GetDailyStats(ctx, &dto.GetRequestDailyStatsRequest{
		LastDates: req.GetLastDates(),
		Gap:       req.GetGap(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetRequestDailyStatsResponse{
			Data: []*hubpb.RequestDailyItem{},
			Gap:  dto.RequestStatsGapDates,
		}, nil
	}

	data := make([]*hubpb.RequestDailyItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.RequestDailyItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}

	return &hubpb.GetRequestDailyStatsResponse{
		Data:      data,
		LastDates: resp.LastDates,
		Gap:       resp.Gap,
	}, nil
}

// @Summary Dashboard thống kê người dùng tương tác
// @Description Tổng tương tác, số active users và top user theo lastDates ngày gần nhất (chỉ request có JWT)
// @Tags RequestStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3)"
// @Success 200 {object} hubpb.GetUserInteractionDashboardResponse
// @Router /v2/hub/request-stats/users/dashboard [get]
func (h *RequestStatsHandler) GetUserInteractionDashboard(
	ctx context.Context,
	req *hubpb.GetUserInteractionDashboardRequest,
) (*hubpb.GetUserInteractionDashboardResponse, error) {
	if req == nil {
		req = &hubpb.GetUserInteractionDashboardRequest{}
	}

	resp, errDTO := h.RequestStatsUsecase.GetUserInteractionDashboard(ctx, &dto.GetUserInteractionDashboardRequest{
		LastDates: req.GetLastDates(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetUserInteractionDashboardResponse{
			Data: []*hubpb.UserInteractionTopItem{},
		}, nil
	}

	data := make([]*hubpb.UserInteractionTopItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.UserInteractionTopItem{
			ProfileId:      item.ProfileID,
			OrganizationId: item.OrganizationID,
			Count:          item.Count,
			FullName:       item.FullName,
			Avatar:         item.Avatar,
			Phone:          item.Phone,
			TickVerified:   item.TickVerified,
		})
	}

	return &hubpb.GetUserInteractionDashboardResponse{
		TotalInteractions: resp.TotalInteractions,
		ActiveUsers:       resp.ActiveUsers,
		Data:              data,
		LastDates:         resp.LastDates,
	}, nil
}

// @Summary Thống kê tương tác user theo ngày/giờ/phút
// @Description Lượt tương tác và active users theo gap=dates|hours|minutes trong lastDates ngày gần nhất. minutes giới hạn lastDates tối đa 2.
// @Tags RequestStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3; gap=minutes tối đa 2)"
// @Param gap query string false "Khoảng gom: dates (mặc định) | hours | minutes"
// @Success 200 {object} hubpb.GetUserInteractionDailyStatsResponse
// @Router /v2/hub/request-stats/users/daily [get]
func (h *RequestStatsHandler) GetUserInteractionDailyStats(
	ctx context.Context,
	req *hubpb.GetUserInteractionDailyStatsRequest,
) (*hubpb.GetUserInteractionDailyStatsResponse, error) {
	if req == nil {
		req = &hubpb.GetUserInteractionDailyStatsRequest{}
	}

	resp, errDTO := h.RequestStatsUsecase.GetUserInteractionDailyStats(ctx, &dto.GetUserInteractionDailyStatsRequest{
		LastDates: req.GetLastDates(),
		Gap:       req.GetGap(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetUserInteractionDailyStatsResponse{
			Data: []*hubpb.UserInteractionDailyItem{},
			Gap:  dto.RequestStatsGapDates,
		}, nil
	}

	data := make([]*hubpb.UserInteractionDailyItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.UserInteractionDailyItem{
			Date:         item.Date,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}

	return &hubpb.GetUserInteractionDailyStatsResponse{
		Data:      data,
		LastDates: resp.LastDates,
		Gap:       resp.Gap,
	}, nil
}

// @Summary Danh sách audit log request đặc biệt (spec_urls)
// @Description Trả từng lần gọi theo timestamp (không gộp count). Lọc path ví dụ /v2/crm/public/content/sitemap.xml
// @Tags RequestStats
// @Produce json
// @Param path query string false "Lọc theo path (vd: /v2/crm/public/content/sitemap.xml)"
// @Param page query int false "Trang (0-based, mặc định 0)"
// @Param size query int false "Số bản ghi mỗi trang (mặc định 20, tối đa 100)"
// @Success 200 {object} hubpb.GetRequestAuditLogsResponse
// @Router /v2/hub/request-stats/audit [get]
func (h *RequestStatsHandler) GetRequestAuditLogs(
	ctx context.Context,
	req *hubpb.GetRequestAuditLogsRequest,
) (*hubpb.GetRequestAuditLogsResponse, error) {
	if req == nil {
		req = &hubpb.GetRequestAuditLogsRequest{}
	}

	dtoReq := &dto.GetRequestAuditLogsRequest{
		Path: req.GetPath(),
	}
	dtoReq.Page = req.GetPage()
	dtoReq.Size = req.GetSize()

	resp, errDTO := h.RequestStatsUsecase.GetAuditLogs(ctx, dtoReq)
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetRequestAuditLogsResponse{
			Data: []*hubpb.RequestAuditLogItem{},
		}, nil
	}

	data := make([]*hubpb.RequestAuditLogItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.RequestAuditLogItem{
			Id:         item.ID,
			Method:     item.Method,
			Path:       item.Path,
			Service:    item.Service,
			StatusCode: item.StatusCode,
			DurationMs: item.DurationMs,
			Timestamp:  item.Timestamp,
		})
	}

	return &hubpb.GetRequestAuditLogsResponse{
		Data:  data,
		Total: resp.Total,
		Page:  resp.Page,
		Size:  resp.Size,
	}, nil
}
