# syntax=docker/dockerfile:1
# Build from repository-root context:
# docker build -f gateway-service/Dockerfile .
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY gateway-service ./gateway-service
WORKDIR /src/gateway-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/gateway .

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/gateway /app/gateway
COPY gateway-service/config /app/config
COPY gateway-service/docs /app/docs
# Gateway hiện ghi một bản log compatibility cạnh stdout. Runtime user cần
# quyền trên đúng thư mục này, không cần quyền root cho toàn container.
RUN chown -R app:app /app
USER app
ENTRYPOINT ["/app/gateway", "http"]
