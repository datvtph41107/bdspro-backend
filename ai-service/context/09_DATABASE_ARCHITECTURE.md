# Database Architecture - Complete Overview

## 📋 Database Strategy

**Database System**: PostgreSQL 14+
**ORM**: GORM v1.25+
**Cache Layer**: Redis 7+
**Pattern**: Database per Service (Microservices Pattern)

---

## 🗄️ Database Per Service

### Strategy
Each microservice has its own isolated database to ensure:
- **Independence**: Services can be deployed independently
- **Scalability**: Each database can be scaled separately
- **Technology Choice**: Different services can use different DB versions
- **Fault Isolation**: Database failure affects only one service

### Database List

| Service | Database Name | Tables | Estimated Size |
|---------|--------------|--------|----------------|
| **Auth** | `auth_db` | 12 | 100MB-1GB |
| **User** | `user_db` | 8 | 500MB-5GB |
| **Organization** | `organization_db` | 25+ | 1GB-10GB |
| **BDSPro** | `bdspro_db` | 50+ | 5GB-50GB |
| **CRM** | `crm_db` | 10+ | 2GB-20GB |
| **Payment** | `payment_db` | 8+ | 1GB-10GB |
| **Notification** | `notification_db` | 7+ | 500MB-5GB |
| **Social** | `social_db` | 5+ | 2GB-20GB |
| **Chat** | `chat_db` | 8+ | 5GB-50GB |
| **Transaction** | `transaction_db` | 7+ | 1GB-10GB |
| **Appointment** | `appointment_db` | 5+ | 100MB-1GB |
| **Marketing** | `marketing_db` | 6+ | 500MB-5GB |
| **Membership** | `membership_db` | 4+ | 100MB-1GB |
| **File** | `file_db` | 3+ | 100MB-1GB |

**Total Databases**: 14
**Total Tables**: 150+
**Estimated Total Size**: 20GB-200GB

---

## 🏗️ Common Schema Patterns

### 1. Base Entity Pattern

All entities inherit from `BaseEntity`:

```sql
-- Common fields in every table
id BIGSERIAL PRIMARY KEY,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW(),
deleted_at TIMESTAMP,                  -- Soft delete
created_by BIGINT,                      -- Audit trail
updated_by BIGINT                       -- Audit trail
```

**Benefits**:
- Automatic timestamps
- Soft delete support
- Audit trail
- Consistent ID strategy

### 2. Owner Pattern

Many entities have ownership tracking:

```sql
owner_id BIGINT NOT NULL,
owner_type INT,                         -- 1:user, 2:group, 3:organization
```

**Use Cases**:
- Products: User or Organization owned
- Posts: Personal or Company posts
- Contacts: Individual or Team contacts
- Assets: User or Organization assets

### 3. Status Pattern

Most entities have status tracking:

```sql
status INT DEFAULT 1,
-- Common status values:
-- 1: active/pending
-- 2: in_progress/processing
-- 3: completed/approved
-- 4: cancelled/rejected
-- 5: archived/deleted
```

### 4. Visibility Pattern

Content visibility control:

```sql
visibility INT DEFAULT 1,
-- 1: public
-- 2: private
-- 3: organization_only
-- 4: group_only
-- 5: friends_only
```

---

## 📊 Core Database Schemas

### Auth Database (`auth_db`)

```sql
-- Main authentication table
CREATE TABLE auth_method (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(20),               -- PHONE, GOOGLE, FACEBOOK, ZALO, ADMIN
    auth_name VARCHAR(255),
    password VARCHAR(255),              -- Hashed
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255),
    avatar VARCHAR(500),
    role_key INT,
    status SMALLINT DEFAULT 1,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_auth_phone ON auth_method(phone);
CREATE INDEX idx_auth_email ON auth_method(email);
CREATE INDEX idx_auth_user ON auth_method(user_id);

-- Session management
CREATE TABLE session (
    id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    session_key VARCHAR(255) UNIQUE,
    platform VARCHAR(50),
    activate BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (auth_id) REFERENCES auth_method(id)
);

-- OTP verification
CREATE TABLE otp (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    otp_code VARCHAR(10) NOT NULL,
    auth_id BIGINT,
    expired_at TIMESTAMP NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Roles & Permissions
CREATE TABLE role (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(100) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE permission (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    module VARCHAR(100),
    action VARCHAR(50),
    is_granted BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (role_id) REFERENCES role(id)
);
```

### Organization Database (`organization_db`)

```sql
-- Organizations
CREATE TABLE organization (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,           -- ORG000001
    tax_code VARCHAR(50) UNIQUE,
    address VARCHAR(500),
    phone VARCHAR(20),
    email VARCHAR(100),
    logo_url VARCHAR(500),
    owner_id BIGINT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Organization members
CREATE TABLE organization_member (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role_id BIGINT,
    status INT DEFAULT 1,
    joined_at TIMESTAMP,
    removed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    UNIQUE(organization_id, user_id)
);

-- Branches
CREATE TABLE organization_branch (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    address VARCHAR(500),
    manager_id BIGINT,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

-- Groups/Teams
CREATE TABLE "group" (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    branch_id BIGINT,
    name VARCHAR(255) NOT NULL,
    owner_id BIGINT NOT NULL,
    chat_id BIGINT,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

-- Deals
CREATE TABLE deal (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,           -- DEAL000001
    name VARCHAR(255) NOT NULL,
    deal_type INT,
    status INT DEFAULT 1,
    total_value DECIMAL(18,2),
    commission_rate DECIMAL(5,2),
    owner_id BIGINT NOT NULL,
    organization_id BIGINT,
    branch_id BIGINT,
    group_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### BDSPro Database (`bdspro_db`)

```sql
-- Products (Real estate properties)
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,           -- SP000001
    name VARCHAR(500),
    description TEXT,
    property_type_id BIGINT,
    
    -- Location
    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    address TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    
    -- Specs
    area DECIMAL(10,2),
    price DECIMAL(18,2),
    bedrooms INT,
    bathrooms INT,
    
    -- Status
    status INT DEFAULT 1,
    transaction_type INT,
    
    -- Ownership
    owner_id BIGINT,
    owner_type INT,
    visibility INT DEFAULT 1,
    
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_location ON product(province_id, district_id);
CREATE INDEX idx_product_type ON product(property_type_id);
CREATE INDEX idx_product_price ON product(price);
CREATE INDEX idx_product_owner ON product(owner_id);

-- Posts (Property listings)
CREATE TABLE post (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    title VARCHAR(500),
    content TEXT,
    post_price DECIMAL(18,2),
    transaction_type INT,
    expired_at TIMESTAMP,
    status INT DEFAULT 1,
    hidden BOOLEAN DEFAULT FALSE,
    num_view BIGINT DEFAULT 0,
    owner_id BIGINT,
    owner_type INT,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_post_product ON post(product_id);
CREATE INDEX idx_post_expired ON post(expired_at);

-- Assets
CREATE TABLE asset (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,           -- AS000001
    product_id BIGINT,
    name VARCHAR(500),
    acquisition_cost DECIMAL(18,2),
    current_value DECIMAL(18,2),
    status INT DEFAULT 1,
    owner_id BIGINT,
    owner_type INT,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Projects
CREATE TABLE project (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(500),
    developer_id BIGINT,
    total_area DECIMAL(18,2),
    total_units INT,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Chat Database (`chat_db`)

```sql
-- Conversations
CREATE TABLE conversation (
    id BIGSERIAL PRIMARY KEY,
    type INT,                          -- 0:private, 1:group, 2:channel
    name VARCHAR(255),
    avatar VARCHAR(500),
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Participants
CREATE TABLE participant (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role INT DEFAULT 0,                -- 0:member, 1:admin, 2:banned
    mute_notif BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES conversation(id),
    UNIQUE(conversation_id, user_id)
);

-- Messages
CREATE TABLE message (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    content TEXT,
    content_type VARCHAR(20),          -- text, image, video, file
    file_url VARCHAR(500),
    is_recall BOOLEAN DEFAULT FALSE,
    reply_to BIGINT,                   -- Reply to message ID
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (conversation_id) REFERENCES conversation(id)
);

CREATE INDEX idx_message_conversation ON message(conversation_id);
CREATE INDEX idx_message_sender ON message(sender_id);

-- Read receipts
CREATE TABLE read_receipt (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    read_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (message_id) REFERENCES message(id),
    UNIQUE(message_id, user_id)
);
```

---

## 🔍 Indexing Strategy

### Primary Indexes
Every table has:
- `PRIMARY KEY` on `id` (clustered index)
- `INDEX` on `deleted_at` (for soft delete queries)

### Foreign Key Indexes
```sql
-- Always index foreign keys
CREATE INDEX idx_table_foreign_id ON table_name(foreign_id);
```

### Search Indexes
```sql
-- For text search
CREATE INDEX idx_product_name ON product USING gin(to_tsvector('english', name));

-- For location search
CREATE INDEX idx_product_location ON product(province_id, district_id, ward_id);

-- For price range search
CREATE INDEX idx_product_price ON product(price);
```

### Composite Indexes
```sql
-- For common query patterns
CREATE INDEX idx_product_owner_status ON product(owner_id, status, deleted_at);
CREATE INDEX idx_post_date_status ON post(created_at DESC, status);
```

---

## 🔐 Data Security

### Sensitive Data Encryption
```sql
-- Passwords are always hashed (bcrypt)
password VARCHAR(255)  -- Stored as hash

-- Sensitive fields use application-level encryption
-- (e.g., tax codes, bank account numbers)
```

### Row-Level Security (Future)
```sql
-- PostgreSQL RLS for multi-tenant data
ALTER TABLE product ENABLE ROW LEVEL SECURITY;

CREATE POLICY product_isolation ON product
    USING (owner_id = current_setting('app.current_user_id')::bigint);
```

---

## 📈 Database Optimization

### Connection Pooling
```go
// GORM configuration
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Query Optimization
```sql
-- Always use indexes for WHERE clauses
SELECT * FROM product 
WHERE province_id = 1        -- Uses index
  AND deleted_at IS NULL;    -- Uses index

-- Avoid N+1 queries - use JOINs or preload
SELECT p.*, m.* FROM product p
LEFT JOIN product_media m ON m.product_id = p.id;
```

### Pagination
```sql
-- Always use LIMIT + OFFSET for large datasets
SELECT * FROM product 
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;
```

---

## 🔄 Database Migrations

### Migration Strategy
```go
// Auto-migration on startup (development)
db.AutoMigrate(
    &Product{},
    &Post{},
    &Asset{},
)

// Manual migrations (production)
// Create SQL files in migrate/ folder
// Apply with migration tools
```

### Sample Migration
```sql
-- migrate/001_create_products.up.sql
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    price DECIMAL(18,2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- migrate/001_create_products.down.sql
DROP TABLE product;
```

---

## 💾 Backup Strategy

### Daily Backups
```bash
# Automated backup script
pg_dump -U postgres -d bdspro_db -F c -f backup_$(date +%Y%m%d).dump

# Retention: 30 days
# Storage: AWS S3 or local
```

### Point-in-Time Recovery
```sql
-- PostgreSQL WAL archiving enabled
archive_mode = on
archive_command = 'cp %p /var/lib/postgresql/archive/%f'
```

---

## 📊 Database Statistics

### Total Metrics
- **Databases**: 14
- **Tables**: 150+
- **Indexes**: 300+
- **Foreign Keys**: 200+
- **Estimated Records**: 1M-100M

### Performance Targets
- **Query Time**: < 10ms (with indexes)
- **Connection Pool**: 10-100 connections
- **Throughput**: 1000+ queries/second
- **Availability**: 99.9%

---

## 🔧 Database Tools

### Development
- **pgAdmin 4**: GUI management
- **DBeaver**: Universal database tool
- **GORM**: ORM for Go

### Monitoring
- **pg_stat_statements**: Query statistics
- **pgBadger**: Log analyzer
- **Prometheus**: Metrics collection

### Backup
- **pg_dump**: Logical backup
- **pg_basebackup**: Physical backup
- **Barman**: Backup and recovery manager

---

**Last Updated**: October 15, 2025
**Schema Version**: v1.0
**Status**: ✅ Production Ready

