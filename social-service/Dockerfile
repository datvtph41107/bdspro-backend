FROM alpine:latest

WORKDIR /

COPY main .

ENV ENV_RUNTIME=release

EXPOSE 8209

CMD ["./main","grpc"]