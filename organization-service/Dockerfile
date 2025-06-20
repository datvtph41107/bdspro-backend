FROM alpine:latest

WORKDIR /

COPY main .

ENV APP_ENV=release

CMD ["./main", "grpc-server"]