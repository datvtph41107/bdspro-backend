package c_config

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

// NewDatabaseConfig tạo cấu hình database
func NewDatabaseConfig() *DatabaseConfig {
	fmt.Printf("URL")
	fmt.Printf(viper.GetString("database.dsn"))
	return &DatabaseConfig{
		DSN: viper.GetString("database.dsn"),
	}
}

var ConfigDB *gorm.DB

// NewDB khởi tạo kết nối database
func NewDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		log.Fatal("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		return nil, err
	}
	ConfigDB = db
	fmt.Println("✅ Kết nối PostgreSQL thành công!")
	return db, nil
}
