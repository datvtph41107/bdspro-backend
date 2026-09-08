ALTER TABLE support_tickets
    ADD COLUMN IF NOT EXISTS handling_team INT NULL;

CREATE INDEX IF NOT EXISTS idx_support_tickets_handling_team
    ON support_tickets(handling_team)
    WHERE deleted_at IS NULL AND handling_team IS NOT NULL;
