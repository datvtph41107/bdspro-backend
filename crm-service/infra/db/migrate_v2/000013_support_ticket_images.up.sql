-- 000013_support_ticket_images.up.sql
ALTER TABLE support_tickets
    ADD COLUMN IF NOT EXISTS images JSONB NOT NULL DEFAULT '[]'::jsonb;
