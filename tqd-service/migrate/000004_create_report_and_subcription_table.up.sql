-- user_subscriptions
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id VARCHAR(255) NOT NULL,
    subscription_scope JSONB NOT NULL,
    trigger_config JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_triggered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_sub_user_target ON user_subscriptions(user_id, target_type, target_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sub_target ON user_subscriptions(target_type, target_id);
CREATE INDEX idx_sub_user_status ON user_subscriptions(user_id, status);

-- user_notifications
CREATE TABLE IF NOT EXISTS user_notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    subscription_id BIGINT,
    change_event_id BIGINT,
    title VARCHAR(255) NOT NULL,
    message TEXT[] NOT NULL,
    notification_type INT NOT NULL,
    severity INT NOT NULL DEFAULT 20,
    action_link TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    channels VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_notif_user_read ON user_notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_notif_expires ON user_notifications(expires_at);

CREATE TABLE IF NOT EXISTS reports (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    report_type INT NOT NULL,
    target_id VARCHAR(255),
    target_snapshot JSONB,
    profile INT NOT NULL,
    format INT NOT NULL DEFAULT 10,
    status INT NOT NULL DEFAULT 10,
    file_url TEXT,
    file_size BIGINT,
    file_hash VARCHAR(64),
    generation_duration_ms INT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_report_user_status ON reports(user_id, status, created_at DESC);
CREATE INDEX idx_report_expires ON reports(expires_at);
