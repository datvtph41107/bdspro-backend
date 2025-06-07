package grpc_service

import (
	_dto "common/dto"
	_models "common/models"
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

type ContactService struct {
	crmpb.UnimplementedContactServiceServer
	UC *usecase.ContactUsecase
}

func NewContactService(uc *usecase.ContactUsecase) *ContactService {
	return &ContactService{UC: uc}
}

// @Summary Lấy danh sách liên hệ
// @Description Lấy danh sách
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /contact [get]
func (s *ContactService) Search(ctx context.Context, req *crmpb.ContactSearchRequest) (*crmpb.ContactListDTO, error) {
	contacts, total, err := s.UC.Search(ctx, dto.ContactSearchDTO{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Text: req.Text,
	})
	if err != nil {
		return nil, err
	}

	contactsPb := make([]*crmpb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = s.DomainToPb(&contact)
	}

	return &crmpb.ContactListDTO{
		Data:          contactsPb,
		TotalElements: uint32(total),
	}, nil
}

func (s *ContactService) Detail(ctx context.Context, req *crmpb.ContactDTO) (*crmpb.ContactDTO, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detail not implemented")
}

// @Summary Tạo liên hệ mới
// @Description Tạo một liên hệ mới với thông tin được cung cấp
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.ContactEntity true "Thông tin cập nhật"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 409 {string} string "Liên hệ đã tồn tại"
// @Router /contact [post]
func (s *ContactService) Create(ctx context.Context, req *crmpb.ContactDTO) (*crmpb.ContactDTO, error) {
	domain := s.PbToDomain(req)
	if domain.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}
	return s.DomainToPb(result), nil
}

// @Summary Cập nhật thông tin liên hệ
// @Description Cập nhật thông tin liên hệ hiện tại
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param request body domain.GroupEntity true "Thông tin cập nhật"
// @Param id path int true "ID của lời mời kết bạn"
// @Security BearerAuth
// @Router /contact/{id} [put]
func (s *ContactService) Update(ctx context.Context, req *crmpb.ContactDTO) (*crmpb.ContactDTO, error) {
	domain := s.PbToDomain(req)
	if domain.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}
	return s.DomainToPb(result), nil
}

// @Summary Xóa liên hệ
// @Description Xóa một liên hệ khỏi hệ thống
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Router /contact/{id} [delete]
func (s *ContactService) Delete(ctx context.Context, req *crmpb.ContactDTO) (*crmpb.ContactDTO, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return s.DomainToPb(nil), nil
}

// --map--

func (s *ContactService) DomainToPb(contact *domain.ContactEntity) *crmpb.ContactDTO {
	return &crmpb.ContactDTO{
		Id:        contact.ID,
		FullName:  contact.FullName,
		Phone:     contact.Phone,
		OwnerId:   contact.OwnerID,
		OwnerType: int32(contact.OwnerType),
		CreatedAt: _utils.FormatTimeToString(contact.CreatedAt),
		CreatedBy: contact.CreatedBy,
	}
}

func (s *ContactService) PbToDomain(pb *crmpb.ContactDTO) *domain.ContactEntity {
	return &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName:  pb.FullName,
		Phone:     pb.Phone,
		OwnerID:   pb.OwnerId,
		OwnerType: enums.EOwnerType(pb.OwnerType),
	}
}
