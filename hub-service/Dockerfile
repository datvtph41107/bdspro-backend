# syntax=docker/dockerfile:1
# Hub phục vụ metadata/location nội bộ cho gateway và file-service.
FROM golang:1.25.10-alpine AS build
WORKDIR /src
COPY shared/common ./shared/common
COPY shared/protobuf ./shared/protobuf
COPY shared/base ./shared/base
COPY hub-service ./hub-service
WORKDIR /src/hub-service
RUN --mount=type=cache,id=qhpro-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=qhpro-go-build,target=/root/.cache/go-build \
    GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/hub-service .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/hub-service /app/hub-service
COPY hub-service/config /app/config
USER app
EXPOSE 8280
ENTRYPOINT ["/app/hub-service", "grpc"]
