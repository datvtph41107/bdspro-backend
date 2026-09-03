-- ============================================================
-- MIGRATION: Add Legal Documents Registry
-- ============================================================
-- Version: v2
-- Author: hieubap
-- Date: 2024-06-01
-- LOGICAL: Registry pattern — single source of truth

-- ============================================================
-- 1. Bảng qh_legal_documents — Registry
-- ============================================================
CREATE TABLE IF NOT EXISTS qh_legal_documents (
    id BIGSERIAL PRIMARY KEY,
    
    -- Identity (INT cho enum)
    document_type INT NOT NULL,
    document_code VARCHAR(100),
    document_number VARCHAR(100),
    
    -- Metadata
    name VARCHAR(255) NOT NULL,
    full_name TEXT,
    description TEXT,
    
    -- Issuance
    issued_by VARCHAR(255) NOT NULL,
    issued_at DATE,
    effective_at DATE,
    expires_at DATE,
    
    -- Lifecycle (INT cho enum)
    status INT NOT NULL DEFAULT 10,
    superseded_by BIGINT,
    supersedes JSONB DEFAULT '[]',
    
    -- File
    file_type VARCHAR(20),
    file_size_kb INT DEFAULT 0,
    page_count INT DEFAULT 0,
    file_url TEXT,
    preview_url TEXT,
    download_url TEXT,
    
    -- Access
    is_public BOOLEAN DEFAULT TRUE,
    requires_auth BOOLEAN DEFAULT FALSE,
    
    -- References
    relevant_articles JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    affected_layer_ids JSONB DEFAULT '[]',
    affected_zone_ids JSONB DEFAULT '[]',
    affected_parcel_ids JSONB DEFAULT '[]',
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0
);

-- ============================================================
-- Indexes
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_type ON qh_legal_documents(document_type);
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_status ON qh_legal_documents(status);
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_issued_at ON qh_legal_documents(issued_at);
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_effective_at ON qh_legal_documents(effective_at);
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_superseded_by ON qh_legal_documents(superseded_by);
CREATE INDEX IF NOT EXISTS idx_qh_legal_docs_deleted_at ON qh_legal_documents(deleted_at);

-- ============================================================
-- 2. Bảng qh_layer_legal_docs — Layer - Legal Document
-- ============================================================
CREATE TABLE IF NOT EXISTS qh_layer_legal_docs (
    layer_id BIGINT NOT NULL,
    legal_document_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (layer_id, legal_document_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_layer_legal_docs_layer ON qh_layer_legal_docs(layer_id);
CREATE INDEX IF NOT EXISTS idx_qh_layer_legal_docs_doc ON qh_layer_legal_docs(legal_document_id);

-- ============================================================
-- 3. Bảng qh_zone_legal_docs — Zone - Legal Document
-- ============================================================
CREATE TABLE IF NOT EXISTS qh_zone_legal_docs (
    zone_id BIGINT NOT NULL,
    legal_document_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (zone_id, legal_document_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_zone_legal_docs_zone ON qh_zone_legal_docs(zone_id);
CREATE INDEX IF NOT EXISTS idx_qh_zone_legal_docs_doc ON qh_zone_legal_docs(legal_document_id);

-- ============================================================
-- 4. Bảng qh_parcel_legal_docs — Parcel - Legal Document
-- ============================================================
CREATE TABLE IF NOT EXISTS qh_parcel_legal_docs (
    parcel_id BIGINT NOT NULL,
    legal_document_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (parcel_id, legal_document_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_parcel_legal_docs_parcel ON qh_parcel_legal_docs(parcel_id);
CREATE INDEX IF NOT EXISTS idx_qh_parcel_legal_docs_doc ON qh_parcel_legal_docs(legal_document_id);

-- ============================================================
-- 5. Seed data — Document types mapping
-- ============================================================
INSERT INTO qh_legal_documents (
    document_type, document_code, document_number, name, full_name,
    issued_by, issued_at, effective_at, status, file_type, is_public
) VALUES 
(
    30, 'LDD-2013', '45/2013/QH13',
    'Luật Đất đai 2013',
    'Luật số 45/2013/QH13 về Đất đai',
    'Quốc hội', '2013-11-29', '2014-07-01',
    10, 'PDF', true
),
(
    10, 'QĐ-4000', '4000/QĐ-UBND',
    'QĐ 4000/QĐ-UBND ngày 15/01/2020',
    'Quyết định phê duyệt đồ án Quy hoạch sử dụng đất đến năm 2030',
    'UBND TP.HCM', '2020-01-15', '2020-02-01',
    10, 'PDF', true
),
(
    40, 'NĐ-43', '43/2014/NĐ-CP',
    'Nghị định 43/2014/NĐ-CP',
    'Nghị định quy định chi tiết thi hành một số điều của Luật Đất đai',
    'Chính phủ', '2014-05-15', '2014-07-01',
    10, 'PDF', true
);