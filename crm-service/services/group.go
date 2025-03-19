package services

import (
	_jwt "common/jwt"
	"crm/models"
	"crm/repo"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupService struct {
	groupRepo *repo.GroupRepo
}

func NewGroupService(groupRepo *repo.GroupRepo) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

// Thêm nhóm
func (s *GroupService) CreateGroup(c *gin.Context) (*models.GroupEntity, error) {
	var group models.GroupEntity
	if err := c.ShouldBindJSON(&group); err != nil {
		return nil, err
	}

	err := s.groupRepo.CreateGroup(c, &group)
	return &group, err
}

// Cập nhật nhóm
func (s *GroupService) UpdateGroup(c *gin.Context) (*models.GroupEntity, error) {
	var group models.GroupEntity
	if err := c.ShouldBindJSON(&group); err != nil {
		return nil, err
	}

	err := s.groupRepo.UpdateGroup(c, &group)
	return &group, err
}

// Xóa nhóm
func (s *GroupService) DeleteGroup(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}

	profileId := _jwt.GetProfileId(c)

	return s.groupRepo.DeleteGroup(&profileId, id)
}

// Lấy danh sách nhóm
func (s *GroupService) ListGroups(c *gin.Context) ([]models.GroupEntity, error) {
	profileId := _jwt.GetProfileId(c)
	return s.groupRepo.ListGroups(profileId)
}

// Lấy thông tin nhóm theo ID
func (s *GroupService) GetGroupByID(c *gin.Context) (*models.GroupEntity, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, errors.New("ID không hợp lệ")
	}

	return s.groupRepo.GetGroupByID(uint(id))
}
