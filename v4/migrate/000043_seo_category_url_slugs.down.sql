UPDATE seo_category SET slug = 'tin-tuc-quy-hoach', updated_at = NOW()
WHERE code = 'planning_news' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'du-an-quy-hoach', updated_at = NOW()
WHERE code = 'planning_project' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'van-ban-phap-ly', updated_at = NOW()
WHERE code = 'legal_document' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'bao-cao-phan-tich', updated_at = NOW()
WHERE code = 'analysis_report' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'vung-quy-hoach', updated_at = NOW()
WHERE code = 'planning_region' AND deleted_at IS NULL;
