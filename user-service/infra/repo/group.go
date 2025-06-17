package repo

import (
	_db "common/db"
	_routes "common/routes"
	"fmt"
	"strconv"
	user_model "user/internal/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupRepo struct {
	DB *gorm.DB
}

// NewGroupRepo creates a new instance of NewGroupRepo
func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{DB: db}
}

// // Search retrieves FriendGroupEntity records based on the provided criteria
// func (repo *GroupRepo) Search(dto *user_model.GroupEntity, pageable *_db.Page) ([]user_model.GroupEntity, error) {
// 	var groups []user_model.GroupEntity

// 	// Building the query
// 	query := repo.DB.Model(&user_model.GroupEntity{}).
// 		Where("deleted_at is null")

// 	// Add name filtering condition if needed
// 	if dto.Name != "" {
// 		query = query.Where("lower(name) LIKE ?", "%"+dto.Name+"%")
// 	}

// 	// Execute the query with pagination
// 	err := query.Offset(int(pageable.Skip)).
// 		Limit(int(pageable.Limit)).
// 		Find(&groups).Error

// 	return groups, err
// }

// // GetByID retrieves a FriendGroupEntity by its ID
// func (repo *GroupRepo) GetByID(id uint64) (*user_model.GroupEntity, error) {
// 	var group user_model.GroupEntity
// 	// Find the group by ID
// 	if err := repo.DB.Where("id = ? AND deleted_at is null", id).First(&group).Error; err != nil {
// 		// If no record is found, return an error
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil // or you could return an error here
// 		}
// 		// For any other error, return it
// 		return nil, err
// 	}
// 	return &group, nil
// }

// type GroupRepo struct {
// 	DB *gorm.DB
// }

// func NewGroupRepo(DB *gorm.DB) *GroupRepo {
// 	return &GroupRepo{DB: DB}
// }

// Thêm nhóm mới
func (repo *GroupRepo) CreateGroup(c *gin.Context, group *user_model.GroupEntity) error {
	return _db.SaveWithAudit(c, &group)
}

// Cập nhật nhóm
func (repo *GroupRepo) UpdateGroup(c *gin.Context, group *user_model.GroupEntity) error {
	// Lấy ID từ path parameter
	idParam := c.Param("id")

	// Chuyển ID từ string sang uint
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return fmt.Errorf("ID không hợp lệ: %v", err)
	}

	// Gán ID vào group
	group.ID = id

	return _db.SaveWithAudit(c, group)
}

// Xóa nhóm
func (repo *GroupRepo) DeleteGroup(profileId *uint64, id uint64) error {
	var group user_model.GroupEntity

	// Tìm nhóm theo ID
	err := repo.DB.First(&group, id).Error
	if err != nil {
		return fmt.Errorf("nhóm không tồn tại hoặc đã bị xóa")
	}

	// Kiểm tra quyền xóa
	if group.CreatedBy != profileId {
		return &_routes.Except{
			Code:    401,
			Message: "bạn không có quyền xóa nhóm này",
		}
	}

	return repo.DB.Delete(&user_model.GroupEntity{}, id).Error
}

// Lấy danh sách nhóm
func (repo *GroupRepo) ListGroups(createdBy uint64) ([]user_model.GroupEntity, error) {
	var groups []user_model.GroupEntity
	err := repo.DB.Where("deleted_at is null and created_by = ?", createdBy).Find(&groups).Error
	return groups, err
}

// Lấy thông tin chi tiết nhóm theo ID
func (repo *GroupRepo) GetGroupByID(id uint) (*user_model.GroupEntity, error) {
	var group user_model.GroupEntity
	err := repo.DB.First(&group, id).Error
	return &group, err
}
