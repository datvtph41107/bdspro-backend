package mapper

import (
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/domain"
)

type RoleMapper struct {
}

func NewRoleMapper() *RoleMapper {
	return &RoleMapper{}
}

func (m *RoleMapper) MapRoleListToItemPb(roles []domain.Role) []*sharepb.Role {
	var pbRoles []*sharepb.Role
	for _, role := range roles {
		pbRoles = append(pbRoles, &sharepb.Role{
			Id:       role.ID,
			RoleName: role.RoleName,
			RoleKey:  role.RoleKey,
		})
	}
	return pbRoles
}

func (m *RoleMapper) MapRoleListToPb(roles []domain.Role) []*sharepb.Role {
	var pbRoles []*sharepb.Role
	for _, role := range roles {
		pbRoles = append(pbRoles, m.MapRoleToPb(&role))
	}
	return pbRoles
}

// mapRoleToPb chuyển domain Role sang protobuf Role
func (m *RoleMapper) MapRoleToPb(role *domain.Role) *sharepb.Role {
	if role == nil {
		return nil
	}

	pbRole := &sharepb.Role{
		Id:              role.ID,
		RoleName:        role.RoleName,
		RoleDescription: role.RoleDescription,
		// Key:             role.Key,
		PermissionIds: role.PermissionIDs,
		// OrganizationId:  role.OrganizationID,
		IsDefault: role.IsDefault,
		// DomainType:      role.DomainType,
		ColorId:     role.ColorID,
		RoleGroupId: role.RoleGroupID,
		AllowAssign: role.AllowAssign,
		RoleKey:     role.RoleKey,
		CreatedAt:   _utils.FormatTimeToString(role.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(role.UpdatedAt),
	}

	if role.RoleGroup != nil {
		pbRole.RoleGroupName = &role.RoleGroup.GroupName
	}

	// if role.ColorID != nil {
	// 	pbRole.ColorId = role.ColorID
	// }

	if role.Color != nil {
		pbRole.Color = &sharepb.Color{
			Id:              role.Color.ID,
			Name:            role.Color.Name,
			ContentColor:    role.Color.ContentColor,
			BackgroundColor: role.Color.BackgroundColor,
			ColorKey:        role.Color.ColorKey,
			HexCode:         role.Color.HexCode,
			CreatedAt:       _utils.FormatTimeToString(role.Color.CreatedAt),
			UpdatedAt:       _utils.FormatTimeToString(role.Color.UpdatedAt),
		}
	}

	if role.Permissions != nil {
		pbRole.Permissions = m.MapPermissionListToPb(role.Permissions)
	}

	return pbRole
}

func (m *RoleMapper) MapPermissionListToPb(permissions []*domain.Permission) []*sharepb.Permission {
	var pbPermissions []*sharepb.Permission
	for _, permission := range permissions {
		pbPermissions = append(pbPermissions, m.MapPermissionToPb(permission))
	}
	return pbPermissions
}

func (m *RoleMapper) MapPermissionToPb(permission *domain.Permission) *sharepb.Permission {
	if permission == nil {
		return nil
	}

	pbPermission := &sharepb.Permission{
		Id:          permission.ID,
		Name:        permission.Name,
		Key:         permission.Key,
		Description: permission.Description,
		Module:      permission.Module,
		ParentId:    permission.ParentID,
		Priority:    permission.Priority,
	}

	return pbPermission
}

func (m *RoleMapper) MapPermissionListWithChild(permission *domain.Permission) *sharepb.Permission {
	if permission == nil {
		return nil
	}

	pbPermission := m.MapPermissionToPb(permission)

	if permission.Children != nil {
		pbPermission.Children = make([]*sharepb.Permission, len(permission.Children))
		for i, child := range permission.Children {
			pbPermission.Children[i] = m.MapPermissionListWithChild(child)
		}
	}

	return pbPermission
}

// MapAuthAdminWithRoleToPb convert domain.AuthAdminWithRole to proto
func (m *RoleMapper) MapAuthAdminWithRoleToPb(admin *domain.AuthAdminWithRole) *userpb.AuthAdminWithRole {
	if admin == nil {
		return nil
	}

	pbAdmin := &userpb.AuthAdminWithRole{
		Id:        admin.ID,
		UserId:    admin.UserID,
		Username:  admin.Username,
		FullName:  admin.FullName,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Avatar:    admin.Avatar,
		RoleKey:   admin.RoleKey,
		Status:    admin.Status,
		CreatedAt: _utils.FormatTimeToString(&admin.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(&admin.UpdatedAt),
	}

	// Map role nếu có
	if admin.Role != nil {
		pbAdmin.Role = m.MapRoleToPb(admin.Role)
	}

	return pbAdmin
}
