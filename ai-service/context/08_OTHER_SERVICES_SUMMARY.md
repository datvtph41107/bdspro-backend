# Other Services - Complete Summary

## 📋 Overview

This document provides a comprehensive overview of the remaining microservices in the BDSPro ecosystem that haven't been detailed individually.

---

## 💼 CRM Service (Port: 50054)

### Purpose
Customer Relationship Management - Manage contacts, leads, sales pipeline

### Core Features
- **Contact Management**: CRUD operations for contacts
- **Lead Tracking**: Track potential customers
- **Pipeline Management**: Sales pipeline stages
- **Stage Management**: Custom pipeline stages
- **Rule Management**: Automation rules
- **Friend Management**: Social connections
- **Following System**: Follow/unfollow contacts
- **Block System**: Block unwanted contacts
- **Sharing**: Share contacts with team

### Key Entities
```
- contact.go              # Customer contact
- lead.go                 # Sales lead
- pipeline.go             # Sales pipeline
- stage.go                # Pipeline stages
- rule.go                 # Automation rules
- friend.go               # Friend connections
- friend_group.go         # Friend grouping
- follow.go               # Follow relationships
- block.go                # Blocked contacts
- sharing.go              # Contact sharing
- invitation_install.go   # App install invitations
```

### Key APIs
```
GET  /v2/crm/contact/list/{ownerOf}/{ownerId}
GET  /v2/crm/contact/detail/{id}
POST /v2/crm/contact
PUT  /v2/crm/contact/edit/{id}
DELETE /v2/crm/contact/{id}
POST /v2/crm/contact/sync
GET  /v2/crm/contact/relation-ship/{id}
PUT  /v2/crm/contact/note/{id}
```

### Database Tables: 10+
### Use Cases: Lead management, sales pipeline, contact sync

---

## 💳 Payment Service (Port: 50055)

### Purpose
Payment processing, wallet management, transaction tracking

### Core Features
- **Wallet Management**: User wallets for virtual currency
- **Wallet Transactions**: Deposits, withdrawals, payments
- **Payment Methods**: Multiple payment methods
- **Transaction Types**: Categorize transactions
- **Bank Integration**: Bank account linking
- **Sepay Integration**: Payment gateway webhook
- **Dashboard**: Payment statistics
- **Reports**: Transaction reports & exports

### Key Entities
```
- wallet.go               # User wallet
- wallet_transaction.go   # Transaction records
- payment_method.go       # Payment methods
- transaction_type.go     # Transaction categories
- withdrawal_request.go   # Withdrawal requests
- deal_transaction.go     # Deal-related payments
```

### Key APIs
```
GET  /v2/payment/wallets/dashboard
GET  /v2/payment/wallets/transactions
POST /v2/payment/wallets/payment
POST /v2/payment/wallets/deposit
GET  /v2/payment/payment-methods
POST /v2/payment/payment-methods
GET  /v2/payment/transaction-types
POST /v2/payment/sepay/webhook
GET  /v2/payment/stats
```

### Special Features
- **Sepay Webhook**: Auto-process bank transfers
- **Wallet System**: Virtual currency management
- **Transaction Tracking**: Detailed transaction history
- **Commission Calculation**: For deals/transactions

### Database Tables: 8+
### Use Cases: Subscription payments, deal transactions, wallet top-up

---

## 🔔 Notification Service (Port: 50056)

### Purpose
Multi-channel notification delivery, activity logging

### Core Features
- **In-App Notifications**: System notifications
- **Push Notifications**: FCM (Firebase Cloud Messaging)
- **Email Notifications**: Email delivery
- **Admin History**: Track admin actions
- **Account Warnings**: Send warnings to users
- **Warning Templates**: Reusable warning templates
- **System Warnings**: Platform-wide alerts
- **Notification Categories**: Categorized notifications

### Key Entities
```
- notification.go         # Notification record
- history.go              # Activity history
- deal_history.go         # Deal-specific history
- account_warning.go      # User warnings
- warning_template.go     # Warning templates
- warning_log.go          # Warning action logs
- system_warning.go       # System-wide warnings
```

### Key APIs
```
POST /v2/notification/new
GET  /v2/notification/count
PUT  /v2/notification/read/{id}
PUT  /v2/notification/read-all
DELETE /v2/notification/remove/{id}
GET  /v2/notification/list
POST /v2/notification/push/token
POST /v2/notification/push/topic
GET  /v2/admin-history/search
POST /v2/notification/{ownerOf}/account-warning/send
GET  /v2/notification/{ownerOf}/account-warning/list
```

### Integration Points
- **Firebase**: Push notifications to mobile
- **Email Service**: SMTP integration
- **SMS Service**: For OTP and alerts
- **All Services**: Receive notification requests

### Database Tables: 7+
### Use Cases: Activity tracking, user notifications, admin actions

---

## 📱 Social Service (Port: 50058)

### Purpose
Social networking features, news feed management

### Core Features
- **News Feed**: Social media-style posts
- **Comments**: Comment on posts
- **Likes**: Like/dislike functionality
- **Shares**: Share posts
- **Reels**: Short video content
- **Visibility Control**: Public/private posts
- **Owner Types**: User, Group, Organization posts
- **Reporting**: Report inappropriate content

### Key Entities
```
- news_feed.go            # News feed posts
- news_feed_media.go      # Post media
- comment.go              # Comments
- like.go                 # Likes/reactions
- report.go               # Content reports
```

### Key APIs
```
POST /v2/social/news-feed/{ownerOf}
GET  /v2/social/news-feed/{ownerOf}/{ownerId}
PUT  /v2/social/news-feed/{id}
GET  /v2/social/news-feed/detail/{newsFeedId}
DELETE /v2/social/news-feed/{id}
POST /v2/social/news-feed/{newsFeedId}/share
PUT  /v2/social/news-feed/{newsFeedId}/visibility
GET  /v2/social/news-feed/global
GET  /v2/social/news-feed/global/reel
```

### Integration Points
- **BDSPro Service**: Share property posts
- **User Service**: Get user profiles
- **Notification Service**: Notify on likes/comments

### Database Tables: 5+
### Use Cases: Social networking, content sharing, engagement

---

## 💬 Chat Service (Port: 50057)

### Purpose
Real-time messaging, conversation management

### Core Features
- **Conversations**: One-on-one and group chats
- **Messages**: Text, images, files, products
- **Message Types**: Text, image, video, file, link
- **Read Receipts**: Track message reading
- **Message Reactions**: React to messages
- **Pinned Messages**: Pin important messages
- **Message Recall**: Delete sent messages
- **Message Forwarding**: Forward to other chats
- **Search**: Search messages and conversations
- **Background Images**: Custom chat backgrounds
- **Mute**: Mute conversation notifications

### Key Entities
```
- conversation.go         # Chat conversation
- message.go              # Chat message
- participant.go          # Conversation members
- message_reaction.go     # Message reactions
- read_receipt.go         # Message read status
- background_image.go     # Chat backgrounds
```

### Key APIs
```
POST /v2/chat/conversations
GET  /v2/chat/conversations
GET  /v2/chat/conversation/{conversationId}
POST /v2/chat/messages
POST /v2/chat/messages/to-receiver
GET  /v2/chat/messages
POST /v2/chat/messages/forward
POST /v2/chat/read
POST /v2/chat/messages/pin
POST /v2/chat/recall
POST /v2/chat/messages/reaction
GET  /v2/chat/unread-count
```

### Conversation Types
- **Private**: One-on-one chat
- **Group**: Multi-user group chat
- **Channel**: Broadcast channel

### Integration Points
- **Organization Service**: Auto-create group chats
- **BDSPro Service**: Share products in chat
- **User Service**: Get user profiles
- **Relay Service**: WebSocket notifications

### Database Tables: 8+
### Use Cases: Team communication, customer chat, support

---

## 🔄 Relay Service (Port: 50064)

### Purpose
WebSocket management, real-time event broadcasting

### Core Features
- **WebSocket Connections**: Manage persistent connections
- **Event Broadcasting**: Push events to clients
- **Room Management**: Organize connections by rooms
- **Presence Tracking**: Online/offline status
- **Connection Pooling**: Efficient connection management

### Use Cases
- Real-time chat notifications
- Live updates for posts/products
- Online status tracking
- Activity notifications

### Integration Points
- **All Services**: Broadcast events via Relay
- **Chat Service**: Real-time message delivery
- **Notification Service**: Push notifications

---

## 🗺️ Map Service (Port: 8210)

### Purpose
Location services, geocoding, nearby search

### Core Features
- **Location Search**: Search by address
- **Nearby Search**: Find nearby properties
- **Geocoding**: Convert address to coordinates
- **Reverse Geocoding**: Convert coordinates to address
- **Distance Calculation**: Calculate distances

### Integration Points
- **BDSPro Service**: Geocode property locations
- **Organization Service**: Office/branch locations

---

## 🔍 Search Service (Port: 8211)

### Purpose
Full-text search, advanced search capabilities

### Core Features
- **Full-Text Search**: Search across multiple fields
- **Fuzzy Search**: Typo-tolerant search
- **Faceted Search**: Filter by multiple criteria
- **Search Suggestions**: Auto-complete
- **Search History**: Track user searches

### Integration Points
- **BDSPro Service**: Property search
- **CRM Service**: Contact search
- **User Service**: User search

---

## 📅 Appointment Service (Port: 50059)

### Purpose
Scheduling and appointment management

### Core Features
- **Appointment CRUD**: Create, manage appointments
- **Calendar Integration**: Calendar view
- **Reminder Notifications**: Appointment reminders
- **Time Slot Management**: Available time slots
- **Recurring Appointments**: Repeating schedules

### Use Cases
- Property viewing appointments
- Customer meetings
- Team meetings

### Integration Points
- **BDSPro Service**: Schedule property viewings
- **CRM Service**: Customer appointments
- **Notification Service**: Send reminders

---

## 📢 Marketing Service (Port: 50060)

### Purpose
Marketing campaigns, lead generation

### Core Features
- **Campaign Management**: Create marketing campaigns
- **Lead Capture**: Capture leads from campaigns
- **Campaign Analytics**: Track campaign performance
- **Email Campaigns**: Email marketing
- **SMS Campaigns**: SMS marketing

### Integration Points
- **CRM Service**: Convert leads to contacts
- **Notification Service**: Send campaign messages
- **BDSPro Service**: Property marketing

---

## 💰 Transaction Service (Port: 50061)

### Purpose
Deal transaction management, contract tracking

### Core Features
- **Transaction CRUD**: Manage transactions
- **Contract Management**: Contract documents
- **Deal Costs**: Transaction costs
- **Transaction Types**: Buy, sell, rent
- **Transaction Status**: Pending, completed, cancelled

### Integration Points
- **Organization Service**: Deal tracking
- **Payment Service**: Payment processing
- **Notification Service**: Transaction updates

---

## 👥 Membership Service (Port: 50063)

### Purpose
Subscription and membership management

### Core Features
- **Membership Plans**: Different tiers (Free, Pro, Enterprise)
- **Plan Features**: Feature access control
- **Subscription Management**: Active subscriptions
- **Billing**: Subscription billing
- **Plan Upgrades**: Upgrade/downgrade plans

### Integration Points
- **Payment Service**: Process subscription payments
- **Organization Service**: Organization subscriptions
- **Auth Service**: Feature access validation

---

## 📁 File Service (Port: 8083)

### Purpose
File upload, storage, and management

### Core Features
- **File Upload**: Multi-file upload
- **Image Processing**: Resize, crop, optimize
- **Video Processing**: Video transcoding (with FFmpeg)
- **File Storage**: Local or cloud storage
- **File Retrieval**: Serve files
- **File Types**: Images, videos, documents

### Key APIs
```
POST /v1/file/upload
GET  /v1/file/load/{filename}
```

### Integration Points
- **All Services**: Upload and retrieve files
- **BDSPro Service**: Property images/videos
- **Social Service**: Post media
- **Chat Service**: Message attachments

---

## 🤖 Assistant Service (Port: TBD)

### Purpose
AI assistant, chatbot functionality

### Core Features (Planned)
- **Chatbot**: AI-powered conversations
- **Property Recommendations**: Smart suggestions
- **Query Understanding**: Natural language processing
- **Auto-Responses**: Automated replies

### Integration Points
- **All Services**: Provide context for AI

---

## 📊 Services Comparison Matrix

| Service | Port | DB Tables | APIs | Primary Tech |
|---------|------|-----------|------|--------------|
| **CRM** | 50054 | 10+ | 8+ | gRPC, PostgreSQL |
| **Payment** | 50055 | 8+ | 18+ | gRPC, PostgreSQL, Sepay |
| **Notification** | 50056 | 7+ | 28+ | gRPC, FCM, PostgreSQL |
| **Social** | 50058 | 5+ | 12+ | gRPC, PostgreSQL |
| **Chat** | 50057 | 8+ | 33+ | gRPC, WebSocket, PostgreSQL |
| **Relay** | 50064 | 2+ | - | WebSocket |
| **Map** | 8210 | 3+ | 5+ | HTTP, Location APIs |
| **Search** | 8211 | - | 3+ | HTTP, Full-text search |
| **Appointment** | 50059 | 5+ | 8+ | gRPC, PostgreSQL |
| **Marketing** | 50060 | 6+ | 10+ | gRPC, PostgreSQL |
| **Transaction** | 50061 | 7+ | 12+ | gRPC, PostgreSQL |
| **Membership** | 50063 | 4+ | 6+ | gRPC, PostgreSQL |
| **File** | 8083 | 3+ | 2 | HTTP, File System |
| **Assistant** | TBD | - | - | AI/ML |

---

## 🔗 Service Dependencies Graph

```
Gateway Service (Entry Point)
    ↓
    ├─→ Auth Service ──────────────────┐
    │                                   │
    ├─→ User Service ←─────────────────┤
    │       ↓                           │
    ├─→ Organization Service            │
    │       ↓                           │
    ├─→ BDSPro Service ─────────────→ Notification Service
    │       ↓                           ↑
    ├─→ CRM Service ────────────────────┤
    │       ↓                           │
    ├─→ Payment Service ────────────────┤
    │       ↓                           │
    ├─→ Social Service ─────────────────┤
    │       ↓                           │
    ├─→ Chat Service ←──────────────────┘
    │       ↓
    ├─→ Transaction Service
    │       ↓
    ├─→ Marketing Service
    │       ↓
    ├─→ Appointment Service
    │       ↓
    ├─→ Membership Service
    │       ↓
    ├─→ File Service
    │       ↓
    ├─→ Map Service
    │       ↓
    ├─→ Search Service
    │       ↓
    └─→ Relay Service (WebSocket)
```

---

## 📈 Service Complexity Ranking

### High Complexity
1. **Organization Service** - Multi-level hierarchy, deals, permissions
2. **BDSPro Service** - Complex domain, many entities
3. **Auth Service** - Security critical, OAuth, RBAC
4. **Chat Service** - Real-time, WebSocket, many features

### Medium Complexity
5. **CRM Service** - Pipeline management, rules
6. **Payment Service** - Financial transactions, integrations
7. **Social Service** - Social features, engagement
8. **Notification Service** - Multi-channel delivery

### Low Complexity
9. **Transaction Service** - Simpler domain
10. **Appointment Service** - Scheduling logic
11. **Marketing Service** - Campaign management
12. **Membership Service** - Subscription logic
13. **File Service** - File operations
14. **Map Service** - Location services
15. **Search Service** - Search operations
16. **Relay Service** - WebSocket relay

---

## 🎯 Service Responsibilities Summary

### User Management
- **Auth Service**: Authentication & authorization
- **User Service**: User profiles & info

### Business Operations
- **Organization Service**: Company structure
- **BDSPro Service**: Real estate management
- **CRM Service**: Customer relationships
- **Transaction Service**: Deal transactions

### Financial
- **Payment Service**: Payments & wallets
- **Membership Service**: Subscriptions

### Communication
- **Chat Service**: Messaging
- **Notification Service**: Notifications
- **Social Service**: Social features
- **Relay Service**: Real-time events

### Support Services
- **File Service**: File management
- **Map Service**: Location services
- **Search Service**: Search functionality
- **Appointment Service**: Scheduling
- **Marketing Service**: Campaigns

---

**Last Updated**: October 15, 2025
**Services Documented**: 14 (out of 19 total)
**Status**: ✅ Complete Summary

