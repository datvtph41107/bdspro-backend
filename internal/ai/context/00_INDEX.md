# BDSPro Microservices - Documentation Index

## 📚 Complete Context Documentation

**Last Updated**: October 18, 2025
**Documentation Version**: v1.1
**Total Services**: 19 microservices
**Status**: ✅ Complete + Quick Reference

---

## ⚡ FASTEST START (NEW!)

### 🚀 **[Quick Reference Guide](23_QUICK_REFERENCE_GUIDE.md)** ⭐⭐⭐
   - **CREATE API trong 7 bước**
   - **All import aliases & conventions**
   - **Repository pattern template**
   - **Common patterns & code templates**
   - **All services ports at a glance**
   - **Special services quick lookup**
   - **Debug checklist**
   - **~5 min scan, instant lookup**
   - 💡 **Bookmark này để tra cứu nhanh!**

---

## 🎯 Essential Reading (Start Here)

1. **[System Overview](00_SYSTEM_OVERVIEW.md)** ⭐
   - Architecture overview
   - All 19 services summary
   - Technology stack summary
   - Business domains
   - ~10 min read

2. **[API Endpoints Complete](01_API_ENDPOINTS_COMPLETE.md)** ⭐
   - All 145+ API endpoints
   - Organized by service
   - Request/response patterns
   - Authentication info
   - ~20 min read

3. **[Development Workflow](07_DEVELOPMENT_WORKFLOW.md)** ⭐
   - How to add new APIs
   - Step-by-step guides
   - Best practices
   - Quick commands
   - ~15 min read

---

## 🔧 Core Services Documentation

### 4. **[Auth Service](02_AUTH_SERVICE.md)**
   - Authentication & Authorization
   - JWT token management
   - OAuth integration (Google, Facebook, Zalo)
   - Role-Based Access Control (RBAC)
   - QR authentication
   - Admin access control
   - **Port**: 50062 (gRPC)
   - **Tables**: 12
   - **Endpoints**: 20
   - ~30 min read

### 5. **[Organization Service](05_ORGANIZATION_SERVICE.md)**
   - Organization management
   - Branch/office management
   - Team/group management
   - Deal management
   - Commission tracking
   - Permission system
   - **Port**: 50052 (gRPC)
   - **Tables**: 25+
   - **Endpoints**: 50+
   - ~25 min read

### 6. **[BDSPro Service](06_BDSPRO_SERVICE.md)**
   - Real estate property management
   - Post/listing management
   - Product management
   - Asset management
   - Project management
   - Property search
   - **Port**: 50053 (gRPC)
   - **Tables**: 50+
   - **Endpoints**: 40+
   - ~30 min read

---

## 📦 Foundation & Shared Components

### 7. **[Shared Common Library](03_SHARED_COMMON_LIBRARY.md)**
   - Base entities (BaseEntity, AuditBase)
   - Common DTOs (Pagable, ErrorDTO, HistoryDTO)
   - CRUD interfaces & implementations
   - Common enumerations
   - Code generators
   - Provider base classes
   - **Used by**: All services
   - ~25 min read

### 8. **[Technology Stack](04_TECHNOLOGY_STACK.md)**
   - Complete tech stack overview
   - Go frameworks (Gin, gRPC)
   - Database (PostgreSQL, Redis)
   - DevOps tools (Docker, Wire, Buf)
   - Why we chose each technology
   - Performance characteristics
   - ~20 min read

### 9. **[Database Architecture](09_DATABASE_ARCHITECTURE.md)**
   - Database per service strategy
   - Common schema patterns
   - Core database schemas
   - Indexing strategy
   - Security & optimization
   - Backup strategy
   - **Databases**: 14
   - **Tables**: 150+
   - ~20 min read

---

## 🌟 Other Services Summary

### 10. **[Other Services Summary](08_OTHER_SERVICES_SUMMARY.md)**
   - CRM Service (Contact & Lead management)
   - Payment Service (Wallet & transactions)
   - Notification Service (Multi-channel notifications)
   - Social Service (News feed & social features)
   - Chat Service (Real-time messaging)
   - Relay Service (WebSocket)
   - Map Service (Location services)
   - Search Service (Full-text search)
   - Appointment Service (Scheduling)
   - Marketing Service (Campaigns)
   - Transaction Service (Deal transactions)
   - Membership Service (Subscriptions)
   - File Service (File management)
   - Assistant Service (AI - planned)
   - ~25 min read

### 11. **[Additional Services Deep Dive](22_ADDITIONAL_SERVICES_DEEP_DIVE.md)** 🆕
   - **TQD Service**: POI management, geolocation
   - **Assistant Service**: AI DeepSeek integration
   - **Payment Service**: Advanced wallet & locking
   - **Membership Service**: Subscription plans
   - **Task Service**: Project management
   - **Detailed implementation patterns**
   - ~30 min read

---

## 📊 Documentation Statistics

### Coverage
- ✅ **Quick Reference Guide**: ⚡ NEW - Fast lookup
- ✅ **System Overview**: Complete
- ✅ **API Documentation**: 145+ endpoints
- ✅ **Core Services**: 3 detailed (Auth, Organization, BDSPro)
- ✅ **Supporting Services**: 14 summarized
- ✅ **Additional Services**: 5 deep dive (TQD, Assistant, Payment, Membership, Task)
- ✅ **Shared Libraries**: Complete
- ✅ **Technology Stack**: Complete
- ✅ **Database Architecture**: Complete
- ✅ **Development Workflow**: Complete

### Metrics
- **Total Documentation Pages**: 23 files
- **Total Services Covered**: 19/19 (100%)
- **Total Lines**: ~7,500+ lines
- **Total Endpoints Documented**: 145+
- **Total Database Tables**: 150+
- **Code Templates**: 10+ ready-to-use
- **Estimated Reading Time**: ~5-6 hours (full read)
- **Quick Lookup Time**: ~5 minutes (Quick Reference)

---

## 🎓 Learning Path

### For New Developers

**Day 1: Quick Start** ⚡ (UPDATED!)
1. **Bookmark** [Quick Reference Guide](23_QUICK_REFERENCE_GUIDE.md) (5 min scan)
2. Read [System Overview](00_SYSTEM_OVERVIEW.md) (10 min)
3. Browse [API Endpoints](01_API_ENDPOINTS_COMPLETE.md) (20 min)
4. **Practice**: Tạo API đầu tiên theo Quick Reference (30 min)

**Day 2: Deep Dive into Core**
1. Study [Auth Service](02_AUTH_SERVICE.md) (30 min)
2. Study [Shared Common Library](03_SHARED_COMMON_LIBRARY.md) (25 min)
3. Practice with [Development Workflow](07_DEVELOPMENT_WORKFLOW.md) (15 min)

**Day 3: Business Logic**
1. Study [Organization Service](05_ORGANIZATION_SERVICE.md) (25 min)
2. Study [BDSPro Service](06_BDSPRO_SERVICE.md) (30 min)
3. Review [Additional Services](22_ADDITIONAL_SERVICES_DEEP_DIVE.md) (30 min)

**Day 4: Data & Practice**
1. Study [Database Architecture](09_DATABASE_ARCHITECTURE.md) (20 min)
2. Practice adding features (2-3 hours)
3. **Always keep Quick Reference open for lookup!**

### For Experienced Developers

**Instant Start** ⚡ (15 min):
1. [Quick Reference Guide](23_QUICK_REFERENCE_GUIDE.md) - 5 min scan
2. [System Overview](00_SYSTEM_OVERVIEW.md) - 10 min skim
3. Start coding with Quick Reference as lookup!

**As-Needed Reference**:
- **Quick Reference**: Patterns, templates, conventions
- **Service Docs**: Detailed business logic
- **Workflow Docs**: Step-by-step when stuck

---

## 🔍 Find What You Need

### Looking for...

**"Cần tra cứu NHANH patterns/conventions?"** ⚡
→ [Quick Reference Guide](23_QUICK_REFERENCE_GUIDE.md) - Instant lookup!

**"How do I add a new API?"** 
→ [Quick Reference](23_QUICK_REFERENCE_GUIDE.md) - 7 bước + code templates
→ [Development Workflow](07_DEVELOPMENT_WORKFLOW.md) - Chi tiết step-by-step

**"Import aliases là gì?"**
→ [Quick Reference](23_QUICK_REFERENCE_GUIDE.md) - Section "Import Aliases"

**"Repository pattern viết như thế nào?"**
→ [Quick Reference](23_QUICK_REFERENCE_GUIDE.md) - Templates + examples

**"Service X có port gì?"**
→ [Quick Reference](23_QUICK_REFERENCE_GUIDE.md) - Bảng tất cả ports

**"Lấy profileId/organizationId ra sao?"**
→ [Quick Reference](23_QUICK_REFERENCE_GUIDE.md) - Code snippets

**"What are all the API endpoints?"**
→ [API Endpoints Complete](01_API_ENDPOINTS_COMPLETE.md) - All 145+ endpoints

**"How does authentication work?"**
→ [Auth Service](02_AUTH_SERVICE.md) - Complete auth flow

**"How do I use the shared library?"**
→ [Shared Common Library](03_SHARED_COMMON_LIBRARY.md) - BaseEntity, CRUD, etc.

**"Payment service hoạt động ra sao?"**
→ [Additional Services Deep Dive](22_ADDITIONAL_SERVICES_DEEP_DIVE.md) - Payment section

**"What database tables exist?"**
→ [Database Architecture](09_DATABASE_ARCHITECTURE.md) - All schemas

**"Which service handles X?"**
→ [Other Services Summary](08_OTHER_SERVICES_SUMMARY.md) - Service responsibilities

---

## 📝 Documentation Conventions

### Symbols
- ⭐ Essential reading
- ✅ Complete documentation
- 🔄 In progress
- 📋 Reference material
- 🎯 Quick guide
- 💡 Best practice
- ⚠️ Important note
- ❌ Common mistake

### Code Blocks
```go
// Go code examples
```

```sql
-- SQL examples
```

```bash
# Shell commands
```

```json
// JSON examples
```

### File References
- Relative paths: `internal/domain/product.go`
- Absolute paths: `/Users/.../go-microservices/...`
- Proto files: `shared/protobuf/schema/user/profile.proto`

---

## 🔗 External Resources

### Official Documentation
- **Go**: https://go.dev/doc/
- **gRPC**: https://grpc.io/docs/
- **GORM**: https://gorm.io/docs/
- **Protocol Buffers**: https://protobuf.dev/
- **Gin**: https://gin-gonic.com/docs/

### Tools
- **Wire**: https://github.com/google/wire
- **Buf**: https://buf.build/docs/
- **Swag**: https://github.com/swaggo/swag
- **Air**: https://github.com/cosmtrek/air

### Repository
- **GitHub**: [Link to repository]
- **Issues**: [Link to issues]
- **Wiki**: [Link to wiki]

---

## 🤝 Contributing to Documentation

### How to Update
1. Edit the relevant `.md` file in `/context/`
2. Follow existing format and structure
3. Update this index if adding new files
4. Commit with descriptive message

### Documentation Standards
- Clear, concise language
- Code examples for complex concepts
- Real-world use cases
- Vietnamese terms where appropriate (BĐS, Tin đăng, etc.)
- Keep documents under 1000 lines
- Update "Last Updated" date

---

## 📞 Support & Questions

### Where to Ask
- **Technical Questions**: Team lead or senior developers
- **Architecture Questions**: Solution architect
- **Business Logic**: Product manager
- **Database Questions**: Database admin

### Common Questions Answered
- [Development Workflow](07_DEVELOPMENT_WORKFLOW.md) - "How do I...?"
- [API Endpoints](01_API_ENDPOINTS_COMPLETE.md) - "Which endpoint...?"
- [Database Architecture](09_DATABASE_ARCHITECTURE.md) - "Where is X stored...?"

---

## 🎯 Quick Reference Cards

### Service Ports
```
Gateway:        8080 (HTTP)
User:           50051 (gRPC)
Organization:   50052 (gRPC)
BDSPro:         50053 (gRPC)
CRM:            50054 (gRPC)
Payment:        50055 (gRPC)
Notification:   50056 (gRPC)
Chat:           50057 (gRPC)
Social:         50058 (gRPC)
Auth:           50062 (gRPC)
File:           8083 (HTTP)
```

### Common Commands
```bash
# Generate protobuf
make buf-<service>

# Generate wire
make wire <service>

# Start service
make <service>-grpc

# Run all services
make start-all
```

### Import Aliases
```go
import (
    _dto "common/domain/dto"
    _models "common/domain/entity"
    _enum "common/domain/enum"
    _err "common/domain/err"
    _crud "common/domain/crud"
    _utils "common/utils"
    _provider "common/provider"
)
```

---

## 📅 Changelog

### v1.0 - October 15, 2025
- ✅ Initial complete documentation
- ✅ All 19 services documented
- ✅ 145+ API endpoints catalogued
- ✅ Database architecture complete
- ✅ Development workflow guide
- ✅ Technology stack overview
- ✅ Shared library documentation

---

## 🎉 Documentation Complete!

**Total Reading Time**: ~3-4 hours for complete understanding
**Coverage**: 100% of microservices
**Status**: Production-ready documentation

**Congratulations!** You now have access to comprehensive documentation covering the entire BDSPro Microservices ecosystem. Use this index to navigate to specific topics as needed.

**Pro Tip**: Bookmark this index and the most relevant service docs for quick reference during development.

---

**Created with ❤️ for the BDSPro Development Team**
**Last Updated**: October 15, 2025

