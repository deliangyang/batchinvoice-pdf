# 登录功能使用说明

## 功能概述

现在程序启动时会首先显示登录界面，要求用户输入邮箱密码。邮箱地址已预设为常量，且会自动识别对应的 IMAP 服务器配置。

## 主要特性

### 1. 固定邮箱地址
- 邮箱地址在 `main.go` 中配置为常量 `DefaultEmail`
- 登录界面中邮箱地址不可编辑，确保安全性

### 2. 自动 IMAP 服务器识别
系统支持自动识别以下邮箱的 IMAP 服务器配置：

**国内邮箱：**
- QQ邮箱 (qq.com, foxmail.com)
- 网易邮箱 (163.com, 126.com, yeah.net, vip.163.com, vip.126.com)
- 新浪邮箱 (sina.com, sina.cn)
- 搜狐邮箱 (sohu.com)
- 移动邮箱 (139.com)
- 联通邮箱 (wo.cn)
- 电信邮箱 (189.cn)
- 阿里云邮箱 (aliyun.com)

**国际邮箱：**
- Gmail (gmail.com)
- Outlook/Hotmail (outlook.com, hotmail.com, live.com)
- Yahoo (yahoo.com)
- iCloud (icloud.com, me.com)
- AOL (aol.com)
- Yandex (yandex.com)
- Mail.com (mail.com)
- GMX (gmx.com)
- Zoho (zoho.com)

### 3. 登录验证
- 输入密码后会自动连接邮箱服务器验证
- 验证成功后自动关闭登录窗口并打开主界面
- 验证失败会显示错误提示

## 配置方法

### 修改默认邮箱地址

编辑 `main.go` 文件中的 `DefaultEmail` 常量：

```go
const (
    // 默认邮箱地址（这里可以配置为你的固定邮箱）
    DefaultEmail = "your-email@gmail.com"  // 修改为你的邮箱
)
```

### 添加自定义 IMAP 服务器

如果你使用的邮箱不在支持列表中，可以在 `src/gui/imap_servers.go` 文件的 `commonIMAPServers` 映射中添加：

```go
var commonIMAPServers = map[string]IMAPServerConfig{
    // 添加你的邮箱配置
    "yourdomain.com": {Server: "imap.yourdomain.com", Port: 993},
    // ... 其他配置
}
```

## 使用流程

1. **启动程序**
   - 运行程序后自动显示登录界面
   - 界面显示固定的邮箱地址
   - 自动识别并显示 IMAP 服务器信息

2. **输入密码**
   - 在密码框中输入邮箱密码或应用专用密码
   - Gmail 等邮箱需要使用"应用专用密码"而非账户密码

3. **登录**
   - 点击"登录"按钮或按回车键
   - 系统自动验证邮箱连接
   - 验证成功后进入主界面

4. **使用主界面**
   - 主界面显示已登录的邮箱信息
   - 只需配置"最大邮件数"和"税率"即可开始提取

## 注意事项

### Gmail 用户
Gmail 需要使用"应用专用密码"：
1. 登录 Google 账户
2. 访问"安全性"页面
3. 启用"两步验证"
4. 生成"应用专用密码"
5. 使用生成的密码登录

### QQ 邮箱用户
QQ 邮箱需要开启 IMAP 服务并获取授权码：
1. 登录 QQ 邮箱网页版
2. 设置 → 账户 → 开启 IMAP/SMTP 服务
3. 生成授权码
4. 使用授权码作为密码登录

### 其他邮箱
大多数邮箱都需要在设置中启用 IMAP 服务，某些邮箱可能需要应用专用密码。

## 文件说明

- `main.go` - 程序入口，配置默认邮箱地址
- `src/gui/login.go` - 登录窗口实现
- `src/gui/imap_servers.go` - IMAP 服务器映射配置
- `src/gui/app.go` - 主窗口实现

## 编译和运行

```bash
# 编译
go build -o batchinvoice-pdf main.go

# 运行
./batchinvoice-pdf
```

或使用提供的脚本：

```bash
# Linux/Mac
./build.sh
./run.sh

# Windows
build.bat
run.bat
```
