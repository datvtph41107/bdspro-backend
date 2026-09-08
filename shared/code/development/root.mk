SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

include shared/code/development/toolchain.versions

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
QHPRO_ACCEPTANCE_USERNAME ?= qhpro_acceptance_0000000000
QHPRO_ACCEPTANCE_PHONE ?= 0390000000
QHPRO_ACCEPTANCE_EMAIL ?= qhpro-0000000000@e.invalid
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
NATIVE_SERVICES := assistant notification organization payment file user auth hub tqd gateway
CONFIG_SERVICES := $(CORE_SERVICES) assistant bdspro chat chat-v1 crm map relay search social
REPOSITORY_MODULES := assistant-service bdspro-service chat-service chat-v1-service \
	crm-service map-service relay-service search-service social-service \
	shared/base shared/code shared/common shared/protobuf
MODULES ?= $(REPOSITORY_MODULES)
DB_SERVICES := user organization payment tqd notification file hub
# Source-level migration authorities. bdspro/crm are not standing services in
# the default Compose topology, but their schema history is still repository-owned.
MIGRATION_SERVICES := $(DB_SERVICES) bdspro crm
COMPOSE_FILE := compose.yaml
ENV_FILE ?= .env
ACCEPTANCE_ENV_FILE ?= integration-test/.env
ENV_FILES := ./.env $(foreach service,$(CORE_SERVICES),./$(service)-service/.env) ./assistant-service/.env
OWNED_ENV_FILES := $(ENV_FILES) ./$(ACCEPTANCE_ENV_FILE)
ROOT_REQUIRED_ENV_KEYS := QHPRO_ENVIRONMENT QHPRO_EXECUTION_MODE JWT_KEY_GENERATE \
	QHPRO_INTERNAL_METADATA_SECRET QHPRO_TRUSTED_METADATA_MODE SERVICE_AUTH_KEY

.PHONY: help setup deps deps-down up down status logs smoke dev migrate test-e2e reset verify accept verify-migrations rebuild bootstrap-admin provision-development-identities doctor configure config-check env-check bootstrap generate build test verify-backend verify-config-isolation generate-backend build-backend test-backend \
	vet-backend race-backend diff-check-backend accept-code-backend migrate-backend verify-runtime-backend docker-backend compose-up compose-up-build compose-down \
	compose-config compose-ps compose-logs integration-up integration-down integration-status integration-logs native-build native-up native-down native-status native-logs native-restart dev-up dev-status dev-logs dev-smoke dev-down dev-reset \
	dev-dependencies dev-dependencies-down \
	acceptance-env-check provision-acceptance-admin provision-acceptance-client runtime-notification-broker-recovery runtime-value-chain-acceptance runtime-commercial-acceptance runtime-admin-acceptance runtime-client-acceptance clean-backend verify-repository-modules verify-non-go-source verify-docs accept-backend

help:
	@printf '%s\n' \
	  'BDSPro Backend' \
	  '' \
	  'make setup                         Chuẩn bị config và pinned toolchain' \
	  'make doctor                        Chẩn đoán toolchain + Docker environment' \
	  'make deps                          Khởi động Postgres, Redis và RabbitMQ local' \
	  'make up                            Build và chạy toàn bộ Go service trực tiếp trên host' \
	  'make status                        Xem hạ tầng Docker + process native' \
	  'make logs [service=payment-service] Theo dõi log process native' \
	  'make smoke                         Kiểm tra public boundary cơ bản' \
	  'make bootstrap-admin               Tạo root operator một lần trên database mới' \
	  'Development accounts: admin/admin123 và 0900000000/client123' \
	  'make test-e2e                      Chạy flow mua gói đến báo cáo' \
	  'make integration-up                Full Compose, chỉ dùng khi kiểm tra image/package' \
	  'make rebuild [service=payment-service] Build lại integration image (một/all)' \
	  'make down                          Dừng process native + hạ tầng, giữ dữ liệu' \
	  'make reset                         Xóa dữ liệu local có xác nhận' \
	  '' \
	  'Daily service: cd <service> && make deps-up && make dev.' \
	  'CI/release: make accept (full source, image, runtime và E2E proof).'

# Public developer interface. The longer targets below remain implementation
# details and compatibility aliases; developers do not need to memorize them.
setup: bootstrap

# This is an explicit operator action, not setup/seed behavior. The command is
# transactionally refused after the first active root operator exists.
bootstrap-admin:
	@$(MAKE) -C user-service bootstrap-admin

# Known credentials are development data, not runtime behavior. Keep them out
# of production migrations and refuse execution unless the repository-owned
# environment is explicitly development.
provision-development-identities: env-check
	@set -a; source "$(ENV_FILE)"; set +a; \
	  test "$${QHPRO_ENVIRONMENT}" = development || { \
	    echo 'development identities are forbidden outside QHPRO_ENVIRONMENT=development' >&2; exit 2; \
	  }; \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" exec -T postgres \
	    psql 'postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable' \
	    < shared/code/development/identities.sql
	@echo 'Development Admin: admin / admin123'
	@echo 'Development Client: 0900000000 (hoặc client) / client123'

# Dependencies are intentionally separate from the service process. Daily Go
# development runs the active service natively; Docker is used for backing
# services and, separately, for whole-system integration proof.
deps: env-check
	@if test -n "$(SERVICE)"; then \
	  $(MAKE) dev-dependencies SERVICE="$(SERVICE)"; \
	else \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --wait postgres redis rabbitmq; \
	fi

deps-down: env-check
	@if test -n "$(SERVICE)"; then \
	  $(MAKE) dev-dependencies-down SERVICE="$(SERVICE)"; \
	else \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" stop postgres redis rabbitmq; \
	fi

up: native-up

down: native-down

status: native-status

logs: native-logs

smoke: dev-smoke

reset: dev-reset

verify: accept-code-backend

accept: accept-backend

verify-migrations:
	@bash shared/code/development/verify-migrations.sh

# With service=..., refresh only that integration container. Without it,
# rebuild the complete Compose candidate.
rebuild:
	@if test -n "$(SERVICE)"; then \
	  case " $(DOCKER_SERVICES) " in *" $(SERVICE) "*) ;; *) echo 'unknown Docker service: $(SERVICE)-service' >&2; exit 2 ;; esac; \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --build --no-deps --wait --wait-timeout 120 "$(SERVICE)"; \
	else \
	  $(MAKE) compose-up-build; \
	fi

test-e2e: native-up runtime-value-chain-acceptance runtime-admin-acceptance

# Daily edit loop: no image build. Each service owns its native run arguments;
# the shared service.mk supplies the common env + Air mechanics.
dev:
	@test -n "$(SERVICE)" || { echo 'usage: make dev service=payment-service' >&2; exit 2; }
	@case " $(CORE_SERVICES) " in *" $(SERVICE) "*) ;; *) echo 'native dev supports core service: $(SERVICE)-service' >&2; exit 2 ;; esac
	@bash shared/code/development/native-stack.sh stop-one "$(SERVICE)"
	@$(MAKE) -C "$(SERVICE)-service" dev

configure:
	@bash shared/code/development/configure.sh

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
	for runtime_command in ffmpeg weasyprint; do \
	  if command -v "$$runtime_command" >/dev/null 2>&1; then \
	    printf '%-18s READY  %s\n' "$$runtime_command" "$$(command -v "$$runtime_command")"; \
	  else \
	    printf '%-18s MISSING (native runtime dependency)\n' "$$runtime_command"; status=1; \
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
	  case ' $(OWNED_ENV_FILES) ' in *" $$found "*) ;; *) echo "File env không có owner: $$found" >&2; exit 2 ;; esac; \
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
	@bash shared/code/development/toolchain.sh bootstrap
	@bash shared/code/development/toolchain.sh verify
	@$(MAKE) config-check
	@echo 'Setup ready. Run make doctor for Docker/environment diagnostics, then make deps.'

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
	  case " $(MIGRATION_SERVICES) " in \
	    *" $(SERVICE) "*) \
	      case "$(SERVICE)" in bdspro|crm) bash shared/code/development/ensure-source-databases.sh ;; esac; \
	      $(MAKE) -C "$(SERVICE)-service" migrateup ;; \
	    *) echo 'service has no canonical database migration: $(SERVICE)-service' >&2; exit 2 ;; \
	  esac; \
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
	  rg -q 'shared/code/development/(go|service)\.mk' "$$service-service/Makefile" || { echo "$$service-service: Makefile không pin repository toolchain" >&2; status=1; }; \
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
verify-backend: env-check verify-config-isolation verify-migrations
	@bash shared/code/development/verify-source-layout.sh
	@status=0; \
	rg -q '^up: native-up$$' shared/code/development/root.mk || \
	  { echo 'Public make up phải chạy application process native' >&2; status=1; }; \
	rg -q '^integration-up: compose-up$$' shared/code/development/root.mk || \
	  { echo 'Full Compose phải là explicit image/package topology' >&2; status=1; }; \
	test -f shared/code/development/identities.sql || \
	  { echo 'Thiếu development identity fixture' >&2; status=1; }; \
	rg -Fq 'test "$$$${QHPRO_ENVIRONMENT}" = development' shared/code/development/root.mk || \
	  { echo 'Development identity fixture phải fail closed ngoài environment development' >&2; status=1; }; \
	rg -Fq "crypt('admin123', gen_salt('bf', 10))" shared/code/development/identities.sql || \
	  { echo 'Development Admin phải được lưu bằng bcrypt, không phải plaintext' >&2; status=1; }; \
	rg -Fq "crypt('client123', gen_salt('bf', 10))" shared/code/development/identities.sql || \
	  { echo 'Development Client phải được lưu bằng bcrypt, không phải plaintext' >&2; status=1; }; \
	if rg -n 'admin123|client123' user-service/migrate "$(COMPOSE_FILE)" user-service/config $(ENV_FILES) --glob '*' >/dev/null; then \
	  echo 'Known development credential không được nằm trong migration, Compose hoặc runtime config' >&2; status=1; \
	fi; \
	rg -q 'QHPRO_EXECUTION_MODE=host' shared/code/development/native-stack.sh || \
	  { echo 'Native supervisor phải khóa host runtime selection' >&2; status=1; }; \
	if rg -n 'docker compose.*(up|restart).*\b(auth|user|organization|payment|tqd|notification|file|gateway|hub|assistant)\b' shared/code/development/native-stack.sh >/dev/null; then \
	  echo 'Native supervisor không được khởi chạy application container' >&2; status=1; \
	fi; \
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
	rg -Fq 'NewCheckoutAuthorizer(database)' user-service/cmd/grpc/access_target.go || \
	  { echo 'Organization checkout phải dùng membership authority do User sở hữu' >&2; status=1; }; \
	if rg -n 'NewSubscriptionCheckoutHandler\(checkoutService, organizationClient\)' user-service --glob '*.go' >/dev/null; then \
	  echo 'Commercial checkout không được gọi Organization RPC legacy' >&2; status=1; \
	fi; \
	if rg -n 'actor\.Role' user-service/infra/handler/grpc/subscription_checkout.go >/dev/null; then \
	  echo 'Checkout không được dùng global Actor.Role làm organization permission' >&2; status=1; \
	fi; \
	count="$$(rg -n 'RegisterAdminOrganizationServiceServer[[:space:]]*\(' user-service --glob '*.go' --glob '!**/*_test.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "Admin Organization phải do User đăng ký đúng một lần, hiện có $$count" >&2; status=1; }; \
	count="$$(rg -n 'RegisterAdminOrganizationServiceHandler' gateway-service/internal/grpcgateway --glob '*.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "Gateway phải đăng ký Admin Organization đúng một lần, hiện có $$count" >&2; status=1; }; \
	count="$$(rg -n 'RegisterInternalOrganizationMembershipServiceServer[[:space:]]*\(' user-service --glob '*.go' --glob '!**/*_test.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "Internal Organization membership phải do User đăng ký đúng một lần, hiện có $$count" >&2; status=1; }; \
	rg -q 'MapOrganizationMembersToContacts' crm-service/infra/handler/contact_handler.go || \
	  { echo 'CRM contact organization projection phải đọc từ User Service' >&2; status=1; }; \
	rg -q 'CheckOrganizationMembers' crm-service/infra/handler/friend_handler.go || \
	  { echo 'CRM membership projection phải đọc từ User Service' >&2; status=1; }; \
	if rg -n 'CheckUsersInOrganization|GetOrganizationMemberByIds' crm-service --glob '*.go' >/dev/null; then \
	  echo 'CRM không được tiếp tục đọc organization membership từ Organization Service legacy' >&2; status=1; \
	fi; \
	if rg -n 'GetAdminRoles' user-service --glob '*.go' >/dev/null; then \
	  echo 'User IAM không được đọc catalog vai trò từ Organization Service legacy' >&2; status=1; \
	fi; \
	if rg -n 'ADMIN_VT_(TAO|SUA|XOA|XEM|GAN)' user-service/infra/handler/role_handler.go >/dev/null; then \
	  echo 'Role handler phải dùng quyền IAM_ROLE_* do User migration sở hữu' >&2; status=1; \
	fi; \
	rg -U -q 'post: "/v2/auth/role/permissions"[[:space:]]+body: "\*"' shared/protobuf/schema/auth/role.proto || \
	  { echo 'Role permission HTTP contract phải decode request body' >&2; status=1; }; \
	count="$$(rg -n 'h\.InternalHandler\.HasPermissions' user-service/infra/handler/role_group_handler.go | wc -l | tr -d ' ')"; \
	test "$$count" = 7 || { echo "Mọi RoleGroup operation phải qua IAM authorization, hiện có $$count/7" >&2; status=1; }; \
	if rg -n 'c\.OrganizationClient\.GetOrganizationMember' user-service/infra/client --glob '*.go' >/dev/null; then \
	  echo 'User auth/token không được đọc organization membership qua RPC legacy' >&2; status=1; \
	fi; \
	rg -q 'Table\("organization_members AS member"\)' user-service/infra/client/organization_provider_client.go || \
	  { echo 'User auth/token phải đọc membership từ User-owned store' >&2; status=1; }; \
	rg -q 'OrganizationService\.ListForProfile' user-service/infra/handler/profile_handler.go || \
	  { echo 'User account switcher phải đọc organization directory từ User-owned store' >&2; status=1; }; \
	if rg -n 'GetOrganizationsByUserId' user-service --glob '*.go' >/dev/null; then \
	  echo 'User profile không được đọc danh sách organization qua RPC legacy' >&2; status=1; \
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
	count="$$(rg -n 'RegisterMapPointServiceServer[[:space:]]*\(' tqd-service/cmd --glob '*.go' --glob '!**/*_test.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "MapPoint phải do TQD đăng ký đúng một lần, hiện có $$count" >&2; status=1; }; \
	count="$$(rg -n 'RegisterMapPointServiceHandler' gateway-service/internal/grpcgateway --glob '*.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "Gateway phải đăng ký MapPoint contract đúng một lần, hiện có $$count" >&2; status=1; }; \
	if rg -n '"map"[[:space:]]*:' gateway-service/internal/v1proxy/proxy.go >/dev/null; then \
	  echo 'MapPoint đã thuộc TQD gRPC; Gateway không được reverse-proxy map-service cũ' >&2; status=1; \
	fi; \
	if rg -n '^  map:' "$(COMPOSE_FILE)" >/dev/null; then \
	  echo 'MapPoint đã thuộc TQD; map-service cũ không được là standing container' >&2; status=1; \
	fi; \
	if rg -n 'OpenFile\(|app/logs|logs/app\.log' gateway-service/cmd --glob '*.go' >/dev/null; then \
	  echo 'Gateway phải log stdout/stderr; không sở hữu hidden log file' >&2; status=1; \
	fi; \
	if rg -n '(_db|common_db)\.NewDB\(' crm-service/cmd --glob '*.go' >/dev/null; then \
	  echo 'CRM cmd không được mở DB thứ hai ngoài Wire composition graph' >&2; status=1; \
	fi; \
	test ! -d relay-service/socket || { echo 'Relay còn package socket song song không có runtime caller' >&2; status=1; }; \
	count="$$(rg -n 'redis\.NewClient\(' relay-service --glob '*.go' | wc -l | tr -d ' ')"; \
	test "$$count" = 1 || { echo "Relay phải có đúng một Redis constructor owner, hiện có $$count" >&2; status=1; }; \
	rg -Fq 'NewWebSocketHandler(chatClient, userClient, notificationClient, redisClient)' relay-service/server/ws_server.go || \
	  { echo 'Relay Redis phải được process root inject vào WebSocket handler' >&2; status=1; }; \
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
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
	  git diff --check; \
	else \
	  echo 'git metadata unavailable; diff-only whitespace check skipped'; \
	fi

verify-non-go-source:
	@bash shared/code/development/verify-python.sh

verify-docs:
	@python3 shared/code/development/verify-docs.py

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
accept-code-backend: generate-backend verify-backend verify-non-go-source verify-docs vet-backend test-backend race-backend build-backend verify-repository-modules diff-check-backend

migrate-backend:
	@bash shared/code/development/ensure-source-databases.sh
	@for service in $(MIGRATION_SERVICES); do $(MAKE) -C "$$service-service" migrateup || exit; done

docker-backend: generate-backend
	@for service in $(DOCKER_SERVICES); do $(MAKE) -C "$$service-service" docker-build || exit; done

# Compose remains the image/packaging integration topology. The normal edit and
# runtime loop below runs application processes directly from source; only
# PostgreSQL, Redis and RabbitMQ stay in Docker.
compose-config: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" config --quiet

# Healthcheck + depends_on là startup readiness owner của Compose. Không thêm
# wait-for.sh/start.sh vào từng image vì sẽ tạo thêm một cơ chế chờ trùng lặp.
compose-up: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --wait --wait-timeout 300

compose-up-build: env-check
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" up -d --build --wait --wait-timeout 300

# Runtime acceptance follows the same topology developers debug: backing
# infrastructure in Docker, application binaries on the host.
verify-runtime-backend: doctor compose-config native-up native-status dev-smoke

compose-down:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" down

compose-ps:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" ps

compose-logs:
	docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" logs -f --tail=200

# Explicit packaging/image integration commands. They are deliberately absent
# from the default developer path so `make up` never hides application code in
# containers.
integration-up: compose-up

integration-down: compose-down

integration-status: compose-ps

integration-logs: compose-logs

# Build every native application once, sequentially, before starting the
# process graph. Air remains the focused one-service edit loop (`make dev`).
native-build:
	@for service in $(NATIVE_SERVICES); do $(MAKE) -C "$$service-service" build || exit; done

native-up: doctor env-check
	@$(MAKE) deps
	@docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" stop $(DOCKER_SERVICES) >/dev/null
	@bash shared/code/development/native-stack.sh stop
	@bash shared/code/development/ensure-source-databases.sh
	@for service in $(DB_SERVICES); do $(MAKE) -C "$$service-service" migrateup || exit; done
	@$(MAKE) provision-development-identities
	@$(MAKE) native-build
	@bash shared/code/development/native-stack.sh start

native-down:
	@bash shared/code/development/native-stack.sh stop
	@$(MAKE) deps-down

native-status:
	@docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" ps postgres redis rabbitmq
	@bash shared/code/development/native-stack.sh status

native-logs:
	@TAIL="$(TAIL)" FOLLOW="$(FOLLOW)" bash shared/code/development/native-stack.sh logs "$(SERVICE)"

native-restart:
	@test -n "$(SERVICE)" || { echo 'usage: make native-restart service=<name>-service' >&2; exit 2; }
	@$(MAKE) -C "$(SERVICE)-service" build
	@bash shared/code/development/native-stack.sh restart "$(SERVICE)"

# Stable developer compatibility names now point to the native topology.
dev-up: native-up

dev-status: native-status

dev-dependencies: env-check
	@test -n "$(SERVICE)" || { echo 'usage: make deps service=<name>-service' >&2; exit 2; }
	@echo 'Deprecated root service router; use: make -C $(SERVICE)-service deps-up' >&2
	@$(MAKE) -C "$(SERVICE)-service" deps-up

dev-dependencies-down:
	@test -n "$(SERVICE)" || { echo 'usage: make deps-down service=<name>-service' >&2; exit 2; }
	@echo 'Deprecated root service router; use: make -C $(SERVICE)-service deps-down' >&2
	@$(MAKE) -C "$(SERVICE)-service" deps-down

dev-logs: native-logs

dev-smoke: env-check
	@set -a; source "$(ENV_FILE)"; set +a; bash shared/code/development/smoke.sh

dev-down: native-down

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

acceptance-env-check: env-check
	@test -f "$(ACCEPTANCE_ENV_FILE)" || { echo 'missing integration-test/.env; run make setup' >&2; exit 2; }
	@set -a; source "$(ENV_FILE)"; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  case "$${QHPRO_ENVIRONMENT:-}" in development|test|ci|acceptance) ;; \
	    *) echo 'acceptance fixtures are forbidden outside development/test/ci/acceptance' >&2; exit 2 ;; \
	  esac; \
	  test -n "$${QHPRO_ACCEPTANCE_ADMIN_PASSWORD:-}" || { echo 'QHPRO_ACCEPTANCE_ADMIN_PASSWORD is required by integration tests' >&2; exit 2; }; \
	  test -n "$${QHPRO_ACCEPTANCE_USER_PASSWORD:-}" || { echo 'QHPRO_ACCEPTANCE_USER_PASSWORD is required by integration tests' >&2; exit 2; }

provision-acceptance-admin: acceptance-env-check
	@set -a; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-admin-test-data
	@bash shared/code/development/native-stack.sh restart auth
	@echo 'Acceptance Admin: qhpro.acceptance.admin (password: QHPRO_ACCEPTANCE_ADMIN_PASSWORD)'

provision-acceptance-client: acceptance-env-check
	@test -n "$(QHPRO_ACCEPTANCE_USERNAME)" || (echo 'QHPRO_ACCEPTANCE_USERNAME is required' >&2; exit 2)
	@test -n "$(QHPRO_ACCEPTANCE_PHONE)" || (echo 'QHPRO_ACCEPTANCE_PHONE is required' >&2; exit 2)
	@test -n "$(QHPRO_ACCEPTANCE_EMAIL)" || (echo 'QHPRO_ACCEPTANCE_EMAIL is required' >&2; exit 2)
	@set -a; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  QHPRO_ACCEPTANCE_USERNAME="$(QHPRO_ACCEPTANCE_USERNAME)" \
	  QHPRO_ACCEPTANCE_PHONE="$(QHPRO_ACCEPTANCE_PHONE)" \
	  QHPRO_ACCEPTANCE_EMAIL="$(QHPRO_ACCEPTANCE_EMAIL)" \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-client-test-data
	@echo 'Acceptance User: $(QHPRO_ACCEPTANCE_USERNAME) (password: QHPRO_ACCEPTANCE_USER_PASSWORD)'

# Failure-model gate: RabbitMQ may restart without taking the Notification API
# down. The next value-chain gate proves that the reconnected consumer still
# transfers the durable Payment event into Notification Inbox/effect state.
runtime-notification-broker-recovery: env-check
	@bash shared/code/development/recovery-notification-broker.sh "$(COMPOSE_FILE)" "$(ENV_FILE)"

# Mỗi lần chạy tạo duy nhất identity acceptance, nhưng toàn bộ business state
# phải đi qua public HTTP contracts. Fixture chỉ sở hữu credential đăng nhập;
# tuyệt đối không ghi order/subscription/quota/report/notification.
runtime-value-chain-acceptance: acceptance-env-check
	@set -euo pipefail; set -a; source "$(ENV_FILE)"; source "$(ACCEPTANCE_ENV_FILE)"; source payment-service/.env; set +a; \
	  run_id="$$(date -u +%s)$$(date -u +%N)"; run_id="$${run_id:0:19}"; \
	  phone_suffix="$${run_id: -7}"; \
	  export QHPRO_ACCEPTANCE_USERNAME="qhpro_acceptance_$${run_id}"; \
	  export QHPRO_ACCEPTANCE_PHONE="039$${phone_suffix}"; \
	  export QHPRO_ACCEPTANCE_EMAIL="qhpro-$${run_id}@e.invalid"; \
	  docker compose --env-file "$(ENV_FILE)" -f "$(COMPOSE_FILE)" --profile acceptance run --rm user-client-test-data; \
	  cd integration-test; \
	  QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test -run '^TestQHPROValueChainRuntime$$' -count=1 -v

# Không gắn acceptance này vào verify-runtime-backend: fresh migrations chỉ seed catalog và
# permission, không được giả lập giao dịch/subscription/usage của profile 1.
# Target PASS khi fixture được tạo qua full commercial flow vẫn đọc được qua
# Gateway và artifact cuối cùng thực sự là PDF do File Service phục vụ.
runtime-commercial-acceptance: acceptance-env-check
	@set -a; source "$(ENV_FILE)"; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialRuntime$$' -count=1 -v

runtime-admin-acceptance: provision-acceptance-admin
	@set -a; source "$(ENV_FILE)"; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialAdminRuntime$$' -count=1 -v

runtime-client-acceptance: provision-acceptance-client
	@set -a; source "$(ENV_FILE)"; source "$(ACCEPTANCE_ENV_FILE)"; set +a; \
	  cd integration-test && QHPRO_ACCEPTANCE_USERNAME="$(QHPRO_ACCEPTANCE_USERNAME)" QHPRO_RUNTIME_SMOKE=1 "$(QHPRO_GO)" test \
	    -run '^TestQHPROCommercialClientLoginRuntime$$' -count=1 -v

clean-backend:
	@for service in $(CORE_SERVICES); do $(MAKE) -C "$$service-service" clean || exit; done

# CLOSED chỉ có nghĩa khi code gate, runtime topology, customer value chain và
# operator IAM/projection cùng PASS trong một lần chạy.
accept-backend: accept-code-backend verify-runtime-backend runtime-notification-broker-recovery runtime-value-chain-acceptance runtime-admin-acceptance
