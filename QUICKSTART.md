# BatchInvoice PDF

## 快速开始

### 1. 安装依赖

```bash
make deps
```

或者直接使用 Go 命令：

```bash
go mod download
```

### 2. 运行开发版本

```bash
make run
```

或者：

```bash
go run main.go
```

### 3. 构建可执行文件

#### 构建当前平台

```bash
make build
```

#### 构建所有平台

```bash
make build-all
```

这将生成以下文件：
- `build/batchinvoice-pdf-windows-amd64.exe` - Windows 64位
- `build/batchinvoice-pdf-macos-amd64` - macOS Intel
- `build/batchinvoice-pdf-macos-arm64` - macOS Apple Silicon
- `build/batchinvoice-pdf-linux-amd64` - Linux 64位

### 4. 使用应用

1. 启动应用程序
2. 填写邮箱配置：
   - **IMAP服务器**：如 `imap.gmail.com`
   - **端口**：通常为 993
   - **用户名**：你的邮箱地址
   - **密码**：邮箱密码或应用专用密码
   - **最大邮件数**：要读取的最近邮件数量（默认20）
   - **税率**：增值税率（默认0.13，即13%）

3. 点击"开始提取"按钮
4. 等待处理完成
5. 查看提取结果
6. 可选：导出为 JSON、HTML 或 Markdown 格式

## 常见邮箱配置

### Gmail
- IMAP服务器: `imap.gmail.com`
- 端口: `993`
- 注意：需要启用"允许不够安全的应用"或使用应用专用密码

### QQ邮箱
- IMAP服务器: `imap.qq.com`
- 端口: `993`
- 注意：需要在QQ邮箱设置中开启IMAP服务并获取授权码

### 163邮箱
- IMAP服务器: `imap.163.com`
- 端口: `993`
- 注意：需要在邮箱设置中开启IMAP服务并获取授权码

### Outlook
- IMAP服务器: `outlook.office365.com`
- 端口: `993`

## 功能说明

### 二维码识别
- 自动从PDF发票中提取二维码
- 支持多页PDF文档
- 自动解析二维码内容

### 税额计算
- 自动计算不含税金额
- 计算税额
- 默认使用13%增值税率
- 可自定义税率

### 数据导出
支持三种导出格式：
- **JSON** - 结构化数据，便于程序处理
- **HTML** - 美观的网页表格，可直接在浏览器中查看
- **Markdown** - 文档格式，可用于文档编写

## 开发说明

### 项目结构

```
batchinvoice-pdf/
├── main.go                 # 程序入口
├── go.mod                  # Go模块定义
├── README.md              # 项目说明
├── QUICKSTART.md          # 快速开始指南
├── Makefile               # 构建工具
├── build.sh               # Linux/Mac构建脚本
├── build.bat              # Windows构建脚本
├── src/
│   ├── gui/              # GUI界面
│   │   ├── app.go       # 主应用窗口
│   │   └── theme.go     # 自定义主题
│   ├── core/             # 核心功能
│   │   ├── email.go     # 邮件处理
│   │   ├── pdf.go       # PDF处理
│   │   ├── qrcode.go    # 二维码识别
│   │   └── invoice.go   # 发票数据模型
│   └── utils/            # 工具函数
│       ├── tax.go       # 税额计算
│       └── export.go    # 数据导出
└── build/                # 构建输出
```

### 技术栈

- **Fyne** - 跨平台GUI框架
- **go-imap** - IMAP邮件客户端
- **pdfcpu** - PDF处理
- **gozxing** - 二维码识别

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

## 故障排查

### 连接邮箱失败
- 检查IMAP服务器地址和端口是否正确
- 确认已开启邮箱的IMAP服务
- 对于Gmail，需要使用应用专用密码
- 对于QQ/163邮箱，需要使用授权码而不是登录密码

### PDF处理失败
- 确认PDF文件格式正确
- 某些加密的PDF可能无法处理
- 尝试用其他PDF阅读器打开确认文件完整性

### 二维码识别不到
- 确认PDF中确实包含二维码
- 检查二维码图像是否清晰
- 部分矢量格式的二维码可能无法识别

## 许可证

MIT License

## 支持

如有问题，请提交 Issue 或 Pull Request。
