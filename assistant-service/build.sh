#!/bin/bash

# Set CGO environment variables for Tesseract and Leptonica
# These must be set before building to allow gosseract to find the headers
export CGO_ENABLED=1
export CGO_CFLAGS="-I/opt/homebrew/include"
export CGO_LDFLAGS="-L/opt/homebrew/lib -ltesseract -lleptonica"

echo "Building assistant-service with CGO..."
echo "CGO_CFLAGS=$CGO_CFLAGS"
echo "CGO_LDFLAGS=$CGO_LDFLAGS"

# Build the application
go build -o assistant-service .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed. Make sure Tesseract and Leptonica are installed:"
    echo "   brew install tesseract leptonica"
    exit 1
fi

