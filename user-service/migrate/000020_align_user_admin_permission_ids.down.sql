BEGIN;

DELETE FROM role_permissions
 WHERE permission_id IN (309, 310);

DELETE FROM group_permissions
 WHERE permission_id IN (309, 310);

DELETE FROM permissions
 WHERE id IN (309, 310)
   AND key IN ('USER_ADMIN_VIEW', 'USER_ADMIN_MANAGE');

COMMIT;
