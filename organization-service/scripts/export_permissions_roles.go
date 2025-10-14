package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"organization/internal/enums"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrganizationPermissionModel struct {
	Id             uint32               `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string               `gorm:"not null;index" json:"name"`
	Key            string               `gorm:"not null;index:idx_org_permissions_key" json:"key"`
	Description    string               `gorm:"type:text" json:"description"`
	ParentId       uint32               `json:"parent_id"`
	PermissionType enums.PermissionType `gorm:"not null;default:10" json:"permission_type"`
	CreatedAt      time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time           `gorm:"index" json:"deleted_at,omitempty"`
	IsDeleted      bool                 `gorm:"default:false" json:"is_deleted"`
}

func (OrganizationPermissionModel) TableName() string {
	return "organization_permissions"
}

type RolePermissionModel struct {
	OrganizationRoleModelId       uint32 `gorm:"primaryKey" json:"organization_role_model_id"`
	OrganizationPermissionModelId uint32 `gorm:"primaryKey" json:"organization_permission_model_id"`
}

func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

type PermissionNode struct {
	Id             uint32               `json:"id"`
	Name           string               `json:"name"`
	Key            string               `json:"key"`
	Description    string               `json:"description"`
	ParentId       uint32               `json:"parent_id"`
	PermissionType enums.PermissionType `json:"permission_type"`
	Childs         []*PermissionNode    `json:"childs,omitempty"`
}

type ExportData struct {
	Permissions     []*PermissionNode     `json:"permissions"`
	RoleGroups      []RolePermissionModel `json:"role_groups"`
	ExportedAt      time.Time             `json:"exported_at"`
	TotalRootPerms  int                   `json:"total_root_permissions"`
	TotalRoleGroups int                   `json:"total_role_groups"`
}

func main() {
	// Đọc connection string từ environment hoặc sử dụng default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default connection string - thay đổi theo config của bạn
		dsn = "host=localhost user=postgres password=postgres dbname=organization_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	}

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// Lấy tất cả permissions (không bị xóa)
	var permissions []OrganizationPermissionModel
	err = db.Where("is_deleted = ? OR is_deleted IS NULL", false).
		Order("parent_id ASC, id ASC").
		Find(&permissions).Error
	if err != nil {
		log.Fatalf("Failed to fetch permissions: %v", err)
	}

	fmt.Printf("Fetched %d permissions\n", len(permissions))

	// Build hierarchy
	permissionMap := make(map[uint32]*PermissionNode)
	var rootPermissions []*PermissionNode

	// Tạo tất cả nodes trước
	for _, perm := range permissions {
		node := &PermissionNode{
			Id:             perm.Id,
			Name:           perm.Name,
			Key:            perm.Key,
			Description:    perm.Description,
			ParentId:       perm.ParentId,
			PermissionType: perm.PermissionType,
			Childs:         []*PermissionNode{},
		}
		permissionMap[perm.Id] = node
	}

	// Build hierarchy: gắn các node vào parent hoặc root
	for _, perm := range permissions {
		node := permissionMap[perm.Id]
		if perm.ParentId == 0 {
			// Root permission
			rootPermissions = append(rootPermissions, node)
		} else {
			// Gắn vào parent
			if parent, exists := permissionMap[perm.ParentId]; exists {
				parent.Childs = append(parent.Childs, node)
			} else {
				// Nếu không tìm thấy parent, coi như root
				rootPermissions = append(rootPermissions, node)
			}
		}
	}

	fmt.Printf("Built hierarchy with %d root permissions\n", len(rootPermissions))

	// Lấy tất cả role_permissions
	var roleGroups []RolePermissionModel
	err = db.Find(&roleGroups).Error
	if err != nil {
		log.Fatalf("Failed to fetch role_permissions: %v", err)
	}

	fmt.Printf("Fetched %d role_groups\n", len(roleGroups))

	// Tạo export data
	exportData := ExportData{
		Permissions:     rootPermissions,
		RoleGroups:      roleGroups,
		ExportedAt:      time.Now(),
		TotalRootPerms:  len(rootPermissions),
		TotalRoleGroups: len(roleGroups),
	}

	// Export ra file JSON
	outputFile := "permissions_roles_export.json"
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("\n✅ Export completed successfully!\n")
	fmt.Printf("📁 Output file: %s\n", outputFile)
	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("   - Root Permissions: %d\n", len(rootPermissions))
	fmt.Printf("   - Total Permissions: %d\n", len(permissions))
	fmt.Printf("   - Role-Permission Mappings: %d\n", len(roleGroups))
	fmt.Printf("   - Exported at: %s\n", exportData.ExportedAt.Format("2006-01-02 15:04:05"))
}
