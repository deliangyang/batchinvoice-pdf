package gui

import "strings"

// IMAPServerConfig IMAP服务器配置
type IMAPServerConfig struct {
	Server string
	Port   int
}

// commonIMAPServers 常见邮箱IMAP服务器映射
var commonIMAPServers = map[string]IMAPServerConfig{
	// 国内邮箱
	"qq.com":      {Server: "imap.qq.com", Port: 993},
	"163.com":     {Server: "imap.163.com", Port: 993},
	"126.com":     {Server: "imap.126.com", Port: 993},
	"yeah.net":    {Server: "imap.yeah.net", Port: 993},
	"vip.163.com": {Server: "imap.vip.163.com", Port: 993},
	"vip.126.com": {Server: "imap.vip.126.com", Port: 993},
	"sina.com":    {Server: "imap.sina.com", Port: 993},
	"sina.cn":     {Server: "imap.sina.cn", Port: 993},
	"sohu.com":    {Server: "imap.sohu.com", Port: 993},
	"139.com":     {Server: "imap.139.com", Port: 993},
	"wo.cn":       {Server: "imap.wo.cn", Port: 993},
	"189.cn":      {Server: "imap.189.cn", Port: 993},
	"aliyun.com":  {Server: "imap.aliyun.com", Port: 993},
	"foxmail.com": {Server: "imap.qq.com", Port: 993},

	// 国际邮箱
	"gmail.com":      {Server: "imap.gmail.com", Port: 993},
	"outlook.com":    {Server: "outlook.office365.com", Port: 993},
	"hotmail.com":    {Server: "outlook.office365.com", Port: 993},
	"live.com":       {Server: "outlook.office365.com", Port: 993},
	"yahoo.com":      {Server: "imap.mail.yahoo.com", Port: 993},
	"icloud.com":     {Server: "imap.mail.me.com", Port: 993},
	"me.com":         {Server: "imap.mail.me.com", Port: 993},
	"aol.com":        {Server: "imap.aol.com", Port: 993},
	"protonmail.com": {Server: "127.0.0.1", Port: 1143}, // ProtonMail需要Bridge
	"yandex.com":     {Server: "imap.yandex.com", Port: 993},
	"mail.com":       {Server: "imap.mail.com", Port: 993},
	"gmx.com":        {Server: "imap.gmx.com", Port: 993},
	"zoho.com":       {Server: "imap.zoho.com", Port: 993},
}

// GetIMAPServer 根据邮箱地址获取IMAP服务器配置
func GetIMAPServer(email string) (IMAPServerConfig, bool) {
	// 提取域名部分
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return IMAPServerConfig{Server: "imap.gmail.com", Port: 993}, false
	}

	domain := strings.ToLower(strings.TrimSpace(parts[1]))

	// 查找对应的IMAP服务器
	if config, ok := commonIMAPServers[domain]; ok {
		return config, true
	}

	// 默认返回Gmail配置
	return IMAPServerConfig{Server: "imap.gmail.com", Port: 993}, false
}

// AddIMAPServer 添加自定义IMAP服务器映射
func AddIMAPServer(domain string, server string, port int) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	commonIMAPServers[domain] = IMAPServerConfig{
		Server: server,
		Port:   port,
	}
}
