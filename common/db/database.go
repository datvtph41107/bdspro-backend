package _db

import (
	"fmt"
	"log"

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
	if err != nil {
		fmt.Println("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		log.Fatal("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		return nil, err
	}
	DB = db
	fmt.Println("✅ Kết nối PostgreSQL thành công!")
	return db, nil
}
