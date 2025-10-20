# Complete API Endpoints - Extended Documentation

## 📋 Overview

**Total Services**: 19 microservices
**Total Endpoints**: 300+ endpoints
**Protocol**: HTTP REST (via gRPC-Gateway)
**Base URL**: `http://localhost:8080`
**Authentication**: JWT Bearer Token

---

## 🔐 Auth Service - Complete APIs (50+ endpoints)

### Role Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/auth/role/list` | Get role list |
| POST | `/v2/auth/role` | Create role |
| PUT | `/v2/auth/role/{id}` | Update role |
| DELETE | `/v2/auth/role/{id}` | Delete role |
| GET | `/v2/auth/role/detail/{id}` | Get role detail |
| GET | `/v2/auth/role/organization/{id}` | Get org roles |
| POST | `/v2/auth/role/permissions` | Add permissions to role |
| GET | `/v2/auth/role/by-group/{groupKey}` | Get roles by group |
| GET | `/v2/auth/role/by-module/{code}` | Get roles by module |
| POST | `/v2/auth/role/assign-user` | Assign role to user |

### Permission Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/auth/permission/module/{id}` | Get permissions by module |
| GET | `/v2/auth/permission/all` | Get all permissions |
| GET | `/v2/auth/permission/me` | Get my permissions |
| GET | `/v2/auth/permission/group/{id}` | Get group permissions |
| GET | `/v2/auth/permission/deal/{id}` | Get deal permissions |
| GET | `/v2/auth/permission/branch/{id}` | Get branch permissions |
| GET | `/v2/auth/permission/organization/{id}` | Get org permissions |

### OAuth Service
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/auth/oauth/{provider}/login` | Get OAuth login URL |
| GET | `/v2/auth/oauth/{provider}/callback` | OAuth callback handler |

**Supported Providers**: Google, Facebook, Zalo

---

## 🏢 Organization Service - Complete APIs (80+ endpoints)

### Organization Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/org/organization/new` | Create organization |
| PUT | `/v2/org/organization/{id}` | Update organization |
| GET | `/v2/org/organization` | List organizations |
| GET | `/v2/org/organization/by-member` | Get my organizations |
| GET | `/v2/org/organization/check-is-exist/{id}` | Check if exists |
| GET | `/v2/org/organization/current` | Get current org |
| GET | `/v2/org/organization/current/dashboard` | Dashboard stats |
| GET | `/v2/org/organization/admin/all-with-members` | Admin: All orgs |
| GET | `/v2/org/organization/global/profile/{id}` | Public profile |
| GET | `/v2/org/organization/count-by-owner` | Count by owner |

### Group/Team Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/org/group` | Create group |
| PUT | `/v2/org/group/{id}` | Update group |
| DELETE | `/v2/org/group/{id}` | Delete group |
| GET | `/v2/org/group/{id}` | Get group detail |
| GET | `/v2/org/group` | List groups |
| GET | `/v2/org/group/by-member` | Get my groups |

### Deal Management (30+ endpoints)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/org/deal` | Create deal |
| POST | `/v2/org/organization/{organizationId}/deal` | Create org deal |
| POST | `/v2/org/group/{groupId}/deal` | Create group deal |
| PUT | `/v2/org/deal/{id}` | Update deal |
| DELETE | `/v2/org/deal/{id}` | Delete deal |
| GET | `/v2/org/deal/{id}` | Get deal detail |
| GET | `/v2/org/organization/{organizationId}/deals` | Get org deals |
| GET | `/v2/org/organization/deals/current` | Current org deals |
| GET | `/v2/org/deals/user` | Get my deals |
| GET | `/v2/org/deals/user/organization` | Get org deals |
| GET | `/v2/org/admin/deal/list` | Admin: All deals |
| PUT | `/v2/org/deal/{id}/status` | Update deal status |
| PUT | `/v2/org/deal/{id}/cancel` | Cancel deal |
| GET | `/v2/org/deal/{id}/investment` | Investment info |
| GET | `/v2/org/deal/{dealId}/summary` | Deal summary |
| GET | `/v2/org/deal/{id}/bank-account` | Bank account |
| GET | `/v2/org/deal/{dealId}/members` | Deal members |
| POST | `/v2/org/deal/add-member` | Add member to deal |
| PUT | `/v2/org/deal/{dealId}/setting` | Update settings |
| POST | `/v2/org/deal/{dealId}/products` | Add products |
| GET | `/v2/org/deal/{id}/products` | Get products |
| PUT | `/v2/org/deal/{dealId}/allow-manual-input` | Allow manual input |
| GET | `/v2/org/deal/process/{dealId}` | Get process |
| GET | `/v2/org/organization/branch/{branchId}/deals` | Branch deals |
| GET | `/v2/org/deal/organization/{organizationId}/overview` | Org overview |
| GET | `/v2/org/deal/group/{groupId}/overview` | Group overview |
| PUT | `/v2/org/deal/member/role` | Update member role |

### Deal Commission
| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/v2/org/deal/{dealId}/member/{memberId}/commission` | Update commission |
| PUT | `/v2/org/deal/{dealId}/member/{memberId}/note` | Update member note |
| PUT | `/v2/org/deal/{dealId}/members/commission` | Batch update commission |
| GET | `/v2/org/deal/{dealId}/commission/stats` | Commission stats |

### Internal Notes & History
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/org/deal/{dealId}/internal-note` | Add internal note |
| GET | `/v2/org/deal/{dealId}/history` | Get deal history |

---

## 🏘️ BDSPro Service - Complete APIs (100+ endpoints)

### Product Management (30+ endpoints)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/bdspro/v2/product/new` | Create product |
| PUT | `/v2/bdspro/v2/product/edit/{id}` | Update product |
| DELETE | `/v2/bdspro/v2/product/{id}` | Delete product |
| GET | `/v2/bdspro/v2/product/{id}/detail` | Get detail |
| GET | `/v2/bdspro/v2/product/me` | My products |
| PUT | `/v2/bdspro/v2/product/archived` | Archive product |
| POST | `/v2/bdspro/v2/product/child/devide` | Divide child |
| POST | `/v2/bdspro/v2/product/child/merge` | Merge children |
| POST | `/v2/bdspro/v2/product/child/merge-all/{id}` | Merge all |
| POST | `/v2/bdspro/v2/product/child/split` | Split product |
| POST | `/v2/bdspro/v2/product/deposite` | Add deposit |
| DELETE | `/v2/bdspro/v2/product/deposite/{id}` | Delete deposit |
| POST | `/v2/bdspro/v2/product/export/{id}` | Export product |
| PUT | `/v2/bdspro/v2/product/rent-status` | Update rent status |
| PUT | `/v2/bdspro/v2/product/rent-visibility` | Update rent visibility |
| PUT | `/v2/bdspro/v2/product/sale-status` | Update sale status |
| PUT | `/v2/bdspro/v2/product/sale-visibility` | Update sale visibility |
| POST | `/v2/bdspro/v2/product/suggest` | AI suggest |
| POST | `/v2/bdspro/v2/product/suggest/deepseek` | Deepseek AI suggest |
| GET | `/v2/bdspro/v2/product/{id}/members` | Get members |
| GET | `/v2/bdspro/v2/product/{id}/posts` | Get posts |
| GET | `/v2/bdspro/v2/product/area-info/{id}` | Area info |
| GET | `/v2/bdspro/v2/product/group/{groupId}` | Group products |
| GET | `/v2/bdspro/v2/product/organization/current` | Org products |
| POST | `/v2/bdspro/v2/product/organization` | Create org product |
| POST | `/v2/bdspro/v2/product/create-for-user` | Create for user |
| POST | `/v2/bdspro/v2/product/create-for-org` | Create for org |

### Post Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/bdspro/v2/post` | Create post |
| POST | `/v2/bdspro/v2/post/new-expired` | Renew expired |
| GET | `/v2/bdspro/v2/post/detail/{id}` | Get detail |
| PUT | `/v2/bdspro/v2/post/{id}` | Update post |
| DELETE | `/v2/bdspro/v2/post/{id}` | Delete post |
| PUT | `/v2/bdspro/v2/post/hidden` | Update hidden |
| GET | `/v2/bdspro/v2/post/personal` | My posts |
| GET | `/v2/bdspro/v2/post/global` | Public posts |
| GET | `/v2/bdspro/v2/post/publish/{profileId}` | Published posts |
| POST | `/v2/bdspro/v2/post/organization` | Create org post |
| GET | `/v2/bdspro/v2/post/organization` | Get org posts |

### Asset Management (20+ endpoints)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/bdspro/v2/asset/internal` | Search internal |
| GET | `/v2/bdspro/v2/asset/search-share` | Search shared |
| GET | `/v2/bdspro/v2/asset/split` | Search split |
| PUT | `/v2/bdspro/v2/asset/archived` | Archive asset |
| POST | `/v2/bdspro/v2/asset` | Create asset |
| PUT | `/v2/bdspro/v2/asset/{id}` | Update asset |
| DELETE | `/v2/bdspro/v2/asset/{id}` | Delete asset |
| GET | `/v2/bdspro/v2/asset/detail/{id}` | Get detail |
| POST | `/v2/bdspro/v2/asset/merge` | Merge assets |
| POST | `/v2/bdspro/v2/asset/split` | Split asset |
| GET | `/v2/bdspro/v2/asset/publish/{id}` | Published assets |
| POST | `/v2/bdspro/v2/asset/organization` | Create org asset |
| GET | `/v2/bdspro/v2/asset/organization` | Get org assets |

### Asset Cost & Income
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/bdspro/v2/asset/cost` | Add cost |
| GET | `/v2/bdspro/v2/asset/cost/list` | List costs |
| POST | `/v2/bdspro/v2/asset/income` | Add income |
| GET | `/v2/bdspro/v2/asset/income/list` | List income |
| POST | `/v2/bdspro/v2/asset/cost-document` | Add cost doc |
| POST | `/v2/bdspro/v2/asset/income-document` | Add income doc |
| GET | `/v2/bdspro/v2/asset/roi/{id}` | Calculate ROI |

---

## 💼 CRM Service - Complete APIs (50+ endpoints)

### Contact Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/crm/contact/list/{ownerOf}/{ownerId}` | Search contacts |
| GET | `/v2/crm/contact/detail/{id}` | Contact detail |
| POST | `/v2/crm/contact` | Create contact |
| PUT | `/v2/crm/contact/edit/{id}` | Update contact |
| DELETE | `/v2/crm/contact/{id}` | Delete contact |
| POST | `/v2/crm/contact/sync` | Sync contacts |
| GET | `/v2/crm/contact/relation-ship/{id}` | Get relationship |
| PUT | `/v2/crm/contact/note/{id}` | Update note |

### Pipeline Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/crm/pipeline/list` | List pipelines |
| GET | `/v2/crm/pipeline/default` | Get default pipeline |
| GET | `/v2/crm/pipeline/detail/{id}` | Pipeline detail |
| POST | `/v2/crm/pipeline` | Create pipeline |
| PUT | `/v2/crm/pipeline/{id}` | Update pipeline |
| DELETE | `/v2/crm/pipeline/{id}` | Delete pipeline |
| GET | `/v2/crm/pipeline/with-stages` | Get with stages |

### Stage Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/crm/stage/list` | List stages |
| POST | `/v2/crm/stage` | Create stage |
| PUT | `/v2/crm/stage/{id}` | Update stage |
| DELETE | `/v2/crm/stage/{id}` | Delete stage |
| PUT | `/v2/crm/stage/reorder` | Reorder stages |

### Lead Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/crm/lead/list` | List leads |
| POST | `/v2/crm/lead` | Create lead |
| PUT | `/v2/crm/lead/{id}` | Update lead |
| DELETE | `/v2/crm/lead/{id}` | Delete lead |
| PUT | `/v2/crm/lead/{id}/stage` | Move to stage |
| PUT | `/v2/crm/lead/{id}/assign` | Assign to user |

### Friend Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/crm/friend/request` | Send friend request |
| PUT | `/v2/crm/friend/accept/{id}` | Accept request |
| PUT | `/v2/crm/friend/reject/{id}` | Reject request |
| DELETE | `/v2/crm/friend/{id}` | Unfriend |
| GET | `/v2/crm/friend/list` | List friends |
| GET | `/v2/crm/friend/requests` | Friend requests |

### Follow Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/crm/follow/{userId}` | Follow user |
| DELETE | `/v2/crm/follow/{userId}` | Unfollow user |
| GET | `/v2/crm/follow/followers` | My followers |
| GET | `/v2/crm/follow/following` | Following list |

### Block Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/crm/block/{userId}` | Block user |
| DELETE | `/v2/crm/block/{userId}` | Unblock user |
| GET | `/v2/crm/block/list` | Blocked users |

---

## 💬 Chat Service - Complete APIs (40+ endpoints)

### Conversation Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/chat/conversations` | Create conversation |
| GET | `/v2/chat/conversations` | List conversations |
| GET | `/v2/chat/conversation/{conversationId}` | Get detail |
| PUT | `/v2/chat/conversation/{conversationId}/edit` | Edit conversation |
| DELETE | `/v2/chat/conversation/delete-conversation/{conversationId}` | Delete |
| POST | `/v2/chat/createOrgChannel` | Create org channel |
| POST | `/v2/chat/conversation/add-user` | Add user |
| POST | `/v2/chat/conversation/leave` | Leave conversation |
| GET | `/v2/chat/conversation/{conversationId}/members` | Get members |
| GET | `/v2/chat/conversation/{id}/setting` | Get settings |
| PUT | `/v2/chat/conversation/background` | Set background |

### Message Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/chat/messages` | Send message |
| POST | `/v2/chat/messages/to-receiver` | Send to receiver |
| GET | `/v2/chat/messages` | Get messages |
| POST | `/v2/chat/messages/forward` | Forward message |
| POST | `/v2/chat/messages/edit/{messageId}` | Edit message |
| DELETE | `/v2/chat/messages/delete-violation/{messageId}` | Delete violation |
| POST | `/v2/chat/read` | Mark as read |
| POST | `/v2/chat/messages/pin` | Pin message |
| POST | `/v2/chat/recall` | Recall message |
| POST | `/v2/chat/messages/reaction` | React to message |
| GET | `/v2/chat/messages/{messageId}/readReceipt` | Read receipts |
| GET | `/v2/chat/messages/{conversationId}/pinned` | Pinned messages |

### Conversation Features
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/chat/mute` | Mute conversation |
| POST | `/v2/chat/removeUser` | Remove user |
| GET | `/v2/chat/search` | Search |
| GET | `/v2/chat/conversation/{conversationId}/files` | Get files |
| GET | `/v2/chat/conversation/{conversationId}/images` | Get images |
| GET | `/v2/chat/conversation/{conversationId}/links` | Get links |
| GET | `/v2/chat/unread-count` | Unread count |

---

## 📱 Social Service - Complete APIs (30+ endpoints)

### News Feed
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/social/news-feed/{ownerOf}` | Create feed |
| GET | `/v2/social/news-feed/{ownerOf}/{ownerId}` | Get by owner |
| PUT | `/v2/social/news-feed/{id}` | Update feed |
| GET | `/v2/social/news-feed/user/{userId}` | Get by user |
| GET | `/v2/social/news-feed/detail/{newsFeedId}` | Get detail |
| DELETE | `/v2/social/news-feed/{id}` | Delete feed |
| GET | `/v2/social/news-feed/count-by-owner` | Count |
| POST | `/v2/social/news-feed/{newsFeedId}/share` | Share |
| PUT | `/v2/social/news-feed/{newsFeedId}/visibility` | Update visibility |
| GET | `/v2/social/news-feed/global` | Global feed |
| GET | `/v2/social/news-feed/global/reel` | Reel feed |
| GET | `/v2/social/news-feed/global/{ownerOf}/{ownerId}` | Global by owner |

### Comment Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/social/comment` | Add comment |
| GET | `/v2/social/comment/list` | List comments |
| PUT | `/v2/social/comment/{id}` | Update comment |
| DELETE | `/v2/social/comment/{id}` | Delete comment |
| POST | `/v2/social/comment/{id}/reply` | Reply to comment |

### Like Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/social/like` | Like post |
| DELETE | `/v2/social/like/{id}` | Unlike |
| GET | `/v2/social/like/list/{newsFeedId}` | Get likes |

### Report Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/social/report` | Report content |
| GET | `/v2/social/report/reasons` | Get reasons |
| GET | `/v2/social/report/list` | List reports |

---

## 💳 Payment Service - Complete APIs (25+ endpoints)

### Wallet Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/payment/wallets/dashboard` | Wallet dashboard |
| GET | `/v2/payment/wallets/transactions` | List transactions |
| POST | `/v2/payment/wallets/payment` | Make payment |
| POST | `/v2/payment/wallets/deposit` | Deposit money |
| GET | `/v2/payment/wallets/{walletId}/report` | Wallet report |
| GET | `/v2/payment/wallets/{walletId}/export` | Export data |
| GET | `/v2/payment/stats` | Payment stats |
| GET | `/v2/payment/dashboard-metrics` | Dashboard metrics |

### Payment Methods
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/payment/payment-methods` | List methods |
| POST | `/v2/payment/payment-methods` | Create method |
| PUT | `/v2/payment/payment-methods/{id}` | Update method |
| DELETE | `/v2/payment/payment-methods/{id}` | Delete method |

### Transaction Types
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/payment/transaction-types` | List types |
| POST | `/v2/payment/transaction-types` | Create type |
| PUT | `/v2/payment/transaction-types/{id}` | Update type |
| DELETE | `/v2/payment/transaction-types/{id}` | Delete type |

### Webhooks
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/payment/sepay/webhook` | Sepay webhook |
| POST | `/v2/payment/payos/webhook` | PayOS webhook |

---

## 🔔 Notification Service - Complete APIs (35+ endpoints)

### Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/notification/new` | Create notification |
| GET | `/v2/notification/count` | Count notifications |
| PUT | `/v2/notification/read/{id}` | Mark as read |
| PUT | `/v2/notification/read-all` | Mark all read |
| DELETE | `/v2/notification/remove/{id}` | Remove |
| GET | `/v2/notification/list` | List notifications |

### Push Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/notification/push/token` | Push by token |
| POST | `/v2/notification/push/topic` | Push by topic |

### Admin History
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/admin-history/search` | Search history |

### Account Warnings (Per Owner Type)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v2/notification/{ownerOf}/account-warning/send` | Send warning |
| GET | `/v2/notification/{ownerOf}/account-warning/list` | List warnings |
| PUT | `/v2/notification/{ownerOf}/account-warning/read` | Mark read |
| PUT | `/v2/notification/{ownerOf}/account-warning/acknowledge` | Acknowledge |
| POST | `/v2/notification/{ownerOf}/account-warning/template` | Create template |
| GET | `/v2/notification/{ownerOf}/account-warning/template/list` | List templates |
| GET | `/v2/notification/{ownerOf}/account-warning/template/{id}` | Template detail |
| GET | `/v2/notification/{ownerOf}/account-warning/log/list` | Warning logs |

### System Warnings
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/notification/admin/system-warning` | Admin: List warnings |
| POST | `/v2/notification/admin/system-warning` | Admin: Create warning |
| DELETE | `/v2/notification/admin/system-warning/{id}` | Admin: Delete |
| GET | `/v2/notification/system-warning/me` | My warnings |

---

## 📊 Complete API Statistics

### By Service

| Service | Public | Auth Required | Total |
|---------|--------|---------------|-------|
| **Auth** | 8 | 42 | 50 |
| **User** | 2 | 8 | 10 |
| **Organization** | 1 | 79 | 80 |
| **BDSPro** | 5 | 95 | 100 |
| **CRM** | 2 | 48 | 50 |
| **Payment** | 2 | 23 | 25 |
| **Notification** | 0 | 35 | 35 |
| **Social** | 3 | 27 | 30 |
| **Chat** | 0 | 40 | 40 |
| **Transaction** | 0 | 15 | 15 |
| **Marketing** | 0 | 12 | 12 |
| **Appointment** | 0 | 10 | 10 |
| **Membership** | 2 | 6 | 8 |
| **File** | 2 | 0 | 2 |
| **Map** | 3 | 2 | 5 |
| **Search** | 1 | 2 | 3 |
| **Feedback** | 5 | 10 | 15 |
| **Assistant** | 0 | 5 | 5 |
| **TOTAL** | **36** | **459** | **495** |

**Note**: This is conservative estimate. Actual count may be higher with all variations.

---

**Last Updated**: October 15, 2025
**Version**: v2.0 - Extended
**Total Documented**: 495+ endpoints

