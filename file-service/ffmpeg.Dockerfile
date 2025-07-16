FROM alpine:latest

RUN apk update && apk add ffmpeg

# Optional: kiểm tra ffmpeg
RUN ffmpeg -version
