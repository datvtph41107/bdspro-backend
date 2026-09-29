# syntax=docker/dockerfile:1
# Build từ repository-root context vì Organization import shared modules.
FROM golang:1.25.10-alpine AS build
WORKDIR /src

COPY infrastructure/runtime ./infrastructure/runtime
COPY proto ./proto
COPY infrastructure/base ./infrastructure/base
COPY internal/chat ./internal/chat
COPY internal/organization-legacy ./internal/organization-legacy

WORKDIR /src/internal/organization-legacy
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/organization .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/organization /app/organization
COPY internal/organization-legacy/config /app/config

USER app
EXPOSE 8207
ENTRYPOINT ["/app/organization", "grpc"]
