-- Align seo_category.slug with public URL <loại> segments:
-- /tin-tuc/<loại>/<id>/<slug>

UPDATE seo_category SET slug = 'tin', updated_at = NOW()
WHERE code = 'planning_news' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'do-an-quy-hoach', updated_at = NOW()
WHERE code = 'planning_project' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'tai-lieu', updated_at = NOW()
WHERE code = 'legal_document' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'bao-cao', updated_at = NOW()
WHERE code = 'analysis_report' AND deleted_at IS NULL;

UPDATE seo_category SET slug = 'vung-quy-hoach', updated_at = NOW()
WHERE code = 'planning_region' AND deleted_at IS NULL;
