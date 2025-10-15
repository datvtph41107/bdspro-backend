FROM alpine:latest

WORKDIR /

COPY main .

COPY config config

ENV ENV_RUNTIME=develop

EXPOSE 8080

CMD ["./main"]