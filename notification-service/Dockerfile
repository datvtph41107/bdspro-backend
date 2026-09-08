# syntax=docker/dockerfile:1
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY notification-service ./notification-service
WORKDIR /src/notification-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' \
      -o /out/notification-service . && \
    GOWORK=off CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' \
      -o /out/payment-event-worker ./cmd/payment-event-worker && \
    GOWORK=off CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' \
      -o /out/delivery-worker ./cmd/delivery-worker

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && \
    addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/notification-service /app/notification-service
COPY --from=build /out/payment-event-worker /app/payment-event-worker
COPY --from=build /out/delivery-worker /app/delivery-worker
COPY notification-service/config /app/config
USER app
ENTRYPOINT ["/app/notification-service"]
CMD ["grpc"]
