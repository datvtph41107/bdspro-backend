BEGIN;

-- Provider/bank descriptions may normalize ASCII case. Order reference
-- identity is therefore case-insensitive even for orders created before the
-- generator switched to canonical uppercase output.
CREATE UNIQUE INDEX IF NOT EXISTS uq_commerce_orders_reference_canonical
    ON commerce_orders (UPPER(reference));

COMMIT;
