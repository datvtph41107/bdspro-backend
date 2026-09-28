BEGIN;

DROP INDEX IF EXISTS idx_demo_deal_log_profile_id;
DROP TABLE IF EXISTS demo_deal_log;

DROP INDEX IF EXISTS idx_demo_appointment_log_profile_id;
DROP TABLE IF EXISTS demo_appointment_log;

COMMIT;
