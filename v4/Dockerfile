# syntax=docker/dockerfile:1
# Build from repository-root context because userv4 imports Shared modules
# through local replace directives in go.mod.
FROM golang:1.25.10-alpine AS build
WORKDIR /src

COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY shared/base ./shared/base
COPY v4 ./v4

WORKDIR /src/v4
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/userv4 ./cmd/grpc && \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/userv4-http ./cmd/http

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd wget && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/userv4 /app/userv4
COPY --from=build /out/userv4-http /app/userv4-http
COPY v4/config /app/config
COPY v4/migrate /app/migrate
ENV USER_MIGRATION_DIR=/app/migrate

USER app
EXPOSE 8211 8004
ENTRYPOINT ["/app/userv4"]
