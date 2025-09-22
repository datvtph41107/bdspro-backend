package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"github.com/jinzhu/copier"
)

type SharingAccessService struct {
	bdspropb.UnimplementedSharingAccessServiceServer
	Usecase            *shared_usecase.SharingAccessUsecase
	Mapper             *mapper.SharingAccessMapper
	Validator          *validator.SharingAccessValidator
	UserClient         *client.UserClient
	OrganizationClient *client.OrganizationClient
}

func NewSharingAccessService(
	uc *shared_usecase.SharingAccessUsecase,
	mapper *mapper.SharingAccessMapper,
	validator *validator.SharingAccessValidator,
	userClient *client.UserClient,
	organizationClient *client.OrganizationClient,
) *SharingAccessService {
	return &SharingAccessService{
		Usecase:            uc,
		Mapper:             mapper,
		Validator:          validator,
		UserClient:         userClient,
		OrganizationClient: organizationClient,
	}
}

func (s *SharingAccessService) BulkSave(ctx context.Context, req *bdspropb.SharingAccessBulkRequest) (*bdspropb.SharingAccessResponse, error) {
	dtos := s.Mapper.SharingAccessBulkToDomain(req)
	if err := s.Validator.ValidateBulkSave(req); err != nil {
		return nil, err
	}

	entities, err := s.Usecase.BulkSave(ctx, req.DomainId, enums.EOwnerOf(req.FromType), dtos)
	if err != nil {
		return nil, err
	}
	pb_entities := []*bdspropb.SharingAccess{}
	copier.Copy(&pb_entities, &entities)
	return &bdspropb.SharingAccessResponse{
		Data: pb_entities,
	}, nil
}
func (s *SharingAccessService) SearchProduct(ctx context.Context, req *bdspropb.SharingAccessRequest) (*bdspropb.SharingAccessResponse, error) {
	searchDTO := s.Mapper.SharingAccessToSearchDTO(req)
	if searchDTO.DomainID == 0 || searchDTO.FromType == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "domainID and domain and fromType are required",
		}
	}
	entities, err := s.Usecase.SearchForProduct(ctx, enums.EOwnerOf(req.FromType), searchDTO)
	if err != nil {
		return nil, err
	}
	pb_entities := s.Mapper.SharingAccessToSearchPb(entities)
	s.UserClient.MapUserToSharingAccess(ctx, pb_entities)
	s.OrganizationClient.MapOrganizationToSharingAccess(ctx, pb_entities)
	s.OrganizationClient.MapGroupToSharingAccess(ctx, pb_entities)
	return &bdspropb.SharingAccessResponse{
		Data: pb_entities,
	}, nil
}

func (s *SharingAccessService) SearchAsset(ctx context.Context, req *bdspropb.SharingAccessRequest) (*bdspropb.SharingAccessResponse, error) {
	searchDTO := s.Mapper.SharingAccessToSearchDTO(req)

	if err := s.Validator.ValidateSearchAsset(req); err != nil {
		return nil, err
	}
	entities, err := s.Usecase.SearchForAsset(ctx,
		enums.EOwnerOf(req.FromType),
		searchDTO,
	)
	if err != nil {
		return nil, err
	}
	result := s.Mapper.SharingAccessToSearchPb(entities)

	s.UserClient.MapUserToSharingAccess(ctx, result)
	return &bdspropb.SharingAccessResponse{
		Data: result,
	}, nil
}

func (s *SharingAccessService) SetCommission(ctx context.Context, req *bdspropb.SharingAccessByDomainRequest) (*bdspropb.SharingAccessResponse, error) {
	if err := s.Validator.ValidateSetCommission(req); err != nil {
		return nil, err
	}

	err := s.Usecase.UpdateCommission(ctx, dto.CommissionUpdate{
		ToID:       req.ToId,
		ToType:     uint64(req.ToType),
		ProductID:  req.DomainId,
		FromType:   uint64(req.FromType),
		Commission: _utils.RoundToMaxDecimals(req.Commission, 5),
	})
	if err != nil {
		return nil, err
	}

	return &bdspropb.SharingAccessResponse{
		Data: nil,
	}, nil
}

func (s *SharingAccessService) ProductDelete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// if err := s.Validator.ValidateDelete(req); err != nil {
	// 	return nil, err
	// }

	err := s.Usecase.ProductDelete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "Xóa thành công",
	}, nil
}
