package handler

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LeadService struct {
	crmpb.UnimplementedLeadServiceServer
	UC           *usecase.LeadUsecase
	UserClient   *client.UserClient
	BdsproClient *client.BdsproClient
	LeadMapper   *mapper.LeadMapper
}

func NewLeadService(uc *usecase.LeadUsecase,
	userClient *client.UserClient,
	bdsproClient *client.BdsproClient,
	leadMapper *mapper.LeadMapper,
) *LeadService {
	return &LeadService{
		UC:           uc,
		UserClient:   userClient,
		BdsproClient: bdsproClient,
		LeadMapper:   leadMapper,
	}
}

// @Summary Chi tiết khách hàng
// @Description Chi tiết khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Response 200 {object} crmpb.LeadDetailDTO
// @Response 404 {object} _routes.Except
// @Router /lead/detail/{id} [get]
func (s *LeadService) Detail(ctx context.Context, req *sharepb.IdRequest) (*crmpb.LeadDetailDTO, error) {
	result, err := s.UC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	leadPb := s.LeadMapper.LeadDetailToPb(result)
	s.BdsproClient.MapProductPbByIds(ctx, leadPb)
	s.UserClient.PbChargePersonToLeadDetail(ctx, leadPb)
	return leadPb, err
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

	leadPb := s.LeadMapper.LeadToPb(result)
	s.UserClient.PbChargePersonAndProfileToLeads(ctx, []*crmpb.LeadDTO{leadPb})

	return leadPb, err
}

// @Summary Tạo khách hàng
// @Description Tạo khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param entity body crmpb.LeadSaveRequest true "Khách hàng"
// @Param ownerOf path sharepb.OwnerOf true "Loại chủ sở hữu"
// @Param ownerId path int true "ID chủ sở hữu"
// @Security BearerAuth
// @Router /lead/{ownerOf}/{ownerId} [post]
func (s *LeadService) Create(ctx context.Context, req *crmpb.LeadSaveRequest) (*crmpb.LeadDTO, error) {
	dto := s.LeadMapper.LeadSaveToDomain(req)
	if dto.Phone == "" {
		return nil, _errors.ReturnError(service.PhoneRequired)
	}
	if dto.StageID == nil {
		return nil, _errors.ReturnError(service.StageIDRequired)
	}

	result, err := s.UC.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.LeadMapper.LeadToPb(result), nil
}

// @Summary Cập nhật khách hàng
// @Description Cập nhật khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body crmpb.LeadSaveRequest true "Khách hàng"
// @Security BearerAuth
// @Router /lead/{id} [put]
func (s *LeadService) Update(ctx context.Context, req *crmpb.LeadSaveRequest) (*crmpb.LeadDTO, error) {
	dto := s.LeadMapper.LeadSaveToDomain(req)
	if req.Id == 0 {
		return nil, _errors.ReturnError(service.IDRequired, _errors.WithPublicMessage("ID is required"))
	}

	result, err := s.UC.Update(ctx, req.Id, dto)
	return s.LeadMapper.LeadToPb(result), err
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

// @Summary Gán liên hệ cho nhân viên phụ trách
// @Description Gán liên hệ cho nhân viên phụ trách
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadAssignPersonRequest true "Gán khách hàng"
// @Security BearerAuth
// @Router /lead/assign-person [post]
func (s *LeadService) AssignPerson(ctx context.Context, req *crmpb.LeadAssignPersonRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.Assign(ctx, req.LeadId, &dto.LeadDTO{
		OwnerID:        &req.OwnerId,
		OwnerType:      enums.EOwnerOf(req.OwnerType),
		ChargePersonId: &req.ChargePersonId,
		ChargeType:     enums.EOwnerOf(req.ChargePersonType),
		AssignNote:     req.NoteAssign,
	})
	if err != nil {
		return nil, err
	}
	return s.LeadMapper.LeadToPb(result), err
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
	return s.LeadMapper.LeadToPb(result), err
}

// @Summary Chuyển trạng thái khách hàng
// @Description Chuyển trạng thái khách hàng
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadStageRequest true "Chuyển trạng thái khách hàng"
// @Security BearerAuth
// @Router /lead/stage [put]
func (s *LeadService) Stage(ctx context.Context, req *crmpb.LeadStageRequest) (*crmpb.LeadDTO, error) {
	if req.LeadId == 0 || req.StageId == 0 {
		return nil, _errors.ReturnError(service.LeadStageFieldsRequired)
	}

	result, err := s.UC.SwitchStage(ctx, req.LeadId, &req.StageId, req.StageNote)
	if err != nil {
		return nil, err
	}
	return s.LeadMapper.LeadToPb(result), nil
}

// @Summary Gán liên hệ vào CRM
// @Description Gán liên hệ vào CRM
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param body body crmpb.LeadAssignCrmRequest true "Gán khách hàng"
// @Security BearerAuth
// @Router /lead/assign-crm [post]
func (s *LeadService) AssignCrm(ctx context.Context, req *crmpb.LeadAssignCrmRequest) (*crmpb.LeadDTO, error) {
	result, err := s.UC.AssignCrm(ctx, req.ContactId, req.PersonChargeId, req.StageId, req.Note)
	if err != nil {
		return nil, err
	}
	return s.LeadMapper.LeadToPb(result), err
}

// @Summary Lấy danh sách khách hàng hiển thị trên dashboard
// @Description Lấy danh sách khách hàng hiển thị trên dashboard
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param query query crmpb.LeadManagerRequest true "Lấy danh sách khách hàng hiển thị trên dashboard"
// @Security BearerAuth
// @Param pipelineId path int true "ID của pipeline"
// @Router /lead/manage/{ownerOf}/{ownerId}/{pipelineId} [get]
func (s *LeadService) Manager(ctx context.Context, req *crmpb.LeadManagerRequest) (*crmpb.LeadManagerResponse, error) {
	dto := dto.LeadSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: 999,
		},
		LastUpdatedAt: _utils.ParseStringToTime(req.LastUpdatedAt),
		PipelineID:    &req.PipelineId,
		OwnerID:       req.OwnerId,
		OwnerOf:       base_enum.EOwnerOf(req.OwnerOf),
	}

	result, total, err := s.UC.GetManager(ctx, dto)

	managerItems := make([]*crmpb.LeadManagerItem, len(result))
	for i, lead := range result {
		managerItems[i] = s.LeadMapper.LeadManagerToPb(&lead)
	}

	s.UserClient.PbLeadToManager(ctx, managerItems)

	return &crmpb.LeadManagerResponse{
		Data:  managerItems,
		Total: uint32(total),
	}, err
}

// @Summary Lấy thời gian hết hạn bảo trì
// @Description Lấy thời gian hết hạn bảo trì
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /lead/care-expired [get]
func (s *LeadService) GetCareExpiredTime(ctx context.Context, req *sharepb.Empty) (*crmpb.CareExpiredResponse, error) {
	count, err := s.UC.GetCareExpiredTime(ctx)
	if err != nil {
		return nil, err
	}
	return &crmpb.CareExpiredResponse{
		Count: count,
	}, nil
}
