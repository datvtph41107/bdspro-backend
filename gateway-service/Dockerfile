FROM alpine:latest

WORKDIR /

COPY main .

COPY docs docs

COPY config config

ENV ENV_RUNTIME=release

EXPOSE 8080

CMD ["./main","http"]