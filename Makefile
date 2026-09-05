SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

include tools/development/toolchain.versions

# Root Make API always executes Go commands with one explicit distribution.
# This prevents a shell-level GOROOT from silently pairing an older compiler
# tree with a newer go binary. Override QHPRO_TOOLCHAIN_ROOT only to relocate
# the same pinned distribution, never to select an implicit PATH toolchain.
QHPRO_GO_VERSION ?= $(GO_VERSION)
QHPRO_TOOLCHAIN_ROOT ?= $(HOME)/.cache/qhpro/toolchain
QHPRO_GOROOT := $(QHPRO_TOOLCHAIN_ROOT)/go$(QHPRO_GO_VERSION)
QHPRO_GO := $(QHPRO_GOROOT)/bin/go
QHPRO_TOOL_BIN := $(QHPRO_TOOLCHAIN_ROOT)/bin
QHPRO_PROTOC_ROOT := $(QHPRO_TOOLCHAIN_ROOT)/protoc$(PROTOC_VERSION)
QHPRO_DOCKER_CONFIG ?= $(CURDIR)/.tmp/docker-public
# Docker Desktop/WSL on developer machines commonly has a small memory budget.
# Keep image builds sequential by default so the canonical full-stack gate does
# not start several Go compilers at once. CI or larger machines may override it.
QHPRO_COMPOSE_PARALLEL_LIMIT ?= 1
TAIL ?= 200
FOLLOW ?= 1
service ?=
role ?=
# Public commands use the readable deployable name (payment-service).  Keep
# accepting the historical short form through SERVICE while old automation is
# migrated.
SERVICE ?= $(patsubst %-service,%,$(service))
ROLE ?= $(role)
export GOROOT := $(QHPRO_GOROOT)
export GOTOOLCHAIN := local
export GOWORK := off
export PATH := $(QHPRO_TOOL_BIN):$(QHPRO_PROTOC_ROOT)/bin:$(QHPRO_GOROOT)/bin:$(PATH)
# Development Compose chỉ dùng public base/runtime images. Isolate Docker CLI
# khỏi credential helper cá nhân (ví dụ desktop.exe không chạy trong
# WSL) nhưng vẫn cho phép CI/developer override bằng QHPRO_DOCKER_CONFIG.
export DOCKER_CONFIG := $(QHPRO_DOCKER_CONFIG)
export COMPOSE_PARALLEL_LIMIT := $(QHPRO_COMPOSE_PARALLEL_LIMIT)

# Đây là owner duy nhất điều phối backend lõi. Mỗi service vẫn tự sở hữu cách
# build binary của mình trong bin/; Makefile gốc chỉ gọi lại contract đó.
CORE_SERVICES := auth user organization payment tqd notification file gateway hub
DOCKER_SERVICES := $(CORE_SERVICES) assistant
CONFIG_SERVICES := $(CORE_SERVICES) assistant bdspro chat chat-v1 crm map relay search social
REPOSITORY_MODULES := assistant-service bdspro-service chat-service chat-v1-service \
	crm-service map-service relay-service search-service social-service \
	shared/base shared/code shared/common shared/protobuf
MODULES ?= $(REPOSITORY_MODULES)
DB_SERVICES := user organization payment tqd notification file hub
COMPOSE_FILE := compose.yaml
ENV_FILE ?= .env
ENV_FILES := ./.env $(foreach service,$(CORE_SERVICES),./$(service)-service/.env) ./assistant-service/.env
ROOT_REQUIRED_ENV_KEYS := QHPRO_ENVIRONMENT QHPRO_EXECUTION_MODE JWT_KEY_GENERATE \
	QHPRO_INTERNAL_METADATA_SECRET QHPRO_TRUSTED_METADATA_MODE SERVICE_AUTH_KEY \
	QHPRO_ACCEPTANCE_ADMIN_PASSWORD QHPRO_ACCEPTANCE_USER_PASSWORD

.PHONY: help setup up down status logs smoke dev migrate test-e2e reset doctor configure config-check env-check bootstrap generate build test verify-backend verify-config-isolation generate-backend build-backend test-backend \
	vet-backend race-backend diff-check-backend accept-code-backend migrate-backend verify-runtime-backend docker-backend compose-up compose-down \
	compose-config compose-ps compose-logs dev-up dev-status dev-logs dev-smoke dev-down dev-reset \
	dev-dependencies dev-dependencies-down \
	provision-acceptance-admin provision-acceptance-client settle-acceptance-order runtime-value-chain-acceptance runtime-commercial-acceptance runtime-admin-acceptance runtime-client-acceptance clean-backend verify-repository-modules accept-backend

help:
	@printf '%s\n' \
	  'BDSPro Backend' \
	  '' \
	  'make setup                         Chuẩn bị config và toolchain' \
	  'make up                            Khởi động toàn backend' \
	  'make status                        Xem service và health' \
	  'make logs [service=payment-service] Theo dõi log' \
	  'make smoke                         Kiểm tra public API cơ bản' \
	  'make dev service=payment-service  Rebuild/restart một service' \
	  'make test [service=payment-service] Chạy test' \
	  'make test-e2e                      Chạy flow mua gói đến báo cáo' \
	  'make migrate [service=payment-service] Chạy migration' \
	  'make generate                      Sinh protobuf và Wire code' \
	  'make down                          Dừng backend, giữ dữ liệu' \
	  'make reset                         Xóa dữ liệu local có xác nhận' \
	  '' \
	  'Tên service: gateway-service, user-service, organization-service,' \
	  'payment-service, tqd-service, notification-service, file-service.'

# Public developer interface. The longer targets below remain implementation
# details and compatibility aliases; developers do not need to memorize them.
setup: bootstrap

up: dev-up

down: dev-down

status: dev-status

logs: dev-logs

smoke: dev-smoke

reset: dev-reset

test-e2e: runtime-value-chain-acceptance runtime-admin-acceptance

dev:
	@test -n "$(SERVICE)" || { echo 'usage: make dev service=payment-service' >&2; exit 2; }
	@case " $(CORE_SERVICES) assistant " in *" $(SERVICE) "*) ;; *) echo 'unknown service: $(SERVICE)-service' >&2; exit 2 ;; esac
	@targets='$(SERVICE)'; \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --build --no-deps $$targets

configure:
	@bash tools/development/configure.sh

config-check: env-check

doctor:
	@mkdir -p "$(QHPRO_DOCKER_CONFIG)"
	@status=0; \
	test -x "$(QHPRO_GO)" || { echo "Pinned Go missing: $(QHPRO_GO)" >&2; status=1; }; \
	for command in go buf wire migrate rg air grpcurl; do \
	  if command -v "$$command" >/dev/null 2>&1; then \
	    printf '%-18s READY  %s\n' "$$command" "$$(command -v "$$command")"; \
	  else \
	    printf '%-18s MISSING\n' "$$command"; status=1; \
	  fi; \
	done; \
	if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then \
	  printf '%-18s READY\n' docker; \
	else \
	  printf '%-18s UNAVAILABLE (daemon/WSL integration)\n' docker; status=1; \
	fi; \
	if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then printf '%-18s READY\n' 'docker compose'; \
	else printf '%-18s UNAVAILABLE\n' 'docker compose'; status=1; fi; \
	required="$$(awk '/^go / { print $$2; exit }' go.work)"; \
	selected="$$($(QHPRO_GO) env GOVERSION 2>/dev/null || true)"; \
	printf '%-18s %s\n' 'go.work' "$$required"; \
	printf '%-18s %s\n' 'go pinned' '$(QHPRO_GO_VERSION)'; \
	printf '%-18s %s\n' 'GOROOT' '$(QHPRO_GOROOT)'; \
	printf '%-18s %s\n' 'go selected' "$$selected"; \
	test "$$selected" = "go$(QHPRO_GO_VERSION)" || { echo 'Go toolchain không đúng phiên bản pin' >&2; status=1; }; \
	case "$$selected" in "go$${required%.*}."*|"go$$required") ;; *) echo 'Go toolchain không tương thích go.work' >&2; status=1 ;; esac; \
	exit $$status

# Root .env thuộc full Compose; mỗi service/.env thuộc process chạy native của
# service đó. Không dùng .env.local/.env.example hoặc file env không có owner.
env-check:
	@for file in $(ENV_FILES); do test -f "$$file" || { echo "Thiếu $$file" >&2; exit 2; }; done
	@extra="$$(find . -type f \( -name '.env.local' -o -name '.env.example' \) ! -path './.git/*' -print)"; \
	test -z "$$extra" || { echo 'Không dùng .env.local/.env.example:' >&2; echo "$$extra" >&2; exit 2; }
	@for found in $$(find . -type f -name '.env' ! -path './.git/*' -print); do \
	  case ' $(ENV_FILES) ' in *" $$found "*) ;; *) echo "File env không có owner: $$found" >&2; exit 2 ;; esac; \
	done
	@for key in $(ROOT_REQUIRED_ENV_KEYS); do \
	  rg -q "^$${key}=.+$$" "$(ENV_FILE)" || { echo "$(ENV_FILE): thiếu hoặc để trống $$key" >&2; exit 2; }; \
	done
	@for service in $(CORE_SERVICES); do \
	  file="./$$service-service/.env"; \
	  for key in QHPRO_SERVICE_ID; do \
	    rg -q "^$${key}=.+$$" "$$file" || { echo "$$file: thiếu hoặc để trống $$key" >&2; exit 2; }; \
	  done; \
	  case "$$service" in \
	    user) required='DATABASE_URL MIGRATION_URL REDIS_ADDRESS USER_REDIS_DB RPC_USER_ADDRESS RPC_AUTH_ADDRESS RPC_ORGANIZATION_ADDRESS RPC_HUB_ADDRESS RPC_BDSPRO_ADDRESS RPC_NOTIFICATION_ADDRESS RPC_PAYMENT_ADDRESS RPC_CHAT_ADDRESS RPC_CRM_ADDRESS QHPRO_USER_OUTBOUND_MESSAGING_ENABLED' ;; \
	    organization) required='ORGANIZATION_DATABASE_URL MIGRATION_URL ORGANIZATION_GRPC_PORT ORGANIZATION_USER_GRPC_ADDRESS ORGANIZATION_BDSPRO_GRPC_ADDRESS ORGANIZATION_NOTIFICATION_GRPC_ADDRESS ORGANIZATION_PAYMENT_GRPC_ADDRESS ORGANIZATION_CHAT_GRPC_ADDRESS ORGANIZATION_TRANSACTION_GRPC_ADDRESS ORGANIZATION_AUTH_GRPC_ADDRESS' ;; \
	    payment) required='PAYMENT_DATABASE_URL MIGRATION_URL PAYMENT_RABBIT_URL PAYMENT_GRPC_ADDRESS QHPRO_USER_GRPC_ADDRESS QHPRO_AUTH_GRPC_ADDRESS QHPRO_NOTIFICATION_GRPC_ADDRESS SEPAY_API_KEY' ;; \
	    tqd) required='TQD_DATABASE_URL MIGRATION_URL TQD_REDIS_ADDR TQD_REDIS_DB TQD_GRPC_PORT QHPRO_USER_GRPC_ADDR QHPRO_AUTH_GRPC_ADDR QHPRO_ASSISTANT_GRPC_ADDR QHPRO_FILE_HTTP_BASE_URL QHPRO_FILE_PUBLIC_BASE_URL QHPRO_REPORT_GENERATOR_MODE TQD_CLASSIFY_ENABLED' ;; \
	    notification) required='NOTIFICATION_DATABASE_DSN MIGRATION_URL NOTIFICATION_REDIS_ADDR NOTIFICATION_REDIS_DB QHPRO_RABBITMQ_URL NOTIFICATION_GRPC_PORT QHPRO_USER_GRPC_ADDR' ;; \
	    file) required='DATABASE_DSN MIGRATION_URL FILE_SIGNATURE_KEY FILE_XOR_CRYPT_KEY AUTH_GRPC_ADDRESS HUB_GRPC_ADDRESS' ;; \
	    hub) required='HUB_DATABASE_URL MIGRATION_URL HUB_REDIS_ADDRESS HUB_REDIS_DB HUB_USER_GRPC_ADDRESS HUB_BDSPRO_GRPC_ADDRESS HUB_NOTIFICATION_GRPC_ADDRESS' ;; \
	    gateway) required='' ;; \
	    auth) required='RPC_USER_ADDRESS' ;; \
	  esac; \
	  for key in $$required; do rg -q "^$${key}=.+$$" "$$file" || { echo "$$file: thiếu hoặc để trống $$key" >&2; exit 2; }; done; \
	  for root_owned in QHPRO_ENVIRONMENT QHPRO_EXECUTION_MODE ENV_RUNTIME JWT_KEY_GENERATE QHPRO_INTERNAL_METADATA_SECRET QHPRO_TRUSTED_METADATA_MODE QHPRO_INTERNAL_METADATA_MAX_AGE QHPRO_INTERNAL_METADATA_CLOCK_SKEW QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS SERVICE_AUTH_KEY; do \
	    ! rg -q "^$${root_owned}=" "$$file" || { echo "$$file: $$root_owned thuộc root .env" >&2; exit 2; }; \
	  done; \
	  done
	@for key in QHPRO_SERVICE_ID QHPRO_AI_PROVIDER_MODE DEEPSEEK_API_KEY DEEPSEEK_BASE_URL DEEPSEEK_MODEL DEEPSEEK_TIMEOUT_SECONDS DEEPSEEK_MAX_TOKENS OPENAI_API_KEY OPENAI_BASE_URL OPENAI_MODEL OPENAI_TIMEOUT_SECONDS OPENAI_MAX_TOKENS GEMINI_API_KEY GEMINI_BASE_URL GEMINI_MODEL GEMINI_TIMEOUT_SECONDS GEMINI_MAX_TOKENS; do \
	  rg -q "^$${key}=" ./assistant-service/.env || { echo "assistant-service/.env: thiếu $$key" >&2; exit 2; }; \
	done
	@for file in $(ENV_FILES); do \
	  duplicates="$$(sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' "$$file" | sort | uniq -d)"; \
	  test -z "$$duplicates" || { echo "$$file có key trùng:" >&2; echo "$$duplicates" >&2; exit 2; }; \
	done
	@! rg -n 'bdspro\.(com|vn)|qhpro\.(vn|com)|14\.225\.|103\.145\.|103\.172\.' $(ENV_FILES) || \
	  { echo 'Develop env không được chứa endpoint production' >&2; exit 2; }
	@echo 'root/service .env contract PASS'

bootstrap: configure
	@bash tools/development/toolchain.sh bootstrap
	@$(MAKE) doctor
	@$(MAKE) config-check
	@$(MAKE) generate-backend

# Short commands are the public API; their implementation remains at the
# explicit owner targets below.
generate: generate-backend
build: build-backend
test:
	@if test -n "$(SERVICE)"; then \
	  case " $(CORE_SERVICES) assistant " in *" $(SERVICE) "*) $(MAKE) -C "$(SERVICE)-service" test ;; *) echo 'unknown service: $(SERVICE)-service' >&2; exit 2 ;; esac; \
	else \
	  $(MAKE) test-backend; \
	fi

migrate:
	@if test -n "$(SERVICE)"; then \
	  case " $(DB_SERVICES) " in *" $(SERVICE) "*) $(MAKE) -C "$(SERVICE)-service" migrateup ;; *) echo 'service has no canonical database migration: $(SERVICE)-service' >&2; exit 2 ;; esac; \
	else \
	  $(MAKE) migrate-backend; \
	fi

# Chặn nhầm kết nối production ngay từ lớp cấu hình được commit. Secret/DSN
# thật chỉ được đưa qua root .env development hoặc secret store của môi trường deploy.
verify-config-isolation:
	@status=0; \
	for service in $(CONFIG_SERVICES); do \
	  test -f "$$service-service/config/runtime.yml" || { echo "$$service-service: thiếu config/runtime.yml" >&2; status=1; }; \
	  test -f "$$service-service/Makefile" || { echo "$$service-service: thiếu developer Make API" >&2; status=1; }; \
	  rg -q 'tools/development/(go|service)\.mk' "$$service-service/Makefile" || { echo "$$service-service: Makefile không pin repository toolchain" >&2; status=1; }; \
	  for retired in local develop app config; do \
	    test ! -f "$$service-service/config/$$retired.yml" || { echo "$$service-service: còn config profile $$retired.yml" >&2; status=1; }; \
	  done; \
	done; \
	for service in $(DOCKER_SERVICES); do \
	  test -f "$$service-service/Dockerfile" || { echo "$$service-service: thiếu Dockerfile" >&2; status=1; continue; }; \
	  rg -q '^FROM golang:$(GO_VERSION)-alpine AS build$$' "$$service-service/Dockerfile" || { echo "$$service-service: Docker build Go không khớp toolchain $(GO_VERSION)" >&2; status=1; }; \
	done; \
	for file in $$(find . -type f -path '*/config/runtime.yml'); do \
	  if rg -n 'bdspro\.(com|vn)|qhpro\.(vn|com)|14\.225\.|103\.145\.|103\.172\.' "$$file" >/dev/null; then \
	    echo "$$file: runtime config chứa endpoint production" >&2; status=1; \
	  fi; \
	done; \
	legacy_importers="$$(rg -l 'ENV_RUNTIME' . --glob '*.go' --glob '!shared/common/configloader/runtime.go' --glob '!shared/common/configloader/runtime_test.go' || true)"; \
	test -z "$$legacy_importers" || { echo 'Production code còn đọc ENV_RUNTIME:' >&2; echo "$$legacy_importers" >&2; status=1; }; \
	test $$status -eq 0 || exit $$status; \
	echo 'runtime config isolation verification PASS'

# Guard có giá trị vận hành được đặt ngay tại owner Makefile, không phụ thuộc
# vào một cây scripts trung gian khó truy vết.
verify-backend: env-check verify-config-isolation
	@status=0; \
	rg -q '^      RPC_USER_ADDRESS: user:8201$$' "$(COMPOSE_FILE)" || \
	  { echo 'Compose Auth phải gọi User qua service DNS user:8201' >&2; status=1; }; \
	for service in $(CORE_SERVICES); do \
	  for retired in scripts experiments; do \
	    test ! -d "$$service-service/$$retired" || { echo "$$service-service: còn thư mục $$retired/ không có runtime owner" >&2; status=1; }; \
	  done; \
	  test -f "$$service-service/config/runtime.yml" || { echo "$$service-service: thiếu config/runtime.yml" >&2; status=1; }; \
	  for retired in local develop; do \
	    test ! -f "$$service-service/config/$$retired.yml" || { echo "$$service-service: còn config profile $$retired.yml" >&2; status=1; }; \
	  done; \
	  for legacy in config.yaml docker_config.yaml; do \
	    test ! -f "$$service-service/config/$$legacy" || { echo "$$service-service: còn config legacy $$legacy" >&2; status=1; }; \
	  done; \
	  rg -q '^tmp_dir = "\.tmp/air"$$' "$$service-service/.air.toml" || { echo "$$service-service: Air phải ghi vào .tmp/air" >&2; status=1; }; \
	done; \
	for service in auth user organization payment tqd notification hub; do \
	  rg -q 'netcat-openbsd' "$$service-service/Dockerfile" || { echo "$$service-service: healthcheck dùng nc nhưng runtime image chưa cài netcat-openbsd" >&2; status=1; }; \
	done; \
	for service in file gateway; do \
	  rg -q 'wget' "$$service-service/Dockerfile" || { echo "$$service-service: healthcheck dùng wget nhưng runtime image chưa cài wget" >&2; status=1; }; \
	done; \
	for service in $(DB_SERVICES); do \
	  migration_dir="$$service-service/migrate"; \
	  test -d "$$migration_dir" || { echo "$$service-service: thiếu migrate/" >&2; status=1; continue; }; \
	  duplicate="$$(find "$$service-service" -type d \( -name migrations -o -name migrate_canonical \) -print -quit)"; \
	  test -z "$$duplicate" || { echo "$$service-service: migration owner trùng $$duplicate" >&2; status=1; }; \
	  for up in "$$migration_dir"/*.up.sql; do \
	    test -e "$$up" || { echo "$$service-service: migrate/ không có file .up.sql" >&2; status=1; break; }; \
	    down="$${up%.up.sql}.down.sql"; \
	    test -f "$$down" || { echo "$$service-service: thiếu $${down#$$service-service/}" >&2; status=1; }; \
	  done; \
	  for down in "$$migration_dir"/*.down.sql; do \
	    test -e "$$down" || continue; \
	    up="$${down%.down.sql}.up.sql"; \
	    test -f "$$up" || { echo "$$service-service: thiếu $${up#$$service-service/}" >&2; status=1; }; \
	  done; \
	  if rg -n '\.AutoMigrate[[:space:]]*\(' "$$service-service/cmd" "$$service-service/initial" "$$service-service/infra/handler" --glob '*.go' --glob '!**/*_test.go' 2>/dev/null; then \
	    echo "$$service-service: serving code không được gọi AutoMigrate trực tiếp" >&2; status=1; \
	  fi; \
	done; \
	count="$$(rg -n 'RegisterAuthInternalServiceServer[[:space:]]*\(' . --glob '*.go' --glob '!shared/protobuf/**' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "AuthInternal phải có đúng một serving owner, hiện có $$count" >&2; status=1; }; \
	if rg -n 'RegisterAuthInternalServiceHandler' gateway-service --glob '*.go' >/dev/null; then \
	  echo 'AuthInternal là service-to-service contract, không được expose qua HTTP Gateway' >&2; status=1; \
	fi; \
	count="$$(rg -n 'RegisterInternalOrganizationServiceServer[[:space:]]*\(' organization-service --glob '*.go' --glob '!**/*_test.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "InternalOrganization phải do Organization Service sở hữu đúng một lần, hiện có $$count" >&2; status=1; }; \
	if rg -n 'RegisterInternalOrganizationServiceHandler' gateway-service --glob '*.go' >/dev/null; then \
	  echo 'InternalOrganization là service-to-service contract, không được expose qua HTTP Gateway' >&2; status=1; \
	fi; \
	rg -q 'AuthorizeOrganizationAction' user-service/infra/handler/grpc/subscription_checkout.go || \
	  { echo 'Organization checkout phải hỏi permission owner tại Organization Service' >&2; status=1; }; \
	if rg -n 'actor\.Role' user-service/infra/handler/grpc/subscription_checkout.go >/dev/null; then \
	  echo 'Checkout không được dùng global Actor.Role làm organization permission' >&2; status=1; \
	fi; \
	rg -q 'ServiceCallerFromContext' organization-service/infrastructure/handler/internal_organization_handler.go || \
	  { echo 'Organization action authorization phải xác thực immediate service caller' >&2; status=1; }; \
	rg -q 'MapWorkspaceService_CreateGeneratedReport_FullMethodName' tqd-service/cmd/operation_target.go || \
	  { echo 'Generated Report charge point phải thuộc MapWorkspaceService.CreateGeneratedReport' >&2; status=1; }; \
	rg -q 'UseCreateReportAdapter' tqd-service/cmd/report_target_start.go || \
	  { echo 'TQD serving process chưa cut over CreateGeneratedReport sang canonical report target' >&2; status=1; }; \
	rg -q 'QHPRO_REPORT_GENERATOR_MODE: production' "$(COMPOSE_FILE)" || \
	  { echo 'TQD container phải chạy production report renderer' >&2; status=1; }; \
	rg -q 'py3-weasyprint' tqd-service/Dockerfile || \
	  { echo 'TQD runtime image thiếu production PDF engine' >&2; status=1; }; \
	rg -q 'font-dejavu' tqd-service/Dockerfile || \
	  { echo 'TQD production PDF image thiếu font Unicode/Vietnamese' >&2; status=1; }; \
	rg -q 'cfg\.HTTPEndpoint' gateway-service/internal/v1proxy/proxy.go || \
	  { echo 'Gateway v1 HTTP proxy không được dùng gRPC endpoint' >&2; status=1; }; \
	if rg -n 'OpenFile\(|app/logs|logs/app\.log' gateway-service/cmd --glob '*.go' >/dev/null; then \
	  echo 'Gateway phải log stdout/stderr; không sở hữu hidden log file' >&2; status=1; \
	fi; \
	rg -q 'tx\.Table\("notification"\)' notification-service/infra/postgres/eventing/payment_completed_store.go || \
	  { echo 'PaymentCompleted phải tạo customer-visible Notification projection trong Inbox transaction' >&2; status=1; }; \
	rg -q 'outboxPublisher\.Run\(actorCtx\)' payment-service/cmd/grpc/runtime.go || \
	  { echo 'Payment process phải sở hữu durable outbox publisher' >&2; status=1; }; \
	if rg -n '^  payment-publisher:' "$(COMPOSE_FILE)" >/dev/null; then \
	  echo 'Payment outbox là component, không phải standing container riêng' >&2; status=1; \
	fi; \
	rg -q 'paymentConsumer\.Run\(ctx\)' notification-service/cmd/grpc_server.go || \
	  { echo 'Notification process phải sở hữu payment event consumer' >&2; status=1; }; \
	if rg -n '^  notification-(payment-worker|delivery-worker):' "$(COMPOSE_FILE)" >/dev/null; then \
	  echo 'Notification consumer/worker là component, không phải standing container riêng' >&2; status=1; \
	fi; \
	if rg -n 'ReportService_CreateReportAsync_FullMethodName' tqd-service/cmd/operation_target.go tqd-service/cmd/report_target*.go >/dev/null; then \
	  echo 'Legacy ReportService.CreateReportAsync không được trở thành commercial charge point' >&2; status=1; \
	fi; \
	for registration in RegisterAuthServiceServer RegisterOAuthServiceServer RegisterPermissionServiceServer RegisterRoleServiceServer RegisterRoleGroupServiceServer RegisterUserInfoServiceServer; do \
	  count="$$(rg -n "$$registration[[:space:]]*\(" user-service --glob '*.go' --glob '!**/*_test.go' | wc -l | tr -d ' ')"; \
	  test "$$count" = 1 || { echo "$$registration phải do User Service sở hữu đúng một lần, hiện có $$count" >&2; status=1; }; \
	done; \
	test $$status -eq 0 || exit $$status; \
	echo 'backend owner verification PASS'

generate-backend:
	@$(MAKE) -C shared/code buf-shared
	@$(MAKE) -C shared/code buf-crm
	@$(MAKE) -C shared/code buf-bdspro
	@$(MAKE) -C shared/code buf-hub
	@$(MAKE) -C shared/code buf-chat
	@$(MAKE) -C shared/code buf-chat-v2
	@$(MAKE) -C shared/code buf-organization
	@$(MAKE) -C shared/code buf-social
	@$(MAKE) -C shared/code buf-assistant
	@$(MAKE) -C shared/code buf-user-service
	@$(MAKE) -C shared/code buf-payment
	@$(MAKE) -C shared/code buf-tqd
	@$(MAKE) -C shared/code buf-notification
	@$(MAKE) -C shared/code buf-file
	@for service in user organization tqd notification; do $(MAKE) -C shared/code wire $$service || exit; done

build-backend:
	@for service in $(CORE_SERVICES); do $(MAKE) -C "$$service-service" build || exit; done

test-backend:
	@for service in $(CORE_SERVICES); do \
	  if [[ "$$service" == payment ]]; then $(MAKE) -C payment-service GO_TEST_FLAGS='-p=1' test || exit; \
	  else $(MAKE) -C "$$service-service" test || exit; fi; \
	done

vet-backend:
	@for service in $(CORE_SERVICES); do $(MAKE) -C "$$service-service" vet || exit; done

race-backend:
	@for service in $(CORE_SERVICES); do $(MAKE) -C "$$service-service" test-race || exit; done

diff-check-backend:
	@git diff --check

# Core stack có Makefile/runtime contract riêng. Gate này bảo đảm phần source
# còn lại và Shared không trở thành compile island chỉ vì chưa nằm trong
# topology Compose commercial mặc định. integration-test là runtime suite và
# được gọi tường minh sau khi full-flow acceptance đã tạo durable fixture.
verify-repository-modules:
	@for module in $(MODULES); do \
		printf '\n[%s] go vet\n' "$$module"; \
		(cd "$$module" && $(QHPRO_GO) vet ./...) || exit; \
		printf '[%s] go test\n' "$$module"; \
		(cd "$$module" && $(QHPRO_GO) test -count=1 ./...) || exit; \
		printf '[%s] go test -race\n' "$$module"; \
		(cd "$$module" && $(QHPRO_GO) test -race -count=1 ./...) || exit; \
		printf '[%s] go build\n' "$$module"; \
		(cd "$$module" && $(QHPRO_GO) build ./...) || exit; \
	done

# Gate này không cần database/broker. Nó cố ý có tên khác
# accept-backend để không thể nhầm compile/test PASS với runtime ACCEPTED.
accept-code-backend: generate-backend verify-backend vet-backend test-backend race-backend build-backend verify-repository-modules diff-check-backend

migrate-backend:
	@for service in $(DB_SERVICES); do $(MAKE) -C "$$service-service" migrateup || exit; done

docker-backend: generate-backend
	@for service in $(DOCKER_SERVICES); do $(MAKE) -C "$$service-service" docker-build || exit; done

# Compose root là owner duy nhất của development stack. Root .env chứa
# secret/wiring theo environment; topology không nhạy cảm nằm trong compose.yaml.
compose-config: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" config --quiet

# Healthcheck + depends_on là startup readiness owner của Compose. Không thêm
# wait-for.sh/start.sh vào từng image vì sẽ tạo thêm một cơ chế chờ trùng lặp.
compose-up: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --build --wait --wait-timeout 300

# Compose sở hữu fresh migrations qua các one-shot *-migrate service và
# dependency health graph. Target chỉ PASS sau khi tất cả process runtime
# trong topology mặc định healthy; không dùng database native ngoài stack.
verify-runtime-backend: doctor compose-config compose-up compose-ps dev-smoke

compose-down:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" down

compose-ps:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" ps

compose-logs:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" logs -f --tail=200

# Stable developer vocabulary. The compose-* targets remain internal aliases
# so existing automation can migrate without creating a second implementation.
dev-up: compose-up

dev-status: compose-ps

dev-dependencies: env-check
	@case "$(SERVICE)" in \
	  user) targets='postgres redis user-migrate' ;; \
	  organization) targets='postgres organization-migrate' ;; \
	  payment) targets='postgres rabbitmq payment-migrate' ;; \
	  tqd) targets='postgres redis tqd-migrate' ;; \
	  notification) targets='postgres redis rabbitmq notification-migrate' ;; \
	  file) targets='postgres file-migrate' ;; \
	  hub) targets='postgres redis hub-migrate' ;; \
	  *) echo 'SERVICE must be one of: user organization payment tqd notification file hub' >&2; exit 2 ;; \
	esac; \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --wait $$targets

dev-dependencies-down:
	@case "$(SERVICE)" in \
	  user) targets='postgres redis' ;; \
	  organization) targets='postgres' ;; \
	  payment) targets='postgres rabbitmq' ;; \
	  tqd) targets='postgres redis' ;; \
	  notification) targets='postgres redis rabbitmq' ;; \
	  file) targets='postgres' ;; \
	  hub) targets='postgres redis' ;; \
	  *) echo 'SERVICE must be one of: user organization payment tqd notification file hub' >&2; exit 2 ;; \
	esac; \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" stop $$targets

dev-logs:
	@target="$(SERVICE)"; \
	set -- logs --tail="$(TAIL)"; \
	if test -n "$(SINCE)"; then set -- "$$@" --since "$(SINCE)"; fi; \
	if test "$(FOLLOW)" != 0; then set -- "$$@" -f; fi; \
	if test -n "$$target"; then set -- "$$@" "$$target"; fi; \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" "$$@"

dev-smoke: env-check
	@set -a; source "$(ENV_FILE)"; set +a; bash tools/development/smoke.sh

dev-down: compose-down

dev-reset:
	@confirmation="$(CONFIRM)"; \
	if test "$$confirmation" != reset-development-data; then \
	  printf '%s\n' \
	    'This removes BDSPro development databases, queues and file volumes.' \
	    'Source code and the pinned toolchain are preserved.'; \
	  if test -t 0; then \
	    printf 'Type reset-development-data to continue: '; \
	    read -r confirmation; \
	  fi; \
	fi; \
	test "$$confirmation" = reset-development-data || { echo 'Reset cancelled.'; exit 2; }; \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" down --volumes --remove-orphans

provision-acceptance-admin: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-admin-test-data
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" restart auth
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --no-deps --no-build --wait --wait-timeout 120 auth
	@echo 'Acceptance Admin: qhpro.acceptance.admin (password: QHPRO_ACCEPTANCE_ADMIN_PASSWORD)'

provision-acceptance-client: env-check
	@test -n "$(QHPRO_ACCEPTANCE_USERNAME)" || (echo 'QHPRO_ACCEPTANCE_USERNAME is required' >&2; exit 2)
	@test -n "$(QHPRO_ACCEPTANCE_PHONE)" || (echo 'QHPRO_ACCEPTANCE_PHONE is required' >&2; exit 2)
	@test -n "$(QHPRO_ACCEPTANCE_EMAIL)" || (echo 'QHPRO_ACCEPTANCE_EMAIL is required' >&2; exit 2)
	QHPRO_ACCEPTANCE_USERNAME="$(QHPRO_ACCEPTANCE_USERNAME)" \
	QHPRO_ACCEPTANCE_PHONE="$(QHPRO_ACCEPTANCE_PHONE)" \
	QHPRO_ACCEPTANCE_EMAIL="$(QHPRO_ACCEPTANCE_EMAIL)" \
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-client-test-data
	@echo 'Acceptance User: qhpro.acceptance.user (password: QHPRO_ACCEPTANCE_USER_PASSWORD)'

# Mỗi lần chạy tạo duy nhất identity acceptance, nhưng toàn bộ business state
# phải đi qua public HTTP contracts. Fixture chỉ sở hữu credential đăng nhập;
# tuyệt đối không ghi order/subscription/quota/report/notification.
runtime-value-chain-acceptance: env-check
	@set -euo pipefail; set -a; source "$(ENV_FILE)"; source payment-service/.env; set +a; \
	  run_id="$$(date -u +%s)$$(date -u +%N)"; run_id="$${run_id:0:19}"; \
	  phone_suffix="$${run_id: -7}"; \
	  export QHPRO_ACCEPTANCE_USERNAME="qhpro_acceptance_$${run_id}"; \
	  export QHPRO_ACCEPTANCE_PHONE="039$${phone_suffix}"; \
	  export QHPRO_ACCEPTANCE_EMAIL="qhpro-$${run_id}@e.invalid"; \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-client-test-data; \
	  cd integration-test; \
	  QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test -run '^TestQHPROValueChainRuntime$$' -count=1 -v

settle-acceptance-order: env-check
	@test -n "$(order_ref)" || (echo 'order_ref is required' >&2; exit 2)
	@test -n "$(amount_minor)" || (echo 'amount_minor is required' >&2; exit 2)
	@set -a; source "$(ENV_FILE)"; source payment-service/.env; set +a; \
	  cd integration-test && "$(QHPRO_GO)" run ./cmd/settle-acceptance-order \
	    -reference "$(order_ref)" -amount-minor "$(amount_minor)"

# Không gắn acceptance này vào verify-runtime-backend: fresh migrations chỉ seed catalog và
# permission, không được giả lập giao dịch/subscription/usage của profile 1.
# Target PASS khi fixture được tạo qua full commercial flow vẫn đọc được qua
# Gateway và artifact cuối cùng thực sự là PDF do File Service phục vụ.
runtime-commercial-acceptance: env-check
	@set -a; source "$(ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialRuntime$$' -count=1 -v

runtime-admin-acceptance: provision-acceptance-admin
	@set -a; source "$(ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialAdminRuntime$$' -count=1 -v

runtime-client-acceptance: provision-acceptance-client
	@set -a; source "$(ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialClientLoginRuntime$$' -count=1 -v

clean-backend:
	@for service in $(CORE_SERVICES); do $(MAKE) -C "$$service-service" clean || exit; done

# CLOSED chỉ có nghĩa khi code gate, runtime topology, customer value chain và
# operator IAM/projection cùng PASS trong một lần chạy.
accept-backend: accept-code-backend verify-runtime-backend runtime-value-chain-acceptance runtime-admin-acceptance
