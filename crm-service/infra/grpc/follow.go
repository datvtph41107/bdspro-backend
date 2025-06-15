package grpc_service

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	pb_shared "pb/types/shared"
)

type FollowService struct {
	UC *usecase.FollowUsecase
	crmpb.UnimplementedFollowServiceServer
}

func NewFollowService(uc *usecase.FollowUsecase) *FollowService {
	return &FollowService{
		UC: uc,
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
func (s *FollowService) ListFollowers(ctx context.Context, req *pb_shared.IdRequest) (*crmpb.ListFollowersResponse, error) {
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
func (s *FollowService) ListFollowing(ctx context.Context, req *pb_shared.IdRequest) (*crmpb.ListFollowingResponse, error) {
	following, err := s.UC.FollowingUser(ctx)
	if err != nil {
		return nil, err
	}

	return &crmpb.ListFollowingResponse{
		Data: mapper.ListFollowersToPb(following),
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
// @Router /follow/{id} [post]
func (s *FollowService) AddFollow(ctx context.Context, req *pb_shared.IdRequest) (*pb_shared.SubmitResponse, error) {
	follow, err := s.UC.FollowUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_shared.SubmitResponse{
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
// @Router /follow/{id} [delete]
func (s *FollowService) RemoveFollow(ctx context.Context, req *pb_shared.IdRequest) (*pb_shared.SubmitResponse, error) {
	follow, err := s.UC.UnfollowUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_shared.SubmitResponse{
		Message: "success",
		Id:      follow.ID,
	}, nil
}
