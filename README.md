# BDSPro Golang Microservices

A comprehensive real estate management system built with Go microservices architecture, featuring gRPC communication, REST APIs, and modern development practices.

## 🏗️ System Architecture

This system consists of multiple microservices, each handling specific business domains:

### Core Services
- **Gateway Service** - API Gateway and routing
- **User Service** - User management and authentication
- **Organization Service** - Organization and business domain management
- **BDSPro Service** - Real estate property management
- **CRM Service** - Customer relationship management
- **Payment Service** - Payment processing and transactions
- **Notification Service** - Push notifications and messaging
- **File Service** - File upload and management
- **Marketing Service** - Marketing campaigns and advertising

### Communication Services
- **Chat Service** - Real-time messaging
- **Relay Service** - WebSocket and real-time communication
- **Social Service** - Social features and interactions

### Support Services
- **Appointment Service** - Scheduling and appointments
- **Membership Service** - Membership and subscription management
- **Task Service** - Task and project management
- **Search Service** - Search functionality
- **Map Service** - Location and mapping services

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make
- Buf CLI
- Air (for hot reloading)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/hieubap/bdspro-golang-microservices.git
cd bdspro-golang-microservices
```

2. **Install dependencies**
```bash
make init
```

3. **Generate protobuf files**
```bash
make buf-all
```

4. **Start all services**
```bash
make start-all
```

## 📋 Makefile Commands

The project includes a comprehensive Makefile for managing all aspects of development and deployment.

### 🔧 Development Commands

#### Service Management
```bash
# Start all services in background
make start-all

# Stop all services
make stop-all

# Start individual services
make gateway-http      # Gateway service (HTTP)
make user-http         # User service (HTTP)
make user-grpc         # User service (gRPC)
make organization      # Organization service
make crm-grpc          # CRM service
make bdspro-grpc       # BDSPro service
make notification      # Notification service
make social            # Social service
make appointment       # Appointment service
make chat-http         # Chat service (HTTP)
make chat-grpc         # Chat service (gRPC)
make relay-grpc        # Relay service
make payment-grpc      # Payment service
make membership-grpc   # Membership service
make file-http         # File service
make marketing-grpc    # Marketing service
```

#### Protobuf Generation
```bash
# Generate all protobuf files
make buf-all

# Generate specific service protobuf files
make buf-user          # User service
make buf-organization  # Organization service
make buf-bdspro        # BDSPro service
make buf-crm           # CRM service
make buf-social        # Social service
make buf-notification  # Notification service
make buf-appointment   # Appointment service
make buf-payment       # Payment service
make buf-chat          # Chat service
make buf-task          # Task service
make buf-marketing     # Marketing service
make buf-shared        # Shared types
```

#### Swagger Documentation
```bash
# Generate swagger documentation for file service
make swag-file

# Merge all swagger files
make merge-swagger
```

#### Wire Dependency Injection
```bash
# Generate wire files for specific service
make wire [service-name]
```

### 🏃‍♂️ Running Services

#### Development Mode (with hot reload)
```bash
# Start individual service with Air (hot reload)
make gateway-http      # Uses Air for automatic restart
make user-http         # Uses Air for automatic restart
make organization      # Uses Air for automatic restart
```

#### Production Mode
```bash
# Run service directly with Go
make crm              # Runs without Air
```

### 🐳 Docker Commands

```bash
# Build and run with Docker Compose
docker-compose up -d

# Stop Docker services
docker-compose down

# View logs
docker-compose logs -f [service-name]
```

## 📁 Project Structure

```
bdspro-golang-microservices/
├── gateway-service/          # API Gateway
├── user-service/            # User management
├── organization-service/    # Organization management
├── bdspro-service/         # Real estate management
├── crm-service/            # Customer relationship
├── payment-service/        # Payment processing
├── notification-service/   # Notifications
├── file-service/          # File management
├── marketing-service/     # Marketing campaigns
├── chat-service/          # Real-time chat
├── relay-service/         # WebSocket relay
├── social-service/        # Social features
├── appointment-service/   # Appointments
├── membership-service/    # Memberships
├── task-service/         # Task management
├── search-service/       # Search functionality
├── map-service/          # Mapping services
├── common/               # Shared utilities
├── protobuf/            # Protocol buffers
│   ├── schema/          # .proto files
│   ├── types/           # Generated Go types
│   └── docs/            # Generated documentation
├── code/                # Code generation scripts
├── docs/                # API documentation
├── Makefile             # Build and development commands
├── docker-compose.yml   # Docker orchestration
└── deploy.sh           # Deployment script
```

## 🔌 API Documentation

### Swagger UI
- **Gateway API**: http://localhost:8080/swagger/index.html
- **File Service**: http://localhost:8081/swagger/index.html

### gRPC Services
All services expose gRPC endpoints on their respective ports. Use tools like:
- **gRPCui**: `grpcui -plaintext localhost:50051`
- **BloomRPC**: For testing gRPC calls
- **Postman**: For REST API testing

## 🛠️ Development Workflow

### 1. Adding a New Service
```bash
# 1. Create service directory
mkdir new-service

# 2. Initialize Go module
cd new-service && go mod init new-service

# 3. Add to Makefile
# Add build and run commands

# 4. Generate protobuf
make buf-[service-name]

# 5. Start development
make [service-name]-grpc
```

### 2. Modifying Protobuf
```bash
# 1. Edit .proto files in protobuf/schema/
# 2. Generate new types
make buf-[service-name]
# 3. Update swagger
make merge-swagger
```

### 3. Database Migrations
```bash
# Run migrations for specific service
cd [service-name] && go run migrate.go
```

## 🔧 Configuration

### Environment Variables
Each service has its own configuration files:
- `config/config.yaml` - Main configuration
- `config/docker_config.yaml` - Docker-specific config
- `.env` files for environment-specific settings

### Database
- **PostgreSQL**: Primary database
- **Redis**: Caching and sessions
- **MongoDB**: Document storage (if needed)

## 🚀 Deployment

### Local Development
```bash
# Start all services
make start-all

# Monitor logs
tail -f */tmp/*.log

# Stop all services
make stop-all
```

### Production Deployment
```bash
# Deploy with Docker
./deploy.sh

# Or use Docker Compose
docker-compose -f docker-compose.prod.yml up -d
```

## 🧪 Testing

### Unit Tests
```bash
# Run tests for specific service
cd [service-name] && go test ./...

# Run all tests
find . -name "*.go" -exec dirname {} \; | sort | uniq | xargs -I {} sh -c 'cd {} && go test ./...'
```

### Integration Tests
```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test ./tests/...
```

## 📊 Monitoring

### Health Checks
- **Gateway**: http://localhost:8080/health
- **Services**: http://localhost:[port]/health

### Metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Port conflicts**: Check if ports are already in use
```bash
lsof -i :8080
```

2. **Protobuf generation errors**: Ensure buf CLI is installed
```bash
go install github.com/bufbuild/buf/cmd/buf@latest
```

3. **Service not starting**: Check logs
```bash
tail -f [service-name]/tmp/*.log
```

4. **Database connection issues**: Verify database is running
```bash
docker-compose ps
```

### Getting Help
- Check the logs in `*/tmp/*.log`
- Review service-specific documentation
- Open an issue on GitHub

## 🔗 Links

- **API Documentation**: [Swagger UI](http://localhost:8080/swagger/index.html)
- **Repository**: [GitHub](https://github.com/hieubap/bdspro-golang-microservices)
- **Issues**: [GitHub Issues](https://github.com/hieubap/bdspro-golang-microservices/issues)

---

**Happy Coding! 🚀** 