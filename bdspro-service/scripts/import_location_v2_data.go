package main

import (
	_models "common/domain/entity"
	"common/logging"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
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
	os.Exit(runImportLocationV2())
}

func runImportLocationV2() int {
	closeLogger, err := logging.Configure("bdspro-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure BDSPro location import logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	// Lấy connection string từ environment hoặc dùng default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=bdspro_service port=5432 sslmode=disable"
		slog.Warn(
			"using default database connection",
			slog.String("environment_variable", "DATABASE_URL"),
		)
	}

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		slog.Error(
			"connect to database failed",
			slog.Any("error", err),
		)
		return 1
	}

	// Auto migrate tables
	slog.Info("running location v2 auto migration")
	if err := db.AutoMigrate(&domain.ProvinceV2{}, &domain.WardV2{}); err != nil {
		slog.Error(
			"migrate location v2 tables failed",
			slog.Any("error", err),
		)
		return 1
	}

	// Import provinces and wards
	if err := importLocationV2Data(db); err != nil {
		slog.Error(
			"import location v2 data failed",
			slog.Any("error", err),
		)
		return 1
	}

	slog.Info("location v2 import completed successfully")
	return 0
}

func importLocationV2Data(db *gorm.DB) error {
	slog.Info("importing location v2 data")

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
	slog.Info("truncating old location v2 data")
	if err := db.Exec("TRUNCATE TABLE ward_v2 RESTART IDENTITY CASCADE").Error; err != nil {
		slog.Warn(
			"truncate location v2 table failed",
			slog.String("table", "ward_v2"),
			slog.Any("error", err),
		)
	}
	if err := db.Exec("TRUNCATE TABLE province_v2 RESTART IDENTITY CASCADE").Error; err != nil {
		slog.Warn(
			"truncate location v2 table failed",
			slog.String("table", "province_v2"),
			slog.Any("error", err),
		)
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
		slog.Info(
			"inserting location v2 provinces",
			slog.Int("count", len(provinces)),
		)
		if err := db.CreateInBatches(provinces, 100).Error; err != nil {
			return fmt.Errorf("failed to insert provinces: %w", err)
		}
		slog.Info(
			"imported location v2 provinces",
			slog.Int("count", len(provinces)),
		)
	}

	// Batch insert wards
	if len(allWards) > 0 {
		slog.Info(
			"inserting location v2 wards",
			slog.Int("count", len(allWards)),
		)
		if err := db.CreateInBatches(allWards, 500).Error; err != nil {
			return fmt.Errorf("failed to insert wards: %w", err)
		}
		slog.Info(
			"imported location v2 wards",
			slog.Int("count", len(allWards)),
		)
	}

	return nil
}
