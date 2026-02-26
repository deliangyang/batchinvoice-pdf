#!/bin/bash

# 开发快速启动脚本

echo "🚀 Starting BatchInvoice PDF..."

# 检查依赖
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# 安装依赖
echo "📦 Installing dependencies..."
go mod download

# 运行程序
echo "▶️  Running application..."
go run main.go
