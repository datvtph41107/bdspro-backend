package postgre

import (
	_errors "common/errors"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.FriendGroupRepo
type FriendGroupRepo struct {
	DB *gorm.DB
}

// NewGroupRepo creates a new instance of NewGroupRepo
func NewGroupRepo(db *gorm.DB) *FriendGroupRepo {
	return &FriendGroupRepo{DB: db}
}

// // Search retrieves FriendGroupEntity records based on the provided criteria
// func (repo *GroupRepo) Search(dto *domain.GroupEntity, pageable *_db.Page) ([]domain.GroupEntity, error) {
// 	var groups []domain.GroupEntity

// 	// Building the query
// 	query := repo.DB.Model(&domain.GroupEntity{}).
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
// func (repo *GroupRepo) GetByID(id uint64) (*domain.GroupEntity, error) {
// 	var group domain.GroupEntity
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
func (repo *FriendGroupRepo) CreateGroup(c context.Context, group *domain.GroupEntity) error {
	return repo.DB.WithContext(c).Create(&group).Error
}

// Cập nhật nhóm
func (repo *FriendGroupRepo) UpdateGroup(c context.Context, id uint64, group *domain.GroupEntity) error {
	// Lấy ID từ path parameter
	// idParam := c.Param("id")

	// Chuyển ID từ string sang uint
	// id, err := strconv.ParseUint(idParam, 10, 64)
	// if err != nil {
	// 	return fmt.Errorf("ID không hợp lệ: %v", err)
	// }

	// Gán ID vào group
	group.ID = id

	return repo.DB.WithContext(c).
		Model(&domain.GroupEntity{}).
		Where("id = ? and deleted_at is null", id).
		Updates(group).Error
}

// Xóa nhóm
func (repo *FriendGroupRepo) DeleteGroup(c context.Context, profileId *uint64, id uint64) error {
	var group domain.GroupEntity

	// Tìm nhóm theo ID
	err := repo.DB.WithContext(c).First(&group, id).Error
	if err != nil {
		return fmt.Errorf("nhóm không tồn tại hoặc đã bị xóa")
	}

	// Kiểm tra quyền xóa
	if group.CreatedBy != profileId {
		return _errors.ReturnError(service.FriendGroupDeleteDenied)
	}

	return repo.DB.WithContext(c).Delete(&domain.GroupEntity{}, id).Error
}

// Lấy danh sách nhóm
func (repo *FriendGroupRepo) ListGroups(c context.Context, createdBy uint64) ([]domain.GroupEntity, error) {
	var groups []domain.GroupEntity
	err := repo.DB.WithContext(c).Where("deleted_at is null and created_by = ?", createdBy).Find(&groups).Error
	return groups, err
}

// Lấy thông tin chi tiết nhóm theo ID
func (repo *FriendGroupRepo) GetGroupByID(c context.Context, id uint64) (*domain.GroupEntity, error) {
	var group domain.GroupEntity
	err := repo.DB.WithContext(c).First(&group, id).Error
	return &group, err
}
