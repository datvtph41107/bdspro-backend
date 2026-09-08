DROP TABLE IF EXISTS seo_user_need_evidence;
DROP TABLE IF EXISTS seo_user_need;
DROP TABLE IF EXISTS seo_audience_segment;

ALTER TABLE seo_objective_evidence
DROP CONSTRAINT IF EXISTS
    uq_seo_objective_evidence_id_objective;