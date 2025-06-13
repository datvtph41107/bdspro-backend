package grpc_service

import (
	_dto "common/dto"
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
)

type SharingAccessService struct {
	UC *usecase.SharingAccessUsecase
	crmpb.UnimplementedSharingAccessServiceServer
}

func NewSharingAccessService(uc *usecase.SharingAccessUsecase) *SharingAccessService {
	return &SharingAccessService{UC: uc}
}

func (s *SharingAccessService) Share(ctx context.Context, req *crmpb.SharingRequest) (*crmpb.SharingResponse, error) {
	receivers := mapper.ListSharingPbToDomain(req.Receivers)
	err := s.UC.BulkShare(ctx, req.ContactId, req.DeletedIds, receivers)
	if err != nil {
		return nil, err
	}
	return &crmpb.SharingResponse{
		Status:  1,
		Message: "Success",
	}, nil
}

func (s *SharingAccessService) List(ctx context.Context, req *crmpb.SharingListRequest) (*crmpb.SharingListResponse, error) {
	query := dto.SharingAccessSearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
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
