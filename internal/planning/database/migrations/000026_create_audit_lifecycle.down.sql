-- Rollback 000019: Remove audit and lifecycle tracking tables
DROP TABLE IF EXISTS qh_lifecycle_transitions;
DROP TABLE IF EXISTS qh_audit_entries;
