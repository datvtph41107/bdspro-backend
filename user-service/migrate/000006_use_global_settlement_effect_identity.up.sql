BEGIN;

-- Payment order_id is a sequence local to the Payment database. It may be
-- reused after an independent restore/recreation while User correctly retains
-- older receipts. Cross-owner idempotency is therefore owned by effect_key,
-- which includes Payment's immutable order reference; order_id remains indexed
-- evidence for operations and audit, not a global uniqueness authority.
ALTER TABLE subscription_settlement_receipts
    DROP CONSTRAINT IF EXISTS subscription_settlement_receipts_order_id_key;

CREATE INDEX IF NOT EXISTS ix_subscription_settlement_receipts_order_id
    ON subscription_settlement_receipts(order_id);

COMMIT;
