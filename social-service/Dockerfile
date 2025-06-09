FROM alpine:latest

WORKDIR /

COPY main .

ENV APP_ENV=release

EXPOSE 8209

CMD ["./main","grpc"]