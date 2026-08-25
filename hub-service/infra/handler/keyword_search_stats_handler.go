package handler

import (
	"context"

	_errors "common/errors"
	"hub/infra/mapper"
	"hub/internal/dto"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type KeywordSearchStatsHandler struct {
	hubpb.UnimplementedKeywordSearchStatsServiceServer
	KeywordSearchMapper  *mapper.KeywordSearchMapper
	KeywordSearchUsecase _usecase.IKeywordSearchUsecase
}

func NewKeywordSearchStatsHandler(
	keywordSearchUsecase _usecase.IKeywordSearchUsecase,
	keywordSearchMapper *mapper.KeywordSearchMapper,
) *KeywordSearchStatsHandler {
	return &KeywordSearchStatsHandler{
		KeywordSearchMapper:  keywordSearchMapper,
		KeywordSearchUsecase: keywordSearchUsecase,
	}
}

// @Summary Thống kê từ khóa tìm kiếm
// @Description Top từ khóa theo lượt tìm (count giảm dần) trong lastDates ngày gần nhất (mặc định 3, giới hạn 1-30)
// @Tags KeywordSearchStats
// @Produce json
// @Param lastDates query int false "Số ngày gần nhất (1-30, mặc định 3)"
// @Param limit query int false "Số từ khóa tối đa (mặc định 50, tối đa 100)"
// @Param source query string false "Lọc theo source: crm | tqd"
// @Success 200 {object} hubpb.GetKeywordSearchStatsResponse
// @Router /v2/hub/keyword-search-stats [get]
func (h *KeywordSearchStatsHandler) GetKeywordSearchStats(
	ctx context.Context,
	req *hubpb.GetKeywordSearchStatsRequest,
) (*hubpb.GetKeywordSearchStatsResponse, error) {
	if req == nil {
		req = &hubpb.GetKeywordSearchStatsRequest{}
	}

	resp, errDTO := h.KeywordSearchUsecase.GetStats(ctx, &dto.GetKeywordSearchStatsRequest{
		LastDates: req.GetLastDates(),
		Limit:     req.GetLimit(),
		Source:    req.GetSource(),
	})
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	if resp == nil {
		return &hubpb.GetKeywordSearchStatsResponse{
			Data: []*hubpb.KeywordSearchItem{},
		}, nil
	}

	data := make([]*hubpb.KeywordSearchItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &hubpb.KeywordSearchItem{
			Keyword: item.Keyword,
			Source:  item.Source,
			Count:   item.Count,
		})
	}

	return &hubpb.GetKeywordSearchStatsResponse{
		Data:          data,
		TotalSearches: resp.TotalSearches,
		LastDates:     resp.LastDates,
		Limit:         resp.Limit,
	}, nil
}
