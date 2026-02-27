# BatchInvoice PDF - 发票二维码批量提取工具

## 功能特性

一个跨平台的桌面应用程序，用于从邮件附件中批量提取发票PDF的二维码信息。

### 主要功能

- 📧 **邮件连接**：支持 IMAP 协议连接各大邮箱（Gmail、QQ、163、Outlook等）
- 📄 **PDF 处理**：自动识别和处理邮件中的发票PDF附件
- 🔍 **二维码识别**：自动提取PDF中的二维码并解析内容
- 💰 **税额计算**：自动计算不含税金额和税额（支持13%增值税）
- 📊 **数据展示**：美观的表格展示发票信息
- 💾 **数据导出**：支持导出为 JSON、HTML、Markdown 格式
- 🖥️ **跨平台**：支持 Windows、macOS 和 Linux

### 技术栈

- **语言**：Go 1.21+
- **GUI 框架**：Fyne v2
- **邮件处理**：go-imap
- **PDF 处理**：pdfcpu
- **二维码识别**：gozxing

## 安装依赖

```bash
go mod download
```

## 编译

### Windows

```bash
go build -o build/batchinvoice-pdf-windows.exe main.go
```

### macOS

```bash
go build -o build/batchinvoice-pdf-macos main.go
```

### Linux

```bash
go build -o build/batchinvoice-pdf-linux main.go
```

若使用 **预编译的 Linux 二进制** 时出现 `GLIBC_2.38 not found` 等错误，说明当前系统 glibc 较旧。请在本机从源码构建（会使用本机 glibc）：

```bash
# 安装构建依赖（Ubuntu/Debian）
sudo apt-get update
sudo apt-get install -y build-essential xorg-dev libgl1-mesa-dev libglu1-mesa-dev libxi-dev

# 在项目根目录执行
go mod download
./build.sh
# 或：go build -o build/batchinvoice-pdf-linux-amd64 main.go
```

构建产物在 `build/` 目录，直接运行即可。

## 使用方法

1. 启动程序
2. 输入邮箱配置信息（IMAP服务器、用户名、密码）
3. 点击"提取发票"按钮
4. 查看提取结果
5. 导出数据（可选）

## 配置说明

### 邮箱配置

- **IMAP服务器**：如 `imap.gmail.com`、`imap.qq.com`
- **端口**：通常为 993（SSL）
- **用户名**：邮箱地址
- **密码**：邮箱密码或应用专用密码

### 常见邮箱服务器

| 邮箱提供商 | IMAP 服务器 | 端口 |
|-----------|-------------|------|
| Gmail | imap.gmail.com | 993 |
| QQ邮箱 | imap.qq.com | 993 |
| 163邮箱 | imap.163.com | 993 |
| Outlook | outlook.office365.com | 993 |

## 项目结构

```
batchinvoice-pdf/
├── main.go                 # 主程序入口
├── src/
│   ├── gui/               # GUI界面
│   │   ├── app.go        # 应用主窗口
│   │   ├── config.go     # 配置界面
│   │   └── result.go     # 结果展示界面
│   ├── core/              # 核心功能
│   │   ├── email.go      # 邮件处理
│   │   ├── pdf.go        # PDF处理
│   │   ├── qrcode.go     # 二维码识别
│   │   └── invoice.go    # 发票数据结构
│   └── utils/             # 工具函数
│       ├── tax.go        # 税额计算
│       └── export.go     # 数据导出
├── resources/             # 资源文件
│   └── icon.png          # 应用图标
└── build/                 # 编译输出目录
```

## 开发说明

### 运行开发版本

```bash
go run main.go
```

### 运行测试

```bash
go test ./...
```

## 许可证

MIT License

## 致谢

基于 email-pdf-qrcode-get Python 技能改造而来。
