package handler

import (
	"context"

	_errors "common/errors"
	"hub/internal/dto"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type SearchLatencyStatsHandler struct {
	hubpb.UnimplementedSearchLatencyStatsServiceServer
	SearchLatencyUsecase _usecase.ISearchLatencyUsecase
}

func NewSearchLatencyStatsHandler(searchLatencyUsecase _usecase.ISearchLatencyUsecase) *SearchLatencyStatsHandler {
	return &SearchLatencyStatsHandler{
		SearchLatencyUsecase: searchLatencyUsecase,
	}
}

// @Summary Thống kê latency API search có text
// @Description Trả avg / min / max / p50 / p95 / p99 / count trong lastDates ngày (mặc định path=/v2/tqd/discovery/search)
// @Tags SearchLatencyStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3)"
// @Param path query string false "Path API search; mặc định /v2/tqd/discovery/search"
// @Success 200 {object} hubpb.GetSearchLatencyDashboardResponse
// @Router /v2/hub/search-latency/dashboard [get]
func (h *SearchLatencyStatsHandler) GetSearchLatencyDashboard(
	ctx context.Context,
	req *hubpb.GetSearchLatencyDashboardRequest,
) (*hubpb.GetSearchLatencyDashboardResponse, error) {
	if req == nil {
		req = &hubpb.GetSearchLatencyDashboardRequest{}
	}
	resp, errDTO := h.SearchLatencyUsecase.GetDashboard(ctx, &dto.GetSearchLatencyDashboardRequest{
		LastDates: req.GetLastDates(),
		Path:      req.GetPath(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		resp = &dto.GetSearchLatencyDashboardResponse{}
	}
	return &hubpb.GetSearchLatencyDashboardResponse{
		Path:          resp.Path,
		TotalRequests: resp.TotalRequests,
		AvgDurationMs: resp.AvgDurationMs,
		MinDurationMs: resp.MinDurationMs,
		MaxDurationMs: resp.MaxDurationMs,
		LastDates:     resp.LastDates,
		P50DurationMs: resp.P50DurationMs,
		P95DurationMs: resp.P95DurationMs,
		P99DurationMs: resp.P99DurationMs,
	}, nil
}

// @Summary Biểu đồ latency search theo ngày/giờ/phút
// @Description gap=dates|hours|minutes; trả count/avg/min/max/p50/p95/p99 để vẽ biểu đồ. minutes giới hạn lastDates tối đa 2.
// @Tags SearchLatencyStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3; minutes max 2)"
// @Param gap query string false "dates | hours | minutes"
// @Param path query string false "Path API search"
// @Success 200 {object} hubpb.GetSearchLatencyDailyStatsResponse
// @Router /v2/hub/search-latency/daily [get]
func (h *SearchLatencyStatsHandler) GetSearchLatencyDailyStats(
	ctx context.Context,
	req *hubpb.GetSearchLatencyDailyStatsRequest,
) (*hubpb.GetSearchLatencyDailyStatsResponse, error) {
	if req == nil {
		req = &hubpb.GetSearchLatencyDailyStatsRequest{}
	}
	resp, errDTO := h.SearchLatencyUsecase.GetDailyStats(ctx, &dto.GetSearchLatencyDailyStatsRequest{
		LastDates: req.GetLastDates(),
		Gap:       req.GetGap(),
		Path:      req.GetPath(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		resp = &dto.GetSearchLatencyDailyStatsResponse{
			Data: []dto.SearchLatencyDailyItem{},
			Gap:  dto.RequestStatsGapDates,
		}
	}

	data := make([]*hubpb.SearchLatencyDailyItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.SearchLatencyDailyItem{
			Date:          item.Date,
			Count:         item.Count,
			AvgDurationMs: item.AvgDurationMs,
			MinDurationMs: item.MinDurationMs,
			MaxDurationMs: item.MaxDurationMs,
			P50DurationMs: item.P50DurationMs,
			P95DurationMs: item.P95DurationMs,
			P99DurationMs: item.P99DurationMs,
		})
	}

	return &hubpb.GetSearchLatencyDailyStatsResponse{
		Data:      data,
		LastDates: resp.LastDates,
		Gap:       resp.Gap,
		Path:      resp.Path,
	}, nil
}
