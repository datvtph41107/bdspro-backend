-- Create directory_suppliers table
CREATE TABLE IF NOT EXISTS directory_suppliers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    contact_person VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(100),
    address VARCHAR(500),
    website VARCHAR(255),
    tax_code VARCHAR(50),
    business_license VARCHAR(100),
    rating DECIMAL(3,2) DEFAULT 0.00 CHECK (rating >= 0 AND rating <= 5),
    is_active BOOLEAN DEFAULT true,
    service_count INTEGER DEFAULT 0,
    contract_count INTEGER DEFAULT 0,
    total_value BIGINT DEFAULT 0,
    join_date TIMESTAMP,
    last_contact_date TIMESTAMP,
    notes TEXT,
    tags JSONB,
    categories JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_code ON directory_suppliers(code);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_email ON directory_suppliers(email);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_tax_code ON directory_suppliers(tax_code);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_is_active ON directory_suppliers(is_active);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_rating ON directory_suppliers(rating);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_created_at ON directory_suppliers(created_at);

-- Create GIN index for JSONB columns for better search performance
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_tags_gin ON directory_suppliers USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_directory_suppliers_categories_gin ON directory_suppliers USING GIN(categories);
