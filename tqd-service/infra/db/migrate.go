package db

import (
	"embed"
	"io/fs"
	"log"
	"path"
	"sort"
	"strings"

	_db "common/db"
)

//go:embed migrations/v3/*.up.sql
var v3UpMigrations embed.FS

// CheckPostGIS kiểm tra PostGIS đã cài trong DB hiện tại (trả version, installed).
func CheckPostGIS() (version string, installed bool) {
	var extVersion string
	err := _db.DB.Raw(
		`SELECT extversion FROM pg_extension WHERE extname = 'postgis'`,
	).Scan(&extVersion).Error
	if err != nil || extVersion == "" {
		return "", false
	}
	var postgisVer string
	_ = _db.DB.Raw(`SELECT PostGIS_Version()`).Scan(&postgisVer).Error
	if postgisVer != "" {
		return postgisVer, true
	}
	return extVersion, true
}

// EnsurePostGIS bật extension PostGIS và unaccent trước khi tạo cột geometry (AutoMigrate / raw SQL).
func EnsurePostGIS() {
	if ver, ok := CheckPostGIS(); ok {
		log.Printf("✅ PostGIS đã có: %s", ver)
	} else {
		log.Println("⚠️  PostGIS chưa bật — thử CREATE EXTENSION postgis ...")
		if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`).Error; err != nil {
			log.Printf("❌ không bật được PostGIS: %v", err)
		}
	}

	log.Println("⚠️  Thử CREATE EXTENSION IF NOT EXISTS unaccent ...")
	if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS unaccent`).Error; err != nil {
		log.Printf("❌ không bật được unaccent: %v", err)
	}

	log.Println("⚠️  Thử CREATE EXTENSION IF NOT EXISTS pg_trgm ...")
	if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		log.Printf("❌ không bật được pg_trgm: %v", err)
	}
}

func ensureSchemaMigrationsTable() {
	if err := _db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`).Error; err != nil {
		log.Printf("⚠️  ensureSchemaMigrationsTable: %v", err)
	}
}

func isMigrationApplied(version string) bool {
	var count int64
	_ = _db.DB.Raw(
		`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version,
	).Scan(&count).Error
	return count > 0
}

func markMigrationApplied(version string) {
	if err := _db.DB.Exec(
		`INSERT INTO schema_migrations(version) VALUES (?) ON CONFLICT (version) DO NOTHING`,
		version,
	).Error; err != nil {
		log.Printf("⚠️  markMigrationApplied(%s): %v", version, err)
	}
}

// runPendingV3Migrations chạy các file migrations/v3/*.up.sql chưa có trong schema_migrations.
func runPendingV3Migrations() {
	ensureSchemaMigrationsTable()

	entries, err := fs.ReadDir(v3UpMigrations, "migrations/v3")
	if err != nil {
		log.Printf("⚠️  runPendingV3Migrations read dir: %v", err)
		return
	}

	var files []string
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(name, ".up.sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.TrimSuffix(name, ".up.sql")
		if isMigrationApplied(version) {
			continue
		}

		content, err := v3UpMigrations.ReadFile(path.Join("migrations/v3", name))
		if err != nil {
			log.Printf("⚠️  runPendingV3Migrations read %s: %v", name, err)
			continue
		}

		log.Printf("▶ applying migration %s ...", version)
		if err := _db.DB.Exec(string(content)).Error; err != nil {
			log.Printf("❌ migration %s failed: %v", version, err)
			continue
		}
		markMigrationApplied(version)
		log.Printf("✅ migration %s applied", version)
	}
}

// ensureDiscoveryParcelSearchIndexes tạo GIN/BTree phục vụ SearchParcels (idempotent).
// Không DROP index ở startup — tránh lock bảng lớn; xem migration 000012.
func ensureDiscoveryParcelSearchIndexes() {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_trgm`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm_nonempty
			ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops)
			WHERE COALESCE(adr_search, '') <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower_nonempty
			ON qh_parcel_info (LOWER(property_code))
			WHERE property_code IS NOT NULL AND property_code <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_map_land
			ON qh_parcel_info (map_number, land_number)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureDiscoveryParcelSearchIndexes: %v", err)
		}
	}
}

// ensureProvinceV2Columns bổ sung cột mới cho province_v2.
// Không AutoMigrate province_v2 — bảng legacy; GORM cố DROP uni_province_v2_code không tồn tại.
func ensureProvinceV2Columns() {
	statements := []string{
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(2000)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS plan_doc_id BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS ward_count INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_role VARCHAR(500) NOT NULL DEFAULT ''`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT ''`,
		`UPDATE province_v2
			SET official_name = full_name
			WHERE (official_name IS NULL OR official_name = '')
			  AND full_name IS NOT NULL
			  AND full_name <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_provinces_qh_order
			ON province_v2(qh_order DESC, full_name ASC)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_geometry ON province_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_seo_id
			ON province_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_plan_doc_id
			ON province_v2(plan_doc_id)
			WHERE plan_doc_id IS NOT NULL AND plan_doc_id > 0`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_deleted_at ON province_v2(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_division_type ON province_v2(division_type)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureProvinceV2Columns: %v", err)
		}
	}
	_ = _db.DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'chk_province_v2_zoom'
				  AND conrelid = 'province_v2'::regclass
			) THEN
				ALTER TABLE province_v2
					ADD CONSTRAINT chk_province_v2_zoom
					CHECK (zoom >= 5 AND zoom <= 22);
			END IF;
		END $$`).Error
}

// ensureWardV2Columns bổ sung cột mới cho ward_v2.
// Không AutoMigrate ward_v2 — bảng legacy; GORM cố DROP uni_ward_v2_code không tồn tại.
func ensureWardV2Columns() {
	statements := []string{
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(2000)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS summary VARCHAR(200)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population BIGINT`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS province_code VARCHAR(20)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_role VARCHAR(500) NOT NULL DEFAULT ''`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT ''`,
		`UPDATE ward_v2
			SET official_name = full_name
			WHERE (official_name IS NULL OR official_name = '')
			  AND full_name IS NOT NULL
			  AND full_name <> ''`,
		`UPDATE ward_v2 w
			SET province_code = p.code
			FROM province_v2 p
			WHERE w.province_id = p.id
			  AND (w.province_code IS NULL OR w.province_code = '')
			  AND p.code IS NOT NULL
			  AND p.code <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_geometry ON ward_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_seo_id
			ON ward_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_deleted_at ON ward_v2(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_division_type ON ward_v2(division_type)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_province_code ON ward_v2(province_code)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_unit_type
			ON ward_v2(unit_type)
			WHERE unit_type > 0`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_data_source
			ON ward_v2(data_source)
			WHERE data_source > 0`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureWardV2Columns: %v", err)
		}
	}
}

// ensurePoisGeom bổ sung geom (Point,4326) cho discovery search / identify / related-entity.
func ensurePoisGeom() {
	statements := []string{
		`ALTER TABLE pois ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326)`,
		`UPDATE pois
			SET geom = ST_SetSRID(
				ST_MakePoint(longitude::double precision, latitude::double precision),
				4326
			)
			WHERE geom IS NULL
			  AND latitude IS NOT NULL
			  AND longitude IS NOT NULL
			  AND latitude BETWEEN -90 AND 90
			  AND longitude BETWEEN -180 AND 180
			  AND NOT (latitude = 0 AND longitude = 0)`,
		`CREATE INDEX IF NOT EXISTS idx_pois_geom
			ON pois USING GIST (geom)
			WHERE deleted_at IS NULL AND geom IS NOT NULL`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_pois_name_trgm
			ON pois USING GIN (LOWER(name) gin_trgm_ops)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_pois_address_trgm
			ON pois USING GIN (LOWER(COALESCE(address, '')) gin_trgm_ops)
			WHERE deleted_at IS NULL`,
		`CREATE OR REPLACE FUNCTION pois_sync_geom_from_latlng()
			RETURNS trigger
			LANGUAGE plpgsql
			AS $$
			BEGIN
			    IF NEW.latitude IS NOT NULL
			       AND NEW.longitude IS NOT NULL
			       AND NEW.latitude BETWEEN -90 AND 90
			       AND NEW.longitude BETWEEN -180 AND 180
			       AND NOT (NEW.latitude = 0 AND NEW.longitude = 0)
			    THEN
			        NEW.geom := ST_SetSRID(
			            ST_MakePoint(NEW.longitude::double precision, NEW.latitude::double precision),
			            4326
			        );
			    ELSIF NEW.geom IS NOT NULL THEN
			        NEW.latitude := ST_Y(NEW.geom);
			        NEW.longitude := ST_X(NEW.geom);
			    END IF;
			    RETURN NEW;
			END;
			$$`,
		`DROP TRIGGER IF EXISTS trg_pois_sync_geom ON pois`,
		`CREATE TRIGGER trg_pois_sync_geom
			BEFORE INSERT OR UPDATE OF latitude, longitude, geom
			ON pois
			FOR EACH ROW
			EXECUTE FUNCTION pois_sync_geom_from_latlng()`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePoisGeom: %v", err)
		}
	}
}

func ensurePlanningProjectClientView() {
	if err := _db.DB.Exec(`
		ALTER TABLE qh_planning_projects
			ADD COLUMN IF NOT EXISTS client_view BOOLEAN NOT NULL DEFAULT true
	`).Error; err != nil {
		log.Printf("⚠️  ensurePlanningProjectClientView: %v", err)
	}
}

func ensurePlanningProjectSeoID() {
	statements := []string{
		`ALTER TABLE qh_planning_projects
			ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_seo_id
			ON qh_planning_projects(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePlanningProjectSeoID: %v", err)
		}
	}
}

func ensurePlanningDocumentSeoID() {
	statements := []string{
		`ALTER TABLE qh_planning_documents
			ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_documents_seo_id
			ON qh_planning_documents(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePlanningDocumentSeoID: %v", err)
		}
	}
}

func ensureLayerColumns() {
	statements := []string{
		`ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS province_id VARCHAR(50)`,
		`ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS ward_id VARCHAR(50)`,
		`CREATE INDEX IF NOT EXISTS idx_qh_layers_seo_id
			ON qh_layers(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
		`CREATE INDEX IF NOT EXISTS idx_qh_layers_province_id
			ON qh_layers(province_id)
			WHERE province_id IS NOT NULL AND province_id <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_qh_layers_ward_id
			ON qh_layers(ward_id)
			WHERE ward_id IS NOT NULL AND ward_id <> ''`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureLayerColumns: %v", err)
		}
	}
}

func ensurePlanningDocumentExt() {
	statements := []string{
		`ALTER TABLE qh_planning_documents
			ADD COLUMN IF NOT EXISTS ext VARCHAR(20) NOT NULL DEFAULT ''`,
		`UPDATE qh_planning_documents
			SET ext = LOWER(regexp_replace(
				COALESCE(NULLIF(BTRIM(relative_path), ''), NULLIF(BTRIM(filepath), ''), ''),
				'^.*\.', ''
			))
			WHERE deleted_at IS NULL
			  AND ext = ''
			  AND COALESCE(NULLIF(BTRIM(relative_path), ''), NULLIF(BTRIM(filepath), '')) ~ '\.[^./\\]+$'`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePlanningDocumentExt: %v", err)
		}
	}
}

func ensurePlanningDocumentLegalMetadata() {
	if err := _db.DB.Exec(`
		ALTER TABLE qh_planning_documents
			ADD COLUMN IF NOT EXISTS indicators JSONB NOT NULL DEFAULT '[]'::jsonb,
			ADD COLUMN IF NOT EXISTS authority VARCHAR(255) NOT NULL DEFAULT '',
			ADD COLUMN IF NOT EXISTS sector VARCHAR(255) NOT NULL DEFAULT '',
			ADD COLUMN IF NOT EXISTS field VARCHAR(255) NOT NULL DEFAULT '',
			ADD COLUMN IF NOT EXISTS signer VARCHAR(255) NOT NULL DEFAULT ''
	`).Error; err != nil {
		log.Printf("⚠️  ensurePlanningDocumentLegalMetadata: %v", err)
	}
}

func ensurePlanningDocumentEvents() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS qh_planning_document_events (
			id BIGSERIAL PRIMARY KEY,
			planning_document_id BIGINT NOT NULL REFERENCES qh_planning_documents(id) ON DELETE CASCADE,
			event_type VARCHAR(80),
			event_name VARCHAR(255) NOT NULL,
			description TEXT,
			event_date TIMESTAMPTZ,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ NULL,
			created_by BIGINT,
			updated_by BIGINT
		)`,
		`ALTER TABLE qh_planning_document_events
			ADD COLUMN IF NOT EXISTS planning_document_id BIGINT,
			ADD COLUMN IF NOT EXISTS event_type VARCHAR(80),
			ADD COLUMN IF NOT EXISTS event_name VARCHAR(255),
			ADD COLUMN IF NOT EXISTS description TEXT,
			ADD COLUMN IF NOT EXISTS event_date TIMESTAMPTZ,
			ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL,
			ADD COLUMN IF NOT EXISTS created_by BIGINT,
			ADD COLUMN IF NOT EXISTS updated_by BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_document_events_doc_timeline
			ON qh_planning_document_events(planning_document_id, event_date DESC, id DESC)
			WHERE deleted_at IS NULL AND planning_document_id IS NOT NULL`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePlanningDocumentEvents: %v", err)
		}
	}
}

func ensureWardNeighbors() {
	statements := []string{
		// ward_v2 legacy thiếu PK → FK ward_neighbors không tạo được
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conrelid = 'ward_v2'::regclass
				  AND contype = 'p'
			) THEN
				ALTER TABLE ward_v2 ADD CONSTRAINT ward_v2_pkey PRIMARY KEY (id);
			END IF;
		END $$`,
		`CREATE TABLE IF NOT EXISTS ward_neighbors (
			id BIGSERIAL PRIMARY KEY,
			ward_id BIGINT NOT NULL REFERENCES ward_v2(id) ON DELETE CASCADE,
			neighbor_ward_id BIGINT NOT NULL REFERENCES ward_v2(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ NULL,
			CONSTRAINT uq_ward_neighbors UNIQUE (ward_id, neighbor_ward_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_neighbors_ward
			ON ward_neighbors(ward_id)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_ward_neighbors_neighbor
			ON ward_neighbors(neighbor_ward_id)
			WHERE deleted_at IS NULL`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureWardNeighbors: %v", err)
		}
	}
}

// ensurePlanningProjectWards tạo bảng nối đồ án ↔ phường (giống qh_planning_project_provinces).
// ward_id VARCHAR(50) — khớp code đang decode compositeId thành string; không FK sang ward_v2(id bigint).
func ensurePlanningProjectWards() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS qh_planning_project_wards (
			id BIGSERIAL PRIMARY KEY,
			planning_project_id BIGINT NOT NULL,
			ward_id VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ NULL,
			created_by BIGINT,
			updated_by BIGINT,
			CONSTRAINT uq_qh_planning_project_wards UNIQUE (planning_project_id, ward_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_project_wards_project
			ON qh_planning_project_wards(planning_project_id)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_project_wards_ward
			ON qh_planning_project_wards(ward_id)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_qh_planning_project_wards_deleted_at
			ON qh_planning_project_wards(deleted_at)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensurePlanningProjectWards: %v", err)
		}
	}
}

// ensureLayerPlanDocRelation tạo bảng quan hệ layer ↔ tài liệu quy hoạch (migration 000015).
func ensureLayerPlanDocRelation() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS layer_plan_doc_relation (
			id                   BIGSERIAL PRIMARY KEY,
			layer_id             BIGINT NOT NULL
				REFERENCES qh_layers(id) ON DELETE CASCADE,
			planning_document_id BIGINT NOT NULL
				REFERENCES qh_planning_documents(id) ON DELETE CASCADE,
			sort_order           INT NOT NULL DEFAULT 0,
			created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at           TIMESTAMPTZ NULL,
			created_by           BIGINT,
			updated_by           BIGINT,
			CONSTRAINT uq_layer_plan_doc_relation
				UNIQUE (layer_id, planning_document_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_layer_plan_doc_relation_layer_order
			ON layer_plan_doc_relation (layer_id, sort_order ASC, id ASC)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_layer_plan_doc_relation_document
			ON layer_plan_doc_relation (planning_document_id)
			WHERE deleted_at IS NULL`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureLayerPlanDocRelation: %v", err)
		}
	}
}

func ensureProAIJobs() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS pro_ai_jobs (
			id                   BIGSERIAL PRIMARY KEY,
			job_type             INT4 NOT NULL DEFAULT 10,
			name                 VARCHAR(255) NOT NULL DEFAULT '',
			source_folder_name   VARCHAR(255),
			process_status       INT4 NOT NULL DEFAULT 10,
			planning_project_id  INT8,
			classify_error       TEXT,
			metadata             JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at           TIMESTAMPTZ,
			updated_at           TIMESTAMPTZ,
			deleted_at           TIMESTAMPTZ,
			created_by           INT8,
			updated_by           INT8,
			completed_at         TIMESTAMPTZ,
			approved_at          TIMESTAMPTZ
		)`,
		`ALTER TABLE pro_ai_jobs ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ`,
		`ALTER TABLE pro_ai_jobs ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_job_type ON pro_ai_jobs(job_type)`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_process_status ON pro_ai_jobs(process_status)`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_planning_project_id ON pro_ai_jobs(planning_project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_deleted_at ON pro_ai_jobs(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_completed_at ON pro_ai_jobs(completed_at)`,
		`CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_approved_at ON pro_ai_jobs(approved_at)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureProAIJobs: %v", err)
		}
	}
}

func AutoMigrate() {
	// Chạy file SQL v3 chưa apply (tracked qua schema_migrations).
	runPendingV3Migrations()

	// ensure*() bổ sung cột/index hot-path idempotent (legacy / không phụ thuộc version file).
	ensureProvinceV2Columns()
	ensureWardV2Columns()
	ensureWardNeighbors()
	ensurePlanningProjectWards()
	ensureLayerPlanDocRelation()
	ensureProAIJobs()
	ensurePoisGeom()
	ensurePlanningProjectClientView()
	ensurePlanningProjectSeoID()
	ensurePlanningDocumentSeoID()
	ensureLayerColumns()
	ensurePlanningDocumentExt()
	ensurePlanningDocumentLegalMetadata()
	ensurePlanningDocumentEvents()
	ensureDiscoveryParcelSearchIndexes()
	// province_v2 / ward_v2 là bảng legacy — chỉ cập nhật schema qua ensure*() ở trên,
	// không dùng GORM AutoMigrate (GORM cố DROP uni_*_code không tồn tại).
}
