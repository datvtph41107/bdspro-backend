package service

import (
	_dto "common/dto"
	"context"
	"noti/infra/mapper"
	"noti/internal/dto"
	"noti/internal/usecase"

	shared_dto "pb/dto"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"
)

type HistoryService struct {
	notificationpb.UnimplementedHistoryServiceServer
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
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/list/{ownerId}/{ownerType} [get]
func (s *HistoryService) Search(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.Search(ctx, req.OwnerId, shared_enum.EOwnerType(req.OwnerType), dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	historiesPb := make([]*notificationpb.HistoryDTO, len(histories))
	for i, history := range histories {
		historiesPb[i] = mapper.HistoryDomainToPb(&history)
	}

	return &notificationpb.HistoryListDTO{
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
// @Param page query notificationpb.HistoryContactRequest false "Trang"
// @Security BearerAuth
// @Router /history/contact/{contactId} [get]
func (s *HistoryService) ContactHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ContactHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *HistoryService) AssetHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.AssetHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *HistoryService) ProductHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ProductHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *HistoryService) ProductChildHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ProductChildHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:          historiesPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *HistoryService) Create(ctx context.Context, req *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error) {
	history, err := s.UC.Create(ctx, &shared_dto.HistoryDTO{
		ActionType: shared_enum.EHistory(req.ActionType),
		TargetId:   req.TargetId,
		OwnerID:    req.OwnerId,
		Note:       req.Note,
		Title:      req.Title,
	})
	if err != nil {
		return nil, err
	}
	return mapper.HistoryDomainToPb(history), nil
}
