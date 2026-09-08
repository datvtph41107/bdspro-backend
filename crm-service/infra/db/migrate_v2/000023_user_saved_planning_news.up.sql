BEGIN;

CREATE TABLE IF NOT EXISTS user_saved_planning_news (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    seo_domain_id BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_user_saved_planning_news UNIQUE (user_id, seo_domain_id)
);

CREATE INDEX IF NOT EXISTS idx_user_saved_planning_news_user
    ON user_saved_planning_news(user_id, updated_at DESC, id DESC)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE user_saved_planning_news IS
    'Danh sách tin quy hoạch trong seo_domain được từng người dùng lưu để xem lại.';
COMMENT ON COLUMN user_saved_planning_news.user_id IS
    'Profile ID của người dùng đã xác thực thực hiện thao tác lưu tin.';
COMMENT ON COLUMN user_saved_planning_news.seo_domain_id IS
    'Định danh tin quy hoạch do CRM sở hữu trong bảng seo_domain.';

COMMIT;
