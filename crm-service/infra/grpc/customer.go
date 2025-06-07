package grpc_service

import (
	_dto "common/dto"
	_models "common/models"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	crmpb.UnimplementedCustomerServiceServer
	UC *usecase.CustomerUsecase
}

func NewCustomerService(uc *usecase.CustomerUsecase) *CustomerService {
	return &CustomerService{UC: uc}
}

// @Summary Lấy danh sách khách hàng
// @Description Lấy danh sách
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ownerId path int true "ID của chủ sở hữu"
// @Param ownerType path enums.EOwnerType true "Loại chủ sở hữu"
// @Param dto query dto.CustomerSearchDTO true "Thông tin tìm kiếm"
// @Router /customer/list/{ownerId}/{ownerType} [get]
func (s *CustomerService) Search(ctx context.Context, req *crmpb.CustomerSearchRequest) (*crmpb.CustomerListDTO, error) {
	steps := make([]enums.EStep, len(req.Steps))
	for i, step := range req.Steps {
		steps[i] = enums.EStep(step)
	}
	sources := make([]enums.ESourceCustomer, len(req.Sources))
	for i, source := range req.Sources {
		sources[i] = enums.ESourceCustomer(source)
	}
	dto := dto.CustomerSearchDTO{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Phone:      req.Phone,
		FullName:   req.FullName,
		Steps:      steps,
		StageIDs:   req.StageIds,
		Sources:    sources,
		ChargeIds:  req.ChargeIds,
		ChargeType: enums.EOwnerType(req.ChargeType),
	}

	result, total, err := s.UC.Search(ctx, req.OwnerId, enums.EOwnerType(req.OwnerType), dto)

	customers := make([]*crmpb.CustomerDTO, len(result))
	for i, customer := range result {
		customers[i] = s.EntityToPb(&customer)
	}
	return &crmpb.CustomerListDTO{
		Data:          customers,
		TotalElements: uint32(total),
	}, err
}
func (s *CustomerService) Detail(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.GetByID(ctx, req.Id)
	return s.EntityToPb(result), err
}
func (s *CustomerService) Create(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	entity := s.PbToEntity(req)

	if entity.OwnerID == nil ||
		entity.OwnerType == 0 ||
		entity.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ownerId and ownerType are required")
	}

	result, err := s.UC.Create(ctx, entity)
	return s.EntityToPb(result), err
}
func (s *CustomerService) Update(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	entity := s.PbToEntity(req)
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID is required",
		}
	}

	result, err := s.UC.Update(ctx, entity.ID, entity)
	return s.EntityToPb(result), err
}
func (s *CustomerService) Delete(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	err := s.UC.Delete(ctx, req.Id)
	return nil, err
}
func (s *CustomerService) Assign(ctx context.Context, req *crmpb.CustomerAssignRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.Assign(ctx, req.CustomerId, &dto.CustomerDTO{
		OwnerID:    &req.OwnerId,
		OwnerType:  enums.EOwnerType(req.OwnerType),
		ChargeId:   &req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		AssignNote: req.NoteAssign,
	})
	return s.EntityToPb(result), err
}
func (s *CustomerService) Note(ctx context.Context, req *crmpb.CustomerNoteRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.UpdateNote(ctx, req.CustomerId, &dto.CustomerDTO{
		Note: req.Note,
	})
	return s.EntityToPb(result), err
}
func (s *CustomerService) Stage(ctx context.Context, req *crmpb.CustomerStageRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.SwitchStage(ctx, req.CustomerId, &req.StageId)
	return s.EntityToPb(result), err
}

// ------------------------------------------------------------
func (s *CustomerService) PbToEntity(req *crmpb.CustomerDTO) *domain.CustomerEntity {
	return &domain.CustomerEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		ContactID:  req.ContactId,
		FullName:   req.FullName,
		Phone:      req.Phone,
		Birthday:   _utils.ParseStringToTime(req.Birthday),
		Source:     enums.ESourceCustomer(req.Source),
		Email:      req.Email,
		Address:    req.Address,
		Note:       req.Note,
		AssignNote: req.AssignNote,
		PipelineID: req.PipelineId,
		StageID:    req.StageId,
		ChargeId:   req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		OwnerID:    req.OwnerId,
		OwnerType:  enums.EOwnerType(req.OwnerType),
	}
}

func (s *CustomerService) EntityToPb(entity *domain.CustomerEntity) *crmpb.CustomerDTO {
	return &crmpb.CustomerDTO{
		Id:         entity.ID,
		ContactId:  entity.ContactID,
		FullName:   entity.FullName,
		Phone:      entity.Phone,
		Source:     int32(entity.Source),
		Email:      entity.Email,
		Address:    entity.Address,
		Birthday:   _utils.FormatTimeToString(entity.Birthday),
		Note:       entity.Note,
		AssignNote: entity.AssignNote,
		PipelineId: entity.PipelineID,
		StageId:    entity.StageID,
		ChargeId:   entity.ChargeId,
		ChargeType: int32(entity.ChargeType),
		OwnerId:    entity.OwnerID,
		OwnerType:  int32(entity.OwnerType),
		// CreatedAt:  _utils.FormatTimeToString(entity.CreatedAt),
	}
}
