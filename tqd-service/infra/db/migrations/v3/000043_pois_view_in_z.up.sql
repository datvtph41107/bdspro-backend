ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS view_in_z INTEGER NOT NULL DEFAULT 16;

UPDATE pois
SET view_in_z = 16
WHERE view_in_z IS NULL OR view_in_z = 0;
