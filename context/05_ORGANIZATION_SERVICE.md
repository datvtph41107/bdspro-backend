# Organization Service - Complete Context

## 📋 Service Overview

**Service Name**: Organization Service
**Port**: 50052 (gRPC)
**Primary Role**: Organization, Team & Deal Management
**Business Domain**: Multi-level organizational hierarchy
**Status**: ✅ Production Ready

---

## 🎯 Core Responsibilities

### 1. Organization Management
- **Organization CRUD**: Create, update, delete organizations
- **Business Domains**: Link organizations to business domains
- **Organization Settings**: Configure organization-level settings
- **Organization Dashboard**: Statistics and metrics

### 2. Branch Management
- **Branch Operations**: Create, manage company branches
- **Branch Members**: Assign members to branches
- **Branch Hierarchy**: Multi-level branch structure
- **Branch Settings**: Branch-specific configurations

### 3. Team/Group Management
- **Group CRUD**: Create, update, delete groups/teams
- **Group Members**: Member assignment and roles
- **Group Chat**: Auto-create chat channels for groups
- **Group Documents**: Document sharing within groups
- **Group Settings**: Team-specific settings
- **Group Notifications**: Team notification preferences

### 4. Deal Management
- **Deal CRUD**: Create, update deals/transactions
- **Deal Members**: Assign members to deals
- **Deal Milestones**: Track deal progress stages
- **Deal Products**: Link products to deals
- **Deal Commission**: Calculate and track commissions
- **Deal Invitation**: Invite collaborators to deals
- **Deal History**: Track all deal activities

### 5. Role & Permission Management
- **Organization Roles**: Define custom roles
- **Permissions**: Fine-grained permission control
- **Role Assignment**: Assign roles to members
- **Permission Checking**: Validate access rights

### 6. Investment Management
- **Investment Tracking**: Track organizational investments
- **Investment Types**: Real estate, stocks, etc.
- **Investment Performance**: ROI calculations

---

## 📁 Service Architecture

### Domain Layer (`internal/domain/entity/`)

**Organization Entities**:
```
- organization.go              # Main organization entity
- organization_member.go       # Organization membership
- organization_branch.go       # Branch/office locations
- organization_branch_member.go # Branch membership
- organization_role.go         # Custom roles
- organization_permission.go   # Permissions
- organization_business_domain.go # Business domain links
- organization_log_activity.go # Activity logs
- bank_account.go             # Organization bank accounts
```

**Group/Team Entities**:
```
- group.go                    # Team/group entity
- group_member.go             # Team membership
- group_chat.go              # Group chat reference
- group_document.go          # Shared documents
- group_setting.go           # Group settings
- group_notification.go      # Notification preferences
- group_log_activity.go      # Team activity logs
```

**Deal Entities**:
```
- deal.go                    # Main deal entity
- deal_of_organization.go    # Org-level deals
- deal_of_group.go          # Group-level deals
- deal_of_branch.go         # Branch-level deals
- deal_member.go            # Deal participants
- deal_milestone.go         # Deal progress stages
- deal_product.go           # Products in deals
- commission_stats.go       # Commission tracking
```

**Supporting Entities**:
```
- business_domain.go        # Business domain reference
- color.go                  # UI color preferences
- investment.go             # Investment records
- internal_note.go          # Private notes
- plan.go                   # Subscription plans
```

---

## 🏢 Organization Hierarchy

```
Organization (Level 1)
├── Organization Member (Users)
├── Organization Role (Roles)
├── Organization Permission (Permissions)
├── Branch (Level 2)
│   ├── Branch Member (Users)
│   └── Deals (Branch-level)
├── Group/Team (Level 3)
│   ├── Group Member (Users)
│   ├── Group Chat (Auto-created)
│   ├── Group Documents
│   └── Deals (Group-level)
└── Deals (Organization-level)
    ├── Deal Members
    ├── Deal Milestones
    ├── Deal Products
    └── Deal Commission
```

---

## 🗄️ Database Schema Highlights

### organization Table
```sql
CREATE TABLE organization (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,            -- Auto-generated: ORG000001
    tax_code VARCHAR(50) UNIQUE,
    business_license_url VARCHAR(500),
    address VARCHAR(500),
    detail_address TEXT,
    phone VARCHAR(20),
    email VARCHAR(100),
    website VARCHAR(255),
    logo_url VARCHAR(500),
    description TEXT,
    founded_at DATE,
    approve_investor BOOLEAN DEFAULT FALSE,
    owner_id BIGINT NOT NULL,           -- Profile ID of owner
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_org_owner ON organization(owner_id);
CREATE INDEX idx_org_code ON organization(code);
```

### organization_member Table
```sql
CREATE TABLE organization_member (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,            -- Profile ID
    role_id BIGINT,                     -- Organization role
    status INT DEFAULT 1,               -- 1:active, 2:inactive, 3:pending
    joined_at TIMESTAMP,
    removed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    UNIQUE(organization_id, user_id)
);

CREATE INDEX idx_org_member_org ON organization_member(organization_id);
CREATE INDEX idx_org_member_user ON organization_member(user_id);
```

### organization_branch Table
```sql
CREATE TABLE organization_branch (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),                   -- Branch code
    address VARCHAR(500),
    phone VARCHAR(20),
    email VARCHAR(100),
    manager_id BIGINT,                  -- Branch manager
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);
```

### group Table
```sql
CREATE TABLE "group" (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    branch_id BIGINT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    avatar VARCHAR(500),
    owner_id BIGINT NOT NULL,           -- Group owner
    chat_id BIGINT,                     -- Auto-created chat room
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (branch_id) REFERENCES organization_branch(id)
);
```

### group_member Table
```sql
CREATE TABLE group_member (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role INT DEFAULT 1,                 -- 1:member, 2:admin, 3:owner
    status INT DEFAULT 1,
    joined_at TIMESTAMP,
    left_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES "group"(id),
    UNIQUE(group_id, user_id)
);
```

### deal Table
```sql
CREATE TABLE deal (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,            -- Auto-generated: DEAL000001
    name VARCHAR(255) NOT NULL,
    description TEXT,
    deal_type INT,                      -- 1:buy, 2:sell, 3:rent
    status INT DEFAULT 1,               -- 1:pending, 2:in_progress, 3:completed, 4:cancelled
    total_value DECIMAL(18,2),
    commission_rate DECIMAL(5,2),
    commission_amount DECIMAL(18,2),
    start_date DATE,
    end_date DATE,
    closed_date DATE,
    owner_id BIGINT NOT NULL,           -- Deal owner
    organization_id BIGINT,
    branch_id BIGINT,
    group_id BIGINT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

CREATE INDEX idx_deal_org ON deal(organization_id);
CREATE INDEX idx_deal_owner ON deal(owner_id);
CREATE INDEX idx_deal_status ON deal(status);
```

### deal_member Table
```sql
CREATE TABLE deal_member (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    member_type INT,                    -- 1:buyer, 2:seller, 3:agent, 4:collaborator
    role INT,                           -- 1:viewer, 2:editor, 3:admin
    commission_rate DECIMAL(5,2),
    commission_amount DECIMAL(18,2),
    status INT DEFAULT 1,
    joined_at TIMESTAMP,
    left_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (deal_id) REFERENCES deal(id),
    UNIQUE(deal_id, user_id)
);
```

### deal_milestone Table
```sql
CREATE TABLE deal_milestone (
    id BIGSERIAL PRIMARY KEY,
    deal_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status INT DEFAULT 1,               -- 1:pending, 2:in_progress, 3:completed
    order_number INT,
    due_date DATE,
    completed_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (deal_id) REFERENCES deal(id)
);

CREATE INDEX idx_milestone_deal ON deal_milestone(deal_id);
CREATE INDEX idx_milestone_order ON deal_milestone(order_number);
```

---

## 🔐 Permission System

### Permission Format
```
{resource}.{action}

Examples:
- organization.view
- organization.create
- organization.update
- organization.delete
- group.view
- group.create
- deal.view
- deal.create
- member.invite
- member.remove
```

### Role Hierarchy
```
1. Owner (Full Control)
   - All permissions
   - Cannot be removed
   
2. Admin
   - Most permissions
   - Cannot delete organization
   - Cannot change owner
   
3. Manager
   - Limited administrative permissions
   - Can manage team members
   - Can create groups
   
4. Member
   - Basic permissions
   - View-only in most cases
   - Can participate in groups
   
5. Viewer
   - Read-only access
   - Cannot modify anything
```

### Permission Checking
```go
func (u *OrganizationUsecase) CheckPermission(
    ctx context.Context,
    orgId uint64,
    resource string,
    action string,
) error {
    profileId := _utils.GetProfileIdWithContext(ctx)
    
    // Get member's role
    member, err := u.memberRepo.GetByOrgAndUser(ctx, orgId, profileId)
    if err != nil {
        return ErrNotMember
    }
    
    // Get role permissions
    permissions, err := u.permissionRepo.GetByRole(ctx, member.RoleId)
    if err != nil {
        return err
    }
    
    // Check if permission exists
    permissionKey := fmt.Sprintf("%s.%s", resource, action)
    hasPermission := false
    for _, p := range permissions {
        if p.Permission == permissionKey {
            hasPermission = true
            break
        }
    }
    
    if !hasPermission {
        return ErrForbidden
    }
    
    return nil
}
```

---

## 🔄 Business Flows

### 1. Create Organization Flow
```
1. User submits organization data
   POST /v2/org/organization/new
   
2. Validate input
   - Check tax code uniqueness
   - Validate business domains
   - Validate owner
   
3. Create organization
   - Generate organization code (ORG000001)
   - Set owner as creator
   - Link business domains
   
4. Create default role
   - Create "Owner" role with full permissions
   
5. Add owner as member
   - Link owner to organization
   - Assign "Owner" role
   
6. Initialize settings
   - Default organization settings
   - Create default color scheme
   
7. Call external services
   - Notify payment service (create wallet)
   - Log history
   
8. Return organization ID
```

### 2. Add Group Member Flow
```
1. Check permissions
   - Validate requester has "group.manage_members" permission
   
2. Validate member
   - Check if user exists (call User Service)
   - Check if already a member
   - Check if organization member
   
3. Create group membership
   - Add to group_member table
   - Set default role (member)
   
4. Add to chat room
   - Call Chat Service to add user to group chat
   
5. Send notification
   - Notify new member
   - Log activity
   
6. Return success
```

### 3. Create Deal Flow
```
1. Validate ownership
   - Check organization/branch/group membership
   - Validate permissions
   
2. Create deal
   - Generate deal code (DEAL000001)
   - Set initial status (pending)
   - Calculate commission if provided
   
3. Add deal owner as member
   - Create deal_member record
   - Role: admin
   - Type: agent
   
4. Link to organization/branch/group
   - Create deal_of_organization record
   - or deal_of_group / deal_of_branch
   
5. Create default milestones
   - Initial contact
   - Negotiation
   - Contract signing
   - Completion
   
6. Notify members
   - Send notifications to involved parties
   
7. Log history
   - Create deal history record
   
8. Return deal details
```

### 4. Deal Commission Calculation
```
1. Get deal total value
2. Get commission rate (from deal or member)
3. Calculate base commission
   commission = total_value * commission_rate
   
4. Distribute to members
   For each deal_member:
     member_commission = commission * member.commission_rate
     
5. Update commission_stats
   - Store individual commissions
   - Update totals
   
6. Log commission changes
```

---

## 📡 External Service Integrations

### 1. User Service
**Purpose**: Get user profiles
**Methods**:
- `GetProfileById(profileId)`: Get user info
- `GetProfileByIds(profileIds)`: Batch get users

**Usage**:
```go
profiles, err := u.userClient.GetProfileByIds(ctx, memberIds)
```

### 2. Auth Service
**Purpose**: Role & permission management
**Methods**:
- `GetRolePermissions(roleId)`: Get role permissions
- `ValidatePermission(userId, permission)`: Check permission

### 3. Chat Service
**Purpose**: Auto-create group chats
**Methods**:
- `CreateGroupChat(name, memberIds)`: Create chat room
- `AddMemberToChat(chatId, userId)`: Add member

**Usage**:
```go
// When creating group
chatId, err := u.chatClient.CreateGroupChat(ctx, group.Name, memberIds)
group.ChatId = chatId
```

### 4. Transaction Service
**Purpose**: Deal transaction tracking
**Methods**:
- `CreateTransaction(dealId, amount)`: Create transaction
- `GetTransactionsByDeal(dealId)`: Get deal transactions

### 5. Notification Service
**Purpose**: Send notifications
**Methods**:
- `SendNotification(userId, title, message)`: Send notification
- `SendBatchNotification(userIds, title, message)`: Batch send

---

## 🎨 Special Features

### 1. Auto-Code Generation
```go
// Organization code: ORG000001
func (u *OrganizationUsecase) generateOrgCode() string {
    count, _ := u.repo.Count(ctx)
    return fmt.Sprintf("ORG%06d", count+1)
}

// Deal code: DEAL000001
func (u *DealUsecase) generateDealCode() string {
    count, _ := u.repo.Count(ctx)
    return fmt.Sprintf("DEAL%06d", count+1)
}
```

### 2. Group Chat Auto-Creation
```go
func (u *GroupUsecase) Create(ctx context.Context, group *Group, memberIds []uint64) error {
    // Create group
    if err := u.repo.Create(ctx, group); err != nil {
        return err
    }
    
    // Auto-create chat room
    chatId, err := u.chatClient.CreateGroupChat(ctx, group.Name, memberIds)
    if err != nil {
        return err
    }
    
    // Link chat to group
    group.ChatId = chatId
    u.repo.Update(ctx, group.ID, group)
    
    return nil
}
```

### 3. Activity Logging
```go
func (u *OrganizationUsecase) LogActivity(
    ctx context.Context,
    orgId uint64,
    actionType int,
    description string,
) error {
    profileId := _utils.GetProfileIdWithContext(ctx)
    
    log := &OrganizationLogActivity{
        OrganizationId: orgId,
        ProfileId:      profileId,
        ActionType:     actionType,
        Description:    description,
        IpAddress:      getIPFromContext(ctx),
        UserAgent:      getUserAgentFromContext(ctx),
    }
    
    return u.logRepo.Create(ctx, log)
}
```

### 4. Commission Distribution
```go
func (u *DealUsecase) CalculateCommission(
    ctx context.Context,
    dealId uint64,
) error {
    // Get deal
    deal, _ := u.repo.GetByID(ctx, dealId)
    
    // Get deal members
    members, _ := u.memberRepo.GetByDeal(ctx, dealId)
    
    totalCommission := deal.TotalValue * deal.CommissionRate / 100
    
    for _, member := range members {
        memberCommission := totalCommission * member.CommissionRate / 100
        
        // Update member commission
        member.CommissionAmount = memberCommission
        u.memberRepo.Update(ctx, member.ID, member)
        
        // Update commission stats
        stats := &CommissionStats{
            DealId:     dealId,
            MemberId:   member.ID,
            Amount:     memberCommission,
            Status:     StatusPending,
        }
        u.statsRepo.Create(ctx, stats)
    }
    
    return nil
}
```

---

## 🚀 API Endpoints Summary

### Organization APIs
- `POST /v2/org/organization/new` - Create organization
- `PUT /v2/org/organization/{id}` - Update organization
- `GET /v2/org/organization` - List organizations
- `GET /v2/org/organization/current` - Get current org
- `GET /v2/org/organization/current/dashboard` - Dashboard stats

### Member Management
- `POST /v2/org/member/invite` - Invite member
- `DELETE /v2/org/member/{id}` - Remove member
- `GET /v2/org/member/list` - List members

### Group APIs
- `POST /v2/org/group` - Create group
- `PUT /v2/org/group/{id}` - Update group
- `GET /v2/org/group/list` - List groups
- `POST /v2/org/group/{id}/member` - Add member

### Deal APIs
- `POST /v2/org/deal` - Create deal
- `PUT /v2/org/deal/{id}` - Update deal
- `GET /v2/org/deal/list` - List deals
- `POST /v2/org/deal/{id}/milestone` - Add milestone
- `POST /v2/org/deal/{id}/member` - Add member

---

## 📊 Key Metrics

### Database Tables: 25+
### Handlers: 15+
### Usecases: 27
### Repositories: 28
### External Clients: 5

---

**Last Updated**: October 15, 2025
**Service Version**: v2
**Status**: ✅ Production Ready

