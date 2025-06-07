package grpc_service

import (
	_dto "common/dto"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	"time"

	crmpb "pb/types/crm"
)

type HistoryService struct {
	crmpb.UnimplementedHistoryServiceServer
	UC *usecase.HistoryUsecase
}

func NewHistoryService(uc *usecase.HistoryUsecase) *HistoryService {
	return &HistoryService{UC: uc}
}

func (s *HistoryService) Search(ctx context.Context, req *crmpb.HistorySearchRequest) (*crmpb.HistoryListDTO, error) {
	histories, total, err := s.UC.Search(ctx, req.OwnerId, enums.EOwnerType(req.OwnerType), dto.HistorySearchDTO{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	historiesPb := make([]*crmpb.HistoryDTO, len(histories))
	for i, history := range histories {
		historiesPb[i] = s.DomainToPb(&history)
	}

	return &crmpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

// --map--

func (s *HistoryService) DomainToPb(history *domain.HistoryEntity) *crmpb.HistoryDTO {
	return &crmpb.HistoryDTO{
		TargetId:   history.TargetId,
		TargetType: int32(history.TargetType),
		ActionType: int32(history.ActionType),
		Note:       history.Note,
		PreStage:   history.PreStage,
		AfterStage: history.AfterStage,
		OwnerId:    history.OwnerID,
		OwnerType:  int32(history.OwnerType),
		CreatedAt:  history.CreatedAt.Format(time.RFC3339),
		CreatedBy:  history.CreatedBy,
	}
}
