-- Intentionally no-op: permission identities may already be assigned to
-- operator roles and rollback must not silently remove durable authority.
SELECT 1;
