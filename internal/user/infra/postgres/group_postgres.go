package postgres

import (
	_db "common/db"
	_errors "common/errors"
	"fmt"
	"strconv"
	"user/internal"
	models "user/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupPostgres struct {
	DB *gorm.DB
}

// NewGroupPostgres creates a new instance of GroupPostgres
func NewGroupPostgres(db *gorm.DB) *GroupPostgres {
	return &GroupPostgres{DB: db}
}

// // Search retrieves FriendGroupEntity records based on the provided criteria
// func (repo *GroupPostgres) Search(dto *models.GroupEntity, pageable *_db.Page) ([]models.GroupEntity, error) {
// 	var groups []models.GroupEntity

// 	// Building the query
// 	query := repo.DB.Model(&models.GroupEntity{}).
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
// func (repo *GroupPostgres) GetByID(id uint64) (*models.GroupEntity, error) {
// 	var group models.GroupEntity
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
func (repo *GroupPostgres) CreateGroup(c *gin.Context, group *models.GroupEntity) error {
	return _db.SaveWithAudit(c, &group)
}

// Cập nhật nhóm
func (repo *GroupPostgres) UpdateGroup(c *gin.Context, group *models.GroupEntity) error {
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
func (repo *GroupPostgres) DeleteGroup(profileId *uint64, id uint64) error {
	var group models.GroupEntity

	// Tìm nhóm theo ID
	err := repo.DB.First(&group, id).Error
	if err != nil {
		return fmt.Errorf("nhóm không tồn tại hoặc đã bị xóa")
	}

	// Kiểm tra quyền xóa
	if group.CreatedBy != profileId {
		return _errors.ReturnError(service.FriendGroupDeleteDenied)
	}

	return repo.DB.Delete(&models.GroupEntity{}, id).Error
}

// Lấy danh sách nhóm
func (repo *GroupPostgres) ListGroups(createdBy uint64) ([]models.GroupEntity, error) {
	var groups []models.GroupEntity
	err := repo.DB.Where("deleted_at is null and created_by = ?", createdBy).Find(&groups).Error
	return groups, err
}

// Lấy thông tin chi tiết nhóm theo ID
func (repo *GroupPostgres) GetGroupByID(id uint) (*models.GroupEntity, error) {
	var group models.GroupEntity
	err := repo.DB.First(&group, id).Error
	return &group, err
}
