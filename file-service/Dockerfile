# syntax=docker/dockerfile:1
FROM golang:1.25.10-alpine AS build
WORKDIR /src

COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY file-service ./file-service

WORKDIR /src/file-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/file-service .

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates ffmpeg tzdata wget && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/file-service /app/file-service
COPY file-service/config /app/config
RUN mkdir -p /app/files && chown -R app:app /app

USER app
EXPOSE 8002
ENTRYPOINT ["/app/file-service"]
