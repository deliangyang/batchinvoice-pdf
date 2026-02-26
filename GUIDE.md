# BatchInvoice PDF - 使用指南

## 🎯 项目说明

**BatchInvoice PDF** 是 `email-pdf-qrcode-get` Python技能 的 **Golang 桌面应用版本**，提供跨平台支持（Windows、macOS、Linux）。

### 核心功能对比

| 功能 | Python 版本 | Golang 桌面版 |
|------|------------|--------------|
| 邮件连接 | ✅ IMAP | ✅ IMAP |
| PDF处理 | ✅ pdf2image | ✅ pdfcpu |
| 二维码识别 | ✅ pyzbar | ✅ gozxing |
| 税额计算 | ✅ 13% | ✅ 可配置 |
| 界面 | ❌ 命令行 | ✅ GUI |
| 跨平台 | ⚠️  需环境配置 | ✅ 原生支持 |

## 🚀 快速开始

### 1. 前置要求

- Go 1.21 或更高版本
- C编译器 (用于CGO)
  - **Linux**: `gcc`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: MinGW-w64 或 TDM-GCC

### 2. 安装依赖

```bash
# Linux
sudo apt install gcc libgl1-mesa-dev xorg-dev

# macOS
xcode-select --install

# Windows
# 下载安装 TDM-GCC: https://jmeubank.github.io/tdm-gcc/
```

### 3. 编译运行

```bash
# 开发模式运行
make run
# 或
./run.sh  # Linux/macOS
run.bat   # Windows

# 构建可执行文件
make build

# 构建所有平台
make build-all
```

## 📋 使用步骤

### 步骤 1: 配置邮箱

启动应用后，填写以下信息：

```
IMAP服务器: imap.gmail.com
端口: 993
用户名: your-email@gmail.com
密码: your-app-password
最大邮件数: 20
税率: 0.13
```

### 步骤 2: 邮箱服务器配置

#### Gmail
1. 开启两步验证
2. 生成应用专用密码
3. 使用应用密码登录

#### QQ邮箱
1. 设置 → 账户 → 开启IMAP服务
2. 获取授权码
3. 使用授权码作为密码

#### 163邮箱
1. 设置 → POP3/SMTP/IMAP
2. 开启IMAP服务
3. 获取授权码

### 步骤 3: 提取发票

1. 点击 "🚀 开始提取" 按钮
2. 等待处理完成
3. 查看提取结果

### 步骤 4: 导出数据

支持导出格式：
- **JSON** - 程序可读
- **HTML** - 美观展示
- **Markdown** - 文档编写

## 🏗️ 项目架构

```
batchinvoice-pdf/
├── main.go                    # 应用入口
├── src/
│   ├── gui/                  # 图形界面
│   │   ├── app.go           # 主窗口
│   │   └── theme.go         # 主题
│   ├── core/                 # 核心逻辑
│   │   ├── email.go         # IMAP邮件
│   │   ├── pdf.go           # PDF处理
│   │   ├── qrcode.go        # 二维码识别
│   │   ├── invoice.go       # 发票模型
│   │   └── parser.go        # 数据解析
│   └── utils/                # 工具函数
│       ├── tax.go           # 税额计算
│       └── export.go        # 数据导出
└── build/                    # 编译输出
```

## 📦 编译产物

运行 `make build-all` 后生成：

```
build/
├── batchinvoice-pdf-windows-amd64.exe   # Windows 64位
├── batchinvoice-pdf-macos-amd64         # macOS Intel
├── batchinvoice-pdf-macos-arm64         # macOS Apple Silicon
└── batchinvoice-pdf-linux-amd64         # Linux 64位
```

## 🔧 开发命令

```bash
# 安装依赖
make deps

# 开发运行
make run

# 格式化代码
make fmt

# 运行测试
make test

# 代码检查
make lint

# 清理构建
make clean

# 查看帮助
make help
```

## ⚡ 性能优化

相比Python版本的优势：

1. **启动速度**: 原生编译，秒级启动
2. **内存占用**: 更低的内存消耗
3. **并发处理**: Go协程并发处理邮件
4. **单文件部署**: 无需Python环境

## 🐛 常见问题

### 编译失败

**问题**: `gcc: command not found`

**解决**:
```bash
# Ubuntu/Debian
sudo apt install build-essential

# macOS
xcode-select --install

# Windows
# 安装 TDM-GCC 或 MinGW-w64
```

### 运行错误

**问题**: IMAP连接失败

**检查项**:
1. IMAP服务器地址是否正确
2. 端口是否为993
3. 用户名和密码是否正确
4. 是否需要应用专用密码
5. 防火墙是否阻止连接

### PDF处理问题

**问题**: 无法提取二维码

**说明**: 当前版本的PDF处理使用pdfcpu，对某些复杂PDF支持有限。建议：

1. 确保PDF格式正确
2. 检查PDF是否包含图像二维码
3. 尝试用PDF阅读器打开验证

## 📚 从Python版本迁移

如果你之前使用Python版本，迁移步骤：

1. **保留原有配置**: 
   - `config.json` 格式兼容
   
2. **数据迁移**:
   - JSON输出格式相同
   - 可直接使用之前的数据

3. **功能对应**:
   - Python `extract_invoice_qrcode.py` → GUI "开始提取"
   - Python `generate_html_table.py` → GUI "导出HTML"
   - Python `generate_markdown_table.py` → GUI "导出Markdown"

## 🎨 界面截图

（界面采用Fyne框架，支持亮色/暗色主题）

### 主界面
- 邮箱配置表单
- 提取按钮
- 结果展示区域

### 结果显示
- 发票清单表格
- 税额计算明细
- 统计汇总信息

## 🔐 安全说明

1. **密码存储**: 应用不会保存密码
2. **数据隐私**: 所有数据本地处理
3. **配置文件**: 不要提交含密码的config.json

建议在 `.gitignore` 中添加：
```gitignore
config.json
*.local.json
```

## 📄 许可证

MIT License

## 🙏 致谢

- 基于 `email-pdf-qrcode-get` Python技能
- 使用 Fyne GUI框架
- 使用 go-imap 处理邮件
- 使用 gozxing 识别二维码
- 使用 pdfcpu 处理PDF

## 📮 反馈

如有问题或建议，欢迎提交 Issue 或 Pull Request。
