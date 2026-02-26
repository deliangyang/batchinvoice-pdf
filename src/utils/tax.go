package utils

import (
	"fmt"
	"math"
)

// CalculateTax 计算税额
// taxRate: 税率（如0.13表示13%）
// totalAmount: 含税总额
// 返回：不含税金额、税额
func CalculateTax(totalAmount float64, taxRate float64) (amountNoTax, taxAmount float64) {
	// 不含税金额 = 含税总额 / (1 + 税率)
	amountNoTax = totalAmount / (1 + taxRate)
	
	// 税额 = 含税总额 - 不含税金额
	taxAmount = totalAmount - amountNoTax
	
	// 四舍五入到分
	amountNoTax = roundToDecimal(amountNoTax, 2)
	taxAmount = roundToDecimal(taxAmount, 2)
	
	return amountNoTax, taxAmount
}

// roundToDecimal 四舍五入到指定小数位
func roundToDecimal(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

// FormatCurrency 格式化货币
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("¥%.2f", amount)
}