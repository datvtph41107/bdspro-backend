package service

import (
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type SharingAccessService struct {
	pb_bdspro.UnimplementedSharingAccessServiceServer
	UC *shared_usecase.SharingAccessUsecase
}

func NewSharingAccessService(uc *shared_usecase.SharingAccessUsecase) *SharingAccessService {
	return &SharingAccessService{
		UC: uc,
	}
}

func (s *SharingAccessService) BulkSave(ctx context.Context, req *pb_bdspro.SharingAccessBulkRequest) (*pb_bdspro.SharingAccessResponse, error) {
	dtos := dto.SharingAccessBulk{}
	copier.Copy(&dtos, &req)
	entities, err := s.UC.BulkSave(ctx, req.DomainId, enums.EOwnerType(req.FromType), dtos)
	if err != nil {
		return nil, err
	}
	pb_entities := []*pb_bdspro.SharingAccess{}
	copier.Copy(&pb_entities, &entities)
	return &pb_bdspro.SharingAccessResponse{
		Data: pb_entities,
	}, nil
}
func (s *SharingAccessService) SearchProduct(ctx context.Context, req *pb_bdspro.SharingAccessRequest) (*pb_bdspro.SharingAccessResponse, error) {
	searchDTO := dto.SharingAccessSearch{}
	copier.Copy(&searchDTO, &req)

	if searchDTO.DomainID == 0 || searchDTO.Domain == 0 || searchDTO.FromType == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "domainID and domain and fromType are required",
		}
	}
	entities, err := s.UC.SearchForProduct(ctx, enums.EOwnerType(req.FromType), searchDTO)
	if err != nil {
		return nil, err
	}
	pb_entities := []*pb_bdspro.SharingAccess{}
	copier.Copy(&pb_entities, &entities)
	return &pb_bdspro.SharingAccessResponse{
		Data: pb_entities,
	}, nil
}
