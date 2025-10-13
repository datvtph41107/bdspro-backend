FROM alpine:latest

WORKDIR /

COPY main .

COPY config config

ENV ENV_RUNTIME=release

CMD ["./main", "grpc"]