package _db

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig chứa thông tin cấu hình DB
type DatabaseConfig struct {
	DSN string
}

var DB *gorm.DB

// NewDB khởi tạo kết nối database
func NewDB() (*gorm.DB, error) {
	dsn := viper.GetString("database.dsn")

	// Cấu hình logger của GORM (in SQL + thời gian + lỗi)
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // logger cơ bản
		logger.Config{
			SlowThreshold: time.Second, // Log query nào chạy quá 1s
			LogLevel:      logger.Info, // Info để hiện tất cả query
			Colorful:      true,        // In màu đẹp hơn trong terminal
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	fmt.Printf("🔗 Kết nối tới: %s\n", dsn)

	if err != nil {
		fmt.Println("❌ Không thể kết nối với cơ sở dữ liệu:", err)
		return nil, err
	}

	// Lấy *sql.DB để cấu hình pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ failed to get db from gorm: %v", err)
	}

	// Cấu hình pool
	sqlDB.SetMaxOpenConns(300)                // Tối đa 300 kết nối mở
	sqlDB.SetMaxIdleConns(200)                // Tối đa 200 kết nối nhàn rỗi
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // Mỗi kết nối tối đa sống 5 phút

	DB = db.Debug() // Bật debug mode: in tất cả query ra console

	fmt.Println("✅ Kết nối PostgreSQL thành công và bật chế độ debug!")
	return db, nil
}
