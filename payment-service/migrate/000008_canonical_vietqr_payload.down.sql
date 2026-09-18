BEGIN;

-- Retired image attempts remain expired; their lossy presentation blob is not
-- restored. Rollback only removes the canonical representation constraint.
ALTER TABLE commerce_payment_attempts
    DROP CONSTRAINT IF EXISTS commerce_payment_attempts_canonical_qr_payload_ck;

COMMIT;
