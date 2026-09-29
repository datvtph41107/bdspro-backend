package database

import (
	"sync"

	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// func NewPostgresDB(dsn string) *gorm.DB {
// 	once.Do(func() {
// 		cfg := &gorm.Config{
// 			Logger: logger.New(
// 				log.New(os.Stdout, "\r\n", log.LstdFlags),
// 				logger.Config{
// 					SlowThreshold:             200 * time.Millisecond,
// 					LogLevel:                  logger.Info,
// 					IgnoreRecordNotFoundError: true,
// 					Colorful:                  true,
// 				},
// 			),
// 		}

// 		db, err := gorm.Open(postgres.New(postgres.Config{
// 			DSN:                  dsn,
// 			PreferSimpleProtocol: true,
// 		}), cfg)

// 		if err != nil {
// 			log.Fatalf("failed to connect to Postgres: %v", err)
// 		}

// 		dbInstance = db
// 	})

// 	return dbInstance
// }
