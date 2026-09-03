-- Intentionally non-destructive.
--
-- Version 000041 adopts tables that may already contain production data and
-- may predate the migration ledger because they were created by GORM.  A down
-- migration cannot prove which rows or tables were introduced by this version,
-- so deleting them would violate ownership and recovery safety.  Rolling back
-- the application is supported because all additions are backward compatible;
-- the schema version can then be moved down/up without data loss.
SELECT 1;
