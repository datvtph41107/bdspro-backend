# Inter-Service Communication Patterns

## 📋 Overview

**Purpose**: How microservices communicate with each other
**Protocol**: gRPC
**Pattern**: Client-based communication through interfaces
**Status**: ✅ Production Implementation

---

## 🔄 Communication Architecture

### Service Communication Flow

```
┌─────────────────────────────────────────────────┐
│         Service A (e.g., Organization)          │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │  Usecase Layer                          │   │
│  │  - Business logic                       │   │
│  │  - Needs user profile data             │   │
│  └────────────┬────────────────────────────┘   │
│               │                                 │
│               │ Uses                            │
│               ▼                                 │
│  ┌─────────────────────────────────────────┐   │
│  │  Provider Interface                     │   │
│  │  (internal/interface/provider)          │   │
│  │                                          │   │
│  │  type IUserClient interface {           │   │
│  │    GetProfileByIds(...) map[uint64]...  │   │
│  │  }                                       │   │
│  └────────────┬────────────────────────────┘   │
│               │                                 │
│               │ Implemented by                  │
│               ▼                                 │
│  ┌─────────────────────────────────────────┐   │
│  │  Client Implementation                  │   │
│  │  (infra/client/user_client.go)         │   │
│  │                                          │   │
│  │  // @bind: internal/interface/provider  │   │
│  │  type UserClient struct {               │   │
│  │    conn *grpc.ClientConn                │   │
│  │  }                                       │   │
│  └────────────┬────────────────────────────┘   │
│               │                                 │
└───────────────┼─────────────────────────────────┘
                │
                │ gRPC Call
                │
                ▼
┌─────────────────────────────────────────────────┐
│         Service B (e.g., User)                  │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │  gRPC Server Handler                    │   │
│  │  (infra/handler/profile_handler.go)     │   │
│  │                                          │   │
│  │  func GetProfileByIds(...) {            │   │
│  │    // Handle request                    │   │
│  │    return profiles                      │   │
│  │  }                                       │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## 🎯 Client Implementation Pattern

### Step 1: Define Interface

**Location**: `<service>/internal/interface/provider/`

```go
// organization-service/internal/interface/user_client.go
package iusecase

import (
    "context"
    sharepb "pb/types/shared"
)

type IUserClient interface {
    // Get single profile
    GetProfileById(ctx context.Context, profileId uint64) (*sharepb.ProfileItem, error)
    
    // Get multiple profiles (batch)
    GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) ([]*sharepb.ProfileItem, error)
    
    // Get as map for easy lookup
    GetMapProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
    
    // Get profiles by phone numbers
    GetProfileByPhones(ctx context.Context, phones []string) (map[string]*sharepb.ProfileItem, error)
}
```

### Step 2: Implement Client

**Location**: `<service>/infra/client/`

```go
// organization-service/infra/client/user_client.go
package client

import (
    "context"
    "organization/internal/interface"
    userpb "pb/types/user"
    sharepb "pb/types/shared"
    "google.golang.org/grpc"
)

// @bind: internal/interface
type UserClient struct {
    conn *grpc.ClientConn
}

func NewUserClient(conn *grpc.ClientConn) iusecase.IUserClient {
    return &UserClient{conn: conn}
}

func (c *UserClient) GetProfileById(ctx context.Context, profileId uint64) (*sharepb.ProfileItem, error) {
    client := userpb.NewProfileServiceClient(c.conn)
    
    resp, err := client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
        Ids: []uint64{profileId},
    })
    
    if err != nil {
        return nil, err
    }
    
    if len(resp.Profiles) == 0 {
        return nil, errors.New("profile not found")
    }
    
    return resp.Profiles[0], nil
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
    client := userpb.NewProfileServiceClient(c.conn)
    
    resp, err := client.GetProfileByIds(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Convert array to map for O(1) lookup
    profileMap := make(map[uint64]*sharepb.ProfileItem)
    for _, profile := range resp.Profiles {
        profileMap[profile.ProfileId] = profile
    }
    
    return profileMap, nil
}

func (c *UserClient) GetProfileByPhones(ctx context.Context, phones []string) (map[string]*sharepb.ProfileItem, error) {
    client := userpb.NewProfileServiceClient(c.conn)
    
    resp, err := client.GetProfileByPhones(ctx, &userpb.GetProfileByPhonesRequest{
        Phones: phones,
    })
    
    if err != nil {
        return nil, err
    }
    
    // Create phone → profile map
    phoneMap := make(map[string]*sharepb.ProfileItem)
    for _, profile := range resp.Profiles {
        if profile.Phone != "" {
            phoneMap[profile.Phone] = profile
        }
    }
    
    return phoneMap, nil
}
```

### Step 3: Inject in Usecase

```go
// organization-service/internal/usecase/organization_usecase.go
type OrganizationUsecase struct {
    repo       repository.OrganizationRepository
    userClient iusecase.IUserClient  // Injected
}

func NewOrganizationUsecase(
    repo repository.OrganizationRepository,
    userClient iusecase.IUserClient,
) OrganizationUsecase {
    return &OrganizationUsecase{
        repo:       repo,
        userClient: userClient,
    }
}
```

### Step 4: Use in Business Logic

```go
func (u *OrganizationUsecase) GetMembers(ctx context.Context, orgId uint64) ([]*dto.MemberWithProfile, error) {
    // Get members from database
    members, _ := u.memberRepo.GetByOrganization(ctx, orgId)
    
    // Extract profile IDs
    profileIds := make([]uint64, len(members))
    for i, member := range members {
        profileIds[i] = member.UserID
    }
    
    // Batch get profiles from User Service
    profilesMap, err := u.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
        Ids: profileIds,
    })
    if err != nil {
        // Graceful degradation - continue without profiles
        profilesMap = make(map[uint64]*sharepb.ProfileItem)
    }
    
    // Combine members with profiles
    result := make([]*dto.MemberWithProfile, len(members))
    for i, member := range members {
        result[i] = &dto.MemberWithProfile{
            Member:  member,
            Profile: profilesMap[member.UserID],  // O(1) lookup
        }
    }
    
    return result, nil
}
```

---

## 🎯 Common Communication Patterns

### Pattern 1: Batch Fetching (Avoid N+1)

```go
// ❌ BAD - N+1 queries (1 query + N service calls)
func (u *Usecase) GetContactsWithProfiles(ctx context.Context) ([]*ContactWithProfile, error) {
    contacts, _ := u.repo.GetAll(ctx)
    
    result := make([]*ContactWithProfile, len(contacts))
    for i, contact := range contacts {
        // N service calls! (if 100 contacts = 100 gRPC calls)
        profile, _ := u.userClient.GetProfileById(ctx, contact.ProfileID)
        result[i] = &ContactWithProfile{
            Contact: contact,
            Profile: profile,
        }
    }
    
    return result, nil
}

// ✅ GOOD - Batch fetching (1 query + 1 service call)
func (u *Usecase) GetContactsWithProfiles(ctx context.Context) ([]*ContactWithProfile, error) {
    contacts, _ := u.repo.GetAll(ctx)
    
    // Extract all profile IDs
    profileIds := make([]uint64, len(contacts))
    for i, contact := range contacts {
        profileIds[i] = contact.ProfileID
    }
    
    // Single batch call
    profilesMap, _ := u.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
        Ids: profileIds,
    })
    
    // Map results
    result := make([]*ContactWithProfile, len(contacts))
    for i, contact := range contacts {
        result[i] = &ContactWithProfile{
            Contact: contact,
            Profile: profilesMap[contact.ProfileID],
        }
    }
    
    return result, nil
}
```

**Performance Impact**:
- Bad: 100 contacts = 1 + 100 = 101 calls (10-500ms total)
- Good: 100 contacts = 1 + 1 = 2 calls (10-20ms total)
- **50-100x faster!**

---

### Pattern 2: Graceful Degradation

```go
// ✅ GOOD - Continue even if external service fails
func (u *OrganizationUsecase) GetMembers(ctx context.Context, orgId uint64) ([]*dto.MemberWithProfile, error) {
    // Get members (critical)
    members, err := u.memberRepo.GetByOrganization(ctx, orgId)
    if err != nil {
        return nil, err  // Fail if can't get members
    }
    
    // Get profiles (nice-to-have)
    profileIds := extractProfileIds(members)
    profilesMap, err := u.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
        Ids: profileIds,
    })
    if err != nil {
        // Log error but continue
        log.Printf("Failed to get profiles: %v", err)
        profilesMap = make(map[uint64]*sharepb.ProfileItem)
    }
    
    // Combine (profiles may be nil)
    result := make([]*dto.MemberWithProfile, len(members))
    for i, member := range members {
        result[i] = &dto.MemberWithProfile{
            Member:  member,
            Profile: profilesMap[member.UserID],  // May be nil
        }
    }
    
    return result, nil
}
```

**Benefit**: Service remains functional even if User Service is down

---

### Pattern 3: Context Propagation

```go
// ✅ GOOD - Pass context through all layers
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, error) {
    // Create product
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, err
    }
    
    // Notify Social Service - context propagated
    if err := u.socialClient.CreateNewsFeed(ctx, &NewsFeed{
        PostId:  product.ID,
        Content: product.Description,
    }); err != nil {
        log.Printf("Failed to create news feed: %v", err)
        // Continue - non-critical
    }
    
    // Notify Notification Service - context propagated
    profileId := _utils.GetProfileIdWithContext(ctx)
    u.notificationClient.Send(ctx, profileId, "Sản phẩm đã được tạo")
    
    return product, nil
}
```

**Benefits**:
- User info automatically available in all services
- Request timeout propagated
- Distributed tracing (future)
- Cancellation support

---

### Pattern 4: Error Handling in Service Calls

```go
// ✅ GOOD - Proper error handling with fallbacks
func (u *ContactUsecase) Sync(ctx context.Context, contacts []Contact) ([]Contact, error) {
    // Extract phones
    phones := make([]string, len(contacts))
    for i, contact := range contacts {
        phones[i] = contact.Phone
    }
    
    // Try to get profiles from User Service
    profileMap, err := u.userClient.GetProfileByPhones(ctx, phones)
    if err != nil {
        // Log error but continue
        log.Printf("Failed to fetch profiles: %v", err)
        profileMap = make(map[string]*sharepb.ProfileItem)
    }
    
    // Link profiles if available
    for i, contact := range contacts {
        if profile, exists := profileMap[contact.Phone]; exists {
            contacts[i].ProfileID = &profile.ProfileId
            contacts[i].FullName = profile.FullName  // Use profile name if available
        }
    }
    
    // Continue with sync
    return u.repo.BulkCreate(ctx, contacts)
}
```

---

## 📡 Real Examples from Codebase

### Example 1: Organization → User Service

**File**: `organization-service/internal/usecase/organization_usecase.go`

```go
func (o *organizationUsecase) GetOrganizationMembers(
    ctx context.Context,
    organizationId uint32,
    offset, limit int,
) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
    // Step 1: Get members from local database
    members, total, err := o.OrganizationMemberRepository.
        FindByOrganizationIdWithPagination(ctx, organizationId, offset, limit, nil, nil)
    if err != nil {
        return nil, 0, err
    }
    
    // Step 2: Extract profile IDs
    profileIds := make([]uint64, len(members))
    for i, member := range members {
        profileIds[i] = uint64(member.UserID)
    }
    
    // Step 3: Batch fetch profiles from User Service
    var profilesMap map[uint64]*sharepb.ProfileItem
    if len(profileIds) > 0 {
        profilesMap, err = o.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
            Ids: profileIds,
        })
        if err != nil {
            // Graceful degradation
            profilesMap = make(map[uint64]*sharepb.ProfileItem)
        }
    } else {
        profilesMap = make(map[uint64]*sharepb.ProfileItem)
    }
    
    // Step 4: Combine members with profiles
    membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
    for i, member := range members {
        membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
            OrganizationMember: member,
            Profile:            profilesMap[uint64(member.UserID)],
        }
    }
    
    return membersWithProfiles, total, nil
}
```

**Communication Details**:
- **From**: Organization Service (Port 50052)
- **To**: User Service (Port 50051)
- **Method**: `ProfileService.GetProfileByIds`
- **Protocol**: gRPC
- **Pattern**: Batch fetching to avoid N+1

---

### Example 2: CRM → User Service

**File**: `crm-service/internal/usecase/contact_usecase.go`

```go
func (s *ContactUsecase) Sync(
    ctx context.Context,
    ownerType base_enum.EOwnerOf,
    contacts []domain.ContactEntity,
) ([]domain.ContactEntity, error) {
    // Step 1: Extract phone numbers
    phones := make([]string, len(contacts))
    for i := 0; i < len(contacts); i++ {
        phones[i] = contacts[i].Phone
    }
    
    // Step 2: Check if phones belong to existing users
    profiles, _ := s.userClient.GetProfileByPhones(ctx, phones)
    
    // Step 3: Link contacts to profiles
    if len(profiles) > 0 {
        for i := 0; i < len(contacts); i++ {
            existedProfile, ok := profiles[contacts[i].Phone]
            if ok {
                contacts[i].ProfileID = &existedProfile.ProfileId
            }
        }
    }
    
    // Step 4: Set ownership
    profileId := _utils.GetProfileIdWithContext(ctx)
    for i := 0; i < len(contacts); i++ {
        contacts[i].OwnerID = profileId
        contacts[i].OwnerOf = ownerType
    }
    
    // Step 5: Bulk insert
    contacts, err := s.repo.Bulk(ctx, contacts)
    if err != nil {
        return nil, err
    }
    
    return contacts, nil
}
```

**Communication Details**:
- **From**: CRM Service (Port 50054)
- **To**: User Service (Port 50051)
- **Method**: `ProfileService.GetProfileByPhones`
- **Purpose**: Link contacts to user accounts
- **Pattern**: Batch lookup by phone numbers

---

### Example 3: Organization → BDSPro Service

**File**: `organization-service/internal/usecase/organization_usecase.go`

```go
func (o *organizationUsecase) GetCurrentDashboard(ctx context.Context) (*organizationpb.DashboardResponse, error) {
    organizationId := _utils.GetOrganizationIdFromContext(ctx)
    
    // Call BDSPro Service for statistics
    bdsproReq := &bdspropb.GetCountByOwnerRequest{
        OwnerOf: sharepb.OwnerOf_organization,
        OwnerId: organizationId,
    }
    
    bdsproResp, err := o.bdsproClient.GetCountByOwner(ctx, bdsproReq)
    if err != nil {
        return nil, err
    }
    
    // Aggregate data
    return &organizationpb.DashboardResponse{
        TotalProduct:     bdsproResp.TotalProduct,
        TotalAsset:       bdsproResp.TotalAsset,
        TotalPost:        bdsproResp.TotalPost,
        TotalProject:     bdsproResp.TotalProject,
        TotalNewsfeed:    0,  // TODO: Call Social Service
        TotalSchedule:    0,  // TODO: Call Appointment Service
        TotalTransaction: 0,  // TODO: Call Transaction Service
    }, nil
}
```

**Communication Details**:
- **From**: Organization Service (Port 50052)
- **To**: BDSPro Service (Port 50053)
- **Method**: `BdsproInternalService.GetCountByOwner`
- **Purpose**: Aggregate statistics for dashboard
- **Pattern**: Internal service methods for stats

---

## 🔧 gRPC Client Configuration

### Connection Setup

**File**: `<service>/wire/wire.go` or initialization

```go
func InitializeGrpcClients() (*Clients, error) {
    // User Service connection
    userConn, err := grpc.Dial(
        "localhost:50051",  // or from config
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithConnectParams(grpc.ConnectParams{
            Backoff: backoff.Config{
                BaseDelay:  1.0 * time.Second,
                Multiplier: 1.6,
                Jitter:     0.2,
                MaxDelay:   120 * time.Second,
            },
            MinConnectTimeout: 20 * time.Second,
        }),
        grpc.WithUnaryInterceptor(clientInterceptor),
    )
    if err != nil {
        return nil, err
    }
    
    // Organization Service connection
    orgConn, err := grpc.Dial("localhost:50052", opts...)
    
    // BDSPro Service connection
    bdsproConn, err := grpc.Dial("localhost:50053", opts...)
    
    return &Clients{
        UserConn:   userConn,
        OrgConn:    orgConn,
        BdsproConn: bdsproConn,
    }, nil
}

// Client interceptor - inject metadata
func clientInterceptor(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error {
    // Extract user info from context
    profileId := _utils.GetProfileIdWithContext(ctx)
    orgId := _utils.GetOrganizationIdFromContext(ctx)
    
    // Create gRPC metadata
    md := metadata.Pairs(
        "profileId", fmt.Sprintf("%d", profileId),
        "organizationId", fmt.Sprintf("%d", orgId),
    )
    
    // Inject into outgoing context
    ctx = metadata.NewOutgoingContext(ctx, md)
    
    // Call service
    return invoker(ctx, method, req, reply, cc, opts...)
}
```

---

## 🔗 Service Dependency Graph

### Primary Dependencies

```
Organization Service depends on:
├── User Service (profiles, user info)
├── BDSPro Service (property stats)
├── Chat Service (create group chats)
├── Transaction Service (deal transactions)
├── Notification Service (notifications)
└── Auth Service (permissions)

BDSPro Service depends on:
├── Organization Service (org info)
├── User Service (owner profiles)
├── Social Service (news feed)
├── Notification Service (notifications)
└── Chat Service (share products)

CRM Service depends on:
├── User Service (profiles)
├── Organization Service (team info)
└── Notification Service (notifications)

Payment Service depends on:
├── User Service (user wallets)
├── Organization Service (org subscriptions)
└── Notification Service (payment alerts)

All Services depend on:
├── Notification Service (notifications)
└── User Service (user profiles)
```

---

## 🚀 Performance Optimizations

### 1. Connection Pooling

```go
// gRPC maintains connection pool automatically
conn, _ := grpc.Dial(
    address,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)

// Reuse connection for multiple calls
client1 := userpb.NewProfileServiceClient(conn)
client2 := userpb.NewAuthServiceClient(conn)  // Same connection
```

### 2. Timeout Configuration

```go
// Set timeout for service calls
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

resp, err := client.GetProfile(ctx, req)
if err != nil {
    if err == context.DeadlineExceeded {
        return nil, errors.New("service timeout")
    }
    return nil, err
}
```

### 3. Retry Mechanism (Future)

```go
// Retry with exponential backoff
func callWithRetry(ctx context.Context, fn func() error) error {
    maxRetries := 3
    backoff := time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(backoff)
            backoff *= 2
        }
    }
    
    return errors.New("max retries exceeded")
}
```

---

## 📊 Communication Statistics

### Service Call Matrix

| From Service | Calls To | Frequency | Critical? |
|-------------|----------|-----------|-----------|
| **Organization** | User | Very High | Yes |
| **Organization** | BDSPro | Medium | No |
| **Organization** | Chat | Low | No |
| **BDSPro** | Organization | Medium | No |
| **BDSPro** | User | High | Yes |
| **BDSPro** | Social | Medium | No |
| **CRM** | User | Very High | Yes |
| **CRM** | Organization | Low | No |
| **All Services** | Notification | High | No |

### Performance Targets

| Metric | Target | Current (Est.) |
|--------|--------|----------------|
| **gRPC Latency** | < 5ms | ~2-5ms |
| **Service Availability** | 99.9% | 99%+ |
| **Connection Pool Size** | 10-100 | 50 |
| **Max Concurrent Calls** | 1000+ | 500+ |

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Production Implementation

