-- Intentionally no-op. These IAM permission identities may already be used by
-- roles/groups and are also part of the canonical base schema. Rolling back the
-- repair must never delete shared authorization state.
SELECT 1;
