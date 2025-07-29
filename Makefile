swag-file:
	cd file-service && swag init --parseDependency --output docs/file
swag:
	./script.sh swag $(word 2, $(MAKECMDGOALS))
wire:
	./script.sh wire $(word 2, $(MAKECMDGOALS))
proto:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path $(word 2, $(MAKECMDGOALS))
run:
	cd $(word 2, $(MAKECMDGOALS))-service && go run main.go $(word 3, $(MAKECMDGOALS))

merge-swagger:
	cd protobuf && go run scripts/merge_swagger.go

buf-crm:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path crm
buf-group:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path group
buf-document:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path document
buf-social:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path social
buf-notification:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path notification
buf-user:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path user
buf-bdspro:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path bdspro
buf-shared:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path shared
buf-org:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path organization
buf-organization:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path organization
buf-appointment:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path appointment
buf-payment:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path payment
buf-chat:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path chat
buf-task:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path task
buf-marketing:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path marketing	
buf-all: buf-social buf-notification buf-user buf-bdspro buf-shared buf-org buf-appointment buf-payment buf-payment buf-crm buf-chat buf-task buf-marketing
	echo "done"	
	
init:
	cd gateway-service && go mod tidy
	cd user-service && go mod tidy
	cd organization-service && go mod tidy
	cd appointment-service && go mod tidy
	cd payment-service && go mod tidy
	cd notification-service && go mod tidy
	cd chat-service && go mod tidy
	cd social-service && go mod tidy
	cd file-service && go mod tidy
	cd membership-service && go mod tidy
	cd relay-service && go mod tidy
	cd bdspro-service && go mod tidy
	cd crm-service && go mod tidy
	cd marketing-service && go mod tidy

gateway-http:
	cd gateway-service && air http
crm-grpc:
	cd crm-service && air grpc
bdspro-grpc:
	cd bdspro-service && air grpc 
user-http:
	cd user-service && air http
user-grpc:
	cd user-service && air grpc
organization:
	cd organization-service && air grpc
notification:
	cd notification-service && air grpc
social:
	cd social-service && air grpc
appointment:
	cd appointment-service && air grpc
chat-http:
	cd chat-service && air grpc
chat-grpc:
	cd chat-service && air grpc
relay-grpc:
	cd relay-service && air
payment-grpc:
	cd payment-service && air grpc
membership-grpc:
	cd membership-service && air grpc
file-http:
	cd file-service && air http
crm:
	cd crm-service && air grpc
marketing-grpc:
	cd marketing-service && air grpc


start-all:
	@echo "Starting all services..."
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
	@echo "Stopping all services..."
	@pkill -f "air" || true
	@pkill -f "go run" || true
	@echo "All services stopped."

build-all:
	cd gateway-service && go build -o gateway-service main.go
	cd user-service && go build -o user-service main.go
	cd organization-service && go build -o organization-service main.go
	cd appointment-service && go build -o appointment-service main.go
	cd payment-service && go build -o payment-service main.go
	cd notification-service && go build -o notification-service main.go
	cd social-service && go build -o social-service main.go
	cd file-service && go build -o file-service main.go
	cd relay-service && go build -o relay-service main.go
	cd bdspro-service && go build -o bdspro-service main.go
	cd crm-service && go build -o crm-service main.go
	cd marketing-service && go build -o marketing-service main.go
	cd membership-service && go build -o membership-service .
