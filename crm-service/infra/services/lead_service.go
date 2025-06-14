package grpc_service

import (
	_dto "common/dto"
	_routes "common/routes"
	"context"
	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LeadService struct {
	crmpb.UnimplementedLeadServiceServer
	UC         *usecase.LeadUsecase
	UserClient *client.UserClient
}

func NewLeadService(uc *usecase.LeadUsecase, userClient *client.UserClient) *LeadService {
	return &LeadService{UC: uc, UserClient: userClient}
}

// @Summary Lấy danh sách khách hàng
// @Description Lấy danh sách
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query crmpb.LeadSearchRequest false "Thông tin tìm kiếm"
// @Router /lead/list [get]
func (s *LeadService) Search(ctx context.Context, req *crmpb.LeadSearchRequest) (*crmpb.LeadListDTO, error) {
	steps := make([]enums.EStep, len(req.Steps))
	for i, step := range req.Steps {
		steps[i] = enums.EStep(step)
	}
	sources := make([]enums.ESourceLead, len(req.Sources))
	for i, source := range req.Sources {
		sources[i] = enums.ESourceLead(source)
	}
	dto := dto.LeadSearchDTO{
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

	leads := make([]*crmpb.LeadDTO, len(result))
	for i, lead := range result {
		leads[i] = mapper.LeadToPb(&lead)
	}
	return &crmpb.LeadListDTO{
		Data:          leads,
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
// @Param body body crmpb.LeadDTO true "Khách hàng"
// @Router /lead/detail/{id} [get]
func (s *LeadService) Detail(ctx context.Context, req *sharepb.IdRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.GetByID(ctx, req.Id)
	return mapper.LeadToPb(result), err
}

// @Summary Lấy khách hàng theo liên hệ
// @Description Lấy khách hàng theo liên hệ
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /lead/by-contact/{id} [get]
func (s *LeadService) GetByContact(ctx context.Context, req *sharepb.IdRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.GetByContact(ctx, req.Id)
	if result == nil {
		return nil, status.Errorf(codes.NotFound, "Lead not found")
	}

	leadPb := mapper.LeadToPb(result)
	s.UserClient.PbUserToLead(ctx, leadPb)

	return leadPb, err
}

// @Summary Tạo khách hàng
// @Description Tạo khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param entity body crmpb.LeadDTO true "Khách hàng"
// @Security BearerAuth
// @Router /lead [post]
func (s *LeadService) Create(ctx context.Context, req *crmpb.LeadDTO) (*crmpb.LeadDTO, error) {
	entity := mapper.LeadToDomain(req)

	// if entity.OwnerID == nil ||
	// 	entity.OwnerType == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "ownerId and ownerType are required")
	// }

	_, err := s.UC.Create(ctx, entity)
	// return mapper.LeadToPb(result), err
	return nil, err
}

// @Summary Cập nhật khách hàng
// @Description Cập nhật khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body crmpb.LeadDTO true "Khách hàng"
// @Security BearerAuth
// @Router /lead/{id} [put]
func (s *LeadService) Update(ctx context.Context, req *crmpb.LeadDTO) (*crmpb.LeadDTO, error) {
	entity := mapper.LeadToDomain(req)
	if req.Id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID is required",
		}
	}

	result, err := s.UC.Update(ctx, entity.ID, entity)
	return mapper.LeadToPb(result), err
}

// @Summary Xóa khách hàng
// @Description Xóa khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /lead/{id} [delete]
func (s *LeadService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.Delete(ctx, req.Id)
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Delete lead success",
	}, err
}

// @Summary Gán khách hàng
// @Description Gán khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadAssignRequest true "Gán khách hàng"
// @Security BearerAuth
// @Router /lead/assign [post]
func (s *LeadService) Assign(ctx context.Context, req *crmpb.LeadAssignRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.Assign(ctx, req.LeadId, &dto.LeadDTO{
		OwnerID:    &req.OwnerId,
		OwnerType:  enums.EOwnerType(req.OwnerType),
		ChargeId:   &req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		AssignNote: req.NoteAssign,
	})
	return mapper.LeadToPb(result), err
}

// @Summary Ghi chú khách hàng
// @Description Ghi chú khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadNoteRequest true "Ghi chú khách hàng"
// @Security BearerAuth
// @Router /lead/note [post]
func (s *LeadService) Note(ctx context.Context, req *crmpb.LeadNoteRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.UpdateNote(ctx, req.LeadId, &dto.LeadDTO{
		Note: req.Note,
	})
	return mapper.LeadToPb(result), err
}

// @Summary Chuyển trạng thái khách hàng
// @Description Chuyển trạng thái khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadStageRequest true "Chuyển trạng thái khách hàng"
// @Security BearerAuth
// @Router /lead/stage [post]
func (s *LeadService) Stage(ctx context.Context, req *crmpb.LeadStageRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.SwitchStage(ctx, req.LeadId, &req.StageId)
	return mapper.LeadToPb(result), err
}
