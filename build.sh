#!/bin/bash

# BatchInvoice PDF - 构建脚本
# 说明：Fyne 使用 CGO 和 OpenGL，无法用 go build 交叉编译到其他 OS。
# 本脚本仅构建当前平台；如需其他平台请在对应系统上执行本脚本，或使用 fyne cross（需 Docker）。

set -e

echo "🚀 Building BatchInvoice PDF..."

# 创建构建目录
mkdir -p build

# 获取版本信息
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

# 仅构建当前平台（Fyne/OpenGL 不支持交叉编译）
GOOS="${GOOS:-$(go env GOOS)}"
GOARCH="${GOARCH:-$(go env GOARCH)}"
echo "📦 Building for current platform: ${GOOS}/${GOARCH}..."

if [ "$GOOS" = "windows" ]; then
  go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-${GOOS}-${GOARCH}.exe main.go
  echo "✅ Build complete: build/batchinvoice-pdf-${GOOS}-${GOARCH}.exe"
else
  go build -ldflags "$LDFLAGS" -o build/batchinvoice-pdf-${GOOS}-${GOARCH} main.go
  echo "✅ Build complete: build/batchinvoice-pdf-${GOOS}-${GOARCH}"
fi

echo ""
echo "🎉 Build completed successfully!"
ls -lh build/
echo ""
echo ""
echo "Run: ./build/batchinvoice-pdf-${GOOS}-${GOARCH}$( [ "$GOOS" = "windows" ] && echo .exe )"
