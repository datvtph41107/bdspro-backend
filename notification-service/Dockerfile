FROM alpine:latest

WORKDIR /

COPY main .

COPY config config

ENV ENV_RUNTIME=develop

CMD ["./main"]