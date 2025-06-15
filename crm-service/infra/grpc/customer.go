package grpc_service

import (
	_dto "common/dto"
	_routes "common/routes"
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	pb_shared "pb/types/shared"

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
// @Param query query crmpb.CustomerSearchRequest false "Thông tin tìm kiếm"
// @Router /customer/list [get]
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
		Pagable: _dto.Pagable{
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
		customers[i] = mapper.CustomerToPb(&customer)
	}
	return &crmpb.CustomerListDTO{
		Data:          customers,
		TotalElements: uint32(total),
	}, err
}

// @Summary Chi tiết khách hàng
// @Description Chi tiết khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Param body body crmpb.CustomerDTO true "Khách hàng"
// @Router /customer/detail/{id} [get]
func (s *CustomerService) Detail(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.GetByID(ctx, req.Id)
	return mapper.CustomerToPb(result), err
}

// @Summary Tạo khách hàng
// @Description Tạo khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param entity body crmpb.CustomerDTO true "Khách hàng"
// @Security BearerAuth
// @Router /customer [post]
func (s *CustomerService) Create(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	entity := mapper.CustomerToDomain(req)

	if entity.OwnerID == nil ||
		entity.OwnerType == 0 ||
		entity.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ownerId and ownerType are required")
	}

	result, err := s.UC.Create(ctx, entity)
	return mapper.CustomerToPb(result), err
}

// @Summary Cập nhật khách hàng
// @Description Cập nhật khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body crmpb.CustomerDTO true "Khách hàng"
// @Security BearerAuth
// @Router /customer/{id} [put]
func (s *CustomerService) Update(ctx context.Context, req *crmpb.CustomerDTO) (*crmpb.CustomerDTO, error) {
	entity := mapper.CustomerToDomain(req)
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID is required",
		}
	}

	result, err := s.UC.Update(ctx, entity.ID, entity)
	return mapper.CustomerToPb(result), err
}

// @Summary Xóa khách hàng
// @Description Xóa khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /customer/{id} [delete]
func (s *CustomerService) Delete(ctx context.Context, req *pb_shared.IdRequest) (*pb_shared.SubmitResponse, error) {
	err := s.UC.Delete(ctx, req.Id)
	return &pb_shared.SubmitResponse{
		Id:      req.Id,
		Message: "Delete customer success",
	}, err
}

// @Summary Gán khách hàng
// @Description Gán khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.CustomerAssignRequest true "Gán khách hàng"
// @Security BearerAuth
// @Router /customer/assign [post]
func (s *CustomerService) Assign(ctx context.Context, req *crmpb.CustomerAssignRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.Assign(ctx, req.CustomerId, &dto.CustomerDTO{
		OwnerID:    &req.OwnerId,
		OwnerType:  enums.EOwnerType(req.OwnerType),
		ChargeId:   &req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		AssignNote: req.NoteAssign,
	})
	return mapper.CustomerToPb(result), err
}

// @Summary Ghi chú khách hàng
// @Description Ghi chú khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.CustomerNoteRequest true "Ghi chú khách hàng"
// @Security BearerAuth
// @Router /customer/note [post]
func (s *CustomerService) Note(ctx context.Context, req *crmpb.CustomerNoteRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.UpdateNote(ctx, req.CustomerId, &dto.CustomerDTO{
		Note: req.Note,
	})
	return mapper.CustomerToPb(result), err
}

// @Summary Chuyển trạng thái khách hàng
// @Description Chuyển trạng thái khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.CustomerStageRequest true "Chuyển trạng thái khách hàng"
// @Security BearerAuth
// @Router /customer/stage [post]
func (s *CustomerService) Stage(ctx context.Context, req *crmpb.CustomerStageRequest) (*crmpb.CustomerDTO, error) {
	result, err := s.UC.SwitchStage(ctx, req.CustomerId, &req.StageId)
	return mapper.CustomerToPb(result), err
}
