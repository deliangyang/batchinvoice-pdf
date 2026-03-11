package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// QRCodeParser 二维码数据解析器
type QRCodeParser struct{}

// NewQRCodeParser 创建二维码解析器
func NewQRCodeParser() *QRCodeParser {
	return &QRCodeParser{}
}

// Parse 解析二维码数据
// 格式示例：01,32, ,26507000000098298342,361.07,20260226, ,1CA3
// 字段说明：
// - 01: 类型码
// - 32: 税务局代码
// - 空: 预留
// - 26507000000098298342: 发票号码
// - 361.07: 金额
// - 20260226: 日期(YYYYMMDD)
// - 空: 预留
// - 1CA3: 校验码
func (p *QRCodeParser) Parse(data string) (invoiceNumber string, amount float64, date time.Time, err error) {
	// 分割数据
	parts := strings.Split(data, ",")

	if len(parts) < 6 {
		return "", 0, time.Time{}, fmt.Errorf("invalid QR code format: not enough fields")
	}

	// 提取发票号码 (第4个字段，索引3)
	if len(parts) > 3 {
		invoiceNumber = strings.TrimSpace(parts[3])
	}

	// 提取金额 (第5个字段，索引4)
	if len(parts) > 4 {
		amountStr := strings.TrimSpace(parts[4])
		amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return "", 0, time.Time{}, fmt.Errorf("invalid amount: %w", err)
		}
	}

	// 提取日期 (第6个字段，索引5)
	if len(parts) > 5 {
		dateStr := strings.TrimSpace(parts[5])
		date, err = parseDate(dateStr)
		if err != nil {
			return "", 0, time.Time{}, fmt.Errorf("invalid date: %w", err)
		}
	}

	return invoiceNumber, amount, date, nil
}

// ParseWithTaxBureau 解析二维码数据（包含税务局代码）
// 返回：发票号码、金额、日期、税务局代码、错误
func (p *QRCodeParser) ParseWithTaxBureau(data string) (invoiceNumber string, amount float64, date time.Time, taxBureauCode string, err error) {
	// 分割数据
	parts := strings.Split(data, ",")

	if len(parts) < 6 {
		return "", 0, time.Time{}, "", fmt.Errorf("invalid QR code format: not enough fields")
	}

	// 提取税务局代码 (第2个字段，索引1)
	if len(parts) > 1 {
		taxBureauCode = strings.TrimSpace(parts[1])
	}

	// 提取发票号码 (第4个字段，索引3)
	if len(parts) > 3 {
		invoiceNumber = strings.TrimSpace(parts[3])
	}

	// 提取金额 (第5个字段，索引4)
	if len(parts) > 4 {
		amountStr := strings.TrimSpace(parts[4])
		amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return "", 0, time.Time{}, "", fmt.Errorf("invalid amount: %w", err)
		}
	}

	// 提取日期 (第6个字段，索引5)
	if len(parts) > 5 {
		dateStr := strings.TrimSpace(parts[5])
		date, err = parseDate(dateStr)
		if err != nil {
			return "", 0, time.Time{}, "", fmt.Errorf("invalid date: %w", err)
		}
	}

	return invoiceNumber, amount, date, taxBureauCode, nil
}

// parseDate 解析日期字符串
// 支持格式：YYYYMMDD, YYYY-MM-DD
func parseDate(dateStr string) (time.Time, error) {
	// 尝试 YYYYMMDD 格式
	if len(dateStr) == 8 && isNumeric(dateStr) {
		return time.Parse("20060102", dateStr)
	}

	// 尝试 YYYY-MM-DD 格式
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr); matched {
		return time.Parse("2006-01-02", dateStr)
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

// isNumeric 检查字符串是否为纯数字
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ExtractInvoiceNumber 从字符串中提取发票号码
func ExtractInvoiceNumber(s string) string {
	// 发票号码通常是一串数字，长度在12-20位之间
	re := regexp.MustCompile(`\d{12,20}`)
	matches := re.FindString(s)
	return matches
}
