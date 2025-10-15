FROM file-ffmpeg

WORKDIR /

COPY docs docs

COPY config config

COPY main .

ENV ENV_RUNTIME=develop

CMD ["./main"]