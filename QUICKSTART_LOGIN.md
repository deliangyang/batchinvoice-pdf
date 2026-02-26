# 快速配置指南

## 第一步：设置默认邮箱地址

打开 `main.go` 文件，修改 `DefaultEmail` 常量：

```go
const (
    // 默认邮箱地址（这里可以配置为你的固定邮箱）
    DefaultEmail = "zhangsan@gmail.com"  // 改成你的邮箱
)
```

## 第二步：编译程序

```bash
# Linux/Mac
go build -o batchinvoice-pdf main.go

# Windows
go build -o batchinvoice-pdf.exe main.go
```

## 第三步：运行程序

```bash
# Linux/Mac
./batchinvoice-pdf

# Windows
batchinvoice-pdf.exe
```

## 第四步：登录

1. 程序启动后会显示登录界面
2. 界面会显示你设置的邮箱地址（不可编辑）
3. 系统会自动识别 IMAP 服务器（如 Gmail 自动配置为 imap.gmail.com:993）
4. 输入邮箱密码或应用专用密码
5. 点击"登录"按钮

## 第五步：使用

登录成功后：
1. 主界面显示已登录的邮箱信息
2. 可以调整"最大邮件数"（默认 20）
3. 可以调整"税率"（默认 0.13）
4. 点击"开始提取"按钮开始处理发票

## 常见邮箱密码获取方式

### Gmail
1. 访问 https://myaccount.google.com/security
2. 启用"两步验证"
3. 在"应用专用密码"中生成新密码
4. 使用生成的 16 位密码登录

### QQ 邮箱
1. 登录 QQ 邮箱网页版
2. 设置 → 账户 → 开启 IMAP/SMTP 服务
3. 按提示发送短信验证
4. 获取授权码（通常是 16 位字符串）
5. 使用授权码作为密码登录

### 163/126 邮箱
1. 登录网易邮箱网页版
2. 设置 → POP3/SMTP/IMAP → 开启 IMAP 服务
3. 设置客户端授权密码
4. 使用授权密码登录

### Outlook/Hotmail
直接使用账户密码即可，如果启用了两步验证则需要应用密码。

## 支持的邮箱类型

系统已内置支持以下邮箱的自动配置：
- Gmail, Outlook, Yahoo, iCloud
- QQ, 163, 126, 新浪, 搜狐
- 阿里云, Foxmail
- 以及更多...

详细列表请查看 [LOGIN_GUIDE.md](LOGIN_GUIDE.md)

## 需要帮助？

查看完整文档：
- [LOGIN_GUIDE.md](LOGIN_GUIDE.md) - 登录功能详细说明
- [README.md](README.md) - 项目总体说明
- [GUIDE.md](GUIDE.md) - 使用指南
