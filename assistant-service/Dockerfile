FROM golang:1.24-bullseye

# Cài các thư viện cần thiết + certificates
RUN apt-get update && apt-get install -y \
    g++ \
    make \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /assistant
# 3️⃣ Copy toàn bộ source code
COPY assistant-service/ /assistant
COPY shared/ /shared

# 2️⃣ Tải các dependency, cache layer này
RUN go mod download

# 4️⃣ Build binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .
