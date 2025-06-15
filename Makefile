
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


init:
	cd gateway-service && go mod tidy
	cd user-service && go mod tidy
	cd organization-service && go mod tidy
	cd appointment-service && go mod tidy

buf-crm:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path crm
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
buf-transaction:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path transaction
buf-group:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path group
buf-payment:
	cd protobuf/schema && buf generate --template ./buf.gen.yaml --path payment

buf-all: buf-social buf-noti buf-user buf-bdspro buf-shared buf-org buf-appointment buf-transaction buf-group buf-payment buf-group buf-payment
	echo "done"	
	

gen:
	cd protobuf/schema && buf generate --template ./notification/buf.gen.yaml --path notification
	cd protobuf/schema && buf generate --template ./bdspro/buf.gen.yaml --path bdspro
	cd protobuf/schema && buf generate --template ./task/buf.gen.yaml --path task
	cd protobuf/schema && buf generate --template ./payment/buf.gen.yaml --path payment
	cd protobuf/schema && buf generate --template ./spec/buf.gen.yaml --path spec
	cd protobuf/schema && buf generate --template ./generic/buf.gen.yaml --path generic
	cd protobuf/schema && buf generate --template ./organization/buf.gen.yaml --path organization

payment:
	cd payment-service && go run main.go http
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
buf-crm:
	cd protobuf/schema && buf generate --template ./crm/buf.gen.yaml --path crm

user-http:
	cd user-service && go run main.go http

chat:
	cd go-chat-service && go run main.go grpc-serve