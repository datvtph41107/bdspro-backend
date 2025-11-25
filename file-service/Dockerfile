FROM debian:stable-slim

# Cài FFmpeg
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY docs docs

COPY config config

COPY main .

ENV ENV_RUNTIME=develop

CMD ["./main","http"]