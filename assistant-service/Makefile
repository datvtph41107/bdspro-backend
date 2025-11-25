.PHONY: build run clean

# Set CGO environment variables for Tesseract and Leptonica
# For Homebrew on macOS (Apple Silicon)
export PKG_CONFIG_PATH=/opt/homebrew/lib/pkgconfig
export CGO_ENABLED=1
# export CGO_CFLAGS=-I/opt/homebrew/include -I/usr/local/include
# export CGO_LDFLAGS=-L/opt/homebrew/lib -ltesseract -lleptonica
export CGO_CFLAGS="-I/opt/homebrew/Cellar/tesseract/5.5.1_1/include -I/opt/homebrew/Cellar/leptonica/1.86.0/include/leptonica"
export CGO_LDFLAGS="-L/opt/homebrew/Cellar/tesseract/5.5.1_1/lib -ltesseract -L/opt/homebrew/opt/libarchive/lib -larchive -lcurl"
# go build

build:
	@echo "Building assistant-service with CGO..."
	@echo "CGO_CFLAGS=$(CGO_CFLAGS)"
	@echo "CGO_LDFLAGS=$(CGO_LDFLAGS)"
	@go build -o assistant-service .

run: build
	@echo "Running assistant-service..."
	@./assistant-service

clean:
	@echo "Cleaning build artifacts..."
	@rm -f assistant-service

test:
	@echo "Running tests..."
	@go test ./...

