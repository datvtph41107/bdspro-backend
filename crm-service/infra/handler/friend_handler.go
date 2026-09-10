package handler

import (
	_dto "common/domain/dto"
	"context"
	"errors"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strings"

	"crm/data"
	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/enums"
	"crm/internal/usecase"
)

type FriendService struct {
	crmpb.UnimplementedFriendServiceServer
	UC                 *usecase.FriendUsecase
	ContactUC          *usecase.ContactUsecase
	userClient         *client.UserClient
	organizationClient *client.OrganizationClient
	FriendMapper       *mapper.FriendMapper
	ContactMapper      *mapper.ContactMapper
}

func NewFriendService(uc *usecase.FriendUsecase,
	contactUC *usecase.ContactUsecase,
	userClient *client.UserClient,
	organizationClient *client.OrganizationClient,
	friendMapper *mapper.FriendMapper,
	contactMapper *mapper.ContactMapper,
) *FriendService {
	return &FriendService{
		UC:                 uc,
		ContactUC:          contactUC,
		userClient:         userClient,
		organizationClient: organizationClient,
		FriendMapper:       friendMapper,
		ContactMapper:      contactMapper,
	}
}

// @Summary Lấy danh sách bạn bè
// @Description Lấy danh sách tất cả bạn bè của người dùng. trong đó truyền: organizationId nếu muốn xem bạn bè có thuộc tổ chức không, groupId nếu muốn xem bạn bè có thuộc nhóm không, dealId nếu muốn xem bạn bè có thuộc thương vụ không
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query commonpb.FriendListRequest false "Tên hoặc số điện thoại"
// @Success 200 {array} data.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/list [get]
func (s *FriendService) FriendList(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendProfileListResponse, error) {
	query := data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	result, total, err := s.UC.FriendList(ctx, query)
	if err != nil {
		return nil, err
	}
	pbs := s.FriendMapper.ListFriendToPb(result)
	if req.DealId != nil {
		s.organizationClient.MapDealMemberToFriend(ctx, *req.DealId, pbs)
	}
	userPbs, _ := s.userClient.FriendToUserPbs(ctx, pbs)
	// s.userClient.PbUserToReceiverFriends(ctx, pbs)

	if req.OrganizationId != nil {
		userIds := make([]uint64, 0, len(userPbs))
		for _, user := range userPbs {
			userIds = append(userIds, user.Id)
		}
		checkUsers, err := s.userClient.CheckOrganizationMembers(ctx, *req.OrganizationId, userIds)
		if err != nil {
			return nil, err
		}
		userMap := make(map[uint64]*userpb.OrganizationMemberCheck)
		for _, user := range checkUsers {
			userMap[user.ProfileId] = user
		}
		for _, user := range userPbs {
			if u, ok := userMap[user.Id]; ok && u != nil {
				user.IsMember = u.IsMember
				user.JoinedAt = u.JoinedAt
			}
		}
	}
	if req.GroupId != nil {
		// s.organizationClient.MapGroupMemberToFriend(ctx, *req.GroupId, userPbs)
	}
	if req.DealId != nil {
		s.organizationClient.MapDealMemberToFriend(ctx, *req.DealId, pbs)
	}
	return &crmpb.FriendProfileListResponse{
		Data:  userPbs,
		Total: int32(total),
	}, err
}

// @Summary Lấy danh sách bạn bè đi kèm thông tin bạn bè có thuộc thương vụ không
// @Description Truyền dealId để xem có phải thành viên thương vụ không, organizationId để xem có phải thành viên tổ chức không, groupId để xem có phải thành viên nhóm không
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query crmpb.FriendListRequest false "Tên hoặc số điện thoại"
// @Success 200 {array} data.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/with-member [get]
// @Deprecated: use ContactListWithMember instead
func (s *FriendService) FriendListWithMember(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	// if req.DealId == nil {
	// 	return nil, errors.New("dealId is required")
	// }
	query := data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	result, total, err := s.UC.FriendList(ctx, query)
	if err != nil {
		return nil, err
	}
	pbs := s.FriendMapper.ListFriendToPb(result)
	s.userClient.PbUserToReceiverFriends(ctx, pbs)

	return &crmpb.FriendListResponse{
		Data:  pbs,
		Total: int32(total),
	}, err
}

// @Summary Lấy danh sách bạn bè đi kèm thông tin bạn bè có thuộc nhóm không
// @Description Lấy danh sách tất cả bạn bè của người dùng đi kèm thông tin bạn bè có thuộc nhóm không
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query crmpb.FriendListRequest false "Tên hoặc số điện thoại"
// @Success 200 {array} data.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/with-lead [get]
func (s *FriendService) FriendListWithLead(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	if req.DealId == nil {
		return nil, errors.New("dealId is required")
	}
	query := data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	result, total, err := s.UC.FriendList(ctx, query)
	if err != nil {
		return nil, err
	}
	pbs := s.FriendMapper.ListFriendToPb(result)
	s.userClient.PbUserToReceiverFriends(ctx, pbs)
	s.organizationClient.MapDealMemberToFriend(ctx, *req.DealId, pbs)
	return &crmpb.FriendListResponse{
		Data:  pbs,
		Total: int32(total),
	}, err
}

// @Summary Lấy danh sách bạn bè đi kèm thông tin bạn bè có thuộc tổ chức không
// @Description Lấy danh sách tất cả bạn bè của người dùng đi kèm thông tin bạn bè có thuộc tổ chức không
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query crmpb.FriendListRequest false "Tên hoặc số điện thoại"
// @Success 200 {array} data.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/with-organization [get]
func (s *FriendService) FriendListWithOrganization(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	if req.OrganizationId == nil {
		return nil, errors.New("organizationId is required")
	}
	query := data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	result, total, err := s.UC.FriendList(ctx, query)
	if err != nil {
		return nil, err
	}
	pbs := s.FriendMapper.ListFriendToPb(result)
	s.userClient.PbUserToFriends(ctx, pbs, false)
	return &crmpb.FriendListResponse{
		Data:  pbs,
		Total: int32(total),
	}, err
}

// @Summary Gửi lời mời kết bạn
// @Description Gửi lời mời kết bạn đến người dùng khác
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receiverId path int true "ID của người nhận lời mời"
// @Success 200 {string} string "Lời mời đã được gửi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người nhận"
// @Router /friend/request/{receiverId} [post]
func (s *FriendService) FriendRequest(ctx context.Context, req *crmpb.FriendRequestDTO) (*crmpb.FriendRequestResponse, error) {
	dto := s.FriendMapper.FriendRequestPbToDomain(req)

	result, err := s.UC.Request(ctx, dto.ReceiverID)
	if err != nil {
		return nil, err
	}
	pb := s.FriendMapper.FriendToPb(result)

	return &crmpb.FriendRequestResponse{
		Data:   pb,
		Status: int32(result.Status),
	}, err
}

// @Summary Đồng ý kết bạn
// @Description Đồng ý lời mời kết bạn đang chờ
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Success 200 {string} string "Lời mời đã được chấp nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /friend/request/{requestId}/accept [put]
func (s *FriendService) FriendAccept(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	result, err := s.UC.ChangeStatus(ctx, req.Id, enums.FriendStatusAccepted)

	return &sharepb.SubmitResponse{
		Id:      uint64(result),
		Message: "Lời mời đã được chấp nhận",
	}, err
}

// @Summary Gán nhóm bạn bè
// @Description Gán nhóm cho bạn bè
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Param groupId path int true "ID của nhóm"
// @Success 200 {string} string "Lời mời đã được chấp nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /friend/group/{requestId}/{groupId} [put]
func (s *FriendService) FriendUpdateGroup(ctx context.Context, req *crmpb.GroupUpdateRequest) (*sharepb.SubmitResponse, error) {
	result, err := s.UC.ChangeGroup(ctx, req.RequestId, req.GroupId)

	return &sharepb.SubmitResponse{
		Id:      uint64(result),
		Message: "Đổi nhóm thành công",
	}, err
}

// @Summary Hủy kết bạn
// @Description Hủy lời mời kết bạn đang chờ
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Success 200 {string} string "Lời mời đã bị từ chối"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /friend/reject/{requestId} [delete]
func (s *FriendService) FriendReject(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	result, err := s.UC.ChangeStatus(ctx, req.Id, enums.FriendStatusReject)

	return &sharepb.SubmitResponse{
		Id:      uint64(result),
		Message: "Lời mời đã bị từ chối",
	}, err
}

// @Summary Lấy danh sách các lời mời kết bạn đã gửi
// @Description Lấy danh sách các lời mời kết bạn mà người dùng đã gửi
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} data.FriendDTO "Danh sách lời mời đã gửi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/request/sent [get]
func (s *FriendService) FriendSent(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	query := data.FriendDTO{
		Text: &req.Text,
	}
	result, err := s.UC.RequestSent(ctx, query)
	pbs := s.FriendMapper.ListFriendToPb(result)
	s.userClient.PbUserToFriends(ctx, pbs, false)

	return &crmpb.FriendListResponse{
		Data:  pbs,
		Total: int32(len(pbs)),
	}, err
}

// @Summary Lấy danh sách các lời mời kết bạn đã nhận
// @Description Lấy danh sách các lời mời kết bạn mà người dùng đã nhận
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} data.FriendDTO "Danh sách lời mời đã nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/request/received [get]
func (s *FriendService) FriendReceived(ctx context.Context, req *sharepb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	query := data.FriendDTO{
		Text: &req.Text,
	}
	result, err := s.UC.RequestReceived(ctx, query)
	pbs := s.FriendMapper.ListFriendToPb(result)
	s.userClient.PbUserToFriends(ctx, pbs, true)

	return &crmpb.FriendListResponse{
		Data:  pbs,
		Total: int32(len(pbs)),
	}, err
}

// @Summary Thông tin số người theo dõi, đang theo dõi và số bạn bè
// @Description Thông tin số người theo dõi, đang theo dõi và số bạn bè
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} data.FollowInfoRes "Số lượng theo dõi và bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/info [get]
func (s *FriendService) FriendInfo(ctx context.Context, req *sharepb.Empty) (*crmpb.FriendInfoResponse, error) {
	result, err := s.UC.FriendInfo(ctx)

	return &crmpb.FriendInfoResponse{
		NumFollower:  int64(result.NumFollower),
		NumFollowing: int64(result.NumFollowing),
		NumFriend:    int64(result.NumFriend),
	}, err
}

// @Summary Hủy kết bạn
// @Description Hủy kết bạn với người dùng khác (không phải từ chối lời mời)
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của người bạn cần hủy kết bạn"
// @Success 200 {string} string "Đã hủy kết bạn thành công"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy quan hệ bạn bè"
// @Router /friend/cancel/{id} [delete]
func (s *FriendService) CancelFriendRequest(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.CancelRequest(ctx, req.Id)

	return &sharepb.SubmitResponse{
		Message: "Đã hủy kết bạn thành công",
		Id:      req.Id,
	}, err
}

func (s *FriendService) FriendNote(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// result, err := s.UC.Note(ctx, req.Id)

	// return &sharepb.SubmitResponse{
	// 	Message: "Lời mời đã bị hủy",
	// 	Id:      req.Id,
	// }, err
	return nil, nil
}
