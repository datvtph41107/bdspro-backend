package handler

import (
	base_enum "base/enum"
	_dto "common/dto"
	"context"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContactService struct {
	crmpb.UnimplementedContactServiceServer
	UC            *usecase.ContactUsecase
	UserClient    *client.UserClient
	ContactMapper *mapper.ContactMapper
	FriendMapper  *mapper.FriendMapper
}

func NewContactService(
	uc *usecase.ContactUsecase,
	userClient *client.UserClient,
	contactMapper *mapper.ContactMapper,
	friendMapper *mapper.FriendMapper,
) *ContactService {
	return &ContactService{
		UC:            uc,
		UserClient:    userClient,
		ContactMapper: contactMapper,
		FriendMapper:  friendMapper,
	}
}

// @Summary Lấy danh sách liên hệ
// @Description Lấy danh sách
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param text query string false "Tên liên hệ"
// @Param ownerType path int true "Loại sở hữu (10: user, 20: nhóm, 30: organization) để 30 nhé"
// @Param friendStatus query []int false "Trạng thái bạn bè, 10: Đã gửi lời mời, 20: Bạn bè, 30: Chưa kết bạn"
// @Param following query bool false "Đang theo dõi"
// @Param roles query []int false "Nhóm quyền: 10: quản trị,20: Nhân viên, 30: khách hàng"
// @Param leaded query bool false "Gán CRM"
// @Param connected query bool false "Kết nối"
// @Param appInstalled query bool false "Đã cài app"
// @Router /contact/list/{ownerOf} [get]
func (s *ContactService) Search(ctx context.Context, req *crmpb.ContactSearchRequest) (*crmpb.ContactListDTO, error) {

	contacts, total, err := s.UC.Search(ctx, base_enum.EOwnerOf(req.OwnerOf), dto.ContactSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:         req.Text,
		FriendStatus: req.FriendStatus,
		Following:    req.Following,
		Roles:        req.Roles,
		Leaded:       req.Leaded,
		Connected:    req.Connected,
		AppInstalled: req.AppInstalled,
		OwnerId:      req.OwnerId,
	})
	if err != nil {
		return nil, err
	}

	contactsPb := s.ContactMapper.ListContactItemToPb(contacts)
	s.UserClient.PbUserToContacts(ctx, contactsPb)

	return &crmpb.ContactListDTO{
		Data:  contactsPb,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy thông tin liên hệ
// @Description Lấy thông tin liên hệ
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Router /contact/detail/{id} [get]
func (s *ContactService) Detail(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ContactDTO, error) {
	contact, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	contactPb := s.ContactMapper.ContactToPb(contact)
	s.UserClient.PbUserToContact(ctx, contactPb)

	return contactPb, nil
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
	domain := s.ContactMapper.PbToContactSave(req)
	if domain.FullName == "" || domain.Phone == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Create(ctx, domain)
	if err != nil {
		return nil, err
	}
	return s.ContactMapper.ContactToPb(result), nil
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
	domain := s.ContactMapper.PbToContactSave(req)
	if domain.FullName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "FullName and Phone are required")
	}
	result, err := s.UC.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}
	return s.ContactMapper.ContactToPb(result), nil
}

// @Summary Xóa liên hệ
// @Description Xóa một liên hệ khỏi hệ thống
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Router /contact/{id} [delete]
func (s *ContactService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
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
	domains := s.ContactMapper.ListContactToDomain(req.Contacts)
	result, err := s.UC.Sync(ctx, base_enum.EOwnerOf(req.OwnerType), domains)
	if err != nil {
		return nil, err
	}
	return &crmpb.ContactSyncResponse{
		Data:  s.ContactMapper.ListContactToPb(result),
		Total: uint32(len(result)),
	}, nil
}

// @Summary Lấy trạng thái quan hệ
// @Description Lấy trạng thái quan hệ
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Response 200 {object} crmpb.RelationShipResponse
// @Router /contact/relation-ship/{id} [get]
func (s *ContactService) RelationShip(ctx context.Context, req *sharepb.IdRequest) (*crmpb.RelationShipResponse, error) {
	result := s.UC.RelationShip(ctx, req.Id)
	return &crmpb.RelationShipResponse{
		FriendStatus: s.FriendMapper.FriendItemToPb(result.FriendStatus),
		Following:    result.Following,
		NumFollower:  result.NumFollower,
		NumFollowing: result.NumFollowing,
		NumFriend:    result.NumFriend,
	}, nil
}

// @Summary Ghi chú liên hệ
// @Description Ghi chú liên hệ
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param contactId path int true "ID của liên hệ"
// @Param note body crmpb.NoteContactDTO true "Ghi chú"
// @Security BearerAuth
// @Router /contact/note/{contactId} [put]
func (s *ContactService) Note(ctx context.Context, req *crmpb.NoteContactDTO) (*crmpb.NoteContactDTO, error) {
	err := s.UC.Note(ctx, req.Id, req.Note)
	if err != nil {
		return nil, err
	}
	return &crmpb.NoteContactDTO{
		Id:   req.Id,
		Note: req.Note,
	}, nil
}
