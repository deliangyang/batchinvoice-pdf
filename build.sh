#!/bin/bash

# BatchInvoice PDF - 构建脚本

set -e

echo "🚀 Building BatchInvoice PDF..."

# 创建构建目录
mkdir -p build

# 获取版本信息
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

# 构建 Windows 版本
echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-windows-amd64.exe main.go
echo "✅ Windows build complete: build/batchinvoice-pdf-windows-amd64.exe"

# 构建 macOS 版本 (Intel)
echo "📦 Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-macos-amd64 main.go
echo "✅ macOS (Intel) build complete: build/batchinvoice-pdf-macos-amd64"

# 构建 macOS 版本 (Apple Silicon)
echo "📦 Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-macos-arm64 main.go
echo "✅ macOS (Apple Silicon) build complete: build/batchinvoice-pdf-macos-arm64"

# 构建 Linux 版本
echo "📦 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-linux-amd64 main.go
echo "✅ Linux build complete: build/batchinvoice-pdf-linux-amd64"

echo ""
echo "🎉 All builds completed successfully!"
echo ""
echo "Build files:"
ls -lh build/

echo ""
echo "To test the build, run:"
echo "  ./build/batchinvoice-pdf-linux-amd64  (Linux)"
echo "  ./build/batchinvoice-pdf-macos-amd64  (macOS Intel)"
echo "  ./build/batchinvoice-pdf-macos-arm64  (macOS Apple Silicon)"
echo "  build\\batchinvoice-pdf-windows-amd64.exe  (Windows)"
