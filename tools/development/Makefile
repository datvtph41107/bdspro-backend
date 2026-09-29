ROOT_DIR := ../..
CD_ROOT := cd $(ROOT_DIR)
SCHEMA := cd ../../proto/schema

# Các target generation là command, không phải file artifact; luôn thực thi
# khi developer gọi lại để đồng bộ source contract.
.PHONY: buf-all buf-shared buf-crm buf-bdspro buf-hub buf-chat buf-chat-v2 buf-organization \
	buf-social buf-assistant buf-user-service buf-auth buf-user buf-payment buf-tqd \
	buf-operation buf-notification buf-transaction buf-file wire proto

AIR_EXE ?= air
SERVICE := $(word 2, $(MAKECMDGOALS))
SOURCE_PATH := $(if $(strip $(SERVICE)),$(shell ./source-path.sh $(SERVICE)),)
ARGS    := $(wordlist 3, 99, $(MAKECMDGOALS))
PRIMARY_GOAL := $(firstword $(MAKECMDGOALS))
PORT = $(shell yq e ".service.port.$(SERVICE)" ../../api/gateway/config/local.yml)

# Debug info
debug-info:
	@echo "SERVICE = $(SERVICE)"
	@echo "PORT = $(PORT)"
	@echo "AIR_EXE = $(AIR_EXE)"

run-dev:
	@export PATH=$$PATH:$$HOME/go/bin && \
	lsof -ti tcp:$(PORT) | xargs kill -9 2>/dev/null || true
	@export PATH=$$PATH:$$HOME/go/bin && \
	cd $(ROOT_DIR)/$(SOURCE_PATH) && $(AIR_EXE) $(ARGS)

%:
	@:	

map:
	cd ../../$(shell ./source-path.sh $(word 2, $(MAKECMDGOALS))) && go run ../../tools/development/script/gen_mapper/v2.go $(word 2, $(MAKECMDGOALS))

dev:
	go run monitor/run-dev.go
swag-file:
	$(CD_ROOT) && cd internal/file && swag init --parseDependency --output docs/file
swag:
	bash ./script.sh swag $(word 2, $(MAKECMDGOALS))
wire:
	@test -n "$(SERVICE)" || (echo "Usage: make wire <service>" >&2; exit 2)
	bash ./script.sh wire $(SERVICE)
proto:
	@test -n "$(SERVICE)" || (echo "Usage: make proto <service>" >&2; exit 2)
	@$(MAKE) buf-$(SERVICE)
run:
	$(CD_ROOT) && cd $(shell ./source-path.sh $(word 2, $(MAKECMDGOALS))) && go run main.go $(word 3, $(MAKECMDGOALS))

merge-swagger:
	$(CD_ROOT) && cd proto && go run scripts/merge_swagger.go

buf-crm:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path crm
buf-group:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path group
buf-document:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path document
buf-social:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path social
buf-notification:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path notification
buf-user:
	@find ../../proto/types/user -type f \
		\( -name '*.pb.go' -o -name '*.pb.gw.go' \) \
		-delete 2>/dev/null || true
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path user
buf-bdspro:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path bdspro
buf-shared:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path shared
# `buf-organization` là tên canonical; giữ `buf-org` chỉ như alias tương thích
# để không tạo thêm một implementation thứ hai.
buf-org: buf-organization
buf-organization:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path organization
buf-appointment:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path appointment
buf-payment:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path payment
# Chat hiện có hai thế hệ schema cùng package protobuf. `chat/v2` vẫn được
# giữ để tương thích source cũ nhưng không được đưa vào output `pb/types/chat`
# (nếu đưa cả cây vào sẽ trùng ChatService/CreateConversationRequest). Gateway
# hiện tại sở hữu contract `chat/*.proto`, vì vậy chỉ generate đúng tập file đó.
buf-chat:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml \
		--path chat/chat.proto \
		--path chat/background.proto \
		--path chat/message.proto
buf-chat-v1:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path chat_v1
buf-task:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path task
buf-transaction:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path transaction
buf-auth:
	@find ../../proto/types/auth -type f \
		\( -name '*.pb.go' -o -name '*.pb.gw.go' \) \
		-delete 2>/dev/null || true
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path auth

buf-user-service: buf-auth buf-user
buf-feedback:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path feedback	
buf-file:
	@find ../../proto/types/file -type f \
		\( -name '*.pb.go' -o -name '*.pb.gw.go' \) \
		-delete 2>/dev/null || true
	@if [ -d ../../proto/schema/file ] && \
		find ../../proto/schema/file -type f -name '*.proto' -print -quit | grep -q .; then \
		$(SCHEMA) && buf generate --template ./buf.gen.yaml --path file; \
	else \
		echo "File protobuf schema absent; generated File outputs remain retired"; \
	fi

buf-operation:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path operation
buf-tqd: buf-operation
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path tqd
buf-assistant:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path assistant	
buf-hub:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path hub	
# Chỉ gọi schema directory còn tồn tại. Các target lịch sử như group/document/
# task không còn source nên không được để fresh clone fail vì path ảo. Chat v1
# không gọi chung vì trùng package với chat contract hiện tại.
buf-all:
	$(MAKE) buf-shared
	$(MAKE) buf-crm
	$(MAKE) buf-bdspro
	$(MAKE) buf-hub
	$(MAKE) buf-chat
	$(MAKE) buf-organization
	$(MAKE) buf-social
	$(MAKE) buf-assistant
	$(MAKE) buf-user-service
	$(MAKE) buf-payment
	$(MAKE) buf-tqd
	$(MAKE) buf-notification
	$(MAKE) buf-transaction
	echo "done"
	
init:
	$(CD_ROOT) && cd api/gateway && go mod tidy
	$(CD_ROOT) && cd internal/user && go mod tidy
	$(CD_ROOT) && cd internal/organization-legacy && go mod tidy
	$(CD_ROOT) && cd appointment-service && go mod tidy
	$(CD_ROOT) && cd internal/payment && go mod tidy
	$(CD_ROOT) && cd internal/notification && go mod tidy
	$(CD_ROOT) && cd internal/chat && go mod tidy
	$(CD_ROOT) && cd internal/social && go mod tidy
	$(CD_ROOT) && cd internal/file && go mod tidy
	$(CD_ROOT) && cd membership-service && go mod tidy
	$(CD_ROOT) && cd internal/realtime-relay && go mod tidy
	$(CD_ROOT) && cd internal/property && go mod tidy
	$(CD_ROOT) && cd internal/crm && go mod tidy
	$(CD_ROOT) && cd internal/authentication && go mod tidy
	$(CD_ROOT) && cd transaction-service && go mod tidy
	
gateway-http:
	$(MAKE) -C ../../api/gateway dev

crm:
	@if [ "$(PRIMARY_GOAL)" = "wire" ] || [ "$(PRIMARY_GOAL)" = "swag" ] || [ "$(PRIMARY_GOAL)" = "map" ] || [ "$(PRIMARY_GOAL)" = "proto" ]; then \
		:; \
	else \
		cd $(ROOT_DIR)/internal/crm && air grpc; \
	fi

bdspro-grpc:
	$(CD_ROOT) && cd internal/property && air grpc

user-grpc:
	$(MAKE) -C ../../internal/user dev

organization-grpc:
	$(MAKE) -C ../../internal/organization-legacy dev

notification-grpc:
	$(MAKE) -C ../../internal/notification dev

social-grpc:
	$(CD_ROOT) && cd internal/social && air grpc

appointment-grpc:
	$(CD_ROOT) && cd appointment-service && air grpc

chat-http:
	$(CD_ROOT) && cd internal/chat && air http

chat-grpc:
	$(CD_ROOT) && cd internal/chat && air grpc

relay-grpc:
	$(CD_ROOT) && cd internal/realtime-relay && air

payment-grpc:
	$(MAKE) -C ../../internal/payment dev

membership-grpc:
	$(CD_ROOT) && cd membership-service && air grpc

file-http:
	$(MAKE) -C ../../internal/file dev

transaction-grpc:
	$(CD_ROOT) && cd transaction-service && air grpc

auth-grpc:
	$(MAKE) -C ../../internal/authentication dev

tqd-http:
	@echo "TQD application HTTP is disabled. Start gateway HTTP and TQD gRPC instead." >&2
	@exit 2

tqd-grpc:
	$(MAKE) -C ../../internal/planning dev

assistant-grpc:
	@export PATH=$$PATH:$$HOME/go/bin && \
	lsof -i :8218 | xargs kill -9 2>/dev/null || true
	@export PATH=$$PATH:$$HOME/go/bin && \
	$(CD_ROOT) && cd internal/assistant && go run .

start-all:
	@echo "Starting all services..."
	$(MAKE) gateway-http &
	$(MAKE) crm &
	$(MAKE) bdspro-grpc &
	$(MAKE) user-grpc &
	$(MAKE) organization-grpc &
	$(MAKE) notification-grpc &
	$(MAKE) social-grpc &
	$(MAKE) appointment-grpc &
	$(MAKE) chat-http &
	$(MAKE) chat-grpc &
	$(MAKE) relay-grpc &
	$(MAKE) payment-grpc &
	$(MAKE) membership-grpc &
	$(MAKE) file-http &
	$(MAKE) transaction-grpc &
	$(MAKE) auth-grpc &
	$(MAKE) tqd-grpc &
	@echo "All services started in background. Use 'make stop-all' to stop them."
	@echo "To see running processes: ps aux | grep air"

stop-all:
	$(CD_ROOT)
	@echo "Stopping all services..."
	@pkill -f "air" || true
	@pkill -f "go run" || true
	@echo "All services stopped."

build-all:
	# Giữ target tương thích, nhưng chuyển quyền tạo artifact về Makefile từng service.
	$(MAKE) -C $(ROOT_DIR) build-backend
