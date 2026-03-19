package main

import (
	_models "common/domain/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bdspro/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type WardJSON struct {
	Name          string `json:"name"`
	Code          int    `json:"code"`
	Codename      string `json:"codename"`
	DivisionType  string `json:"division_type"`
	ShortCodename string `json:"short_codename"`
}

type ProvinceV2JSON struct {
	Name         string     `json:"name"`
	Code         int        `json:"code"`
	Codename     string     `json:"codename"`
	DivisionType string     `json:"division_type"`
	PhoneCode    int        `json:"phone_code"`
	Wards        []WardJSON `json:"wards"`
}

func main() {
	// Lấy connection string từ environment hoặc dùng default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=bdspro_service port=5432 sslmode=disable"
		log.Println("Using default database connection. Set DATABASE_URL environment variable to override.")
	}

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate tables
	log.Println("Running auto migration...")
	if err := db.AutoMigrate(&domain.ProvinceV2{}, &domain.WardV2{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Import provinces and wards
	if err := importLocationV2Data(db); err != nil {
		log.Fatalf("Failed to import location v2 data: %v", err)
	}

	log.Println("✅ Import completed successfully!")
}

func importLocationV2Data(db *gorm.DB) error {
	log.Println("Importing location v2 data...")

	// Đọc file JSON
	data, err := ioutil.ReadFile("../../shared/code/locationv2.json")
	if err != nil {
		return fmt.Errorf("failed to read locationv2.json: %w", err)
	}

	var provincesJSON []ProvinceV2JSON
	if err := json.Unmarshal(data, &provincesJSON); err != nil {
		return fmt.Errorf("failed to parse locationv2.json: %w", err)
	}

	// Xóa dữ liệu cũ (nếu muốn)
	log.Println("Truncating old data...")
	if err := db.Exec("TRUNCATE TABLE ward_v2 RESTART IDENTITY CASCADE").Error; err != nil {
		log.Printf("Warning: Failed to truncate ward_v2 table: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE province_v2 RESTART IDENTITY CASCADE").Error; err != nil {
		log.Printf("Warning: Failed to truncate province_v2 table: %v", err)
	}

	// Chuyển đổi và insert provinces
	var provinces []domain.ProvinceV2
	var allWards []domain.WardV2
	wardCount := 0

	for _, p := range provincesJSON {
		province := domain.ProvinceV2{
			BaseEntity: _models.BaseEntity{},
			Name:       p.Name,
			// Code:         p.Code,
			Codename:     p.Codename,
			DivisionType: p.DivisionType,
			// PhoneCode:    p.PhoneCode,
		}
		provinces = append(provinces, province)

		// Chuẩn bị wards cho province này
		for _, w := range p.Wards {
			ward := domain.WardV2{
				BaseEntity: _models.BaseEntity{},
				Name:       w.Name,
				// Code:          w.Code,
				Codename:      w.Codename,
				DivisionType:  w.DivisionType,
				ShortCodename: w.ShortCodename,
				// ProvinceCode:  p.Code, // Link tới province bằng code
			}
			allWards = append(allWards, ward)
			wardCount++
		}
	}

	// Batch insert provinces
	if len(provinces) > 0 {
		log.Printf("Inserting %d provinces...", len(provinces))
		if err := db.CreateInBatches(provinces, 100).Error; err != nil {
			return fmt.Errorf("failed to insert provinces: %w", err)
		}
		log.Printf("✅ Imported %d provinces", len(provinces))
	}

	// Batch insert wards
	if len(allWards) > 0 {
		log.Printf("Inserting %d wards...", len(allWards))
		if err := db.CreateInBatches(allWards, 500).Error; err != nil {
			return fmt.Errorf("failed to insert wards: %w", err)
		}
		log.Printf("✅ Imported %d wards", len(allWards))
	}

	return nil
}
