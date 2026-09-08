package handler

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	_provider "common/domain/provider"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"time"

	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContactService struct {
	crmpb.UnimplementedContactServiceServer
	UC                 *usecase.ContactUsecase
	UserClient         *client.UserClient
	ChatClient         _provider.ChatProvider
	BdsproClient       *client.BdsproClient
	OrganizationClient *client.OrganizationClient
	ContactMapper      *mapper.ContactMapper
	FriendMapper       *mapper.FriendMapper
	SyncProvider       *_utils.SyncUtil
}

func NewContactService(
	uc *usecase.ContactUsecase,
	userClient *client.UserClient,
	chatClient _provider.ChatProvider,
	bdsproClient *client.BdsproClient,
	organizationClient *client.OrganizationClient,
	contactMapper *mapper.ContactMapper,
	friendMapper *mapper.FriendMapper,
	SyncProvider *_utils.SyncUtil,
) *ContactService {
	return &ContactService{
		UC:                 uc,
		UserClient:         userClient,
		ChatClient:         chatClient,
		BdsproClient:       bdsproClient,
		OrganizationClient: organizationClient,
		ContactMapper:      contactMapper,
		FriendMapper:       friendMapper,
		SyncProvider:       SyncProvider,
	}
}

func (s *ContactService) GetInterestedContactsByProduct(
	ctx context.Context,
	req *crmpb.ContactProductInterestRequest,
) (*crmpb.ContactProductInterestResponse, error) {
	items, total, err := s.UC.GetInterestedContactsByProduct(ctx,
		&dto.ContactProductInterestedFilter{
			ProductID: req.ProductId,
			Pagable: _dto.Pagable{
				Page: req.Page,
				Size: req.Size,
			},
		},
	)
	if err != nil {
		return nil, _errors.InternalServerException("GetInterestedContactsByProduct Error: ", err.Error())
	}

	return &crmpb.ContactProductInterestResponse{
		Data:  s.ContactMapper.ContactProductInterestedToPb(items),
		Total: uint32(total),
	}, nil
}

func (s *ContactService) PinContact(ctx context.Context, req *crmpb.PinRequest) (*crmpb.PinResponse, error) {
	err := s.UC.PinContact(ctx, s.ContactMapper.PinReqPbToDTO(req))
	return &crmpb.PinResponse{Success: err == nil}, err
}

// func (s *ContactService) Unpin(ctx context.Context, req *crmpb.PinRequest) (*crmpb.PinResponse, error) {
// 	err := s.UC.Unpin(ctx, s.ContactMapper.PinReqPbToDTO(req))
// 	return &crmpb.PinResponse{Success: err == nil}, err
// }

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
		FriendStatus: _utils.ParseInt32Array(req.FriendStatus),
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

	// Nếu có conversationId, map thông tin member vào contact
	if req.ConversationId != nil && *req.ConversationId > 0 {
		// todo: map participant to contact
		// s.ChatClient.MapParticipantToContact(ctx, *req.ConversationId, contactsPb)
	}

	return &crmpb.ContactListDTO{
		Data:  contactsPb,
		Total: uint32(total),
	}, nil
}

func (s *ContactService) ListMe(ctx context.Context, req *sharepb.SyncRequest) (*crmpb.ContactListDTO, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyContactMe, 0)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &crmpb.ContactListDTO{}, nil
	}
	contacts, total, err := s.UC.Search(ctx, base_enum.EOwnerOfMember, dto.ContactSearchDTO{})
	if err != nil {
		return nil, err
	}

	contactsPb := s.ContactMapper.ListContactItemToPb(contacts)
	s.UserClient.PbUserToContacts(ctx, contactsPb)

	// Nếu có conversationId, map thông tin member vào contact
	// if req.ConversationId != nil && *req.ConversationId > 0 {
	// 	s.ChatClient.MapParticipantToContact(ctx, *req.ConversationId, contactsPb)
	// }

	// if req.EndId == nil && len(contacts) > 0 {
	// 	s.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("u:contact:%d", profileId), contacts[0].UpdatedAt.UnixMilli())
	// }

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
func (s *ContactService) Detail(ctx context.Context, req *sharepb.SyncRequest) (*crmpb.ContactDTO, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyContactDetail, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &crmpb.ContactDTO{}, nil
	}

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
	s.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("time:contact:%d", req.Id), time.Now().UnixMilli())
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
	s.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("time:contact:%d", req.Id), time.Now().UnixMilli())
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
		FriendStatus: &sharepb.FriendItem{
			Id:         result.FriendStatus.ID,
			ReceiverId: result.FriendStatus.ReceiverID,
			CreatedBy:  &result.FriendStatus.CreatedBy,
			Status:     result.FriendStatus.Status,
		},
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
	s.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("time:contact:%d", req.Id), time.Now().UnixMilli())
	return &crmpb.NoteContactDTO{
		Id:   req.Id,
		Note: req.Note,
	}, nil
}

// @Summary Kiểm tra số điện thoại liên hệ
// @Description Kiểm tra xem số điện thoại đã tồn tại trong danh bạ chưa
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param phone path string true "Số điện thoại cần kiểm tra"
// @Param ownerOf path int true "Loại sở hữu (10: user, 20: nhóm, 30: organization)"
// @Param ownerId path int true "ID của owner"
// @Security BearerAuth
// @Success 200 {object} crmpb.CheckPhoneResponse
// @Router /contact/phone/{ownerOf}/{ownerId}/{phone} [get]
func (s *ContactService) CheckPhone(ctx context.Context, req *crmpb.CheckPhoneRequest) (*crmpb.ContactDTO, error) {
	if req.Phone == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Số điện thoại không được để trống")
	}

	contact, existed, err := s.UC.SearchByPhone(ctx, base_enum.EOwnerOf(req.OwnerOf), req.OwnerId, req.Phone)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, nil
	}

	contactPb := s.ContactMapper.ContactToPb(contact)
	s.UserClient.PbUserToContact(ctx, contactPb)

	return contactPb, nil
}

// @Summary Lấy danh sách contact quan tâm tới sản phẩm
// @Description Lấy danh sách contact quan tâm tới một sản phẩm cụ thể
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param productId path int true "ID của sản phẩm"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Security BearerAuth
// @Router /contact/by-product/{id} [get]
func (s *ContactService) GetListContactByProductId(ctx context.Context, req *sharepb.SyncRequest) (*crmpb.ContactListDTO, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductContacts, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated && req.Page == 0 {
		return &crmpb.ContactListDTO{}, nil
	}
	contacts, total, err := s.UC.GetListContactByProductId(ctx, req.Id, dto.ContactSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	// if req.Page == 0 && len(contacts) > 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if contacts[0].ContactProductUpdatedAt != nil {
	// 		time = contacts[0].ContactProductUpdatedAt.UnixMilli()
	// 	}
	// 	s.SyncProvider.PutTimestamp(ctx, key, time)
	// }
	contactsPb := s.ContactMapper.ListContactItemToPb(contacts)
	s.UserClient.PbUserToContacts(ctx, contactsPb)

	return &crmpb.ContactListDTO{
		Data:  contactsPb,
		Total: uint32(total),
	}, nil
}

// @Summary Cập nhật danh sách sản phẩm quan tâm cho contact
// @Description Cập nhật danh sách sản phẩm quan tâm cho một contact (customer)
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param contactId path int true "ID của contact"
// @Param body body crmpb.UpdateProductIdsRequest true "Danh sách product IDs"
// @Security BearerAuth
// @Router /contact/{contactId}/products [put]
func (s *ContactService) UpdateProductIds(ctx context.Context, req *crmpb.UpdateProductIdsRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.UpdateProductIds(ctx, req.ContactId, req.ProductIds)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.ContactId,
		Message: "Cập nhật sản phẩm quan tâm thành công",
	}, nil
}

// @Summary Lấy danh sách sản phẩm quan tâm của contact
// @Description Lấy danh sách sản phẩm quan tâm của một contact với phân trang
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param contactId path int true "ID của contact"
// @Param page query int false "Số trang (mặc định: 1)"
// @Param size query int false "Số lượng sản phẩm trên mỗi trang (mặc định: 20)"
// @Security BearerAuth
// @Router /contact/detail/{contactId}/products [get]
func (s *ContactService) GetContactProducts(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.SearchResponse, error) {
	productIds, total, err := s.UC.GetContactProducts(ctx, req.Id, _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	})
	if err != nil {
		return nil, err
	}

	print("productIds", len(productIds))
	// Lấy chi tiết products từ bdspro service
	products, err := s.BdsproClient.GetProductByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	return &bdspropb.SearchResponse{
		Data:          products,
		TotalElements: total,
	}, nil
}

// @Summary Lấy danh sách contact kèm thông tin member
// @Description Lấy danh sách contact (không check friend) kèm thông tin profile và member (deal, organization, group) nếu có
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query string false "Tên hoặc số điện thoại"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param dealId query int false "ID của thương vụ"
// @Param organizationId query int false "ID của tổ chức"
// @Param groupId query int false "ID của nhóm"
// @Success 200 {object} crmpb.FriendListResponse "Danh sách contact"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /contact/with-member [get]
func (s *ContactService) ContactListWithMember(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.ContactListDTO, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, errors.New("profileId not found")
	}

	// Lấy danh sách contact (không check friend status)
	contacts, total, err := s.UC.Search(ctx, base_enum.EOwnerOfMember, dto.ContactSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:         req.Text,
		OwnerId:      profileId,
		FriendStatus: _utils.ParseInt32Array(req.FriendStatus),
	})
	if err != nil {
		return nil, err
	}

	// Convert contact sang ContactDTO
	contactsPb := s.ContactMapper.ListContactItemToPb(contacts)

	// Map profile info vào contact
	s.UserClient.PbUserToContacts(ctx, contactsPb)

	// Convert ContactDTO sang Friend format để trả về (vì response type là FriendListResponse)
	// Chỉ lấy contact + profile, không check friend status
	contactPbs := make([]*sharepb.ContactDTO, 0, len(contactsPb))
	for _, contact := range contactsPb {
		friend := s.contactToFriend(ctx, contact, profileId)
		contactPbs = append(contactPbs, friend)
	}

	// Map member info nếu có (deal, organization, group)
	if req.DealId != nil {
		s.BdsproClient.MapDealMemberToContacts(ctx, *req.DealId, contactPbs)
	}
	if req.OrganizationId != nil {
		if err := s.UserClient.MapOrganizationMembersToContacts(ctx, *req.OrganizationId, contactPbs); err != nil {
			return nil, err
		}
	}
	if req.GroupId != nil {
		s.OrganizationClient.MapGroupMemberToContact(ctx, *req.GroupId, contactPbs)
	}

	// Map productUser nếu có shareProductId
	if req.ShareProductId != nil && *req.ShareProductId > 0 {
		productUsers, err := s.BdsproClient.GetProductUsersByProductId(ctx, *req.ShareProductId)
		if err == nil && len(productUsers) > 0 {
			// Tạo map profileId -> ProductUser để map nhanh
			productUserMap := make(map[uint64]*sharepb.ProductUser)
			for _, productUser := range productUsers {
				if productUser.OriginProfileId != nil {
					productUserMap[*productUser.OriginProfileId] = productUser
				}
			}

			// Map productUser vào contact theo profileId
			for _, contact := range contactPbs {
				if contact.OriginProfileId != nil {
					if productUser, ok := productUserMap[*contact.OriginProfileId]; ok {
						contact.ProductUser = productUser
					}
				}
			}
		}
	}

	return &crmpb.ContactListDTO{
		Data:  contactPbs,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách contact là bạn bè
// @Description Lấy danh sách contact mà là bạn bè (đã chấp nhận lời mời kết bạn)
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query string false "Tên hoặc số điện thoại"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} crmpb.ContactListDTO "Danh sách contact là bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /contact/friends [get]
func (s *ContactService) ContactListFriends(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.ContactListDTO, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, errors.New("profileId not found")
	}

	// Lấy danh sách contact với filter FriendStatus = 20 (Accepted - bạn bè)
	contacts, total, err := s.UC.Search(ctx, base_enum.EOwnerOfMember, dto.ContactSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:         req.Text,
		FriendStatus: []int32{20}, // Chỉ lấy contact là bạn bè (FriendStatusAccepted = 20)
		OwnerId:      profileId,
	})
	if err != nil {
		return nil, err
	}

	// Convert contact sang ContactDTO
	contactsPb := s.ContactMapper.ListContactItemToPb(contacts)

	// Map profile info vào contact
	s.UserClient.PbUserToContacts(ctx, contactsPb)

	return &crmpb.ContactListDTO{
		Data:  contactsPb,
		Total: uint32(total),
	}, nil
}

// contactToFriend converts ContactDTO to Friend format (chỉ lấy contact + profile, không check friend)
func (s *ContactService) contactToFriend(ctx context.Context, contact *sharepb.ContactDTO, currentProfileId uint64) *sharepb.ContactDTO {
	friend := &crmpb.Friend{
		CreatedBy: currentProfileId,
		Status:    0, // Không check friend status
		CreatedAt: contact.CreatedAt,
	}

	// Set receiverId từ profileId của contact
	if contact.ProfileId != nil {
		friend.ReceiverId = *contact.ProfileId
		friend.ReceiverUser = contact.ProfileInfo
	}

	// Set createdUser từ profile của user hiện tại
	// currentProfiles, _ := s.UserClient.GetProfileByIds(ctx, []uint64{currentProfileId})
	// if len(currentProfiles) > 0 {
	// 	friend.CreatedUser = currentProfiles[0]
	// }

	return contact
}
