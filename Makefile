
swag:
	./script.sh swag $(word 2, $(MAKECMDGOALS))
wire:
	./script.sh wire $(word 2, $(MAKECMDGOALS))
proto:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path $(word 2, $(MAKECMDGOALS))
run:
	cd $(word 2, $(MAKECMDGOALS))-service && go run main.go $(word 3, $(MAKECMDGOALS))
gateway:
	cd gateway-service && go run main.go http
crm:
	cd crm-service && go run main.go grpc
org:
	cd organization-service && go run main.go grpc-server

init:
	cd gateway-service && go mod tidy
	cd user-service && go mod tidy
	cd organization-service && go mod tidy
	cd appointment-service && go mod tidy
	cd payment-service && go mod tidy
	cd notification-service && go mod tidy
	cd chat-service && go mod tidy

buf-crm:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path crm
buf-group:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path group
buf-document:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path document
buf-social:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path social
buf-noti:
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
buf-all: buf-social buf-noti buf-user buf-bdspro buf-shared buf-org buf-appointment buf-payment buf-payment buf-crm buf-chat
	echo "done"	
	

bdspro:
	cd bdspro-service && go run main.go grpc
user:
	cd user-service && go run main.go grpc
admin:
	cd admin-service && go run main.go grpc
organization:
	cd organization-service && go run main.go grpc-server
noti:
	cd notification-service && go run main.go grpc
social:
	cd social-service && go run main.go grpc
appointment:
	cd appointment-service && go run main.go grpc-server
user-http:
	cd user-service && go run main.go http
chat:
	cd chat-service && go run main.go grpc
swag-merge:
	cd code/script/merge_swagger && go run runtime.go
chat-http:
	cd chat-service && go run main.go http-serve
relay:
	cd relay-service && go run main.go grpc-server
chat-swagger:
	cd chat-service && go run main.go swagger-serve
payment:
	cd payment-service && go run main.go grpc-server
chat-grpc:
	cd chat-service && go run main.go grpc-serve
membership:
	cd membership-service && go run main.go grpc-server

start-all:
	make buf-all
	make swag-merge
	make chat-http
	make chat-ws
	make chat-swagger
	make payment
	make org
	make user-http
	make chat