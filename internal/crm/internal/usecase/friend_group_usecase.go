package usecase

import (
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
)

type FriendGroupUsecase struct {
	friendGroupRepo repo.FriendGroupRepo
}

func NewFriendGroupUsecase(friendGroupRepo repo.FriendGroupRepo) *FriendGroupUsecase {
	return &FriendGroupUsecase{
		friendGroupRepo: friendGroupRepo,
	}
}

// Thêm nhóm
func (s *FriendGroupUsecase) CreateGroup(c context.Context, group domain.GroupEntity) (*domain.GroupEntity, error) {
	err := s.friendGroupRepo.CreateGroup(c, &group)
	return &group, err
}

// Cập nhật nhóm
func (s *FriendGroupUsecase) UpdateGroup(c context.Context, id uint64, group *domain.GroupEntity) (*domain.GroupEntity, error) {
	err := s.friendGroupRepo.UpdateGroup(c, id, group)
	return group, err
}

// Xóa nhóm
func (s *FriendGroupUsecase) DeleteGroup(c context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.friendGroupRepo.DeleteGroup(c, &profileId, id)
}

// Lấy danh sách nhóm
func (s *FriendGroupUsecase) ListGroups(c context.Context) ([]domain.GroupEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.friendGroupRepo.ListGroups(c, profileId)
}

// Lấy thông tin nhóm theo ID
func (s *FriendGroupUsecase) GetGroupByID(c context.Context, id uint64) (*domain.GroupEntity, error) {
	return s.friendGroupRepo.GetGroupByID(c, id)
}