package handler

import (
	"context"
	"crm/infra/client"
	"crm/infra/mapper"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type FollowService struct {
	crmpb.UnimplementedFollowServiceServer
	UC         *usecase.FollowUsecase
	userClient *client.UserClient
}

func NewFollowService(uc *usecase.FollowUsecase, userClient *client.UserClient) *FollowService {
	return &FollowService{
		UC:         uc,
		userClient: userClient,
	}
}

// @Summary Lấy danh sách người theo dõi của người dùng
// @Description Lấy danh sách những người theo dõi người dùng cụ thể
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /follow/followers [get]
func (s *FollowService) ListFollowers(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ListFollowersResponse, error) {
	followers, err := s.UC.FollowerUser(ctx)
	if err != nil {
		return nil, err
	}

	return &crmpb.ListFollowersResponse{
		Data: mapper.ListFollowersToPb(followers),
	}, nil
}

// @Summary Lấy danh sách người mà người dùng đang theo dõi
// @Description Lấy danh sách những người mà người dùng đang theo dõi
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /follow/following [get]
func (s *FollowService) ListFollowing(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ListFollowingResponse, error) {
	following, err := s.UC.FollowingUser(ctx)
	if err != nil {
		return nil, err
	}

	profiles, _ := s.userClient.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: following})

	return &crmpb.ListFollowingResponse{
		Data: profiles.Profiles,
	}, nil
}

// @Summary Thêm người vào danh sách theo dõi
// @Description Thêm một người vào danh sách theo dõi của người dùng
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID của người dùng"
// @Success 200 {string} string "Đã thêm người vào danh sách theo dõi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng hoặc người dùng"
// @Router /follow/{profileId} [post]
func (s *FollowService) AddFollow(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	follow, err := s.UC.FollowUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "success",
		Id:      follow.ID,
	}, nil
}

// @Summary Xóa người khỏi danh sách theo dõi
// @Description Xóa một người khỏi danh sách theo dõi của người dùng
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID của người dùng"
// @Success 200 {string} string "Đã xóa người khỏi danh sách theo dõi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng hoặc người dùng"
// @Router /follow/{profileId} [delete]
func (s *FollowService) RemoveFollow(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.UnfollowUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "success",
		Id:      req.Id,
	}, nil
}