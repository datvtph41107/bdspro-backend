package grpc_service

import (
	"context"
	"crm/data"
	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/enums"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharedpb "pb/types/shared"
	"strings"
)

type FriendService struct {
	crmpb.UnimplementedFriendServiceServer
	UC         *usecase.FriendUsecase
	userClient *client.UserClient
}

func NewFriendService(uc *usecase.FriendUsecase, userClient *client.UserClient) *FriendService {
	return &FriendService{
		UC:         uc,
		userClient: userClient,
	}
}

// @Summary Lấy danh sách bạn bè
// @Description Lấy danh sách tất cả bạn bè của người dùng
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query crmpb.FriendListRequest false "Tên hoặc số điện thoại"
// @Success 200 {array} data.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /friend/list [get]
func (s *FriendService) FriendList(ctx context.Context, req *crmpb.FriendListRequest) (*crmpb.FriendProfileListResponse, error) {
	query := data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
	}
	result, total, err := s.UC.FriendList(ctx, query)
	if err != nil {
		return nil, err
	}
	pbs := mapper.ListFriendToPb(result)
	userPbs, err := s.userClient.FriendToUserPbs(ctx, pbs)
	if err != nil {
		return nil, err
	}

	return &crmpb.FriendProfileListResponse{
		Data:  userPbs,
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
	dto := mapper.FriendRequestPbToDomain(req)

	result, err := s.UC.Request(ctx, dto.ReceiverID)
	if err != nil {
		return nil, err
	}
	pb := mapper.FriendToPb(result)

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
func (s *FriendService) FriendAccept(ctx context.Context, req *sharedpb.IdRequest) (*sharedpb.SubmitResponse, error) {
	result, err := s.UC.ChangeStatus(ctx, req.Id, enums.FriendStatusAccepted)

	return &sharedpb.SubmitResponse{
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
func (s *FriendService) FriendUpdateGroup(ctx context.Context, req *crmpb.GroupUpdateRequest) (*sharedpb.SubmitResponse, error) {
	result, err := s.UC.ChangeGroup(ctx, req.RequestId, req.GroupId)

	return &sharedpb.SubmitResponse{
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
// @Router /friend/request/{requestId}/reject [delete]
func (s *FriendService) FriendReject(ctx context.Context, req *sharedpb.IdRequest) (*sharedpb.SubmitResponse, error) {
	result, err := s.UC.ChangeStatus(ctx, req.Id, enums.FriendStatusReject)

	return &sharedpb.SubmitResponse{
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
func (s *FriendService) FriendSent(ctx context.Context, req *crmpb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	query := data.FriendDTO{
		Text: &req.Text,
	}
	result, err := s.UC.RequestSent(ctx, query)
	pbs := mapper.ListFriendToPb(result)
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
func (s *FriendService) FriendReceived(ctx context.Context, req *crmpb.FriendListRequest) (*crmpb.FriendListResponse, error) {
	query := data.FriendDTO{
		Text: &req.Text,
	}
	result, err := s.UC.RequestReceived(ctx, query)
	pbs := mapper.ListFriendToPb(result)
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
func (s *FriendService) FriendInfo(ctx context.Context, req *sharedpb.Empty) (*crmpb.FriendInfoResponse, error) {
	result, err := s.UC.FriendInfo(ctx)

	return &crmpb.FriendInfoResponse{
		NumFollower:  int64(result.NumFollower),
		NumFollowing: int64(result.NumFollowing),
		NumFriend:    int64(result.NumFriend),
	}, err
}
