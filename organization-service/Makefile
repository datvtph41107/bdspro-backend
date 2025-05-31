.PHONY: build
build:
	go build -o organization-service main.go


.PHONY: swagger
swagger:
	go run main.go merge-swagger

.PHONY: swagger
run:
	go run main.go http-gateway

