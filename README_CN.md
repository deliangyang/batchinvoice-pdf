# 📊 BatchInvoice PDF

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**跨平台的发票二维码批量提取工具**

基于 `email-pdf-qrcode-get` Python 技能改造的 Golang 桌面应用

[快速开始](#快速开始) • [功能特性](#功能特性) • [下载](#下载) • [文档](#文档)

</div>

---

## ✨ 功能特性

### 📧 邮件处理
- ✅ 支持 IMAP 协议连接主流邮箱（Gmail、QQ、163、Outlook）
- ✅ 自动读取最近 N 封邮件（可配置）
- ✅ 智能识别发票 PDF 附件

### 🔍 智能识别
- ✅ 自动提取 PDF 中的二维码
- ✅ 解析发票号码、金额、日期
- ✅ 识别发票来源（沃尔玛、盒马等）

### 💰 税额计算
- ✅ 自动计算不含税金额
- ✅ 精确计算税额（支持自定义税率）
- ✅ 默认 13% 增值税率

### 📊 数据展示
- ✅ 美观的 GUI 界面
- ✅ 实时提取进度显示
- ✅ 发票列表展示

### 💾 数据导出
- ✅ JSON 格式（程序可读）
- ✅ HTML 格式（浏览器查看）
- ✅ Markdown 格式（文档编写）

### 🖥️ 跨平台
- ✅ Windows 64位
- ✅ macOS Intel
- ✅ macOS Apple Silicon
- ✅ Linux 64位

## 🚀 快速开始

### 前置要求

- Go 1.21+ （仅开发需要）
- C 编译器（仅编译需要）
  - Linux: `gcc`
  - macOS: Xcode Command Line Tools
  - Windows: MinGW-w64

### 下载预编译版本

访问 [Releases](https://github.com/yourusername/batchinvoice-pdf/releases) 下载对应平台的可执行文件。

### 从源码编译

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/batchinvoice-pdf.git
cd batchinvoice-pdf

# 2. 安装依赖
make deps

# 3. 编译
make build          # 编译当前平台
make build-all      # 编译所有平台

# 4. 运行
./build/batchinvoice-pdf-linux-amd64  # Linux
./build/batchinvoice-pdf-macos-amd64  # macOS
build\batchinvoice-pdf-windows-amd64.exe  # Windows
```

### 开发模式运行

```bash
make run
# 或
./run.sh    # Linux/macOS
run.bat     # Windows
```

## 📖 使用示例

### 1️⃣ 配置邮箱

启动应用后，填写邮箱信息：

```
IMAP 服务器: imap.gmail.com
端口: 993
用户名: your-email@gmail.com
密码: your-app-password
最大邮件数: 20
税率: 0.13
```

### 2️⃣ 提取发票

1. 点击 "🚀 开始提取" 按钮
2. 等待处理完成
3. 查看提取结果

### 3️⃣ 导出数据

选择导出格式：
- JSON - 结构化数据
- HTML - 美观展示
- Markdown - 文档格式

## 🏗️ 项目结构

```
batchinvoice-pdf/
├── main.go                 # 主程序入口
├── src/
│   ├── gui/               # GUI 界面
│   ├── core/              # 核心逻辑
│   └── utils/             # 工具函数
├── resources/             # 资源文件
├── build/                 # 编译输出
└── docs/                  # 文档
```

## 📚 文档

- [快速开始指南](QUICKSTART.md)
- [详细使用指南](GUIDE.md)
- [更新日志](CHANGELOG.md)

## 🔧 开发

```bash
make deps       # 安装依赖
make run        # 开发运行
make test       # 运行测试
make fmt        # 格式化代码
make lint       # 代码检查
make clean      # 清理构建
make help       # 查看帮助
```

## 📊 性能对比

| 指标 | Python 版本 | Golang 版本 |
|------|------------|-------------|
| 启动时间 | ~3-5秒 | <1秒 |
| 内存占用 | ~150MB | ~50MB |
| 可执行文件大小 | N/A (需环境) | ~40MB |
| 跨平台部署 | ❌ 需配置环境 | ✅ 单文件运行 |
| GUI 界面 | ❌ 命令行 | ✅ 原生 GUI |

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- 基于 `email-pdf-qrcode-get` Python 技能
- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- [go-imap](https://github.com/emersion/go-imap) - IMAP 客户端
- [gozxing](https://github.com/makiuchi-d/gozxing) - 二维码识别
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF 处理

## 📮 联系方式

- Issue: [GitHub Issues](https://github.com/yourusername/batchinvoice-pdf/issues)
- Email: your-email@example.com

---

<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐️ Star！**

Made with ❤️ by [Your Name]

</div>
