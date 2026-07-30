-- =====================================================
-- MIGRATION DOWN: DROP qh_labels
-- =====================================================

DROP TRIGGER IF EXISTS trg_qh_labels_updated_at ON qh_labels;
DROP FUNCTION IF EXISTS update_qh_labels_updated_at();
DROP TABLE IF EXISTS qh_labels;
