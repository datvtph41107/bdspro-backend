-- Cấp quyền cho user_owner trên các bảng auth/user mà app cần (OAuth callback, login, session...).
-- Chạy bằng superuser/owner (vd: admin) trên DB production:
--   psql -h 103.172.239.89 -U admin -d db_user -f 011_grant_user_owner_auth_tables.sql

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE
    admin_access_control,
    auth_config,
    auth_devices,
    auth_method,
    group_permissions,
    permissions,
    role_groups,
    role_permissions,
    role_profiles,
    roles,
    user_auth,
    user_info,
    user_otp,
    user_pin,
    user_session,
    user_status
TO user_owner;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO user_owner;

-- Bảng tạo sau này bởi admin cũng tự có quyền cho user_owner
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO user_owner;

ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO user_owner;
