package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type GroupModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_groups_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_groups_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	Name        string `gorm:"not null;index:idx_groups_name"`
	Description string
	AvatarUrl   string
	Status      entity.GroupStatus `gorm:"default:'active'"`

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (GroupModel) TableName() string {
	return "groups"
}

func GroupModelToEntity(model *GroupModel) *entity.Group {
	return &entity.Group{
		Id:          model.Id,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
		Name:        model.Name,
		Description: model.Description,
		AvatarUrl:   model.AvatarUrl,
		Status:      model.Status,
	}
}

func GroupEntityToModel(entity *entity.Group) *GroupModel {
	return &GroupModel{
		Id:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		Name:        entity.Name,
		Description: entity.Description,
		AvatarUrl:   entity.AvatarUrl,
		Status:      entity.Status,
	}
}

// @bind: organization/internal/domain/repository.GroupRepository
type GroupPostgresRepository struct {
	db *gorm.DB
}

func NewGroupPostgresRepository(db *gorm.DB) *GroupPostgresRepository {
	return &GroupPostgresRepository{db: db}
}

func (r *GroupPostgresRepository) Create(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	model := GroupEntityToModel(group)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *GroupPostgresRepository) Update(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	model := GroupEntityToModel(group)
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *GroupPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Group, error) {
	var model *GroupModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *GroupPostgresRepository) List(ctx context.Context, page, size int, currentUserId uint32) ([]*entity.Group, uint32, error) {
	var models []*GroupModel
	if err := r.db.WithContext(ctx).Offset(page*size).Limit(size).Where("created_by = ? AND is_deleted = ?", currentUserId, false).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.Group, len(models))
	for i, model := range models {
		entities[i] = GroupModelToEntity(model)
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Where("created_by = ? AND is_deleted = ?", currentUserId, false).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

func (r *GroupPostgresRepository) FindByUserId(ctx context.Context, userId uint32) ([]*entity.Group, error) {
	var groups []*GroupModel
	if err := r.db.WithContext(ctx).
		Joins("LEFT JOIN group_members ON groups.id = group_members.group_id").
		Where("groups.created_by = ? OR group_members.user_id = ? AND groups.is_deleted = ?", userId, userId, false).
		Distinct().
		Order("groups.created_at DESC").
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return GroupModelsToEntities(groups), nil
}

func GroupModelsToEntities(models []*GroupModel) []*entity.Group {
	entities := make([]*entity.Group, len(models))
	for i, model := range models {
		entities[i] = GroupModelToEntity(model)
	}
	return entities
}

func (r *GroupPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*entity.Group, error) {
	var models []*GroupModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	return GroupModelsToEntities(models), nil
}

func (r *GroupPostgresRepository) FindByUserIdWithDetails(ctx context.Context, userId uint32) ([]*dto.GroupWithDetails, error) {
	var results []struct {
		Id                        uint32             `gorm:"column:id"`
		CreatedAt                 time.Time          `gorm:"column:created_at"`
		UpdatedAt                 time.Time          `gorm:"column:updated_at"`
		CreatedBy                 uint32             `gorm:"column:created_by"`
		UpdatedBy                 uint32             `gorm:"column:updated_by"`
		Name                      string             `gorm:"column:name"`
		Description               string             `gorm:"column:description"`
		AvatarUrl                 string             `gorm:"column:avatar_url"`
		Status                    entity.GroupStatus `gorm:"column:status"`
		MemberCount               uint32             `gorm:"column:member_count"`
		Role                      string             `gorm:"column:role"`
		RoleName                  string             `gorm:"column:role_name"`
		LatestNotificationContent string             `gorm:"column:latest_notification_content"`
	}

	query := r.db.WithContext(ctx).
		Table("groups g").
		Select(`
			g.id, g.created_at, g.updated_at, g.created_by, g.updated_by,
			g.name, g.description, g.avatar_url, g.status,
			COALESCE(member_counts.count, 0) as member_count,
			COALESCE(gm.role, CASE WHEN g.created_by = ? THEN 'admin' ELSE 'member' END) as role,
			COALESCE(org_role.name, CASE WHEN g.created_by = ? THEN 'Admin' ELSE 'Thành viên' END) as role_name,
			gn.content as latest_notification_content
		`, userId, userId).
		Joins("LEFT JOIN group_members gm ON g.id = gm.group_id AND gm.user_id = ? AND gm.is_deleted = false", userId).
		Joins("LEFT JOIN organization_roles org_role ON gm.role = org_role.key").
		Joins(`
			LEFT JOIN (
				SELECT group_id, COUNT(*) as count 
				FROM group_members 
				WHERE is_deleted = false 
				GROUP BY group_id
			) member_counts ON g.id = member_counts.group_id
		`).
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (group_id) group_id, content
				FROM group_notifications 
				WHERE is_deleted = false 
				ORDER BY group_id, created_at DESC
			) gn ON g.id = gn.group_id
		`).
		Where("(g.created_by = ? OR gm.group_id IS NOT NULL) AND g.is_deleted = ?", userId, false).
		Order("g.created_at DESC")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	groupsWithDetails := make([]*dto.GroupWithDetails, len(results))
	for i, result := range results {
		group := &entity.Group{
			Id:          result.Id,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			CreatedBy:   result.CreatedBy,
			UpdatedBy:   result.UpdatedBy,
			Name:        result.Name,
			Description: result.Description,
			AvatarUrl:   result.AvatarUrl,
			Status:      result.Status,
		}

		// Map status to Vietnamese
		statusMap := map[entity.GroupStatus]string{
			entity.GroupStatusActive:    "Đang hoạt động",
			entity.GroupStatusPaused:    "Tạm ngưng",
			entity.GroupStatusDisbanded: "Đã giải tán",
		}

		// Determine user role
		userRole := result.RoleName
		if userRole == "" {
			// Map role string to Vietnamese
			roleMap := map[string]string{
				"admin":  "Admin",
				"member": "Thành viên",
			}
			userRole = roleMap[result.Role]
			if userRole == "" {
				if result.CreatedBy == userId {
					userRole = "Admin"
				} else {
					userRole = "Thành viên"
				}
			}
		}

		// Determine additional info
		additionalInfo := "Bạn chưa mời thành viên nào"
		if result.LatestNotificationContent != "" {
			additionalInfo = result.LatestNotificationContent
		}

		groupsWithDetails[i] = &dto.GroupWithDetails{
			Group:          group,
			MemberCount:    result.MemberCount,
			UserRole:       userRole,
			GroupStatus:    statusMap[result.Status],
			AdditionalInfo: additionalInfo,
		}
	}

	return groupsWithDetails, nil
}

func (r *GroupPostgresRepository) ListWithDetails(ctx context.Context, page, size int, currentUserId uint32) ([]*dto.GroupWithDetails, uint32, error) {
	var results []struct {
		Id                        uint32             `gorm:"column:id"`
		CreatedAt                 time.Time          `gorm:"column:created_at"`
		UpdatedAt                 time.Time          `gorm:"column:updated_at"`
		CreatedBy                 uint32             `gorm:"column:created_by"`
		UpdatedBy                 uint32             `gorm:"column:updated_by"`
		Name                      string             `gorm:"column:name"`
		Description               string             `gorm:"column:description"`
		AvatarUrl                 string             `gorm:"column:avatar_url"`
		Status                    entity.GroupStatus `gorm:"column:status"`
		MemberCount               uint32             `gorm:"column:member_count"`
		RoleId                    *uint64            `gorm:"column:role_id"`
		RoleName                  string             `gorm:"column:role_name"`
		LatestNotificationContent string             `gorm:"column:latest_notification_content"`
	}

	query := r.db.WithContext(ctx).
		Debug().
		Table("groups g").
		Select(`
			g.id, g.created_at, g.updated_at, g.created_by, g.updated_by,
			g.name, g.description, g.avatar_url, g.status,
			COALESCE(member_counts.count, 0) as member_count,
			gm.role_id,
			COALESCE(org_role.name, CASE WHEN g.created_by = ? THEN 'Admin' ELSE 'Thành viên' END) as role_name,
			gn.content as latest_notification_content
		`, currentUserId).
		Joins("LEFT JOIN group_members gm ON g.id = gm.group_id AND gm.user_id = ? AND gm.is_deleted = false", currentUserId).
		Joins("LEFT JOIN organization_roles org_role ON gm.role = org_role.key").
		Joins(`
			LEFT JOIN (
				SELECT group_id, COUNT(*) as count 
				FROM group_members 
				WHERE is_deleted = false 
				GROUP BY group_id
			) member_counts ON g.id = member_counts.group_id
		`).
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (group_id) group_id, content
				FROM group_notifications 
				WHERE is_deleted = false 
				ORDER BY group_id, created_at DESC
			) gn ON g.id = gn.group_id
		`).
		Where("gm.user_id = ? AND g.is_deleted = ?", currentUserId, false).
		Order("g.created_at DESC")

	// Count total
	var count int64
	countQuery := r.db.WithContext(ctx).
		Table("groups").
		Where("created_by = ? AND is_deleted = ?", currentUserId, false)
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Offset(page * size).Limit(size).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	groupsWithDetails := make([]*dto.GroupWithDetails, len(results))
	for i, result := range results {
		group := &entity.Group{
			Id:          result.Id,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			CreatedBy:   result.CreatedBy,
			UpdatedBy:   result.UpdatedBy,
			Name:        result.Name,
			Description: result.Description,
			AvatarUrl:   result.AvatarUrl,
			Status:      result.Status,
		}

		// Map status to Vietnamese
		statusMap := map[entity.GroupStatus]string{
			entity.GroupStatusActive:    "Đang hoạt động",
			entity.GroupStatusPaused:    "Tạm ngưng",
			entity.GroupStatusDisbanded: "Đã giải tán",
		}

		// Determine user role
		// userRole := result.RoleName
		// if userRole == "" {
		// 	// Map role string to Vietnamese
		// 	roleMap := map[string]string{
		// 		"admin":  "Admin",
		// 		"member": "Thành viên",
		// 	}
		// 	userRole = roleMap[result.Role]
		// 	if userRole == "" {
		// 		if result.CreatedBy == currentUserId {
		// 			userRole = "Admin"
		// 		} else {
		// 			userRole = "Thành viên"
		// 		}
		// 	}
		// }

		// Determine additional info
		additionalInfo := "Bạn chưa mời thành viên nào"
		if result.LatestNotificationContent != "" {
			additionalInfo = result.LatestNotificationContent
		}

		groupsWithDetails[i] = &dto.GroupWithDetails{
			Group:          group,
			MemberCount:    result.MemberCount,
			RoleId:         result.RoleId,
			GroupStatus:    statusMap[result.Status],
			AdditionalInfo: additionalInfo,
			// UserRole:       result.RoleName,
		}
	}

	return groupsWithDetails, uint32(count), nil
}
