# Changelog

All notable changes to BatchInvoice PDF will be documented in this file.

## [1.0.0] - 2026-02-26

### 🎉 初始版本

基于 `email-pdf-qrcode-get` Python 技能改造的 Golang 桌面应用。

### ✨ 新增功能

#### 核心功能
- **IMAP 邮件连接**: 支持 Gmail、QQ、163、Outlook 等主流邮箱
- **PDF 附件处理**: 自动识别和提取邮件中的 PDF 发票附件
- **二维码识别**: 使用 gozxing 识别 PDF 中的二维码
- **数据解析**: 自动解析发票号码、金额、日期等信息
- **税额计算**: 自动计算不含税金额和税额（支持自定义税率，默认13%）
- **来源识别**: 自动识别发票来源（沃尔玛、盒马等）

#### 用户界面
- **跨平台 GUI**: 基于 Fyne 框架的现代化桌面界面
- **配置表单**: 友好的邮箱配置界面
- **实时反馈**: 提取进度提示
- **结果展示**: 美观的发票列表展示
- **自定义主题**: 紫色渐变主题设计

#### 数据导出
- **JSON 导出**: 结构化数据，便于程序处理
- **HTML 导出**: 美观的网页表格，支持浏览器直接查看
- **Markdown 导出**: 文档格式，便于编写报告

#### 开发工具
- **多平台构建**: 支持 Windows、macOS (Intel/ARM)、Linux
- **构建脚本**: 
  - `build.sh` - Linux/macOS 构建脚本
  - `build.bat` - Windows 构建脚本
  - `Makefile` - 统一的构建工具
- **快速运行**:
  - `run.sh` - Linux/macOS 开发运行
  - `run.bat` - Windows 开发运行

### 📦 项目结构

```
batchinvoice-pdf/
├── main.go                    # 主程序入口
├── src/
│   ├── gui/                  # GUI 界面层
│   │   ├── app.go           # 主应用窗口
│   │   └── theme.go         # 自定义主题
│   ├── core/                 # 核心业务逻辑
│   │   ├── email.go         # IMAP 邮件处理
│   │   ├── pdf.go           # PDF 文件处理
│   │   ├── qrcode.go        # 二维码识别
│   │   ├── invoice.go       # 发票数据模型
│   │   └── parser.go        # 数据解析器
│   └── utils/                # 工具函数
│       ├── tax.go           # 税额计算
│       └── export.go        # 数据导出
├── resources/                # 资源文件
├── build/                    # 编译输出目录
└── [配置和文档文件...]
```

### 🔧 技术栈

- **语言**: Go 1.21+
- **GUI 框架**: Fyne v2.4+
- **邮件处理**: 
  - `github.com/emersion/go-imap` - IMAP 客户端
  - `github.com/emersion/go-message` - 邮件消息解析
- **PDF 处理**: 
  - `github.com/pdfcpu/pdfcpu` - PDF 文档处理
- **二维码识别**: 
  - `github.com/makiuchi-d/gozxing` - ZXing Go 实现

### 📝 文档

- `README.md` - 项目概述和功能说明
- `QUICKSTART.md` - 快速开始指南
- `GUIDE.md` - 详细使用指南和开发文档
- `config.example.json` - 配置文件示例

### 🚀 与 Python 版本对比

#### 优势
- ✅ **跨平台**: 原生编译，无需 Python 环境
- ✅ **启动速度**: 秒级启动 vs Python 数秒
- ✅ **内存占用**: 更低的内存消耗
- ✅ **GUI 界面**: 友好的图形界面 vs 命令行
- ✅ **单文件部署**: 一个可执行文件即可运行
- ✅ **并发处理**: Go 协程并发处理邮件

#### 功能保持
- ✅ 相同的邮件处理能力
- ✅ 相同的二维码识别逻辑
- ✅ 相同的税额计算公式
- ✅ 兼容的数据格式

### 🔮 未来计划

#### v1.1.0 (计划中)
- [ ] 改进 PDF 渲染（集成 pdfium 或 poppler）
- [ ] 支持批量导出为 Excel
- [ ] 添加发票数据库存储
- [ ] 支持发票分类和标签
- [ ] 添加邮件过滤规则

#### v1.2.0 (计划中)
- [ ] 支持自动定时扫描
- [ ] 添加邮件通知功能
- [ ] 云端数据同步
- [ ] 多语言支持 (i18n)
- [ ] 插件系统

### 📄 许可证

MIT License - 详见 LICENSE 文件

### 🙏 致谢

感谢 `email-pdf-qrcode-get` Python 技能提供的核心思路和实现参考。

---

## 版本说明

- **Semantic Versioning**: 采用语义化版本 (MAJOR.MINOR.PATCH)
- **发布周期**: 主版本不定期，次版本2-3个月，补丁版本随时

## 如何贡献

欢迎提交 Issue 和 Pull Request！

请参考 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。
