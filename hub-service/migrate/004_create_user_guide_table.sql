-- Migration: Create user guide table
-- Created: 2025-10-28

-- Create tb_user_guide table
CREATE TABLE IF NOT EXISTS tb_user_guide (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    group_key VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

-- Create indexes for tb_user_guide
CREATE INDEX IF NOT EXISTS idx_user_guide_deleted_at ON tb_user_guide(deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_guide_group_key ON tb_user_guide(group_key);
CREATE INDEX IF NOT EXISTS idx_user_guide_title ON tb_user_guide(title);
CREATE INDEX IF NOT EXISTS idx_user_guide_created_at ON tb_user_guide(created_at DESC);

-- Comments
COMMENT ON TABLE tb_user_guide IS 'Stores user guide and help documentation';
COMMENT ON COLUMN tb_user_guide.title IS 'Title of the user guide';
COMMENT ON COLUMN tb_user_guide.description IS 'Detailed description/content of the user guide';
COMMENT ON COLUMN tb_user_guide.group_key IS 'Group key to categorize user guides';

