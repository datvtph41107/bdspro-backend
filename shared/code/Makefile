ROOT_DIR := ../..
CD_ROOT := cd $(ROOT_DIR)
SCHEMA := cd ../protobuf/schema

# Các target generation là command, không phải file artifact; luôn thực thi
# khi developer gọi lại để đồng bộ source contract.
.PHONY: buf-all buf-shared buf-crm buf-bdspro buf-hub buf-chat buf-chat-v2 buf-organization \
	buf-social buf-assistant buf-user-service buf-auth buf-user buf-payment buf-tqd \
	buf-operation buf-notification buf-transaction buf-file wire proto

AIR_EXE ?= air
SERVICE := $(word 2, $(MAKECMDGOALS))
ARGS    := $(wordlist 3, 99, $(MAKECMDGOALS))
PRIMARY_GOAL := $(firstword $(MAKECMDGOALS))
PORT = $(shell yq e ".service.port.$(SERVICE)" ../../gateway-service/config/local.yml)

# Debug info
debug-info:
	@echo "SERVICE = $(SERVICE)"
	@echo "PORT = $(PORT)"
	@echo "AIR_EXE = $(AIR_EXE)"

run-dev:
	@export PATH=$$PATH:$$HOME/go/bin && \
	lsof -ti tcp:$(PORT) | xargs kill -9 2>/dev/null || true
	@export PATH=$$PATH:$$HOME/go/bin && \
	cd $(ROOT_DIR)/$(SERVICE)-service && $(AIR_EXE) $(ARGS)

%:
	@:	

map:
	cd ../../$(word 2, $(MAKECMDGOALS))-service && go run ../shared/code/script/gen_mapper/v2.go $(word 2, $(MAKECMDGOALS))

dev:
	go run monitor/run-dev.go
swag-file:
	$(CD_ROOT) && cd file-service && swag init --parseDependency --output docs/file
swag:
	./script.sh swag $(word 2, $(MAKECMDGOALS))
wire:
	@test -n "$(SERVICE)" || (echo "Usage: make wire <service>" >&2; exit 2)
	./script.sh wire $(SERVICE)
proto:
	@test -n "$(SERVICE)" || (echo "Usage: make proto <service>" >&2; exit 2)
	@$(MAKE) buf-$(SERVICE)
run:
	$(CD_ROOT) && cd $(word 2, $(MAKECMDGOALS))-service && go run main.go $(word 3, $(MAKECMDGOALS))

merge-swagger:
	$(CD_ROOT) && cd shared/protobuf && go run scripts/merge_swagger.go

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
	@find ../protobuf/types/user -type f \
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
buf-marketing:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path marketing	
buf-transaction:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path transaction
buf-auth:
	@find ../protobuf/types/auth -type f \
		\( -name '*.pb.go' -o -name '*.pb.gw.go' \) \
		-delete 2>/dev/null || true
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path auth

buf-user-service: buf-auth buf-user
buf-feedback:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path feedback	
buf-file:
	@find ../protobuf/types/file -type f \
		\( -name '*.pb.go' -o -name '*.pb.gw.go' \) \
		-delete 2>/dev/null || true
	@if [ -d ../protobuf/schema/file ] && \
		find ../protobuf/schema/file -type f -name '*.proto' -print -quit | grep -q .; then \
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
	$(CD_ROOT) && cd gateway-service && go mod tidy
	$(CD_ROOT) && cd user-service && go mod tidy
	$(CD_ROOT) && cd organization-service && go mod tidy
	$(CD_ROOT) && cd appointment-service && go mod tidy
	$(CD_ROOT) && cd payment-service && go mod tidy
	$(CD_ROOT) && cd notification-service && go mod tidy
	$(CD_ROOT) && cd chat-service && go mod tidy
	$(CD_ROOT) && cd social-service && go mod tidy
	$(CD_ROOT) && cd file-service && go mod tidy
	$(CD_ROOT) && cd membership-service && go mod tidy
	$(CD_ROOT) && cd relay-service && go mod tidy
	$(CD_ROOT) && cd bdspro-service && go mod tidy
	$(CD_ROOT) && cd crm-service && go mod tidy
	$(CD_ROOT) && cd marketing-service && go mod tidy
	$(CD_ROOT) && cd auth-service && go mod tidy
	$(CD_ROOT) && cd transaction-service && go mod tidy
	
gateway-http:
	$(MAKE) -C ../../gateway-service dev

crm:
	@if [ "$(PRIMARY_GOAL)" = "wire" ] || [ "$(PRIMARY_GOAL)" = "swag" ] || [ "$(PRIMARY_GOAL)" = "map" ] || [ "$(PRIMARY_GOAL)" = "proto" ]; then \
		:; \
	else \
		cd $(ROOT_DIR)/crm-service && air grpc; \
	fi

bdspro-grpc:
	$(CD_ROOT) && cd bdspro-service && air grpc

user-grpc:
	$(MAKE) -C ../../user-service dev

organization-grpc:
	$(MAKE) -C ../../organization-service dev

notification-grpc:
	$(MAKE) -C ../../notification-service dev

social-grpc:
	$(CD_ROOT) && cd social-service && air grpc

appointment-grpc:
	$(CD_ROOT) && cd appointment-service && air grpc

chat-http:
	$(CD_ROOT) && cd chat-service && air http   

chat-grpc:
	$(CD_ROOT) && cd chat-service && air grpc

relay-grpc:
	$(CD_ROOT) && cd relay-service && air

payment-grpc:
	$(MAKE) -C ../../payment-service dev

membership-grpc:
	$(CD_ROOT) && cd membership-service && air grpc

file-http:
	$(MAKE) -C ../../file-service dev

marketing-grpc:
	$(CD_ROOT) && cd marketing-service && air grpc

transaction-grpc:
	$(CD_ROOT) && cd transaction-service && air grpc

auth-grpc:
	$(MAKE) -C ../../auth-service dev

tqd-http:
	@echo "TQD application HTTP is disabled. Start gateway HTTP and TQD gRPC instead." >&2
	@exit 2

tqd-grpc:
	$(MAKE) -C ../../tqd-service dev

assistant-grpc:
	@export PATH=$$PATH:$$HOME/go/bin && \
	lsof -i :8218 | xargs kill -9 2>/dev/null || true
	@export PATH=$$PATH:$$HOME/go/bin && \
	$(CD_ROOT) && cd assistant-service && go run .

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
	$(MAKE) marketing-grpc &
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
