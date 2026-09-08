# Business Flows & Complex Logic

## 📋 Overview

**Purpose**: Document complex business workflows and logic
**Scope**: Multi-service orchestration, deal management, commission calculation
**Status**: ✅ Production Implementation

---

## 🏢 Deal Management Flow (Complete Lifecycle)

### Phase 1: Deal Creation

```
┌─────────────────────────────────────────────────┐
│ 1. User Creates Deal                            │
│    POST /v2/org/deal                           │
│    {                                            │
│      name: "Bán căn hộ Vista Verde",           │
│      targetProfit: 500000000,  # 500M VND      │
│      dealType: 1,  # Bán                       │
│      products: [productId1, productId2],       │
│      members: [memberId1, memberId2],          │
│      customers: [customerId1]                  │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 2. Organization Service - Validation           │
│    ✓ Check user has permission                 │
│    ✓ Validate organization/group exists        │
│    ✓ Validate products exist                   │
│    ✓ Validate members are in org              │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 3. Create Deal Record (Transaction)            │
│    BEGIN TRANSACTION                            │
│                                                 │
│    3.1. Create deal                            │
│         - Generate code: DEAL000001            │
│         - Set status: Pending                  │
│         - Set owner from context               │
│                                                 │
│    3.2. Link to organization/group             │
│         - Create deal_of_organization          │
│         - or deal_of_group                     │
│                                                 │
│    3.3. Add deal owner as member               │
│         - Create deal_member                   │
│         - Role: Owner (400)                    │
│         - IsOwner: true                        │
│                                                 │
│    3.4. Add other members                      │
│         - Create deal_member for each          │
│         - Role: Member (430)                   │
│         - Set commission if provided           │
│                                                 │
│    3.5. Link products                          │
│         - Create deal_product records          │
│                                                 │
│    3.6. Link customers                         │
│         - Create deal_customer records         │
│                                                 │
│    3.7. Create default milestones              │
│         - Initial contact                      │
│         - Negotiation                          │
│         - Contract preparation                 │
│         - Signing                              │
│         - Completion                           │
│                                                 │
│    COMMIT TRANSACTION                           │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 4. Post-Creation Actions (Async)               │
│                                                 │
│    4.1. Call Chat Service                      │
│         - Create deal discussion room          │
│         - Add all members to chat              │
│                                                 │
│    4.2. Call Notification Service              │
│         - Notify all members                   │
│         - Notify customers                     │
│                                                 │
│    4.3. Log to History Service                 │
│         - Action: CREATE_DEAL                  │
│         - Record all details                   │
│                                                 │
│    4.4. Update Organization Stats              │
│         - Increment deal count                 │
│         - Update dashboard                     │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 5. Return Deal ID                              │
│    Response: { id: 123, code: "DEAL000001" }  │
└─────────────────────────────────────────────────┘
```

---

### Phase 2: Deal Progress Tracking

```
┌─────────────────────────────────────────────────┐
│ 1. Update Deal Milestone                       │
│    PUT /v2/org/deal/milestone/{id}/status      │
│    { status: "COMPLETED" }                     │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 2. Update Milestone Status                     │
│    - Mark milestone as completed               │
│    - Set completed_date                        │
│    - Check if last milestone                   │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 3. Auto-Progress Deal Status                   │
│    IF all milestones completed:                │
│      - Update deal status to COMPLETED         │
│      - Calculate actual profit                 │
│      - Trigger commission distribution         │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 4. Log Activity                                │
│    - Record milestone completion               │
│    - Notify all members                        │
└─────────────────────────────────────────────────┘
```

---

### Phase 3: Commission Calculation & Distribution

```
┌─────────────────────────────────────────────────┐
│ Deal Completed Event                            │
│ Status: COMPLETED                               │
│ Total Amount: 1,000,000,000 VND                │
│ Target Profit: 500,000,000 VND                 │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 1. Calculate Actual Profit                     │
│                                                 │
│    SELECT SUM(amount) FROM transaction          │
│    WHERE deal_id = 123 AND type = 'INCOME'     │
│    Result: total_income = 1,200,000,000        │
│                                                 │
│    SELECT SUM(amount) FROM transaction          │
│    WHERE deal_id = 123 AND type = 'COST'       │
│    Result: total_cost = 600,000,000            │
│                                                 │
│    actual_profit = total_income - total_cost   │
│                  = 1,200,000,000 - 600,000,000 │
│                  = 600,000,000 VND              │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 2. Get Deal Members with Commission Settings   │
│                                                 │
│    Member 1 (Owner):                           │
│      - CommissionType: PERCENTAGE (10)         │
│      - CommissionValue: 50%                    │
│                                                 │
│    Member 2 (Admin):                           │
│      - CommissionType: PERCENTAGE (10)         │
│      - CommissionValue: 30%                    │
│                                                 │
│    Member 3 (Member):                          │
│      - CommissionType: FIXED (20)              │
│      - CommissionValue: 50,000,000 VND         │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 3. Calculate Individual Commissions            │
│                                                 │
│    Member 1:                                   │
│      600,000,000 * 50% = 300,000,000 VND      │
│                                                 │
│    Member 2:                                   │
│      600,000,000 * 30% = 180,000,000 VND      │
│                                                 │
│    Member 3:                                   │
│      50,000,000 VND (fixed)                    │
│                                                 │
│    Total Distributed:                          │
│      300M + 180M + 50M = 530,000,000 VND      │
│                                                 │
│    Remaining:                                  │
│      600M - 530M = 70,000,000 VND (reserve)   │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 4. Create Commission Records (Transaction)     │
│                                                 │
│    BEGIN TRANSACTION                            │
│                                                 │
│    FOR each member:                            │
│      INSERT INTO commission_stats (             │
│        deal_id,                                │
│        member_id,                              │
│        commission_amount,                      │
│        status: PENDING                         │
│      )                                          │
│                                                 │
│    UPDATE deal SET                             │
│      total_commission = 530,000,000,           │
│      remaining_profit = 70,000,000             │
│    WHERE id = 123                              │
│                                                 │
│    COMMIT TRANSACTION                           │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 5. Payment Processing                          │
│                                                 │
│    FOR each member with commission:            │
│      Call Payment Service:                     │
│      - Create payment transaction              │
│      - Deduct from org wallet                  │
│      - Credit to member wallet                 │
│                                                 │
│      Update commission status: PAID            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ 6. Notifications & History                     │
│    - Notify each member of commission          │
│    - Log commission distribution               │
│    - Send email receipts                       │
└─────────────────────────────────────────────────┘
```

---

## 🏘️ Property Listing Flow (Post → Product → Asset)

### Complete Property Management Lifecycle

```
┌─────────────────────────────────────────────────┐
│ Step 1: Create Product (Real Estate Property)  │
│                                                 │
│    POST /v2/bdspro/v2/product/new              │
│    {                                            │
│      name: "Căn hộ 2PN Vista Verde",          │
│      propertyTypeId: 3,  # Apartment           │
│      area: 75.5,  # m²                         │
│      price: 3500000000,  # 3.5B VND           │
│      bedrooms: 2,                              │
│      bathrooms: 2,                             │
│      provinceId: 1,  # TP.HCM                  │
│      districtId: 10,  # Quận 2                 │
│      amenities: [1, 2, 5],  # Pool, Gym, Park  │
│      media: [                                  │
│        {url: "image1.jpg", type: "image"},    │
│        {url: "image2.jpg", type: "image"}     │
│      ]                                          │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 2: Product Created                        │
│                                                 │
│    BEGIN TRANSACTION                            │
│      1. Generate code: SP000001                │
│      2. Set owner from context                 │
│      3. Insert product                         │
│      4. Insert product_media                   │
│      5. Link amenities (many-to-many)         │
│      6. Calculate price_per_m2                 │
│    COMMIT                                       │
│                                                 │
│    Result: Product ID = 456                    │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 3: Create Post (Listing)                  │
│                                                 │
│    POST /v2/bdspro/v2/post                     │
│    {                                            │
│      productId: 456,                           │
│      title: "Cần bán căn hộ 2PN Vista Verde", │
│      content: "Căn góc, view đẹp...",          │
│      postPrice: 3600000000,  # Can differ      │
│      transactionType: 1,  # Sell               │
│      numDate: 30,  # Active for 30 days        │
│      packageVisible: 2,  # VIP listing         │
│      visibility: 1  # Public                   │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 4: Post Created & Published               │
│                                                 │
│    BEGIN TRANSACTION                            │
│      1. Create post                            │
│      2. Set expired_at = NOW() + 30 days       │
│      3. Link to product                        │
│      4. Set status: ACTIVE                     │
│    COMMIT                                       │
│                                                 │
│    Side effects (async):                       │
│      - Create news feed (Social Service)       │
│      - Update search index (Search Service)    │
│      - Notify followers                        │
│                                                 │
│    Result: Post ID = 789                       │
│           Post URL: /v2/bdspro/v2/post/789    │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 5: Post Expiration (Scheduled Job)        │
│                                                 │
│    Daily job checks expired posts:             │
│    SELECT * FROM post                          │
│    WHERE expired_at < NOW()                    │
│      AND status = ACTIVE                       │
│                                                 │
│    For each expired post:                      │
│      - Update status: EXPIRED                  │
│      - Update hidden: true                     │
│      - Notify owner                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 6: Renew Post                             │
│                                                 │
│    POST /v2/bdspro/v2/post/new-expired         │
│    { postId: 789, numDay: 30 }                 │
│                                                 │
│    Actions:                                    │
│      - Update expired_at = NOW() + 30 days     │
│      - Set status: ACTIVE                      │
│      - Set hidden: false                       │
│      - Increment num_date counter              │
│                                                 │
│    Payment (future):                           │
│      - Charge user for renewal                 │
│      - Deduct from wallet                      │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 7: Create Asset (After Sale)              │
│                                                 │
│    POST /v2/bdspro/v2/asset                    │
│    {                                            │
│      productId: 456,  # Link to product        │
│      purchasePrice: 3600000000,                │
│      purchaseDate: "2025-10-15",               │
│      legalStatus: 10,  # Sổ đỏ                 │
│      legalItems: [                             │
│        {                                        │
│          certificateNumber: "SH123456",        │
│          issueDate: "2025-10-15",              │
│          holderName: "Nguyễn Văn A"            │
│        }                                        │
│      ]                                          │
│    }                                            │
│                                                 │
│    Result: Asset ID = 999                      │
│            Code: AS000001                      │
│                                                 │
│    Now hierarchy exists:                       │
│    Post(789) → Product(456) → Asset(999)       │
└─────────────────────────────────────────────────┘
```

---

## 💰 Asset Cost & Income Tracking Flow

### Monthly Asset Management

```
┌─────────────────────────────────────────────────┐
│ Month 1: Add Operating Costs                   │
│                                                 │
│    POST /v2/bdspro/v2/asset/cost               │
│    {                                            │
│      assetId: 999,                             │
│      costTypeId: 1,  # Bảo trì                 │
│      amount: 5000000,  # 5M VND                │
│      costDate: "2025-11-01",                   │
│      description: "Sơn lại tường"              │
│    }                                            │
│                                                 │
│    Actions:                                    │
│      - Create asset_cost record                │
│      - Update asset.total_cost += 5M           │
│      - Update asset.current_value -= 5M        │
│      - Log cost history                        │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Month 1: Add Rental Income                     │
│                                                 │
│    POST /v2/bdspro/v2/asset/income             │
│    {                                            │
│      assetId: 999,                             │
│      incomeTypeId: 1,  # Cho thuê              │
│      amount: 15000000,  # 15M VND/month        │
│      incomeDate: "2025-11-01",                 │
│      description: "Tiền thuê tháng 11"         │
│    }                                            │
│                                                 │
│    Actions:                                    │
│      - Create asset_income record              │
│      - Update asset.total_income += 15M        │
│      - Calculate ROI                           │
│      - Log income history                      │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Calculate ROI                                  │
│                                                 │
│    Total Cost:                                 │
│      Purchase: 3,600,000,000                   │
│      Operating: 5,000,000                      │
│      Total: 3,605,000,000 VND                  │
│                                                 │
│    Total Income: 15,000,000 VND                │
│                                                 │
│    ROI (annual, projected):                    │
│      Monthly income: 15M                       │
│      Annual income: 15M * 12 = 180M            │
│      ROI = (180M / 3605M) * 100                │
│          = 4.99% per year                      │
│                                                 │
│    GET /v2/bdspro/v2/asset/roi/{id}            │
│    Response: { roi: 4.99, period: "annual" }   │
└─────────────────────────────────────────────────┘
```

---

## 📊 Dashboard Aggregation Flow

### Organization Dashboard

```
┌─────────────────────────────────────────────────┐
│ User Requests Dashboard                         │
│ GET /v2/org/organization/current/dashboard     │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Organization Service - Orchestration           │
│                                                 │
│    organizationId = GetOrganizationIdFromContext(ctx)
│                                                 │
│    Parallel calls to multiple services:        │
│    ┌───────────────────────────────────┐       │
│    │ Call BDSPro Service               │       │
│    │ GetCountByOwner(orgId)            │       │
│    │ → totalProduct, totalPost,        │       │
│    │   totalAsset, totalProject        │       │
│    └───────────────────────────────────┘       │
│    ┌───────────────────────────────────┐       │
│    │ Call Social Service               │       │
│    │ CountNewsFeed(orgId)              │       │
│    │ → totalNewsfeed                   │       │
│    └───────────────────────────────────┘       │
│    ┌───────────────────────────────────┐       │
│    │ Call Transaction Service          │       │
│    │ CountTransactions(orgId)          │       │
│    │ → totalTransaction                │       │
│    └───────────────────────────────────┘       │
│    ┌───────────────────────────────────┐       │
│    │ Call Appointment Service          │       │
│    │ CountAppointments(orgId)          │       │
│    │ → totalSchedule                   │       │
│    └───────────────────────────────────┘       │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Aggregate Results                              │
│                                                 │
│    {                                            │
│      totalProduct: 245,                        │
│      totalPost: 189,                           │
│      totalAsset: 45,                           │
│      totalProject: 12,                         │
│      totalNewsfeed: 567,                       │
│      totalSchedule: 23,                        │
│      totalTransaction: 89,                     │
│      totalClicked: 12456,  # From analytics    │
│      totalViewHome: 8934   # From analytics    │
│    }                                            │
└─────────────────────────────────────────────────┘
```

---

## 👥 Contact Sync Flow (CRM Integration)

### Sync Phone Contacts to CRM

```
┌─────────────────────────────────────────────────┐
│ Mobile App Syncs Contacts                      │
│                                                 │
│    POST /v2/crm/contact/sync                   │
│    {                                            │
│      contacts: [                               │
│        {phone: "0901234567", name: "John"},   │
│        {phone: "0907654321", name: "Jane"},   │
│        ... (100 contacts)                      │
│      ]                                          │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 1: Extract Phone Numbers                 │
│    phones = ["0901234567", "0907654321", ...]  │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 2: Call User Service (Batch)             │
│                                                 │
│    profiles = userClient.GetProfileByPhones(   │
│      ctx, phones                               │
│    )                                            │
│                                                 │
│    Result: Map[phone → profile]                │
│    {                                            │
│      "0901234567": {                           │
│        profileId: 123,                         │
│        fullName: "John Doe",                   │
│        avatar: "..."                           │
│      },                                         │
│      "0907654321": null  # Not registered      │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 3: Link Contacts to Profiles             │
│                                                 │
│    FOR each contact:                           │
│      IF profile exists for phone:              │
│        contact.profileId = profile.profileId   │
│        contact.fullName = profile.fullName     │
│        contact.avatar = profile.avatar         │
│        contact.hasApp = true                   │
│      ELSE:                                      │
│        contact.hasApp = false                  │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 4: Upsert Contacts (Transaction)         │
│                                                 │
│    BEGIN TRANSACTION                            │
│                                                 │
│    INSERT INTO contact (phone, full_name, ...)  │
│    VALUES ...                                   │
│    ON CONFLICT (phone, owner_id) DO UPDATE     │
│    SET full_name = EXCLUDED.full_name,         │
│        profile_id = EXCLUDED.profile_id,       │
│        updated_at = NOW()                      │
│                                                 │
│    COMMIT                                       │
│                                                 │
│    Result: 100 contacts synced                 │
│            - 45 linked to app users            │
│            - 55 not registered                 │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 5: Auto-Create Friend Relationships      │
│                                                 │
│    FOR each contact with profileId:            │
│      Check if already friends                  │
│      IF not:                                    │
│        Send friend request                     │
│        OR auto-add as contact                  │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 6: Return Sync Results                   │
│                                                 │
│    {                                            │
│      data: [...contacts with profileIds...],   │
│      total: 100,                               │
│      linkedToApp: 45,                          │
│      notRegistered: 55                         │
│    }                                            │
└─────────────────────────────────────────────────┘
```

---

## 🔄 Group Creation with Auto-Chat

### Create Group → Auto-Create Chat Room

```
┌─────────────────────────────────────────────────┐
│ Step 1: Create Group/Team                      │
│                                                 │
│    POST /v2/org/group                          │
│    {                                            │
│      name: "Sales Team Alpha",                 │
│      description: "...",                       │
│      members: [userId1, userId2, userId3]      │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 2: Organization Service - Create Group    │
│                                                 │
│    BEGIN TRANSACTION                            │
│                                                 │
│    2.1. Create group record                    │
│         INSERT INTO "group" (name, ...)        │
│         Result: groupId = 10                   │
│                                                 │
│    2.2. Add group members                      │
│         FOR each member:                       │
│           INSERT INTO group_member             │
│           (group_id, user_id, role)            │
│                                                 │
│    2.3. Call Chat Service                      │
│         chatId = chatClient.CreateGroupChat(   │
│           name: "Sales Team Alpha",            │
│           members: [userId1, userId2, userId3] │
│         )                                       │
│         Result: chatId = 555                   │
│                                                 │
│    2.4. Link chat to group                     │
│         UPDATE "group"                         │
│         SET chat_id = 555                      │
│         WHERE id = 10                          │
│                                                 │
│    COMMIT TRANSACTION                           │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 3: Chat Service - Create Conversation    │
│                                                 │
│    BEGIN TRANSACTION                            │
│                                                 │
│    3.1. Create conversation                    │
│         INSERT INTO conversation (              │
│           type: GROUP,                         │
│           name: "Sales Team Alpha",            │
│           created_by: currentUserId            │
│         )                                       │
│         Result: conversationId = 555           │
│                                                 │
│    3.2. Add participants                       │
│         FOR each member:                       │
│           INSERT INTO participant (            │
│             conversation_id: 555,              │
│             user_id: userId,                   │
│             role: MEMBER                       │
│           )                                     │
│                                                 │
│    COMMIT TRANSACTION                           │
│                                                 │
│    Return conversationId to Organization Service
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Step 4: Send Welcome Messages                 │
│                                                 │
│    Chat Service:                               │
│      - Send system message to group            │
│        "Group created by [User]"               │
│                                                 │
│    Notification Service:                       │
│      - Notify all members                      │
│        "You've been added to Sales Team Alpha" │
│                                                 │
│    Relay Service:                              │
│      - Broadcast real-time update              │
│      - Members see new group immediately       │
└─────────────────────────────────────────────────┘
```

---

## 🔀 Product Split & Merge

### Split Large Land into Smaller Lots

```
┌─────────────────────────────────────────────────┐
│ Input: Parent Product                          │
│   ID: 100                                      │
│   Name: "Đất 1000m² Quận 2"                   │
│   Area: 1000 m²                                │
│   Price: 10,000,000,000 VND                    │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ User Initiates Split                           │
│                                                 │
│    POST /v2/bdspro/v2/product/child/split      │
│    {                                            │
│      parentId: 100,                            │
│      children: [                               │
│        { name: "Lô 1", area: 200, price: 2B }, │
│        { name: "Lô 2", area: 300, price: 3B }, │
│        { name: "Lô 3", area: 500, price: 5B }  │
│      ]                                          │
│    }                                            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Validation                                     │
│                                                 │
│    ✓ Check ownership                           │
│    ✓ Validate total area:                      │
│      200 + 300 + 500 = 1000 ✓                  │
│    ✓ Check parent not already split            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Execute Split (Transaction)                    │
│                                                 │
│    BEGIN TRANSACTION                            │
│                                                 │
│    1. Create child products                    │
│       INSERT INTO product (                    │
│         parent_id: 100,                        │
│         code: SP000101,                        │
│         name: "Lô 1",                          │
│         area: 200,                             │
│         price: 2000000000                      │
│       )                                         │
│       ... repeat for Lô 2, Lô 3                │
│                                                 │
│    2. Update parent product                    │
│       UPDATE product SET                       │
│         has_children = true,                   │
│         available_area = 0,                    │
│         status = SPLIT                         │
│       WHERE id = 100                           │
│                                                 │
│    3. Create split history                     │
│       INSERT INTO product_history (            │
│         product_id: 100,                       │
│         action_type: SPLIT,                    │
│         child_count: 3,                        │
│         detail: "Split into 3 lots"            │
│       )                                         │
│                                                 │
│    COMMIT TRANSACTION                           │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│ Result: Product Hierarchy                      │
│                                                 │
│    Product #100 (Parent - SPLIT)               │
│    ├── Product #101 (Lô 1 - 200m²)            │
│    ├── Product #102 (Lô 2 - 300m²)            │
│    └── Product #103 (Lô 3 - 500m²)            │
│                                                 │
│    Now each lot can be:                        │
│      - Listed separately (create posts)        │
│      - Sold independently                      │
│      - Tracked as separate assets              │
└─────────────────────────────────────────────────┘
```

---

## 📈 Performance Insights

### Typical Request Latencies

| Flow | Steps | Services | Latency |
|------|-------|----------|---------|
| **Create Deal** | 7 | 3 services | 100-200ms |
| **Dashboard** | 5 | 5 services | 50-150ms |
| **Contact Sync** | 4 | 2 services | 200-500ms |
| **Product Split** | 3 | 1 service | 50-100ms |

### Transaction Sizes

| Flow | DB Operations | Avg Time |
|------|---------------|----------|
| **Create Deal** | 10-15 INSERTs | 50ms |
| **Product Split** | 5-10 INSERTs | 30ms |
| **Contact Sync** | 100 UPSERTs | 200ms |

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Production Flows

