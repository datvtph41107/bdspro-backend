ROOT_DIR := ../..
CD_ROOT := cd $(ROOT_DIR)
SCHEMA := cd ../protobuf/schema

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
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path $(word 2, $(MAKECMDGOALS))
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
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path user
buf-bdspro:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path bdspro
buf-shared:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path shared
buf-org:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path organization
buf-organization:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path organization
buf-appointment:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path appointment
buf-payment:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path payment
buf-chat:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path chat
buf-chat-v1:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path chat_v1
buf-task:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path task
buf-marketing:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path marketing	
buf-transaction:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path transaction
buf-auth:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path auth	
buf-feedback:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path feedback	
buf-tqd:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path tqd
buf-assistant:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path assistant	
buf-hub:
	$(SCHEMA) && buf generate --template ./buf.gen.yaml --path hub	
buf-all:
	$(MAKE) buf-crm
	$(MAKE) buf-group
	$(MAKE) buf-document
	$(MAKE) buf-social
	$(MAKE) buf-notification
	$(MAKE) buf-user
	$(MAKE) buf-bdspro
	$(MAKE) buf-shared
	$(MAKE) buf-org
	$(MAKE) buf-organization
	$(MAKE) buf-appointment
	$(MAKE) buf-payment
	$(MAKE) buf-chat
	$(MAKE) buf-task
	$(MAKE) buf-marketing
	$(MAKE) buf-transaction
	$(MAKE) buf-auth
	$(MAKE) buf-feedback
	$(MAKE) buf-tqd
	$(MAKE) buf-assistant
	$(MAKE) buf-hub
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
	$(CD_ROOT) && cd gateway-service && air http

crm:
	@if [ "$(PRIMARY_GOAL)" = "wire" ] || [ "$(PRIMARY_GOAL)" = "swag" ] || [ "$(PRIMARY_GOAL)" = "map" ] || [ "$(PRIMARY_GOAL)" = "proto" ]; then \
		:; \
	else \
		cd $(ROOT_DIR)/crm-service && air grpc; \
	fi

bdspro-grpc:
	$(CD_ROOT) && cd bdspro-service && air grpc

user-http:
	$(CD_ROOT) && cd user-service && air http

user-grpc:
	$(CD_ROOT) && cd user-service && air grpc

organization-grpc:
	$(CD_ROOT) && cd organization-service && air grpc-server

notification-grpc:
	$(CD_ROOT) && cd notification-service && air grpc

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
	$(CD_ROOT) && cd payment-service && air grpc

membership-grpc:
	$(CD_ROOT) && cd membership-service && air grpc

file-http:
	$(CD_ROOT) && cd file-service && air http

marketing-grpc:
	$(CD_ROOT) && cd marketing-service && air grpc

transaction-grpc:
	$(CD_ROOT) && cd transaction-service && air grpc

auth-grpc:
	# .air.toml đã có args_bin=["grpc"] — chỉ cần `air` (không dùng -server=grpc; không truyền thêm grpc)
	$(CD_ROOT) && cd auth-service && air

tqd-http:
	$(CD_ROOT) && cd tqd-service && air -- -server=http

tqd-grpc:
	@export PATH=$$PATH:$$HOME/go/bin && \
	lsof -i :$(PORT) | xargs kill -9 2>/dev/null || true
	@export PATH=$$PATH:$$HOME/go/bin && \
	$(CD_ROOT) && cd tqd-service && $(AIR_EXE) -- -server=grpc

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
	$(MAKE) user-http &
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
	$(MAKE) tqd-http &
	$(MAKE) tqd-grpc &
	@echo "All services started in background. Prefer: ../../scripts/run-support-gis.sh start"
	@echo "Use 'make stop-all' to stop them."
	@echo "To see running processes: ps aux | grep air"
	@echo "To see logs: tail -f /tmp/bdspro-support-gis/*.log"

stop-all:
	$(CD_ROOT)
	@echo "Stopping all services..."
	@pkill -f "air" || true
	@pkill -f "go run" || true
	@echo "All services stopped."

build-all:
	$(CD_ROOT) && cd gateway-service && go build -o gateway-service main.go
	$(CD_ROOT) && cd user-service && go build -o user-service main.go
	$(CD_ROOT) && cd organization-service && go build -o organization-service main.go
	$(CD_ROOT) && cd appointment-service && go build -o appointment-service main.go
	$(CD_ROOT) && cd payment-service && go build -o payment-service main.go
	$(CD_ROOT) && cd notification-service && go build -o notification-service main.go
	$(CD_ROOT) && cd social-service && go build -o social-service main.go
	$(CD_ROOT) && cd file-service && go build -o file-service main.go
	$(CD_ROOT) && cd relay-service && go build -o relay-service main.go
	$(CD_ROOT) && cd bdspro-service && go build -o bdspro-service main.go
	$(CD_ROOT) && cd crm-service && go build -o crm-service main.go
	$(CD_ROOT) && cd marketing-service && go build -o marketing-service main.go
	# $(CD_ROOT) && cd membership-service && go build -o membership-service .
	$(CD_ROOT) && cd auth-service && go build -o auth-service .
	$(CD_ROOT) && cd transaction-service && go build -o transaction-service .
	$(CD_ROOT) && cd feedback-service && go build -o feedback-service .
	$(CD_ROOT) && cd tqd-service && go build -o tqd-service .
	$(CD_ROOT) && cd assistant-service && go build -o assistant-service .
