BEGIN;

DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE key IN ('IAM_ROLE_VIEW', 'IAM_ROLE_MANAGE')
);
DELETE FROM permissions WHERE key IN ('IAM_ROLE_VIEW', 'IAM_ROLE_MANAGE');

COMMIT;
