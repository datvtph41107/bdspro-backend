package handler

import (
	"bdspro/internal/job"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	_utils "common/utils"
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
	productUsecase   *shared_usecase.ProductUsecase
}

func NewBdsproDashboardHandler(assetStatsJob *job.AssetDashboardStatsJob, productStatsJob *job.ProductDashboardStatsJob, postStatsJob *job.PostDashboardStatsJob, dashboardUsecase *usecases.DashboardUsecase, productUsecase *shared_usecase.ProductUsecase) *BdsproDashboardHandler {
	return &BdsproDashboardHandler{
		assetStatsJob:    assetStatsJob,
		productStatsJob:  productStatsJob,
		postStatsJob:     postStatsJob,
		dashboardUsecase: dashboardUsecase,
		productUsecase:   productUsecase,
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

// @Summary Lấy danh sách sản phẩm gợi ý hôm nay
// @Description Lấy ngẫu nhiên 5 sản phẩm để gợi ý
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param req query bdspropb.SuggestTodayRequest true "Request"
// @Router /v2/bdspro/dashboard/suggest/today [get]
func (h *BdsproDashboardHandler) GetSuggestToday(ctx context.Context, req *bdspropb.SuggestTodayRequest) (*bdspropb.SuggestTodayResponse, error) {
	products, err := h.productUsecase.GetSupportToday(ctx)
	if err != nil {
		return nil, err
	}

	items := []*bdspropb.SuggestTodayItem{}
	for _, product := range products {
		var imageUrl string
		var price float64

		// Lấy ảnh đầu tiên từ MediaList
		if len(product.MediaList) > 0 {
			imageUrl = product.MediaList[0].MediaURL
		}

		// Lấy giá từ Price
		if product.Price != nil {
			if product.Price.SalePrice != nil {
				price = *product.Price.SalePrice
			} else if product.Price.RentPrice != nil {
				price = *product.Price.RentPrice
			}
		}

		items = append(items, &bdspropb.SuggestTodayItem{
			ImageUrl:   imageUrl,
			Title:      product.Name,
			Price:      price,
			LinkId:     product.ID,
			DomainType: 10, // 10 = Product domain type
		})
	}

	return &bdspropb.SuggestTodayResponse{
		Data: items,
	}, nil
}

// @Summary Đếm số tin đã đăng, sản phẩm quan tâm và đã lưu của user
// @Description Lấy thống kê số lượng tin đã đăng, sản phẩm quan tâm và đã lưu
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param req query bdspropb.GetCountOfUserRequest true "Request"
// @Router /v2/bdspro/dashboard/count-of-user [get]
func (h *BdsproDashboardHandler) GetCountOfUser(ctx context.Context, req *bdspropb.GetCountOfUserRequest) (*bdspropb.GetCountOfUserResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)

	counts, err := h.dashboardUsecase.GetCountOfUser(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return &bdspropb.GetCountOfUserResponse{
		TotalPosted:            counts.TotalPosted,
		TotalInterestedProduct: counts.TotalInterestedProduct,
		TotalSavedProduct:      counts.TotalSavedProduct,
	}, nil
}
