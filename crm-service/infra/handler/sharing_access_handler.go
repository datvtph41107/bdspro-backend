package handler

import (
	_dto "common/domain/dto"
	"context"
	crmpb "pb/types/crm"

	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
)

type SharingAccessService struct {
	UC *usecase.SharingAccessUsecase
	crmpb.UnimplementedSharingAccessServiceServer
}

func NewSharingAccessService(uc *usecase.SharingAccessUsecase) *SharingAccessService {
	return &SharingAccessService{UC: uc}
}

// @Summary Chia sẻ liên hệ
// @Description Chia sẻ liên hệ
// @Tags Chia sẻ liên hệ
// @Accept json
// @Produce json
// @Param request body crmpb.SharingRequest true "Thông tin chia sẻ"
// @Security BearerAuth
// @Router /sharing/share [post]
func (s *SharingAccessService) Share(ctx context.Context, req *crmpb.SharingRequest) (*crmpb.SharingResponse, error) {
	receivers := mapper.ListSharingPbToDomain(req.Receivers)
	err := s.UC.BulkShare(ctx,
		req.ContactId,
		req.Permissions,
		req.Note,
		enums.EOwnerOf(req.OwnerType),
		req.DeletedIds,
		receivers,
	)
	if err != nil {
		return nil, err
	}
	return &crmpb.SharingResponse{
		Status:  1,
		Message: "Success",
	}, nil
}

// @Summary Lấy danh sách chia sẻ
// @Description Lấy danh sách chia sẻ
// @Tags Chia sẻ liên hệ
// @Accept json
// @Produce json
// @Param request query crmpb.SharingListRequest true "Thông tin chia sẻ"
// @Security BearerAuth
// @Param contactId path int true "ID của liên hệ"
// @Router /sharing/list/{contactId} [get]
func (s *SharingAccessService) List(ctx context.Context, req *crmpb.SharingListRequest) (*crmpb.SharingListResponse, error) {
	query := dto.SharingAccessSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	sharing, total, err := s.UC.List(ctx, req.ContactId, &query)
	if err != nil {
		return nil, err
	}
	pbs := mapper.ListSharingToPb(sharing)

	return &crmpb.SharingListResponse{
		Data:  pbs,
		Total: int32(total),
	}, nil
}