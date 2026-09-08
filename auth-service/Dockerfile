# syntax=docker/dockerfile:1
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY auth-service ./auth-service
WORKDIR /src/auth-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth-service .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && \
    addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/auth-service /app/auth-service
COPY auth-service/config /app/config
USER app
EXPOSE 8216
ENTRYPOINT ["/app/auth-service"]
