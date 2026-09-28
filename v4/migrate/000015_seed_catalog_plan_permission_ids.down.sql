BEGIN;

DELETE FROM role_permissions
 WHERE permission_id IN (302, 303, 304);

DELETE FROM group_permissions
 WHERE permission_id IN (302, 303, 304);

DELETE FROM permissions
 WHERE id IN (302, 303, 304)
   AND key IN ('CATALOG_PLAN_VIEW', 'CATALOG_PLAN_MANAGE', 'CATALOG_PLAN_PUBLISH');

COMMIT;
