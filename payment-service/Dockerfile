# syntax=docker/dockerfile:1
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY shared/base ./shared/base
COPY payment-service ./payment-service
WORKDIR /src/payment-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/payment .

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/payment /app/payment
COPY payment-service/config /app/config
USER app
EXPOSE 8205
ENTRYPOINT ["/app/payment", "grpc"]
