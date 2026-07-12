package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendGroupService struct {
	UC *usecase.FriendGroupUsecase
	crmpb.UnimplementedFriendGroupServiceServer
}

func NewFriendGroupService(uc *usecase.FriendGroupUsecase) *FriendGroupService {
	return &FriendGroupService{
		UC: uc,
	}
}

// Tạo nhóm
// @Summary Tạo nhóm
// @Description Tạo nhóm
// @Tags Nhóm bạn bè
// @Accept json
// @Produce json
// @Param request body crmpb.FriendGroup true "Thông tin nhóm"
// @Success 200 {object} crmpb.FriendGroup "Nhóm đã được tạo"
// @Router /friend-group [post]
func (s *FriendGroupService) Create(ctx context.Context, req *crmpb.FriendGroup) (*crmpb.FriendGroup, error) {
	group := mapper.PbToFriendGroup(req)
	group, err := s.UC.CreateGroup(ctx, *group)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Create not implemented")
	}
	return mapper.FriendGroupToPb(group), nil
}

// Cập nhật nhóm
// @Summary Cập nhật nhóm
// @Description Cập nhật nhóm
// @Tags Nhóm bạn bè
// @Accept json
// @Produce json
// @Param request body crmpb.FriendGroup true "Thông tin nhóm"
// @Success 200 {object} crmpb.FriendGroup "Nhóm đã được cập nhật"
// @Router /friend-group/{id} [put]
func (s *FriendGroupService) Update(ctx context.Context, req *crmpb.FriendGroup) (*crmpb.FriendGroup, error) {
	group := mapper.PbToFriendGroup(req)
	group, err := s.UC.UpdateGroup(ctx, req.Id, group)
	if err != nil {
		return nil, err
	}
	return mapper.FriendGroupToPb(group), nil
}

// Xóa nhóm
// @Summary Xóa nhóm
// @Description Xóa nhóm
// @Tags Nhóm bạn bè
// @Accept json
// @Produce json
// @Param request body crmpb.FriendGroup true "Thông tin nhóm"
// @Success 200 {object} crmpb.FriendGroup "Nhóm đã được xóa"
// @Router /friend-group/{id} [delete]
func (s *FriendGroupService) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := s.UC.DeleteGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Message: "success",
		Id:      req.Id,
	}, nil
}

// Lấy danh sách nhóm
// @Summary Lấy danh sách nhóm
// @Description Lấy danh sách nhóm
// @Tags Nhóm bạn bè
// @Accept json
// @Produce json
// @Param request body crmpb.FriendGroup true "Thông tin nhóm"
// @Success 200 {object} crmpb.ListFriendGroupsResponse "Danh sách nhóm"
// @Router /friend-group/list [get]
func (s *FriendGroupService) List(ctx context.Context, req *sharepb.IdRequest) (*crmpb.ListFriendGroupsResponse, error) {
	groups, err := s.UC.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	return &crmpb.ListFriendGroupsResponse{
		Data: mapper.ListFriendGroupToPb(groups),
	}, nil
}
