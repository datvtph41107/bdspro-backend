package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"

	"notification/infra/mapper"
	"notification/internal/dto"
	"notification/internal/usecase"

	shared_dto "pb/dto"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HistoryHandler struct {
	notificationpb.UnimplementedHistoryServiceServer
	UC                 *usecase.HistoryUsecase
	Mapper             *mapper.HistoryMapper
	HistoryAuthUsecase *usecase.HistoryAuthUsecase
	HistoryAuthMapper  *mapper.HistoryAuthMapper
}

func NewHistoryHandler(
	uc *usecase.HistoryUsecase,
	mapper *mapper.HistoryMapper,
	historyAuthUsecase *usecase.HistoryAuthUsecase,
	historyAuthMapper *mapper.HistoryAuthMapper,
) *HistoryHandler {
	return &HistoryHandler{
		UC:                 uc,
		Mapper:             mapper,
		HistoryAuthUsecase: historyAuthUsecase,
		HistoryAuthMapper:  historyAuthMapper,
	}
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
func (s *HistoryHandler) Search(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.Search(ctx, req.OwnerId, shared_enum.EOwnerType(req.OwnerType), dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	historiesPb := make([]*notificationpb.HistoryDTO, len(histories))
	for i, history := range histories {
		historiesPb[i] = s.Mapper.HistoryDomainToPb(&history)
	}

	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của CRM
// @Description Lấy lịch sử của CRM
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param crmId path int true "ID của CRM"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/crm/{crmId} [get]
func (s *HistoryHandler) CrmHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.CrmHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	historiesPb := s.Mapper.HistoryToPbs(histories)
	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của liên hệ
// @Description Lấy lịch sử của liên hệ
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param contactId path int true "ID của liên hệ"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/contact/{contactId} [get]
func (s *HistoryHandler) ContactHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ContactHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := s.Mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

func (s *HistoryHandler) AssetHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.AssetHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := s.Mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

func (s *HistoryHandler) ProductHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ProductHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := s.Mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của người dùng
// @Description Lấy lịch sử của người dùng
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param page query notificationpb.HistoryAuthSearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/auth [get]
func (s *HistoryHandler) AuthHistory(ctx context.Context, req *notificationpb.HistoryAuthSearchRequest) (*notificationpb.HistoryAuthListResponse, error) {
	if s.HistoryAuthUsecase == nil || s.HistoryAuthMapper == nil {
		return nil, status.Error(codes.Unavailable, "history auth service not available")
	}

	// if req.GetUserId() == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "userId is required")
	// }

	// var organizationID *uint64
	// if req.OrganizationId != nil {
	// 	val := req.GetOrganizationId()
	// 	organizationID = &val
	// }

	var success *bool
	if req.Success != nil {
		val := req.GetSuccess()
		success = &val
	}

	searchDTO := dto.HistoryAuthSearchDTO{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		ActionType: req.GetActionType(),
		Channel:    req.GetChannel(),
		DeviceID:   req.GetDeviceId(),
		Success:    success,
	}
	// UserID:         req.GetUserId(),
	// OrganizationID: organizationID,

	if from := req.GetFromDate(); from != "" {
		if parsed := _utils.ParseStringToTime(from); parsed != nil {
			searchDTO.FromDate = parsed
		}
	}

	if to := req.GetToDate(); to != "" {
		if parsed := _utils.ParseStringToTime(to); parsed != nil {
			searchDTO.ToDate = parsed
		}
	}

	histories, total, err := s.HistoryAuthUsecase.SearchHistoryAuth(ctx, searchDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to load security history: %v", err))
	}

	data := s.HistoryAuthMapper.DomainsToPb(histories)
	page := searchDTO.Pagable.GetPage()
	size := searchDTO.Pagable.GetSize()
	return &notificationpb.HistoryAuthListResponse{
		Data:  data,
		Total: uint64(total),
		Page:  page,
		Size:  size,
	}, nil
}

// @Summary Lấy lịch sử của sản phẩm con
// @Description Lấy lịch sử của sản phẩm con
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param productId path int true "ID của sản phẩm"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/product/child/{productId} [get]
func (s *HistoryHandler) ProductChildHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.ProductChildHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}

	historiesPb := s.Mapper.HistoryToPbs(histories)

	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của quy tắc
// @Description Lấy lịch sử của quy tắc
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param ruleEventId path int true "ID của quy tắc"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/rule-event/{ruleEventId} [get]
func (s *HistoryHandler) RuleEventHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.RuleEventHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	historiesPb := s.Mapper.HistoryToPbs(histories)
	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

func (s *HistoryHandler) Create(ctx context.Context, req *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error) {
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
	return s.Mapper.HistoryDomainToPb(history), nil
}

func (s *HistoryHandler) CreateInternal(ctx context.Context, req *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error) {
	history, err := s.UC.CreateInternal(ctx, &shared_dto.HistoryDTO{
		ActionType: shared_enum.EHistory(req.ActionType),
		TargetId:   req.TargetId,
		OwnerID:    req.OwnerId,
		Note:       req.Note,
		Title:      req.Title,
	})
	if err != nil {
		return nil, err
	}
	return s.Mapper.HistoryDomainToPb(history), nil
}

func (s *HistoryHandler) SearchInternal(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.SearchInternal(ctx, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	historiesPb := s.Mapper.HistoryToPbs(histories)
	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của chiến dịch quảng cáo
// @Description Lấy lịch sử của chiến dịch quảng cáo
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param campaignId path int true "ID của chiến dịch"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/campaign/{campaignId} [get]
func (s *HistoryHandler) CampaignHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.CampaignHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	historiesPb := s.Mapper.HistoryToPbs(histories)
	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy lịch sử của gói quảng cáo
// @Description Lấy lịch sử của gói quảng cáo
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param packageId path int true "ID của gói quảng cáo"
// @Param page query notificationpb.HistorySearchRequest false "Trang"
// @Security BearerAuth
// @Router /history/package/{packageId} [get]
func (s *HistoryHandler) PackageHistory(ctx context.Context, req *notificationpb.HistorySearchRequest) (*notificationpb.HistoryListDTO, error) {
	histories, total, err := s.UC.PackageHistory(ctx, req.Id, dto.HistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	historiesPb := s.Mapper.HistoryToPbs(histories)
	return &notificationpb.HistoryListDTO{
		Data:  historiesPb,
		Total: uint32(total),
	}, nil
}
