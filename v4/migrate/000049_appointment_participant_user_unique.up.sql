-- Remove orphan rows from failed inserts (appointment_id = 0).
DELETE FROM appointment_participants
WHERE appointment_id IS NULL OR appointment_id = 0;

-- Unique (appointment_id, user_id) for participant upsert (ON CONFLICT).
DROP INDEX IF EXISTS idx_appointment_contact_unique;
DROP INDEX IF EXISTS idx_appointment_user_unique;

CREATE UNIQUE INDEX IF NOT EXISTS idx_appointment_user_unique
ON appointment_participants (appointment_id, user_id);

-- Resync serial after manual/seed inserts with explicit ids.
SELECT setval(
    pg_get_serial_sequence('appointment_participants', 'id'),
    COALESCE((SELECT MAX(id) FROM appointment_participants), 1)
);
