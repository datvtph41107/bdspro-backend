.PHONY: help run test race vet check smoke

help:
	@printf '%s\n' \
	  'run          Run the canonical HTTP Gateway' \
	  'test         Run Gateway tests' \
	  'race         Run Gateway tests with race detector' \
	  'vet          Run go vet' \
	  'check        Run architecture guards + tests + vet' \
	  'smoke        Run focused Gateway architecture tests'

run:
	go run . http

test:
	go test -count=1 ./...

race:
	go test -count=1 -race ./...

vet:
	go vet ./...

check:
	cd .. && bash scripts/architecture/verify-foundation-gateway.sh .
	cd .. && bash scripts/architecture/verify-foundation-gateway-runtime.sh .
	$(MAKE) test
	$(MAKE) vet

smoke:
	go test -count=1 ./config ./internal/authroutes ./internal/grpcgateway ./internal/grpcmetadata ./internal/httpauth ./internal/httperror ./internal/httpmiddleware ./internal/httpresponse ./internal/tqdmultipart ./internal/upstream ./internal/v1proxy ./cmd/http
