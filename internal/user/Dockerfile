# syntax=docker/dockerfile:1
# Build from repository-root context because user-service imports Shared modules
# through local replace directives in go.mod.
FROM golang:1.25.10-alpine AS build
WORKDIR /src

COPY infrastructure/runtime ./infrastructure/runtime
COPY proto ./proto
COPY infrastructure/base ./infrastructure/base
COPY internal/user ./internal/user

WORKDIR /src/internal/user
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/user ./cmd/grpc

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/user /app/user
COPY internal/user/config /app/config

USER app
EXPOSE 8201
ENTRYPOINT ["/app/user"]
