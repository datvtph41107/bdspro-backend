package handler

// type HistoryService struct {
// 	Usecase *shared_usecase.HistoryUsecase
// }

// func NewHistoryService(Usecase *shared_usecase.HistoryUsecase) *HistoryService {
// 	return &HistoryService{
// 		Usecase: Usecase,
// 	}
// }

// // @Summary Lịch sử tách gộp tài sản danh cho tài sản
// // @Tags User: Tài sản
// // @Produce json
// // @Param body query dto.AssetHistoryDTO true "body"
// // @Param id path uint64 true "ID"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/asset/history/:id [get]
// func (s *AssetService) History(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.HistoryResponse, error) {
// 	body := dto.AssetHistoryDTO{}
// 	if err := copier.Copy(&body, req); err != nil {
// 		return nil, err
// 	}

// 	result, total, err := s.UC.History(ctx, req.Id, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	history := make([]*bdspropb.History, len(*result))
// 	copier.Copy(&history, &result)
// 	s.ProfileClient.MapProfileToHistory(ctx, history)

// 	return &bdspropb.HistoryResponse{
// 		Data:  history,
// 		Total: uint64(total),
// 	}, nil
// }
