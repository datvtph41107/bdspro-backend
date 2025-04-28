package postgres

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	postgresInstance, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Ghi log ra terminal
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Đánh dấu query chậm
				LogLevel:                  logger.Info,            // Debug mode
				IgnoreRecordNotFoundError: true,                   // Không log lỗi RecordNotFound
				Colorful:                  true,                   // Hiển thị màu sắc
			},
		),
	})
	if err != nil {
		postgresInstance = nil
	}
	return postgresInstance, err
}
