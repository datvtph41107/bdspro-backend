.PHONY: build
build:
	go build -o organization-service main.go


.PHONY: swagger
swagger:
	go run main.go merge-swagger

.PHONY: run
run: swagger
	go run main.go http-gateway

