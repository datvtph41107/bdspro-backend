BEGIN;

DELETE FROM role_permissions
 WHERE permission_id IN (299, 300, 301);

DELETE FROM group_permissions
 WHERE permission_id IN (299, 300, 301);

DELETE FROM permissions
 WHERE id IN (299, 300, 301)
   AND key IN (
       'COMMERCIAL_USAGE_VIEW',
       'COMMERCIAL_USAGE_RECONCILE',
       'COMMERCIAL_SUBSCRIPTION_VIEW'
   );

COMMIT;
