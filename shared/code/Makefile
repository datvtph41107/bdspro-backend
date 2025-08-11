ROOT := cd ../..

swag-file:
	$(ROOT) && cd file-service && swag init --parseDependency --output docs/file
swag:
	$(ROOT) && ./script.sh swag $(word 2, $(MAKECMDGOALS))
wire:
	$(ROOT) && ./script.sh wire $(word 2, $(MAKECMDGOALS))
proto:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path $(word 2, $(MAKECMDGOALS))
run:
	$(ROOT) && cd $(word 2, $(MAKECMDGOALS))-service && go run main.go $(word 3, $(MAKECMDGOALS))

merge-swagger:
	$(ROOT) && cd protobuf && go run scripts/merge_swagger.go

buf-crm:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path crm
buf-group:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path group
buf-document:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path document
buf-social:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path social
buf-notification:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path notification
buf-user:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path user
buf-bdspro:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path bdspro
buf-shared:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path shared
buf-org:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path organization
buf-organization:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path organization
buf-appointment:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path appointment
buf-payment:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path payment
buf-chat:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path chat
buf-task:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path task
buf-marketing:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path marketing	
buf-transaction:
	$(ROOT) && cd protobuf/schema && buf generate --template ./buf.gen.yaml --path transaction
buf-all: buf-social buf-notification buf-user buf-bdspro buf-shared buf-org buf-appointment buf-payment buf-payment buf-crm buf-chat buf-task buf-marketing
	echo "done"	
	
init:
	$(ROOT) && cd gateway-service && go mod tidy
	$(ROOT) && cd user-service && go mod tidy
	$(ROOT) && cd organization-service && go mod tidy
	$(ROOT) && cd appointment-service && go mod tidy
	$(ROOT) && cd payment-service && go mod tidy
	$(ROOT) && cd notification-service && go mod tidy
	$(ROOT) && cd chat-service && go mod tidy
	$(ROOT) && cd social-service && go mod tidy
	$(ROOT) && cd file-service && go mod tidy
	$(ROOT) && cd membership-service && go mod tidy
	$(ROOT) && cd relay-service && go mod tidy
	$(ROOT) && cd bdspro-service && go mod tidy
	$(ROOT) && cd crm-service && go mod tidy
	$(ROOT) && cd marketing-service && go mod tidy

gateway:
	$(ROOT) && cd gateway-service && go run main.go http
crm-grpc:
	$(ROOT) && cd crm-service && air grpc
bdspro-grpc:
	$(ROOT) && cd bdspro-service && air grpc
user-http:
	$(ROOT) && cd user-service && air http
user-grpc:
	$(ROOT) && cd user-service && air grpc
organization-grpc:
	$(ROOT) && cd organization-service && air grpc-server
notification-grpc:
	$(ROOT) && cd gateway-service && go run main.go http
bdspro:
	$(ROOT) && cd bdspro-service && go run main.go grpc
user:
	$(ROOT) && cd user-service && go run main.go grpc
organization:
	$(ROOT) && cd organization-service && go run main.go grpc-server
notification:
	$(ROOT) && cd notification-service && air grpc
social-grpc:
	$(ROOT) && cd social-service && air grpc
appointment-grpc:
	$(ROOT) && cd appointment-service && air grpc
appointment:
	$(ROOT) && cd appointment-service && go run main.go grpc
chat-http:
	$(ROOT) && cd chat-service && air grpc
chat-grpc:
	$(ROOT) && cd chat-service && air grpc
relay-grpc:
	$(ROOT) && cd relay-service && air
payment-grpc:
	$(ROOT) && cd payment-service && air grpc
membership-grpc:
	$(ROOT) && cd membership-service && air grpc
file-http:
	$(ROOT) && cd file-service && air http
crm:
	$(ROOT) && cd crm-service && go run main.go grpc
marketing-grpc:
	$(ROOT) && cd marketing-service && air grpc
payment:
	$(ROOT) && cd payment-service && go run main.go grpc
transaction:
	$(ROOT) && cd transaction-service && go run main.go grpc

start-all:
	@echo "Starting all services..."
	$(ROOT) &
	$(MAKE) gateway-http &
	$(MAKE) crm-grpc &
	$(MAKE) bdspro-grpc &
	$(MAKE) user-http &
	$(MAKE) user-grpc &
	$(MAKE) organization &
	$(MAKE) notification &
	$(MAKE) social &
	$(MAKE) appointment &
	$(MAKE) chat-http &
	$(MAKE) chat-grpc &
	$(MAKE) relay-grpc &
	$(MAKE) payment-grpc &
	$(MAKE) membership-grpc &
	$(MAKE) file-http &
	@echo "All services started in background. Use 'make stop-all' to stop them."
	@echo "To see running processes: ps aux | grep air"
	@echo "To see logs: tail -f */tmp/*.log"

stop-all:
	$(ROOT)
	@echo "Stopping all services..."
	@pkill -f "air" || true
	@pkill -f "go run" || true
	@echo "All services stopped."

build-all:
	$(ROOT) && cd gateway-service && go build -o gateway-service main.go
	$(ROOT) && cd user-service && go build -o user-service main.go
	$(ROOT) && cd organization-service && go build -o organization-service main.go
	$(ROOT) && cd appointment-service && go build -o appointment-service main.go
	$(ROOT) && cd payment-service && go build -o payment-service main.go
	$(ROOT) && cd notification-service && go build -o notification-service main.go
	$(ROOT) && cd social-service && go build -o social-service main.go
	$(ROOT) && cd file-service && go build -o file-service main.go
	$(ROOT) && cd relay-service && go build -o relay-service main.go
	$(ROOT) && cd bdspro-service && go build -o bdspro-service main.go
	$(ROOT) && cd crm-service && go build -o crm-service main.go
	$(ROOT) && cd marketing-service && go build -o marketing-service main.go
	$(ROOT) && cd membership-service && go build -o membership-service .
