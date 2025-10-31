-- Migration for FAQs table
CREATE TABLE IF NOT EXISTS faqs (
    id BIGSERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    group_key VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_faqs_group_key ON faqs (group_key);
CREATE INDEX IF NOT EXISTS idx_faqs_deleted_at ON faqs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_faqs_question ON faqs USING gin(to_tsvector('english', question));

-- Trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_faqs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_faqs_updated_at
BEFORE UPDATE ON faqs
FOR EACH ROW
EXECUTE FUNCTION update_faqs_updated_at();

COMMENT ON TABLE faqs IS 'Bảng lưu trữ câu hỏi thường gặp (FAQs)';
COMMENT ON COLUMN faqs.question IS 'Câu hỏi';
COMMENT ON COLUMN faqs.answer IS 'Câu trả lời';
COMMENT ON COLUMN faqs.group_key IS 'Nhóm phân loại FAQ (post, contact, product, asset, general...)';

