package core

import (
	"time"
)

// Invoice 发票信息结构
type Invoice struct {
	Name          string    `json:"name"`            // 姓名
	InvoiceNumber string    `json:"invoice_number"`  // 发票号码
	QRCodeData    string    `json:"qrcode_data"`     // 二维码原始数据
	Source        string    `json:"source"`          // 来源（如：沃尔玛、盒马）
	Date          time.Time `json:"date"`            // 日期
	TotalAmount   float64   `json:"total_amount"`    // 含税总额
	AmountNoTax   float64   `json:"amount_no_tax"`   // 不含税金额
	TaxAmount     float64   `json:"tax_amount"`      // 税额
	TaxBureauCode string    `json:"tax_bureau_code"` // 税务局代码
	TaxBureauName string    `json:"tax_bureau_name"` // 税务局名称
	EmailSubject  string    `json:"email_subject"`   // 邮件主题
	FileName      string    `json:"file_name"`       // PDF文件名
	EmailDate     string    `json:"email_date"`      // 邮件收件时间
	PDFData       []byte    `json:"-"`               // PDF原始数据（不导出到JSON）
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
		MaxEmails: 30,
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
	name     string
	keywords []string
}

var (
	sourceList = []sourcer{
		{
			name:     "沃尔玛",
			keywords: []string{"沃尔玛", "Walmart", "walmart"},
		},
		{
			name:     "星巴克",
			keywords: []string{"星巴克", "Starbucks", "starbucks"},
		},
		{
			name:     "家乐福",
			keywords: []string{"家乐福", "Carrefour", "carrefour"},
		},
		{
			name:     "大润发",
			keywords: []string{"大润发", "Dalat", "dalat"},
		},
		{
			name:     "华润万家",
			keywords: []string{"华润万家", "Walmart", "walmart"},
		},
		{
			name:     "永辉超市",
			keywords: []string{"永辉超市", "Yonghui", "yonghui"},
		},
		{
			name:     "山姆会员店",
			keywords: []string{"山姆会员店", "Sam's Club", "sam's club"},
		},
		{
			name:     "麦德龙",
			keywords: []string{"麦德龙", "Metro", "metro"},
		},
		{
			name:     "盒马",
			keywords: []string{"盒马", "Hema", "hema"},
		},
		{
			name:     "京东",
			keywords: []string{"京东", "JD", "jd", "jd.com", "JD.COM", "京东商城", "京东集团", "jingdong", "Jingdong"},
		},
		{
			name:     "淘宝",
			keywords: []string{"淘宝", "Taobao", "taobao"},
		},
		{
			name:     "天猫",
			keywords: []string{"天猫", "Tmall", "tmall"},
		},
		{
			name:     "拼多多",
			keywords: []string{"拼多多", "Pinduoduo", "pinduoduo"},
		},
		{
			name:     "唯品会",
			keywords: []string{"唯品会", "Vipshop", "vipshop"},
		},
		{
			name:     "苏宁易购",
			keywords: []string{"苏宁易购", "Suning", "suning"},
		},
		{
			name:     "国美",
			keywords: []string{"国美", "Gome", "gome"},
		},
		{
			name:     "餐饮",
			keywords: []string{"餐饮", "Catering", "catering"},
		},
		{
			name:     "住宿",
			keywords: []string{"住宿", "Accommodation", "accommodation"},
		},
		{
			name:     "交通",
			keywords: []string{"交通", "Transportation", "transportation"},
		},
		{
			name:     "其他",
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

// 税务局代码到名称的映射
var taxBureauMap = map[string]string{
	"01": "国家税务总局北京市税务局",
	"02": "国家税务总局天津市税务局",
	"03": "国家税务总局河北省税务局",
	"04": "国家税务总局山西省税务局",
	"05": "国家税务总局内蒙古自治区税务局",
	"06": "国家税务总局辽宁省税务局",
	"07": "国家税务总局吉林省税务局",
	"08": "国家税务总局黑龙江省税务局",
	"09": "国家税务总局上海市税务局",
	"10": "国家税务总局江苏省税务局",
	"11": "国家税务总局浙江省税务局",
	"12": "国家税务总局安徽省税务局",
	"13": "国家税务总局福建省税务局",
	"14": "国家税务总局江西省税务局",
	"15": "国家税务总局山东省税务局",
	"16": "国家税务总局河南省税务局",
	"17": "国家税务总局湖北省税务局",
	"18": "国家税务总局湖南省税务局",
	"19": "国家税务总局广东省税务局",
	"20": "国家税务总局广西壮族自治区税务局",
	"21": "国家税务总局海南省税务局",
	"22": "国家税务总局重庆市税务局",
	"23": "国家税务总局四川省税务局",
	"24": "国家税务总局贵州省税务局",
	"25": "国家税务总局云南省税务局",
	"26": "国家税务总局西藏自治区税务局",
	"27": "国家税务总局陕西省税务局",
	"28": "国家税务总局甘肃省税务局",
	"29": "国家税务总局青海省税务局",
	"30": "国家税务总局宁夏回族自治区税务局",
	"31": "国家税务总局新疆维吾尔自治区税务局",
	"32": "国家税务总局大连市税务局",
	"33": "国家税务总局宁波市税务局",
	"34": "国家税务总局厦门市税务局",
	"35": "国家税务总局青岛市税务局",
	"36": "国家税务总局深圳市税务局",
}

// GetTaxBureauName 根据税务局代码获取税务局名称
func GetTaxBureauName(code string) string {
	if name, ok := taxBureauMap[code]; ok {
		return name
	}
	// 如果找不到对应的税务局名称，返回代码本身
	if code != "" {
		return "税务局代码:" + code
	}
	return "未知"
}
