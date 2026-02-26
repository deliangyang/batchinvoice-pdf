package core

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

func init() {
	// 注册中文字符集支持
	charset.RegisterEncoding("gb2312", simplifiedchinese.GB18030)
	charset.RegisterEncoding("gbk", simplifiedchinese.GB18030)
	charset.RegisterEncoding("gb18030", simplifiedchinese.GB18030)
	charset.RegisterEncoding("big5", traditionalchinese.Big5)
}

// EmailClient IMAP邮件客户端
type EmailClient struct {
	config *Config
	client *client.Client
}

// NewEmailClient 创建邮件客户端
func NewEmailClient(config *Config) *EmailClient {
	return &EmailClient{
		config: config,
	}
}

// Connect 连接到邮箱服务器
func (e *EmailClient) Connect() error {
	addr := fmt.Sprintf("%s:%d", e.config.IMAPServer, e.config.Port)
	log.Printf("Connecting to %s...", addr)

	c, err := client.DialTLS(addr, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	if err := c.Login(e.config.Username, e.config.Password); err != nil {
		c.Logout()
		return fmt.Errorf("failed to login: %w", err)
	}

	e.client = c
	log.Println("Successfully connected to email server")
	return nil
}

// Disconnect 断开连接
func (e *EmailClient) Disconnect() {
	if e.client != nil {
		e.client.Logout()
		log.Println("Disconnected from email server")
	}
}

// FetchRecentEmails 获取最近的邮件（从新到旧，限制最近31天）
func (e *EmailClient) FetchRecentEmails(maxEmails int) ([]*EmailMessage, error) {
	// 选择收件箱
	mbox, err := e.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}

	if mbox.Messages == 0 {
		return []*EmailMessage{}, nil
	}

	// 搜索最近31天的邮件
	thirtyOneDaysAgo := time.Now().AddDate(0, 0, -31)
	criteria := imap.NewSearchCriteria()
	criteria.Since = thirtyOneDaysAgo

	log.Printf("Searching emails since %s (total messages in mailbox: %d)", thirtyOneDaysAgo.Format("2006-01-02"), mbox.Messages)

	// 执行搜索
	seqNums, err := e.client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %w", err)
	}

	if len(seqNums) == 0 {
		log.Println("No emails found in the last 31 days")
		return []*EmailMessage{}, nil
	}

	log.Printf("Found %d emails in the last 31 days", len(seqNums))

	// 按序号从大到小排序（最新的在前）
	sort.Slice(seqNums, func(i, j int) bool {
		return seqNums[i] > seqNums[j]
	})

	// 如果邮件数超过 maxEmails，只取最新的 maxEmails 封
	if len(seqNums) > maxEmails {
		seqNums = seqNums[:maxEmails]
	}

	log.Printf("Fetching %d emails (seq range: %d to %d)", len(seqNums), seqNums[0], seqNums[len(seqNums)-1])

	// 创建序列集
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	// 获取邮件
	messages := make(chan *imap.Message, len(seqNums))
	done := make(chan error, 1)
	go func() {
		done <- e.client.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchRFC822,
		}, messages)
	}()

	// 处理邮件，使用 map 存储以便按序号排序
	emailMap := make(map[uint32]*EmailMessage)
	var resultSeqNums []uint32

	for msg := range messages {
		emailMsg, err := e.parseMessage(msg)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}
		emailMap[msg.SeqNum] = emailMsg
		resultSeqNums = append(resultSeqNums, msg.SeqNum)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// 按序号从大到小排序（最新的邮件在前）
	sort.Slice(resultSeqNums, func(i, j int) bool {
		return resultSeqNums[i] > resultSeqNums[j]
	})

	// 按排序后的序号构建结果
	var result []*EmailMessage
	for _, seqNum := range resultSeqNums {
		result = append(result, emailMap[seqNum])
	}

	log.Printf("✓ Successfully fetched %d emails from the last 31 days (newest first)", len(result))

	return result, nil
}

// EmailMessage 邮件消息
type EmailMessage struct {
	Subject     string
	From        string
	Date        string
	Attachments []Attachment
}

// Attachment 附件
type Attachment struct {
	Filename string
	Data     []byte
}

// parseMessage 解析邮件消息
func (e *EmailClient) parseMessage(msg *imap.Message) (*EmailMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		return nil, fmt.Errorf("message body is nil")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail reader: %w", err)
	}

	header := mr.Header
	subject, _ := header.Subject()
	from, _ := header.AddressList("From")
	date, _ := header.Date()

	// 转换为东8区时间并格式化
	loc, _ := time.LoadLocation("Asia/Shanghai")
	emailDate := date.In(loc).Format("2006-01-02 15:04:05")

	emailMsg := &EmailMessage{
		Subject:     subject,
		Date:        emailDate,
		Attachments: []Attachment{},
	}

	if len(from) > 0 {
		emailMsg.From = from[0].Address
	}

	// 解析附件
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 某些邮件的 multipart 边界可能不标准，导致解析错误
			// 这种情况下我们已经获取了能获取的附件，可以安全退出
			if strings.Contains(err.Error(), "multipart") || strings.Contains(err.Error(), "EOF") {
				break
			}
			log.Printf("Failed to read part: %v", err)
			continue
		}

		switch h := p.Header.(type) {
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
				data, err := io.ReadAll(p.Body)
				if err != nil {
					log.Printf("Failed to read attachment: %v", err)
					continue
				}
				emailMsg.Attachments = append(emailMsg.Attachments, Attachment{
					Filename: filename,
					Data:     data,
				})
				log.Printf("Found PDF attachment: %s (%d bytes)", filename, len(data))
			}
		}
	}

	return emailMsg, nil
}
