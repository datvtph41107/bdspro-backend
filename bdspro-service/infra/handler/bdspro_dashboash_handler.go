package handler

import (
	"bdspro/internal/job"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"
)

type BdsproDashboardHandler struct {
	bdspropb.UnimplementedBdsproDashboardServiceServer
	assetStatsJob    *job.AssetDashboardStatsJob
	productStatsJob  *job.ProductDashboardStatsJob
	postStatsJob     *job.PostDashboardStatsJob
	dashboardUsecase *usecases.DashboardUsecase
}

func NewBdsproDashboardHandler(assetStatsJob *job.AssetDashboardStatsJob, productStatsJob *job.ProductDashboardStatsJob, postStatsJob *job.PostDashboardStatsJob, dashboardUsecase *usecases.DashboardUsecase) *BdsproDashboardHandler {
	return &BdsproDashboardHandler{
		assetStatsJob:    assetStatsJob,
		productStatsJob:  productStatsJob,
		postStatsJob:     postStatsJob,
		dashboardUsecase: dashboardUsecase,
	}
}
func (h *BdsproDashboardHandler) GetCashProfit(ctx context.Context, req *bdspropb.GetCashProfitRequest) (*bdspropb.GetCashProfitResponse, error) {
	return nil, nil
}

func (h *BdsproDashboardHandler) GetCashProfitTime(ctx context.Context, req *bdspropb.GetCashProfitRequest) (*bdspropb.GetCashProfitTimeResponse, error) {
	return nil, nil
}

func (h *BdsproDashboardHandler) GetAssetStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, err
	}
	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, err
	}
	stats, stats2, err := h.assetStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, err
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

func (h *BdsproDashboardHandler) GetProductStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, err
	}
	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, err
	}
	stats, stats2, err := h.productStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, err
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

func (h *BdsproDashboardHandler) GetPostStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, err
	}
	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, err
	}
	stats, stats2, err := h.postStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, err
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

func (h *BdsproDashboardHandler) GetCountPostByTime(ctx context.Context, req *bdspropb.GetCountPostByTimeRequest) (*bdspropb.GetCountPostByTimeResponse, error) {
	from, err := time.Parse(time.DateOnly, req.From)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse(time.DateOnly, req.To)
	if err != nil {
		return nil, err
	}
	countPostByTime, err := h.dashboardUsecase.GetCountPostByTime(ctx, from, to)
	if err != nil {
		return nil, err
	}
	countPostByTimeResponse := make([]*bdspropb.CountPostByTime, len(countPostByTime))
	for i, count := range countPostByTime {
		countPostByTimeResponse[i] = &bdspropb.CountPostByTime{
			Date:  count.Date.Format(time.DateOnly),
			Count: count.Count,
		}
	}
	return &bdspropb.GetCountPostByTimeResponse{Data: countPostByTimeResponse}, nil
}
