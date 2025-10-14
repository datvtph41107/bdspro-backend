package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrganizationPermissionModel struct {
	Id             uint32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"not null;index" json:"name"`
	Key            string     `gorm:"not null;index:idx_org_permissions_key" json:"key"`
	Description    string     `gorm:"type:text" json:"description"`
	ParentId       uint32     `json:"parentId"`
	PermissionType uint32     `gorm:"not null;default:10" json:"permissionType"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"index" json:"deletedAt,omitempty"`
	IsDeleted      bool       `gorm:"default:false" json:"isDeleted"`
}

func (OrganizationPermissionModel) TableName() string {
	return "organization_permissions"
}

type ExportData struct {
	Permissions      []OrganizationPermissionModel `json:"permissions"`
	ExportedAt       time.Time                     `json:"exportedAt"`
	TotalPermissions int                           `json:"totalPermissions"`
}

func main() {
	// Đọc connection string từ environment hoặc sử dụng default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default connection string - thay đổi theo config của bạn
		dsn = "host=localhost user=postgres password=postgres dbname=organization_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	}

	fmt.Println("🔌 Connecting to database...")
	fmt.Printf("   DSN: %s\n\n", dsn)

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Connected to database successfully!")

	// Lấy tất cả permissions
	var permissions []OrganizationPermissionModel
	fmt.Println("\n📊 Fetching permissions from organization_permissions table...")

	err = db.Where("is_deleted = ? OR is_deleted IS NULL", false).
		Order("parent_id ASC, id ASC").
		Find(&permissions).Error

	if err != nil {
		log.Fatalf("❌ Failed to fetch permissions: %v", err)
	}

	fmt.Printf("✅ Fetched %d permissions\n", len(permissions))

	// Tạo export data
	exportData := ExportData{
		Permissions:      permissions,
		ExportedAt:       time.Now(),
		TotalPermissions: len(permissions),
	}

	// Export ra file JSON
	outputFile := "permissions_export.json"
	fmt.Printf("\n💾 Saving to file: %s\n", outputFile)

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		log.Fatalf("❌ Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write file: %v", err)
	}

	fmt.Printf("\n✅ Export completed successfully!\n")
	fmt.Printf("📁 Output file: %s\n", outputFile)
	fmt.Printf("📊 Total permissions: %d\n", exportData.TotalPermissions)
	fmt.Printf("🕐 Exported at: %s\n", exportData.ExportedAt.Format("2006-01-02 15:04:05"))

	// Hiển thị một vài permissions đầu tiên
	if len(permissions) > 0 {
		fmt.Printf("\n📋 Sample permissions:\n")
		limit := 5
		if len(permissions) < limit {
			limit = len(permissions)
		}
		for i := 0; i < limit; i++ {
			p := permissions[i]
			fmt.Printf("   [%d] %s (key: %s, parent_id: %d)\n", p.Id, p.Name, p.Key, p.ParentId)
		}
		if len(permissions) > 5 {
			fmt.Printf("   ... and %d more\n", len(permissions)-5)
		}
	}
}
