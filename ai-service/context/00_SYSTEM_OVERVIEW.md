# BDSPro Microservices - System Overview

## 📋 Tổng Quan Hệ Thống

**Dự án**: BDSPro (Bất Động Sản Pro) - Comprehensive Real Estate Management System
**Kiến trúc**: Microservices Architecture với Go
**Giao thức**: gRPC, HTTP REST API
**Framework**: Gin (HTTP), gRPC-Gateway, Wire (DI)
**Database**: PostgreSQL, Redis, MongoDB
**Documentation**: Swagger/OpenAPI
**Authentication**: JWT (JSON Web Token)

---

## 🏗️ Kiến Trúc Tổng Thể

### Clean Architecture
Hệ thống tuân theo **Clean Architecture** với phân tách rõ ràng:
- **Domain Layer**: Business entities & core logic
- **Interface Layer**: Contracts (repositories, providers)
- **Usecase Layer**: Business logic & orchestration
- **Infrastructure Layer**: Technical implementations (DB, HTTP, gRPC, clients)

### Microservices Pattern
- **API Gateway**: Single entry point cho tất cả requests
- **Service Communication**: gRPC inter-service communication
- **Service Discovery**: Dynamic service routing
- **Shared Libraries**: Common domain models & utilities

---

## 🎯 Core Services (19 Services)

### 1. **Gateway Service** (Port 8080)
- **Vai trò**: API Gateway - single entry point
- **Chức năng**: 
  - Request routing
  - JWT authentication & authorization
  - Protocol translation (HTTP ↔ gRPC)
  - Load balancing
  - Swagger documentation aggregation
- **Endpoints**: `/v1/*`, `/v2/*`, `/swagger/*`

### 2. **Auth Service** (gRPC: 50062)
- **Vai trò**: Authentication & Authorization
- **Chức năng**:
  - User authentication (OTP, Email, Phone)
  - OAuth2 integration (Google, Facebook, Zalo)
  - Session management
  - Role & Permission management
  - Admin access control
- **Key Features**: 
  - Multi-provider authentication
  - JWT token generation & validation
  - RBAC (Role-Based Access Control)

### 3. **User Service** (HTTP: 8081, gRPC: 50051)
- **Vai trò**: User & Profile Management
- **Chức năng**:
  - User profile CRUD
  - Profile customization
  - User rating & reviews
  - User reporting
- **Dual Mode**: HTTP + gRPC servers

### 4. **Organization Service** (gRPC: 50052)
- **Vai trò**: Organization & Team Management
- **Chức năng**:
  - Organization management
  - Branch management
  - Team/Group management
  - Deal & transaction tracking
  - Member management
- **Complex Domain**: Multi-level hierarchy (Organization → Branch → Team → Member)

### 5. **BDSPro Service** (gRPC: 50053)
- **Vai trò**: Real Estate Property Management
- **Chức năng**:
  - Property listing (Posts)
  - Product management
  - Market data
  - Property search & filter
  - Property analytics
- **Core Business**: Main real estate operations

### 6. **CRM Service** (gRPC: 50054)
- **Vai trò**: Customer Relationship Management
- **Chức năng**:
  - Contact management
  - Lead tracking
  - Customer lifecycle management
  - Sales pipeline
- **Integration**: Works with Organization, User, BDSPro services

### 7. **Payment Service** (gRPC: 50055)
- **Vai trò**: Payment & Billing
- **Chức năng**:
  - Payment processing
  - Wallet management
  - Bank integration
  - Transaction history
  - Subscription billing
- **Critical**: Handles financial transactions

### 8. **Notification Service** (gRPC: 50056)
- **Vai trò**: Notification & Messaging
- **Chức năng**:
  - Push notifications (FCM)
  - Email notifications
  - SMS notifications
  - In-app notifications
  - Notification history
- **Multi-channel**: Email, SMS, Push, In-app

### 9. **Chat Service** (HTTP: 8082, gRPC: 50057)
- **Vai trò**: Real-time Chat & Messaging
- **Chức năng**:
  - One-on-one chat
  - Group chat
  - Message history
  - Real-time messaging
- **Technology**: WebSocket for real-time

### 10. **Social Service** (gRPC: 50058)
- **Vai trò**: Social Features
- **Chức năng**:
  - News feed
  - Comments & Likes
  - Social interactions
  - Activity tracking
- **Engagement**: User engagement features

### 11. **File Service** (HTTP: 8083)
- **Vai trò**: File Upload & Management
- **Chức năng**:
  - File upload (images, videos, documents)
  - File storage
  - File retrieval
  - Image processing (resize, crop)
- **Storage**: File system / Cloud storage

### 12. **Appointment Service** (gRPC: 50059)
- **Vai trò**: Appointment & Scheduling
- **Chức năng**:
  - Schedule appointments
  - Calendar management
  - Reminder notifications
- **Use Case**: Property viewing appointments

### 13. **Marketing Service** (gRPC: 50060)
- **Vai trò**: Marketing Campaigns
- **Chức năng**:
  - Campaign management
  - Lead generation
  - Marketing analytics
- **Integration**: Works with CRM

### 14. **Transaction Service** (gRPC: 50061)
- **Vai trò**: Transaction Management
- **Chức năng**:
  - Transaction tracking
  - Deal management
  - Transaction history
- **Business Critical**: Property transaction lifecycle

### 15. **Membership Service** (gRPC: 50063)
- **Vai trò**: Subscription & Membership
- **Chức năng**:
  - Membership plans
  - Subscription management
  - Feature access control
- **Monetization**: Premium features

### 16. **Relay Service** (gRPC: 50064)
- **Vai trò**: WebSocket Relay
- **Chức năng**:
  - WebSocket management
  - Real-time event broadcasting
- **Real-time**: Event-driven communication

### 17. **Map Service** (HTTP/gRPC: 8210)
- **Vai trò**: Location & Mapping
- **Chức năng**:
  - Location search
  - Nearby properties
  - Map integration
- **Geo**: Geospatial queries

### 18. **Search Service** (HTTP: 8211)
- **Vai trò**: Search & Indexing
- **Chức năng**:
  - Full-text search
  - Advanced search
  - Search indexing
- **Performance**: Fast search capabilities

### 19. **Feedback Service** (gRPC: TBD)
- **Vai trò**: User Feedback
- **Chức năng**:
  - Feedback collection
  - Rating & reviews
  - Feedback analytics

---

## 📚 Shared Libraries

### 1. **shared/common** - Common Domain
- **Location**: `/shared/common/`
- **Chức năng**:
  - Base entities (BaseEntity, AuditBase)
  - Common DTOs (Pagable, ErrorDTO, HistoryDTO, NotificationDTO)
  - Common enums (Context keys, History actions, Notification types)
  - CRUD base implementations (ICrudRepo, BaseUsecase)
  - Code generators (Product code, Asset code, Post code)
  - Base repository provider (CrudRepo)

### 2. **shared/protobuf** - Protocol Buffers
- **Location**: `/shared/protobuf/`
- **Chức năng**:
  - Protocol buffer definitions (.proto files)
  - Generated gRPC code
  - API contracts between services
  - Swagger/OpenAPI schemas

### 3. **shared/code** - Code Generation
- **Location**: `/shared/code/`
- **Chức năng**:
  - Wire generation scripts
  - Protobuf generation scripts
  - Service management scripts
  - Build automation

---

## 🔄 Communication Patterns

### 1. Client → Gateway → Microservice
```
Client (HTTP/REST)
    ↓
Gateway Service (Port 8080)
    ↓ JWT Auth
    ↓ gRPC
Microservice (gRPC Port)
    ↓
Response (JSON)
```

### 2. Service-to-Service Communication
```
Service A (gRPC Client)
    ↓ gRPC
Service B (gRPC Server)
```

### 3. Event-Driven (via Relay Service)
```
Service A
    ↓ Event
Relay Service (WebSocket)
    ↓ Broadcast
Connected Clients
```

---

## 🔐 Security & Authentication

### JWT Authentication Flow
1. **Login**: User authenticates via Auth Service
2. **Token Generation**: Auth Service generates JWT
3. **Token Storage**: Client stores JWT
4. **API Request**: Client sends JWT in Authorization header
5. **Gateway Validation**: Gateway validates JWT
6. **Context Injection**: Gateway injects user info into gRPC metadata
7. **Service Access**: Microservice extracts user info from context

### Role-Based Access Control (RBAC)
- **Roles**: Admin, User, Member, etc.
- **Permissions**: Granular permissions per resource
- **Role Groups**: Group multiple roles
- **Access Control**: Enforced at Gateway & Service level

---

## 💾 Data Management

### Databases
- **PostgreSQL**: Primary database for all services
- **Redis**: Caching, session storage
- **MongoDB**: Document storage (optional)

### Data Patterns
- **Repository Pattern**: Database abstraction
- **Unit of Work**: Transaction management
- **CQRS**: Command-Query separation (implicit)
- **Soft Delete**: BaseEntity với DeletedAt field

### Audit Trail
- **AuditBase**: Tracks CreatedBy, UpdatedBy
- **History Service**: Logs all important actions
- **Timestamps**: CreatedAt, UpdatedAt auto-managed

---

## 🚀 Deployment Architecture

### Development
```
Local Machine
├── Gateway (Port 8080)
├── Services (Various Ports)
├── PostgreSQL (Port 5432)
└── Redis (Port 6379)
```

### Production (Docker)
```
Docker Network: internal
├── gateway (Port 8080)
├── user (Port 50051)
├── organization (Port 50052)
├── bdspro_grpc (Port 50053)
├── auth (Port 50062)
├── postgres (Port 5432)
└── redis (Port 6379)
```

---

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP), gRPC
- **ORM**: GORM
- **DI**: Wire
- **Validation**: go-playground/validator

### Communication
- **gRPC**: Inter-service communication
- **Protocol Buffers**: Message serialization
- **gRPC-Gateway**: HTTP to gRPC translation
- **WebSocket**: Real-time communication

### Documentation
- **Swagger**: API documentation
- **OpenAPI**: API specification
- **Protobuf Docs**: Generated from .proto files

### DevOps
- **Docker**: Containerization
- **Docker Compose**: Local orchestration
- **Air**: Hot reload for development
- **Make**: Build automation

---

## 📊 Key Design Patterns

### 1. Clean Architecture
- Dependency Inversion
- Layer separation
- Interface-based programming

### 2. Repository Pattern
- Data access abstraction
- Testability
- Database independence

### 3. Dependency Injection
- Wire for compile-time DI
- Constructor injection
- Interface dependencies

### 4. API Gateway Pattern
- Single entry point
- Cross-cutting concerns (auth, logging)
- Protocol translation

### 5. CQRS (Implicit)
- Separate read/write models
- Optimized queries
- Domain-driven design

### 6. Event Sourcing (Partial)
- History tracking
- Audit logging
- Event broadcasting

---

## 🎯 Business Domains

### Real Estate Core
- **Properties**: Posts, Products, Markets
- **Transactions**: Deals, Contracts
- **Search**: Property discovery

### Customer Management
- **CRM**: Contacts, Leads, Customers
- **Organization**: Companies, Branches, Teams
- **Users**: Profiles, Ratings

### Platform Services
- **Communication**: Chat, Notifications
- **Social**: News Feed, Interactions
- **Files**: Document management

### Business Operations
- **Payments**: Transactions, Billing
- **Membership**: Subscriptions
- **Marketing**: Campaigns, Analytics

---

## 📈 Scalability Considerations

### Horizontal Scaling
- **Stateless Services**: Easy to replicate
- **Load Balancing**: Gateway distributes load
- **Database Pooling**: Connection management

### Performance Optimization
- **Caching**: Redis for frequently accessed data
- **Database Indexing**: Optimized queries
- **gRPC**: Efficient binary protocol
- **Connection Pooling**: Reuse connections

### Resilience
- **Circuit Breaker**: Prevent cascading failures
- **Retry Logic**: Handle transient failures
- **Timeout Management**: Prevent hanging requests
- **Graceful Degradation**: Fallback mechanisms

---

## 🔍 Monitoring & Observability

### Logging
- Structured logging
- Request tracing
- Error tracking

### Health Checks
- Service health endpoints
- Database connectivity
- External service status

### Metrics (Planned)
- Request rate
- Response time
- Error rate
- Resource usage

---

## 📝 Development Workflow

### 1. Add New Feature
1. Define in `.proto` file
2. Generate protobuf (`make buf-<service>`)
3. Implement in handler (infra/handler)
4. Write usecase logic (internal/usecase)
5. Create repository (infra/postgre)
6. Add Swagger comments
7. Generate wire (`make wire <service>`)
8. Test & verify

### 2. Service Communication
1. Define interface (internal/interface/provider)
2. Implement client (infra/client)
3. Inject via Wire
4. Call from usecase

### 3. Database Changes
1. Update domain entity
2. Create migration
3. Update repository
4. Test queries

---

## 🎓 Code Organization Principles

### SOLID Principles
- **S**: Single Responsibility (each layer has one job)
- **O**: Open/Closed (extend via interfaces)
- **L**: Liskov Substitution (interface implementations)
- **I**: Interface Segregation (specific interfaces)
- **D**: Dependency Inversion (depend on abstractions)

### Clean Code
- Meaningful names
- Small functions
- Single responsibility
- Comments for Swagger/Wire
- Consistent formatting

---

## 🚧 Current Status & Future Plans

### Implemented ✅
- Core microservices architecture
- API Gateway with routing
- JWT authentication
- gRPC communication
- Database integration
- Wire dependency injection
- Swagger documentation
- Docker deployment

### In Progress 🔄
- Advanced monitoring
- Performance optimization
- Enhanced security
- More integrations

### Planned 📋
- Kubernetes deployment
- Service mesh (Istio)
- Distributed tracing
- Advanced analytics
- Mobile app support

---

**Last Updated**: October 15, 2025
**Total Services**: 19 microservices
**Total Lines of Code**: ~100,000+ lines
**Architecture**: Clean Architecture + Microservices
**Status**: Production-ready with ongoing enhancements

