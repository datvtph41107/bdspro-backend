BEGIN;

DROP INDEX IF EXISTS ix_subscription_settlement_receipts_order_id;

ALTER TABLE subscription_settlement_receipts
    ADD CONSTRAINT subscription_settlement_receipts_order_id_key UNIQUE (order_id);

COMMIT;
