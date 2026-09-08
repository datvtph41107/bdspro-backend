# syntax=docker/dockerfile:1
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY tqd-service ./tqd-service
WORKDIR /src/tqd-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o /out/tqd-service \
    .

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd py3-weasyprint font-dejavu && \
    fc-cache -f && \
    addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/tqd-service /app/tqd-service
COPY tqd-service/config /app/config
COPY tqd-service/template /app/template
# Một số import/PMTiles compatibility flow ghi artifact dưới /app/files. Tạo
# runtime directory có owner rõ trước khi hạ quyền, không chạy container root.
RUN mkdir -p /app/files && chown -R app:app /app
USER app
ENTRYPOINT ["/app/tqd-service"]
CMD ["-server=grpc"]
