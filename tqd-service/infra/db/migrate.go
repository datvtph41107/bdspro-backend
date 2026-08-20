package db

import (
	"log"

	_db "common/db"
)

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
		`ALTER TABLE province_v2
			ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10,
			ADD COLUMN IF NOT EXISTS description VARCHAR(255),
			ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326),
			ADD COLUMN IF NOT EXISTS seo_id BIGINT,
			ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS plan_doc_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_provinces_qh_order
			ON province_v2(qh_order DESC, full_name ASC)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_geometry ON province_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_seo_id
			ON province_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_plan_doc_id
			ON province_v2(plan_doc_id)
			WHERE plan_doc_id IS NOT NULL AND plan_doc_id > 0`,
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
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(255)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_geometry ON ward_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_seo_id
			ON ward_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
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

func AutoMigrate() {
	// EnsurePostGIS()
	ensureProvinceV2Columns()
	ensureWardV2Columns()
	ensurePoisGeom()
	ensurePlanningProjectClientView()
	ensurePlanningDocumentSeoID()
	ensurePlanningDocumentLegalMetadata()
	ensurePlanningDocumentEvents()
	ensureDiscoveryParcelSearchIndexes()
	// province_v2 / ward_v2 là bảng legacy — chỉ cập nhật schema qua ensure*() ở trên,
	// không dùng GORM AutoMigrate (GORM cố DROP uni_*_code không tồn tại).
	// _db.DB.AutoMigrate(
	// 	&qh_domain.Parcel{},
	// 	&qh_domain.QHLayer{},
	// 	// &qh_domain.Parcel{},
	// )
}
