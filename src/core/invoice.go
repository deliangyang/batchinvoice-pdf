package core

import (
	"time"
)

// Invoice 发票信息结构
type Invoice struct {
	Name          string    `json:"name"`           // 姓名
	InvoiceNumber string    `json:"invoice_number"` // 发票号码
	QRCodeData    string    `json:"qrcode_data"`    // 二维码原始数据
	Source        string    `json:"source"`         // 来源（如：沃尔玛、盒马）
	Date          time.Time `json:"date"`           // 日期
	TotalAmount   float64   `json:"total_amount"`   // 含税总额
	AmountNoTax   float64   `json:"amount_no_tax"`  // 不含税金额
	TaxAmount     float64   `json:"tax_amount"`     // 税额
	EmailSubject  string    `json:"email_subject"`  // 邮件主题
	FileName      string    `json:"file_name"`      // PDF文件名
	EmailDate     string    `json:"email_date"`     // 邮件收件时间
	PDFData       []byte    `json:"-"`              // PDF原始数据（不导出到JSON）
}

// ExtractionResult 提取结果
type ExtractionResult struct {
	TotalEmails   int       `json:"total_emails"`   // 总邮件数
	ProcessedPDFs int       `json:"processed_pdfs"` // 处理的PDF数量
	QRCodesFound  int       `json:"qrcodes_found"`  // 找到的二维码数量
	Invoices      []Invoice `json:"invoices"`       // 发票列表
	ExtractedAt   time.Time `json:"extracted_at"`   // 提取时间
}

// Config 邮箱配置
type Config struct {
	IMAPServer string  `json:"imap_server"` // IMAP服务器地址
	Port       int     `json:"port"`        // 端口
	Username   string  `json:"username"`    // 用户名
	Password   string  `json:"password"`    // 密码
	MaxEmails  int     `json:"max_emails"`  // 最大邮件数
	TaxRate    float64 `json:"tax_rate"`    // 税率（默认0.13）
}

// NewConfig 创建默认配置
func NewConfig() *Config {
	return &Config{
		Port:      993,
		MaxEmails: 20,
		TaxRate:   0.13,
	}
}

// ParseQRCodeData 解析二维码数据
// 格式示例：01,32, ,26507000000098298342,361.07,20260226, ,1CA3
func ParseQRCodeData(data string) (invoiceNumber string, amount float64, date time.Time, err error) {
	parser := NewQRCodeParser()
	return parser.Parse(data)
}

type sourcer struct {
	name string
	keywords []string
}
var (
	sourceList = []sourcer{
		{
			name: "沃尔玛",
			keywords: []string{"沃尔玛", "Walmart", "walmart"},
		},
		{
			name: "星巴克",
			keywords: []string{"星巴克", "Starbucks", "starbucks"},
		},
		{
			name: "家乐福",
			keywords: []string{"家乐福", "Carrefour", "carrefour"},
		},
		{
			name: "大润发",
			keywords: []string{"大润发", "Dalat", "dalat"},
		},
		{
			name: "华润万家",
			keywords: []string{"华润万家", "Walmart", "walmart"},
		},
		{
			name: "永辉超市",
			keywords: []string{"永辉超市", "Yonghui", "yonghui"},
		},
		{
			name: "山姆会员店",
			keywords: []string{"山姆会员店", "Sam's Club", "sam's club"},
		},
		{
			name: "麦德龙",
			keywords: []string{"麦德龙", "Metro", "metro"},
		},
		{
			name: "盒马",
			keywords: []string{"盒马", "Hema", "hema"},
		},
		{
			name: "京东",
			keywords: []string{"京东", "JD", "jd"},
		},
		{
			name: "淘宝",
			keywords: []string{"淘宝", "Taobao", "taobao"},
		},
		{
			name: "天猫",
			keywords: []string{"天猫", "Tmall", "tmall"},
		},
		{
			name: "拼多多",
			keywords: []string{"拼多多", "Pinduoduo", "pinduoduo"},
		},
		{
			name: "唯品会",
			keywords: []string{"唯品会", "Vipshop", "vipshop"},	
		},
		{
			name: "苏宁易购",
			keywords: []string{"苏宁易购", "Suning", "suning"},
		},
		{
			name: "国美",
			keywords: []string{"国美", "Gome", "gome"},
		},
		{
			name:"餐饮",
			keywords: []string{"餐饮", "Catering", "catering"},
		},
		{
			name: "住宿",
			keywords: []string{"住宿", "Accommodation", "accommodation"},
		},
		{
			name: "交通",
			keywords: []string{"交通", "Transportation", "transportation"},
		},
		{
			name: "其他",
			keywords: []string{"其他", "Other", "other"},
		},
	}
)

// DetectSource 根据邮件主题或文件名检测发票来源
func DetectSource(subject, filename string) string {
	for _, source := range sourceList {
		if containsAny(subject, source.keywords) ||
			containsAny(filename, source.keywords) {
			return source.name
		}
	}
	return "其他"
}

// containsAny 检查字符串是否包含任意一个子串
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

// contains 简单的字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstr(s, substr))
}

func hasSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
