FROM alpine:latest

WORKDIR /

COPY main .

COPY docs docs

ENV ENV_RUNTIME=release

EXPOSE 8080

CMD ["./main","http"]