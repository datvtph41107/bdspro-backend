
wire:
	cd script && go run gen_wire.go && cd ../cmd/http && wire

http:
	go run main.go http

grpc:
	go run main.go grpc
	
swag:
	swag init --parseDependency

migrateup:
	migrate -path infra/db/migrations/v2 -database "postgresql://hieubap:Matkhau.1@14.225.210.29:5432/db_bdspro?sslmode=disable" -verbose up

migratedown:
	migrate -path infra/db/migrations/v2 -database "postgresql://hieubap:Matkhau.1@14.225.210.29:5432/db_bdspro?sslmode=disable" -verbose down
