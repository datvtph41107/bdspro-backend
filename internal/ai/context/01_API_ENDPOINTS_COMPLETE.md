# BDSPro Microservices - Complete API Endpoints

**Base URL**: `http://localhost:8080` (Gateway)
**Protocol**: HTTP REST → gRPC Gateway
**Authentication**: JWT Bearer Token (except public routes)

---

## 🔐 Auth Service (Port: 50062)

### OTP Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/otp/request` | Request OTP | ❌ |
| POST | `/v2/auth/otp/verify` | Verify OTP | ❌ |
| POST | `/v2/auth/otp/resend` | Resend OTP | ❌ |

### Token Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/token/refresh` | Refresh access token | ❌ |
| DELETE | `/v2/auth/logout` | Logout user | ✅ |

### Account Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| DELETE | `/v2/auth/delete` | Delete account | ✅ |
| PUT | `/v2/auth/restore` | Restore account | ❌ |
| PUT | `/v2/auth/restore-deleted` | Restore deleted account | ❌ |

### Admin Access Control
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/admin-access` | Create admin access | ✅ |
| PUT | `/v2/auth/admin-access/{id}` | Update admin access | ✅ |
| DELETE | `/v2/auth/admin-access/{id}` | Delete admin access | ✅ |
| GET | `/v2/auth/admin-access` | Get admin access list | ✅ |
| POST | `/v2/auth/admin-access/validate` | Validate admin access | ✅ |
| POST | `/v2/auth/admin-access/log` | Create admin access log | ✅ |
| GET | `/v2/auth/admin-access/log` | Get admin access logs | ✅ |

### QR Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/qr/init` | Initialize QR auth | ❌ |
| POST | `/v2/auth/qr/confirm` | Confirm QR auth | ✅ |
| POST | `/v2/auth/qr/status` | Check QR status | ❌ |

### Organization
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/switch/organization` | Switch organization | ✅ |

### Admin Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/auth/admin/login` | Admin login | ❌ |

---

## 👤 User Service (Port: 50051)

### Profile Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/user/profile/info/{id}` | Get user profile by ID | ✅ |
| GET | `/v2/user/profile/info/me` | Get current user profile | ✅ |
| GET | `/v2/user/profile/accounts/me` | List my accounts | ✅ |
| PUT | `/v2/user/profile/set-info` | Update profile info | ✅ |
| GET | `/v2/user/profile/search/public` | Search public profiles | ✅ |
| GET | `/v2/user/profile/find-by-phone` | Find user by phone | ✅ |

---

## 🏢 Organization Service (Port: 50052)

### Organization Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/org/organization/new` | Create organization | ✅ |
| PUT | `/v2/org/organization/{id}` | Update organization | ✅ |
| GET | `/v2/org/organization` | Get organizations | ✅ |
| GET | `/v2/org/organization/by-member` | Get orgs by user | ✅ |
| GET | `/v2/org/organization/check-is-exist/{id}` | Check if org exists | ✅ |
| GET | `/v2/org/organization/current` | Get current organization | ✅ |
| GET | `/v2/org/organization/current/dashboard` | Get dashboard | ✅ |
| GET | `/v2/org/organization/admin/all-with-members` | Admin: Get all orgs with members | ✅ |
| GET | `/v2/org/organization/global/profile/{id}` | Get global profile | ✅ |
| GET | `/v2/org/organization/count-by-owner` | Count by owner | ✅ |

---

## 🏘️ BDSPro Service (Port: 50053)

### Post Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/bdspro/v2/post` | Create post | ✅ |
| POST | `/v2/bdspro/v2/post/new-expired` | Create expired post | ✅ |
| GET | `/v2/bdspro/v2/post/detail/{id}` | Get post detail | ✅ |
| PUT | `/v2/bdspro/v2/post/{id}` | Update post | ✅ |
| DELETE | `/v2/bdspro/v2/post/{id}` | Delete post | ✅ |
| PUT | `/v2/bdspro/v2/post/hidden` | Update hidden status | ✅ |
| GET | `/v2/bdspro/v2/post/personal` | Get personal posts | ✅ |
| GET | `/v2/bdspro/v2/post/global` | Get global posts | ✅ |
| GET | `/v2/bdspro/v2/post/publish/{profileId}` | Get published posts | ✅ |

### Organization Posts
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/bdspro/v2/post/organization` | Create org post | ✅ |
| GET | `/v2/bdspro/v2/post/organization` | Get org posts | ✅ |

---

## 💼 CRM Service (Port: 50054)

### Contact Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/crm/contact/list/{ownerOf}/{ownerId}` | Search contacts | ✅ |
| GET | `/v2/crm/contact/detail/{id}` | Get contact detail | ✅ |
| POST | `/v2/crm/contact` | Create contact | ✅ |
| PUT | `/v2/crm/contact/edit/{id}` | Update contact | ✅ |
| DELETE | `/v2/crm/contact/{id}` | Delete contact | ✅ |
| POST | `/v2/crm/contact/sync` | Sync contacts | ✅ |
| GET | `/v2/crm/contact/relation-ship/{id}` | Get relationship | ✅ |
| PUT | `/v2/crm/contact/note/{id}` | Update contact note | ✅ |

---

## 💳 Payment Service (Port: 50055)

### Wallet Dashboard
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/payment/wallets/dashboard` | Get wallet dashboard | ✅ |
| GET | `/v2/payment/stats` | Get payment stats | ✅ |
| GET | `/v2/payment/dashboard-metrics` | Get dashboard metrics | ✅ |

### Wallet Transactions
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/payment/wallets/transactions` | Get wallet transactions | ✅ |
| POST | `/v2/payment/wallets/payment` | Make payment | ✅ |
| POST | `/v2/payment/wallets/deposit` | Deposit money | ✅ |

### Payment Methods
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/payment/payment-methods` | Get payment methods | ✅ |
| POST | `/v2/payment/payment-methods` | Create payment method | ✅ |
| PUT | `/v2/payment/payment-methods/{id}` | Update payment method | ✅ |
| DELETE | `/v2/payment/payment-methods/{id}` | Delete payment method | ✅ |

### Reports
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/payment/wallets/{walletId}/report` | Get wallet report | ✅ |
| GET | `/v2/payment/wallets/{walletId}/export` | Export wallet data | ✅ |

### Transaction Types
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/payment/transaction-types` | Get transaction types | ✅ |
| POST | `/v2/payment/transaction-types` | Create transaction type | ✅ |
| PUT | `/v2/payment/transaction-types/{id}` | Update transaction type | ✅ |
| DELETE | `/v2/payment/transaction-types/{id}` | Delete transaction type | ✅ |

### Webhooks
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/payment/sepay/webhook` | Sepay webhook | ❌ |

---

## 🔔 Notification Service (Port: 50056)

### Notification Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/notification/new` | Create notification | ✅ |
| GET | `/v2/notification/count` | Count notifications | ✅ |
| PUT | `/v2/notification/read/{id}` | Mark as read | ✅ |
| PUT | `/v2/notification/read-all` | Mark all as read | ✅ |
| DELETE | `/v2/notification/remove/{id}` | Remove notification | ✅ |
| GET | `/v2/notification/list` | List notifications | ✅ |

### Push Notifications
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/notification/push/token` | Push by token | ✅ |
| POST | `/v2/notification/push/topic` | Push by topic | ✅ |

### Admin History
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/admin-history/search` | Search admin history | ✅ |

### Account Warnings
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/notification/{ownerOf}/account-warning/send` | Send account warning | ✅ |
| GET | `/v2/notification/{ownerOf}/account-warning/list` | Get warning list | ✅ |
| PUT | `/v2/notification/{ownerOf}/account-warning/read` | Mark warning as read | ✅ |
| PUT | `/v2/notification/{ownerOf}/account-warning/acknowledge` | Acknowledge warning | ✅ |

### Warning Templates
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/notification/{ownerOf}/account-warning/template` | Create template | ✅ |
| GET | `/v2/notification/{ownerOf}/account-warning/template/list` | Get template list | ✅ |
| GET | `/v2/notification/{ownerOf}/account-warning/template/{id}` | Get template detail | ✅ |
| GET | `/v2/notification/{ownerOf}/account-warning/log/list` | Get warning logs | ✅ |

### System Warnings (Admin)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/notification/admin/system-warning` | Get admin warnings | ✅ |
| POST | `/v2/notification/admin/system-warning` | Create system warning | ✅ |
| DELETE | `/v2/notification/admin/system-warning/{id}` | Delete system warning | ✅ |

### System Warnings (User)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/notification/system-warning/me` | Get my system warnings | ✅ |

---

## 📱 Social Service (Port: 50058)

### News Feed Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/social/news-feed/{ownerOf}` | Create news feed | ✅ |
| GET | `/v2/social/news-feed/{ownerOf}/{ownerId}` | Get news feed by owner | ✅ |
| PUT | `/v2/social/news-feed/{id}` | Update news feed | ✅ |
| GET | `/v2/social/news-feed/user/{userId}` | Get feed by user | ✅ |
| GET | `/v2/social/news-feed/detail/{newsFeedId}` | Get feed detail | ✅ |
| DELETE | `/v2/social/news-feed/{id}` | Delete news feed | ✅ |
| GET | `/v2/social/news-feed/count-by-owner` | Count by owner | ✅ |
| POST | `/v2/social/news-feed/{newsFeedId}/share` | Share news feed | ✅ |
| PUT | `/v2/social/news-feed/{newsFeedId}/visibility` | Update visibility | ✅ |

### Global News Feed
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v2/social/news-feed/global` | Get global feed | ✅ |
| GET | `/v2/social/news-feed/global/reel` | Get reel feed | ✅ |
| GET | `/v2/social/news-feed/global/{ownerOf}/{ownerId}` | Get global feed by owner | ✅ |

---

## 💬 Chat Service (Port: 50057)

### Conversation Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/chat/conversations` | Create conversation | ✅ |
| GET | `/v2/chat/conversations` | Get conversations | ✅ |
| GET | `/v2/chat/conversation/{conversationId}` | Get conversation by ID | ✅ |
| PUT | `/v2/chat/conversation/{conversationId}/edit` | Edit conversation | ✅ |
| DELETE | `/v2/chat/conversation/delete-conversation/{conversationId}` | Delete conversation | ✅ |
| POST | `/v2/chat/createOrgChannel` | Create org channel | ✅ |
| POST | `/v2/chat/conversation/add-user` | Add user to conversation | ✅ |
| POST | `/v2/chat/conversation/leave` | Leave conversation | ✅ |

### Message Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/chat/messages` | Send message | ✅ |
| POST | `/v2/chat/messages/to-receiver` | Send to receiver | ✅ |
| GET | `/v2/chat/messages` | Get messages | ✅ |
| POST | `/v2/chat/messages/forward` | Forward message | ✅ |
| POST | `/v2/chat/read` | Mark as read | ✅ |
| POST | `/v2/chat/messages/pin` | Pin message | ✅ |
| POST | `/v2/chat/recall` | Recall message | ✅ |
| POST | `/v2/chat/removeUser` | Remove user | ✅ |
| POST | `/v2/chat/messages/edit/{messageId}` | Edit message | ✅ |
| DELETE | `/v2/chat/messages/delete-violation/{messageId}` | Delete violation message | ✅ |
| POST | `/v2/chat/messages/reaction` | React to message | ✅ |

### Conversation Features
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v2/chat/mute` | Mute conversation | ✅ |
| GET | `/v2/chat/search` | Search conversations/messages | ✅ |
| GET | `/v2/chat/conversation/{conversationId}/members` | Get participants | ✅ |
| GET | `/v2/chat/messages/{messageId}/readReceipt` | Get read receipts | ✅ |
| GET | `/v2/chat/messages/{conversationId}/pinned` | Get pinned messages | ✅ |
| GET | `/v2/chat/conversation/{id}/setting` | Get conversation settings | ✅ |
| GET | `/v2/chat/conversation/{conversationId}/files` | Get files | ✅ |
| GET | `/v2/chat/conversation/{conversationId}/images` | Get images | ✅ |
| GET | `/v2/chat/conversation/{conversationId}/links` | Get links | ✅ |
| PUT | `/v2/chat/conversation/background` | Set background image | ✅ |
| GET | `/v2/chat/unread-count` | Get unread count | ✅ |

---

## 📊 Summary Statistics

### Total Endpoints by Service
| Service | Public | Authenticated | Total |
|---------|--------|---------------|-------|
| **Auth Service** | 6 | 14 | 20 |
| **User Service** | 0 | 6 | 6 |
| **Organization Service** | 0 | 10 | 10 |
| **BDSPro Service** | 0 | 10 | 10 |
| **CRM Service** | 0 | 8 | 8 |
| **Payment Service** | 1 | 17 | 18 |
| **Notification Service** | 0 | 28 | 28 |
| **Social Service** | 0 | 12 | 12 |
| **Chat Service** | 0 | 33 | 33 |
| **TOTAL** | **7** | **138** | **145** |

---

## 🔑 Authentication

### Public Routes (No Token Required)
```
/v2/auth/otp/*
/v2/auth/token/refresh
/v2/auth/restore
/v2/auth/restore-deleted
/v2/auth/qr/init
/v2/auth/qr/status
/v2/auth/admin/login
/v2/payment/sepay/webhook
```

### Protected Routes (Token Required)
All other routes require JWT Bearer token in `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

---

## 📝 Request/Response Patterns

### Pagination Request
```json
{
  "page": 1,
  "size": 20,
  "sort": "createdAt",
  "text": "search keyword"
}
```

### Pagination Response
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "size": 20
}
```

### Success Response
```json
{
  "success": true,
  "message": "Operation successful"
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error description"
}
```

---

## 🔗 Service Dependencies

### Inter-Service Communication
- **Auth Service** ← Gateway (authentication)
- **User Service** ← Auth, Organization, BDSPro, CRM, Social, Chat
- **Organization Service** ← BDSPro, CRM, Payment, Notification
- **BDSPro Service** ← Social, Notification
- **Payment Service** ← Organization, Notification
- **Notification Service** ← All services (for notifications)
- **Chat Service** ← BDSPro (product attachments)

---

## 🌐 CORS Configuration
```
Allowed Origins:
- http://localhost:3000
- http://localhost:3001
- http://localhost:8000

Allowed Methods:
- GET, POST, PUT, DELETE, OPTIONS

Allowed Headers:
- Origin, Content-Type, Authorization, Lang

Credentials: Allowed
Max Age: 12 hours
```

---

**Last Updated**: October 15, 2025
**API Version**: v2
**Base Path**: `/v2/`
**Total Endpoints**: 145+

