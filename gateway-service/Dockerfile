FROM alpine:latest

WORKDIR /

ENV ENV_RUNTIME=develop

COPY config config

COPY .env.backend.local .env.backend.local

COPY main .
