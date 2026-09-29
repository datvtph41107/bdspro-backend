BEGIN;

CREATE TABLE IF NOT EXISTS notification (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,
  owner_of BIGINT NOT NULL DEFAULT 0,
  owner_id BIGINT NOT NULL DEFAULT 0,
  title VARCHAR(255) NOT NULL DEFAULT '',
  message VARCHAR[] NOT NULL DEFAULT '{}',
  attach_data VARCHAR[] NOT NULL DEFAULT '{}',
  target_id BIGINT NOT NULL DEFAULT 0,
  notification_type BIGINT NOT NULL DEFAULT 0,
  avatar VARCHAR(255) NOT NULL DEFAULT '',
  visible_at TIMESTAMPTZ,
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_notification_owner ON notification(owner_id, owner_of);
CREATE INDEX IF NOT EXISTS idx_notification_target ON notification(target_id);
CREATE INDEX IF NOT EXISTS idx_notification_type ON notification(notification_type);
CREATE INDEX IF NOT EXISTS idx_notification_deleted ON notification(deleted_at);
CREATE INDEX IF NOT EXISTS idx_notification_visible ON notification(visible_at);

CREATE TABLE IF NOT EXISTS notification_history (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  owner_of BIGINT NOT NULL DEFAULT 0,
  owner_id BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_notification_history_owner ON notification_history(owner_id, owner_of);

CREATE TABLE IF NOT EXISTS noti_histories (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  target_id BIGINT NOT NULL DEFAULT 0,
  target_type BIGINT NOT NULL DEFAULT 0,
  action_type BIGINT NOT NULL DEFAULT 0,
  title TEXT NOT NULL DEFAULT '',
  note TEXT[] NOT NULL DEFAULT '{}',
  pre_stage TEXT NOT NULL DEFAULT '',
  after_stage TEXT NOT NULL DEFAULT '',
  owner_id BIGINT,
  owner_type BIGINT NOT NULL DEFAULT 0,
  is_internal BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_noti_histories_target ON noti_histories(target_id, target_type);

CREATE TABLE IF NOT EXISTS admin_histories (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  target_id BIGINT NOT NULL DEFAULT 0,
  target_type BIGINT NOT NULL DEFAULT 0,
  action_type BIGINT NOT NULL DEFAULT 0,
  title TEXT NOT NULL DEFAULT '',
  note TEXT[] NOT NULL DEFAULT '{}',
  pre_stage TEXT NOT NULL DEFAULT '',
  after_stage TEXT NOT NULL DEFAULT '',
  admin_id BIGINT NOT NULL DEFAULT 0,
  owner_id BIGINT,
  owner_of BIGINT NOT NULL DEFAULT 0,
  is_internal BOOLEAN NOT NULL DEFAULT FALSE,
  admin_role TEXT NOT NULL DEFAULT '',
  ip_address TEXT NOT NULL DEFAULT '',
  user_agent TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_admin_histories_target ON admin_histories(target_id, target_type);
CREATE INDEX IF NOT EXISTS idx_admin_histories_admin ON admin_histories(admin_id);

CREATE TABLE IF NOT EXISTS deal_history (
  id BIGSERIAL PRIMARY KEY,
  deal_id BIGINT NOT NULL,
  actor_id BIGINT NOT NULL,
  actor_name VARCHAR(255) NOT NULL DEFAULT '',
  actor_avatar VARCHAR(500) NOT NULL DEFAULT '',
  action_type VARCHAR(50) NOT NULL,
  action_name VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_deal_history_deal ON deal_history(deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_history_actor ON deal_history(actor_id);
CREATE INDEX IF NOT EXISTS idx_deal_history_action ON deal_history(action_type);

CREATE TABLE IF NOT EXISTS history_auth (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  user_id BIGINT NOT NULL,
  organization_id BIGINT,
  action_type VARCHAR(100) NOT NULL,
  action_name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  success BOOLEAN NOT NULL DEFAULT TRUE,
  reason VARCHAR(255) NOT NULL DEFAULT '',
  ip_address VARCHAR(45) NOT NULL DEFAULT '',
  user_agent VARCHAR(500) NOT NULL DEFAULT '',
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  performed_by BIGINT,
  session_id VARCHAR(255) NOT NULL DEFAULT '',
  channel VARCHAR(100) NOT NULL DEFAULT '',
  device_id VARCHAR(255) NOT NULL DEFAULT '',
  location VARCHAR(255) NOT NULL DEFAULT '',
  additional_note VARCHAR(500) NOT NULL DEFAULT '',
  source_service VARCHAR(255) NOT NULL DEFAULT 'auth-service'
);
CREATE INDEX IF NOT EXISTS idx_history_auth_user ON history_auth(user_id);
CREATE INDEX IF NOT EXISTS idx_history_auth_org ON history_auth(organization_id);
CREATE INDEX IF NOT EXISTS idx_history_auth_performed_by ON history_auth(performed_by);

CREATE TABLE IF NOT EXISTS tb_history (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT '',
  image VARCHAR(255) NOT NULL DEFAULT '',
  target VARCHAR(255) NOT NULL DEFAULT '',
  error TEXT NOT NULL DEFAULT '',
  response TEXT NOT NULL DEFAULT '',
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS person_configs (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  user_id BIGINT NOT NULL,
  key VARCHAR(255) NOT NULL,
  checked BOOLEAN NOT NULL DEFAULT FALSE,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  channel BIGINT NOT NULL DEFAULT 10,
  CONSTRAINT uq_person_configs_user_key_channel UNIQUE(user_id, key, channel)
);
CREATE INDEX IF NOT EXISTS idx_person_configs_user ON person_configs(user_id);

CREATE TABLE IF NOT EXISTS property_histories (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  subject_id BIGINT NOT NULL,
  action BIGINT NOT NULL,
  actor_id BIGINT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  occurred_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_property_histories_subject ON property_histories(subject_id);
CREATE INDEX IF NOT EXISTS idx_property_histories_actor ON property_histories(actor_id);
CREATE INDEX IF NOT EXISTS idx_property_histories_occurred ON property_histories(occurred_at);

CREATE TABLE IF NOT EXISTS account_warning_templates (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  name VARCHAR(255) NOT NULL,
  warning_type VARCHAR(50) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  severity VARCHAR(20) NOT NULL DEFAULT 'medium',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_by_user BIGINT NOT NULL DEFAULT 0,
  updated_by_user BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS account_warnings (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  target_id BIGINT NOT NULL,
  target_type VARCHAR(20) NOT NULL,
  warning_type VARCHAR(50) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  severity VARCHAR(20) NOT NULL DEFAULT 'medium',
  status VARCHAR(20) NOT NULL DEFAULT 'sent',
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  read_at TIMESTAMPTZ,
  sent_by BIGINT NOT NULL,
  sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  email_sent BOOLEAN NOT NULL DEFAULT FALSE,
  email_sent_at TIMESTAMPTZ,
  related_id BIGINT,
  related_type VARCHAR(50) NOT NULL DEFAULT '',
  expires_at TIMESTAMPTZ,
  template_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_account_warnings_target ON account_warnings(target_id);
CREATE INDEX IF NOT EXISTS idx_account_warnings_template ON account_warnings(template_id);

CREATE TABLE IF NOT EXISTS account_warning_logs (
  id BIGSERIAL PRIMARY KEY,
  created_by BIGINT, updated_by BIGINT,
  created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  warning_id BIGINT NOT NULL,
  target_id BIGINT NOT NULL,
  action VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL,
  reason VARCHAR(500) NOT NULL DEFAULT '',
  performed_by BIGINT NOT NULL,
  performed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  ip_address VARCHAR(45) NOT NULL DEFAULT '',
  user_agent VARCHAR(500) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_account_warning_logs_warning ON account_warning_logs(warning_id);
CREATE INDEX IF NOT EXISTS idx_account_warning_logs_target ON account_warning_logs(target_id);

COMMIT;
