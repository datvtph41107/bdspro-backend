DROP INDEX IF EXISTS idx_support_tickets_handling_team;
ALTER TABLE support_tickets DROP COLUMN IF EXISTS handling_team;
