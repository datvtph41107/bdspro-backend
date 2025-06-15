package grpc_service

import (
	_dto "common/dto"
	"context"
	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	pb_shared "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContactService struct {
	crmpb.UnimplementedContactServiceServer
	UC         *usecase.ContactUsecase
	UserClient *client.UserClient
}

func NewContactService(uc *usecase.ContactUsecase, userClient *client.UserClient) *ContactService {
	return &ContactService{UC: uc, UserClient: userClient}
}

// @Summary Lấy danh sách liên hệ
// @Description Lấy danh sách
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ownerType path int true "Loại sở hữu (10: user, 20: nhóm, 30: organization)"
// @Router /contact/list/{ownerType} [get]
func (s *ContactService) Search(ctx context.Context, req *crmpb.ContactSearchRequest) (*crmpb.ContactListDTO, error) {

	contacts, total, err := s.UC.Search(ctx, enums.EOwnerType(req.OwnerType), dto.ContactSearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Text: req.Text,
	})
	if err != nil {
		return nil, err
	}

	contactsPb := mapper.ListContactItemToPb(contacts)
	s.UserClient.PbUserToContact(ctx, contactsPb)

	return &crmpb.ContactListDTO{
		Data:          contactsPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *ContactService) Detail(ctx context.Context, req *pb_shared.IdRequest) (*crmpb.ContactDTO, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detail not implemented")
}

// @Summary Tạo liên hệ mới
// @Description Tạo một liên hệ mới với thông tin được cung cấp
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body crmpb.ContactSaveDTO true "Thông tin cập nhật"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 409 {string} string "Liên hệ đã tồn tại"
// @Router /contact [post]
func (s *ContactService) Create(ctx context.Context, req *crmpb.ContactSaveDTO) (*crmpb.ContactDTO, error) {
	domain := mapper.PbToContactSave(req)
	if domain.FullName == "" || domain.Phone == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}
	return mapper.ContactToPb(result), nil
}

// @Summary Cập nhật thông tin liên hệ
// @Description Cập nhật thông tin liên hệ hiện tại
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param request body crmpb.ContactSaveDTO true "Thông tin cập nhật"
// @Param id path int true "ID của lời mời kết bạn"
// @Security BearerAuth
// @Router /contact/{id} [put]
func (s *ContactService) Update(ctx context.Context, req *crmpb.ContactSaveDTO) (*crmpb.ContactDTO, error) {
	domain := mapper.PbToContactSave(req)
	if domain.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}
	return mapper.ContactToPb(result), nil
}

// @Summary Xóa liên hệ
// @Description Xóa một liên hệ khỏi hệ thống
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Router /contact/{id} [delete]
func (s *ContactService) Delete(ctx context.Context, req *pb_shared.IdRequest) (*pb_shared.SubmitResponse, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb_shared.SubmitResponse{
		Id:      req.Id,
		Message: "Delete contact success",
	}, nil
}

// @Summary Đồng bộ danh bạ
// @Description Đồng bộ danh bạ
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param body body crmpb.ContactSyncRequest true "Danh bạ"
// @Security BearerAuth
// @Router /contact/sync [post]
func (s *ContactService) Sync(ctx context.Context, req *crmpb.ContactSyncRequest) (*crmpb.ContactSyncResponse, error) {
	domains := mapper.ListContactToDomain(req.Contacts)
	result, err := s.UC.Sync(ctx, enums.EOwnerType(req.OwnerType), domains)
	if err != nil {
		return nil, err
	}
	return &crmpb.ContactSyncResponse{
		Data:          mapper.ListContactToPb(result),
		TotalElements: uint32(len(result)),
	}, nil
}
