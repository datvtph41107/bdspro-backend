package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"tqd/internal/domain/jsonb"
	qh_domain "tqd/internal/domain/qh"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GeoJSONFeature đại diện cho một feature trong NDJSON
type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   map[string]interface{} `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// mapping land use từ 59 loại khác nhau về code chuẩn
var landUseMapping = map[string]string{
	// Đất cây xanh
	"Đất cây xanh": "DCX",
	"đất cây xanh": "DCX",
	"Đất cây xanh cách ly, phòng hộ, vườn ươm, nghiên *": "DCX",
	"Công viên cây xanh nghĩa trang":                     "DCX",
	"Đất cây xanh,công viên vui chơi giải trí":           "DCX",
	"Đất Cây Xanh": "DCX",
	// Sông hồ
	"Sông hồ":                   "DSH",
	"sông hồ":                   "DSH",
	"đất sông hồ":               "DSH",
	"Mặt nước, sông, hồ, mương": "DSH",
	"Mặt nước":                  "DSH",
	"Sông Hồ":                   "DSH",
	// Đất làng xóm
	"Đất làng xóm":                     "DLX",
	"đất làng xóm":                     "DLX",
	"Đất làng xóm ven đô, ngoại thành": "DLX",
	"Đất làng xóm đô thị hóa":          "DLX",
	"đất láng xóm":                     "DLX",
	"Đất ở làng xóm":                   "ONT",
	// Đất ở đô thị
	"Đất ở đô thị":            "ODT",
	"đất ở đô thị":            "ODT",
	"đất đô thị":              "ODT",
	"Đất đô thị":              "ODT",
	"Đất đô thị sinh thái":    "ODT",
	"Đất đơn vị ở mới":        "ODT",
	"Đất dân dụng hiện trạng": "ODT",
	// Đất công cộng
	"Đất công cộng":           "DCC",
	"đất công cộng":           "DCC",
	"Đất công cộng thành phố": "DCC",
	"đất công cộng thành phố": "DCC",
	"Đất công cộng ở":         "DCC",
	"Đất công cộng đô thị":    "DCC",
	"Đất cồng cộng đô thị":    "DCC",
	"Đất công cộng, hỗn hợp (dịch vụ, thương mại, văn *": "DCC",
	// Đất di tích
	"Đất di tích":                      "DDT",
	"đất di tích":                      "DDT",
	"Đất di tích, công trình tôn giáo": "DDT",
	// Đất giáo dục
	"Đất giáo dục": "DGD",
	"đất giáo dục": "DGD",
	"Đất GIáo Dục": "DGD",
	// Đất lâm nghiệp
	"Đất lâm nghiệp": "LNP",
	"đất lâm nghiệp": "LNP",
	// Đất nông nghiệp
	"Đất nông nghiệp":          "NNP",
	"đất nông nghiệp":          "NNP",
	"nông nghiệp":              "NNP",
	"nông nghiệp thủy sản":     "NNP",
	"Đất nông nghiệp thủy sản": "NNP",
	// Đất hỗn hợp
	"Đất hỗn hợp": "DHH",
	"đất hỗn hợp": "DHH",
	// Đất hạ tầng
	"Đất hạ tầng kỹ thuật":        "DHT",
	"đất hạ tầng kỹ thuật":        "DHT",
	"Đất hạ tầng":                 "DHT",
	"đất hạ tầng":                 "DHT",
	"Đất hạ tầng kĩ thuật":        "DHT",
	"Đất công trình đầu mối HTKT": "DHT",
	// Bãi đỗ xe
	"Bãi đỗ xe": "BDX",
	"Bãi đõ xe": "BDX",
	"Bãi Đỗ Xe": "BDX",
	"Để xe":     "BDX",
	// Đất công nghiệp
	"Đất công nghiệp":                 "DCN",
	"đất công nghiệp":                 "DCN",
	"Đất công nghiệp - công nghệ cao": "DCN",
	// Đất an ninh quốc phòng
	"Đất an ninh quốc phòng":   "DAN",
	"đất an ninh quốc phòng":   "DAN",
	"Đất an ninh, quốc phòng":  "DAN",
	"Đất quốc phòng":           "DAN",
	"đất an ninh - quốc phòng": "DAN",
	// Đất nghĩa trang
	"Đất nghĩa trang": "DNT",
	"đất nghĩa trang": "DNT",
	// Đất du lịch
	"Đất du lịch": "DDL",
	"đất du lịch": "DDL",
	// Đất dự trữ
	"Đất dự trữ": "DDR",
	"đất dự trữ": "DDR",
	// Mặc định
	"": "UNKNOWN",
}

// normalizeLandUseCode trả về code chuẩn từ tên loại đất
func normalizeLandUseCode(raw string) string {
	if raw == "" {
		return "UNKNOWN"
	}
	if code, ok := landUseMapping[raw]; ok {
		return code
	}
	return "UNKNOWN"
}

// parseColor chuyển interface thành uint8
func parseColor(val interface{}) uint8 {
	if val == nil {
		return 200
	}
	switch v := val.(type) {
	case float64:
		return uint8(v)
	case int:
		return uint8(v)
	case json.Number:
		f, _ := v.Float64()
		return uint8(f)
	}
	return 200
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// sanitizeName tạo tên layer từ do_an (loại bỏ ký tự đặc biệt)
func sanitizeName(name string) string {
	re := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"(", "",
		")", "",
		",", "",
		".", "",
		"-", "_",
	)
	return re.Replace(name)
}

func main() {
	// Command line arguments
	dataDir := flag.String("dir", "", "Directory containing NDJSON files")
	dbHost := flag.String("db-host", "localhost", "Database host")
	dbPort := flag.String("db-port", "5432", "Database port")
	dbUser := flag.String("db-user", "postgres", "Database user")
	dbPassword := flag.String("db-password", "", "Database password")
	dbName := flag.String("db-name", "db_tqd", "Database name")
	batchSize := flag.Int("batch", 200, "Batch size for insert")
	flag.Parse()

	if *dataDir == "" {
		log.Fatal("--dir is required")
	}

	// Kết nối database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Tự động migrate các bảng (chỉ lần đầu)
	db.AutoMigrate(&qh_domain.QHLayer{}, &qh_domain.QHRegion{})

	// Map lưu do_an -> layer_id
	layerMap := make(map[string]uint64)

	// Tạo layer mặc định cho những feature không có do_an
	defaultLayer := &qh_domain.QHLayer{
		Name:        "default_planning",
		DisplayName: "Quy hoạch chung (không xác định)",
		Type:        3,
		Status:      10,
		Visible:     1,
	}
	db.Where("name = ?", defaultLayer.Name).FirstOrCreate(defaultLayer)
	layerMap[""] = defaultLayer.ID
	log.Printf("Default layer ID: %d", defaultLayer.ID)

	// Tìm tất cả file NDJSON
	files, err := filepath.Glob(filepath.Join(*dataDir, "*.ndjson"))
	if err != nil {
		log.Fatalf("Failed to list NDJSON files: %v", err)
	}
	log.Printf("Found %d NDJSON files", len(files))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // xử lý đồng thời 4 file
	var mu sync.Mutex
	totalSuccess := 0
	totalFeatures := 0
	startTime := time.Now()

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			log.Printf("Processing: %s", filepath.Base(filePath))

			f, err := os.Open(filePath)
			if err != nil {
				log.Printf("Cannot open file %s: %v", filePath, err)
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			// 20MB buffer
			scanner.Buffer(make([]byte, 10*1024*1024), 50*1024*1024)

			var batch []*qh_domain.QHRegion
			fileTotal := 0
			fileSuccess := 0

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				fileTotal++

				var feature GeoJSONFeature
				if err := json.Unmarshal([]byte(line), &feature); err != nil {
					log.Printf("JSON parse error in %s line %d: %v", filepath.Base(filePath), fileTotal, err)
					continue
				}

				// Lấy do_an để xác định layer
				doAn := ""
				if val, ok := feature.Properties["do_an"]; ok && val != nil {
					doAn = strings.TrimSpace(val.(string))
				}

				mu.Lock()
				layerID, exists := layerMap[doAn]
				if !exists {
					// Tạo layer mới
					newLayer := &qh_domain.QHLayer{
						Name:        sanitizeName(doAn),
						DisplayName: doAn,
						Type:        3,
						Status:      10,
						Visible:     1,
					}
					if err := db.Create(newLayer).Error; err != nil {
						log.Printf("Failed to create layer for %s: %v", doAn, err)
						mu.Unlock()
						continue
					}
					layerID = newLayer.ID
					layerMap[doAn] = layerID
					log.Printf("Created new layer: %s (ID=%d)", doAn, layerID)
				}
				mu.Unlock()

				// Lấy loại đất
				landUseRaw := ""
				if val, ok := feature.Properties["loai_dat_quy_hoach"]; ok && val != nil {
					landUseRaw = strings.TrimSpace(val.(string))
				}
				landUseCode := normalizeLandUseCode(landUseRaw)

				// Màu sắc
				colorRed := parseColor(feature.Properties["red"])
				colorGreen := parseColor(feature.Properties["green"])
				colorBlue := parseColor(feature.Properties["blue"])

				// Văn bản pháp lý
				legalDoc := ""
				if val, ok := feature.Properties["quyet_dinh"]; ok && val != nil {
					legalDoc = val.(string)
				}

				// Geometry dưới dạng JSON string
				geomJSON, err := json.Marshal(feature.Geometry)
				if err != nil {
					log.Printf("Marshal geometry error line %d: %v", fileTotal, err)
					continue
				}

				// Raw properties & geometry JSONB
				propsJSON, _ := json.Marshal(feature.Properties)

				// Tạo StyleConfig
				styleConfigMap := jsonb.JSONB{
					"fillColor":   fmt.Sprintf("#%02X%02X%02X", colorRed, colorGreen, colorBlue),
					"fillOpacity": 0.6,
					"strokeColor": "#000000",
					"strokeWidth": 1,
				}

				// Tạo region object
				region := &qh_domain.QHRegion{
					LayerID:       layerID,
					Name:          landUseRaw,
					DisplayName:   capitalize(landUseRaw),
					LandUseCode:   landUseCode,
					LandUseName:   landUseRaw,
					LegalDoc:      legalDoc,
					PlanningName:  doAn,
					Geometry:      geomJSON, // lưu GeoJSON bytes
					ColorRed:      colorRed,
					ColorGreen:    colorGreen,
					ColorBlue:     colorBlue,
					StyleConfig:   styleConfigMap,
					RawProperties: make(jsonb.JSONB),
					RawGeometry:   make(jsonb.JSONB),
					Status:        10,
					Version:       1,
					IsLatest:      true,
				}

				// Unmarshal raw properties
				json.Unmarshal(propsJSON, &region.RawProperties)
				var geomMap map[string]interface{}
				json.Unmarshal(geomJSON, &geomMap)
				region.RawGeometry = geomMap

				batch = append(batch, region)
				fileSuccess++

				if len(batch) >= *batchSize {
					if err := insertBatch(db, batch); err != nil {
						log.Printf("Batch insert error: %v", err)
					} else {
						mu.Lock()
						totalSuccess += len(batch)
						mu.Unlock()
					}
					batch = nil
				}
			}

			// Insert phần còn lại
			if len(batch) > 0 {
				if err := insertBatch(db, batch); err != nil {
					log.Printf("Final batch insert error: %v", err)
				} else {
					mu.Lock()
					totalSuccess += len(batch)
					mu.Unlock()
				}
			}

			mu.Lock()
			totalFeatures += fileTotal
			mu.Unlock()
			log.Printf("File %s: %d/%d regions imported", filepath.Base(filePath), fileSuccess, fileTotal)

			if err := scanner.Err(); err != nil {
				log.Printf("Scanner error: %v", err)
			}
		}(file)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	log.Printf("\n============================================================")
	log.Printf("IMPORT SUMMARY")
	log.Printf("============================================================")
	log.Printf("Total files: %d", len(files))
	log.Printf("Total features read: %d", totalFeatures)
	log.Printf("Successfully imported: %d", totalSuccess)
	log.Printf("Time taken: %v", elapsed)
	log.Printf("============================================================")

	// Kiểm tra tổng số regions trong DB
	var count int64
	db.Model(&qh_domain.QHRegion{}).Where("deleted_at IS NULL").Count(&count)
	log.Printf("Total regions in qh_regions: %d", count)
}

// insertBatch chèn batch regions bằng raw SQL với ST_GeomFromGeoJSON
func insertBatch(db *gorm.DB, batch []*qh_domain.QHRegion) error {
	if len(batch) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, r := range batch {
			sql := `
                INSERT INTO qh_regions (
                    layer_id, name, display_name, land_use_code, land_use_name,
                    legal_doc, planning_name, geometry, color_red, color_green, color_blue,
                    style_config, raw_properties, raw_geometry, status, version, is_latest
                ) VALUES (
                    $1, $2, $3, $4, $5, $6, $7, ST_GeomFromGeoJSON($8), $9, $10, $11, $12, $13, $14, $15, $16, $17
                )
            `
			styleJSON, _ := json.Marshal(r.StyleConfig)
			rawPropsJSON, _ := json.Marshal(r.RawProperties)
			rawGeomJSON, _ := json.Marshal(r.RawGeometry)

			err := tx.Exec(sql,
				r.LayerID,
				r.Name,
				r.DisplayName,
				r.LandUseCode,
				r.LandUseName,
				r.LegalDoc,
				r.PlanningName,
				string(r.Geometry), // GeoJSON string
				r.ColorRed,
				r.ColorGreen,
				r.ColorBlue,
				string(styleJSON),
				string(rawPropsJSON),
				string(rawGeomJSON),
				r.Status,
				r.Version,
				r.IsLatest,
			).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
