gen:
	cd protobuf/schema && buf generate --template ./notification/buf.gen.yaml --path notification
	cd protobuf/schema && buf generate --template ./bdspro/buf.gen.yaml --path bdspro
	cd protobuf/schema && buf generate --template ./task/buf.gen.yaml --path task
	cd protobuf/schema && buf generate --template ./payment/buf.gen.yaml --path payment
	cd protobuf/schema && buf generate --template ./spec/buf.gen.yaml --path spec
	cd protobuf/schema && buf generate --template ./generic/buf.gen.yaml --path generic

payment:
	cd payment-service && go run main.go http
gateway:
	cd gateway-service && go run main.go http
bdspro:
	cd gateway-service && go run main.go grpc
user:
	cd user-service && go run main.go http
admin:
	cd admin-service && go run main.go grpc
