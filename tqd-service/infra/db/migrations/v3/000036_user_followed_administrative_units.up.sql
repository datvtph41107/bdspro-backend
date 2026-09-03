CREATE TABLE IF NOT EXISTS user_followed_administrative_units (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    composite_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_followed_administrative_units
        UNIQUE (user_id, composite_id)
);

CREATE INDEX IF NOT EXISTS idx_user_followed_administrative_units_user
    ON user_followed_administrative_units(user_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_user_followed_administrative_units_composite
    ON user_followed_administrative_units(composite_id, user_id);

COMMENT ON TABLE user_followed_administrative_units IS
    'User follow đơn vị hành chính theo compositeId ({sourceId}{1=tỉnh|2=phường}). Unfollow xóa cứng.';
