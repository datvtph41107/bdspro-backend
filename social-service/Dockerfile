FROM alpine:latest

WORKDIR /

COPY main .

COPY config config

ENV ENV_RUNTIME=develop

EXPOSE 8209

CMD ["./main","grpc"]