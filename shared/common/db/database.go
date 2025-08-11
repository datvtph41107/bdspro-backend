package _db

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig chứa thông tin cấu hình DB
type DatabaseConfig struct {
	DSN string
}

var DB *gorm.DB

// NewDB khởi tạo kết nối database
func NewDB() (*gorm.DB, error) {
	dsn := viper.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("Kết nối tới: %s", dsn)
	if err != nil {
		fmt.Println("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		return nil, err
	}
	// Lấy *sql.DB để set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get db from gorm: %v", err)
	}

	// Cấu hình pool
	sqlDB.SetMaxOpenConns(300) // Tối đa 50 kết nối đồng thời mở được
	sqlDB.SetMaxIdleConns(200) // Tối đa 25 kết nối nhàn rỗi giữ sẵn
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	// sqlDB.SetConnMaxLifetime(time.Minute * 5)

	DB = db
	fmt.Println("✅ Kết nối PostgreSQL thành công!")
	return db, nil
}
