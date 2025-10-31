-- Migration for user_guide_steps table
CREATE TABLE IF NOT EXISTS user_guide_steps (
    id BIGSERIAL PRIMARY KEY,
    user_guide_id BIGINT NOT NULL,
    step_order INT NOT NULL,
    image TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    
    CONSTRAINT fk_user_guide
        FOREIGN KEY(user_guide_id) 
        REFERENCES user_guides(id)
        ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_guide_steps_user_guide_id ON user_guide_steps (user_guide_id);
CREATE INDEX IF NOT EXISTS idx_user_guide_steps_step_order ON user_guide_steps (step_order);
CREATE INDEX IF NOT EXISTS idx_user_guide_steps_deleted_at ON user_guide_steps (deleted_at);

-- Trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_user_guide_steps_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_guide_steps_updated_at
BEFORE UPDATE ON user_guide_steps
FOR EACH ROW
EXECUTE FUNCTION update_user_guide_steps_updated_at();

COMMENT ON TABLE user_guide_steps IS 'Bảng lưu trữ các bước trong user guide';
COMMENT ON COLUMN user_guide_steps.user_guide_id IS 'ID của user guide (foreign key)';
COMMENT ON COLUMN user_guide_steps.step_order IS 'Thứ tự bước (1, 2, 3...)';
COMMENT ON COLUMN user_guide_steps.image IS 'URL hình ảnh cho bước này';
COMMENT ON COLUMN user_guide_steps.content IS 'Nội dung mô tả bước này';

