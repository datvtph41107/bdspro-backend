# Giai đoạn 1: Build binary
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build binary tĩnh
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/binary .

# Giai đoạn 2: Chạy ứng dụng
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /app/binary .

# Đảm bảo binary có quyền thực thi
RUN chmod +x /app/binary

CMD ["/app/binary", "http-serve"]
