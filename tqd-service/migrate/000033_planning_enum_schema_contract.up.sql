BEGIN;

-- Align the physical planning schema with the Go domain and Proto contract.
-- Existing non-null UUID jurisdiction values cannot be mapped safely to the
-- bigint qh_jurisdictions identity and therefore fail closed for manual repair.
DO $$
DECLARE
    column_type TEXT;
BEGIN
    SELECT data_type INTO column_type
    FROM information_schema.columns
    WHERE table_schema='public' AND table_name='qh_planning_projects' AND column_name='jurisdiction_id';

    IF column_type = 'uuid' THEN
        IF EXISTS (SELECT 1 FROM qh_planning_projects WHERE jurisdiction_id IS NOT NULL) THEN
            RAISE EXCEPTION 'qh_planning_projects.jurisdiction_id contains UUID values; map them to qh_jurisdictions bigint IDs before migration 000007';
        END IF;
        ALTER TABLE qh_planning_projects
            ALTER COLUMN jurisdiction_id TYPE BIGINT USING NULL::BIGINT;
    END IF;
END $$;

ALTER TABLE qh_planning_projects
    ALTER COLUMN planning_type TYPE INT4 USING (
        CASE LOWER(BTRIM(planning_type::text))
            WHEN 'qh sử dụng đất' THEN 10 WHEN 'su-dung-dat' THEN 10 WHEN 'used' THEN 10
            WHEN 'qh chung' THEN 20 WHEN 'quy-hoach-chung' THEN 20 WHEN 'general' THEN 20
            WHEN 'qh phân khu' THEN 30 WHEN 'quy-hoach-phan-khu' THEN 30 WHEN 'urban' THEN 30
            WHEN 'qh chi tiết' THEN 40 WHEN 'quy-hoach-chi-tiet' THEN 40 WHEN 'detail' THEN 40
            ELSE CASE WHEN planning_type::text ~ '^[0-9]+$' THEN planning_type::text::INT4 ELSE 20 END
        END
    ),
    ALTER COLUMN planning_level TYPE INT4 USING (
        CASE LOWER(BTRIM(planning_level::text))
            WHEN 'quốc gia' THEN 10 WHEN 'quoc-gia' THEN 10 WHEN 'national' THEN 10
            WHEN 'khu vực' THEN 15 WHEN 'khu-vuc' THEN 15 WHEN 'regional' THEN 15
            WHEN 'tỉnh' THEN 20 WHEN 'tinh' THEN 20 WHEN 'provincial' THEN 20
            WHEN 'huyện' THEN 30 WHEN 'huyen' THEN 30 WHEN 'district' THEN 30
            WHEN 'xã' THEN 40 WHEN 'xa' THEN 40 WHEN 'ward' THEN 40
            WHEN 'dự án' THEN 45 WHEN 'du-an' THEN 45 WHEN 'project' THEN 45
            ELSE CASE WHEN planning_level::text ~ '^[0-9]+$' THEN planning_level::text::INT4 ELSE 20 END
        END
    );

ALTER TABLE qh_planning_documents
    ALTER COLUMN document_type TYPE INT4 USING (
        CASE LOWER(BTRIM(document_type::text))
            WHEN 'pháp lý' THEN 10 WHEN 'phap-ly' THEN 10 WHEN 'legal' THEN 10
            WHEN 'thuyết minh' THEN 20 WHEN 'thuyet-minh' THEN 20 WHEN 'explanation' THEN 20
            WHEN 'bản đồ' THEN 30 WHEN 'ban-do' THEN 30 WHEN 'map' THEN 30
            WHEN 'cad' THEN 40
            WHEN 'gis' THEN 50 WHEN 'gis gốc' THEN 50
            WHEN 'phụ lục' THEN 60 WHEN 'phu-luc' THEN 60 WHEN 'appendix' THEN 60
            WHEN 'khác' THEN 70 WHEN 'khac' THEN 70 WHEN 'other' THEN 70
            ELSE CASE WHEN document_type::text ~ '^[0-9]+$' THEN document_type::text::INT4 ELSE 70 END
        END
    );

ALTER TABLE qh_planning_projects
    DROP CONSTRAINT IF EXISTS chk_qh_planning_project_type,
    DROP CONSTRAINT IF EXISTS chk_qh_planning_project_level;
ALTER TABLE qh_planning_projects
    ADD CONSTRAINT chk_qh_planning_project_type CHECK (planning_type IN (10,20,30,40)),
    ADD CONSTRAINT chk_qh_planning_project_level CHECK (planning_level IN (10,15,20,30,40,45));

ALTER TABLE qh_planning_documents
    DROP CONSTRAINT IF EXISTS chk_qh_planning_document_type;
ALTER TABLE qh_planning_documents
    ADD CONSTRAINT chk_qh_planning_document_type CHECK (document_type IN (10,20,30,40,50,60,70));

COMMIT;
