package grpc_service

import (
	_dto "common/dto"
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"

	crmpb "pb/types/crm"
)

type HistoryService struct {
	crmpb.UnimplementedHistoryServiceServer
	UC *usecase.HistoryUsecase
}

func NewHistoryService(uc *usecase.HistoryUsecase) *HistoryService {
	return &HistoryService{UC: uc}
}

// @Summary Lấy lịch sử của người dùng
// @Description Lấy lịch sử của người dùng
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param ownerId path int true "ID của người dùng"
// @Param ownerType path int true "Loại người dùng"
// @Param page query crmpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/list/{ownerId}/{ownerType} [get]
func (s *HistoryService) Search(ctx context.Context, req *crmpb.HistorySearchRequest) (*crmpb.HistoryListDTO, error) {
	histories, total, err := s.UC.Search(ctx, req.OwnerId, enums.EOwnerType(req.OwnerType), dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
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
		historiesPb[i] = mapper.HistoryDomainToPb(&history)
	}

	return &crmpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của liên hệ
// @Description Lấy lịch sử của liên hệ
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param contactId path int true "ID của liên hệ"
// @Param page query crmpb.HistoryContactRequest false "Trang"
// @Security BearerAuth
// @Router /history/contact/{contactId} [get]
func (s *HistoryService) ContactHistory(ctx context.Context, req *crmpb.HistoryContactRequest) (*crmpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ContactHistory(ctx, req.ContactId, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := mapper.HistoryToPbs(histories)

	return &crmpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}
