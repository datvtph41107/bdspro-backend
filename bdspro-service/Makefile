
wire:
	cd script && go run gen_wire.go && cd ../cmd/http && wire

http:
	go run main.go http

grpc:
	go run main.go grpc
	
swag:
	swag init --parseDependency
	