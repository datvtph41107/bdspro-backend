-- 000013_support_ticket_images.down.sql
ALTER TABLE support_tickets DROP COLUMN IF EXISTS images;
