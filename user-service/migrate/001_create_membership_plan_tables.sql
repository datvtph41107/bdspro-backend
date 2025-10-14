-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id BIGSERIAL PRIMARY KEY,
    package_name VARCHAR(255) NOT NULL,
    duration BIGINT NOT NULL, -- số ngày
    price BIGINT NOT NULL,
    create_product INT DEFAULT 0,
    post_limit BIGINT DEFAULT 0,
    public_product_limit INT DEFAULT 0,
    create_product_ai INT DEFAULT 0,
    friend_limit INT DEFAULT 0,
    customer_limit INT DEFAULT 0,
    join_org INT DEFAULT 0,
    create_organization INT DEFAULT 0,
    org_member_limit INT DEFAULT 0,
    org_product_limit INT DEFAULT 0,
    color VARCHAR(50),
    features TEXT[],
    visibility BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create memberships table
CREATE TABLE IF NOT EXISTS memberships (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL UNIQUE,
    plan_id BIGINT NOT NULL REFERENCES plans(id),
    finished_date TIMESTAMP,
    order_id BIGINT,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_memberships_profile_id ON memberships(profile_id);
CREATE INDEX IF NOT EXISTS idx_memberships_plan_id ON memberships(plan_id);
CREATE INDEX IF NOT EXISTS idx_memberships_deleted_at ON memberships(deleted_at);
CREATE INDEX IF NOT EXISTS idx_plans_visibility ON plans(visibility);
CREATE INDEX IF NOT EXISTS idx_plans_deleted_at ON plans(deleted_at);

-- Insert default free plan
INSERT INTO plans (
    package_name, 
    duration, 
    price, 
    create_product, 
    post_limit, 
    public_product_limit,
    create_product_ai,
    friend_limit,
    customer_limit,
    join_org,
    create_organization,
    org_member_limit,
    org_product_limit,
    color,
    features,
    visibility
) VALUES (
    'Gói miễn phí',
    365, -- 1 năm
    0,
    10, -- tạo tối đa 10 sản phẩm
    5, -- post 5 bài
    3, -- public 3 sản phẩm
    0, -- không dùng AI
    50, -- 50 bạn bè
    100, -- 100 khách hàng
    3, -- tham gia 3 tổ chức
    0, -- không tạo tổ chức
    0,
    0,
    '#6B7280',
    ARRAY['Quản lý sản phẩm cơ bản', 'Quản lý khách hàng', 'Đăng bài giới thiệu'],
    true
) ON CONFLICT DO NOTHING;

