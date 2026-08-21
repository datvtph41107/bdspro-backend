BEGIN;

CREATE TABLE IF NOT EXISTS user_followed_seo_category (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    seo_category_id BIGINT NOT NULL REFERENCES seo_category(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_user_followed_seo_category UNIQUE (user_id, seo_category_id)
);

CREATE INDEX IF NOT EXISTS idx_user_followed_seo_category_user
    ON user_followed_seo_category(user_id, updated_at DESC, id DESC)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE user_followed_seo_category IS
    'Danh sách chủ đề (seo_category) mà từng người dùng theo dõi.';
COMMENT ON COLUMN user_followed_seo_category.user_id IS
    'Profile ID của người dùng đã xác thực thực hiện thao tác theo dõi.';
COMMENT ON COLUMN user_followed_seo_category.seo_category_id IS
    'Định danh chủ đề trong bảng seo_category.';

COMMIT;
