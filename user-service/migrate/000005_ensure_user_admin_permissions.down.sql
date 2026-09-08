-- Intentionally no-op. User Admin IAM identities may already be assigned to
-- operator roles; deleting them during rollback would remove durable authority.
SELECT 1;
