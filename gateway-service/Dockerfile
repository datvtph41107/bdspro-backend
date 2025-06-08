FROM alpine:latest

WORKDIR /

COPY main .

COPY docs docs

ENV APP_ENV=release

EXPOSE 8080

CMD ["./main","http"]