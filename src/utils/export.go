package utils

import (
	"batchinvoice-pdf/src/core"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"
)

// ExportToJSON 导出为JSON
func ExportToJSON(result *core.ExtractionResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

// ExportToHTML 导出为HTML
func ExportToHTML(result *core.ExtractionResult) ([]byte, error) {
	tmpl := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>发票二维码提取结果</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            margin: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        h1 {
            color: #667eea;
            text-align: center;
            margin-bottom: 10px;
        }
        .summary {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
        }
        .summary p {
            margin: 5px 0;
            font-size: 16px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #dee2e6;
        }
        th {
            background: #667eea;
            color: white;
            font-weight: 600;
        }
        tr:hover {
            background: #f8f9fa;
        }
        .amount-no-tax { color: #28a745; font-weight: bold; }
        .tax-amount { color: #fd7e14; font-weight: bold; }
        .total-amount { color: #dc3545; font-weight: bold; }
        .source-walmart { color: #0071ce; }
        .source-hema { color: #ff6a00; }
        .footer {
            text-align: center;
            color: #6c757d;
            font-size: 14px;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 发票二维码提取结果</h1>
        
        <div class="summary">
            <p><strong>提取时间：</strong>{{ .ExtractedAt }}</p>
            <p><strong>总邮件数：</strong>{{ .TotalEmails }}</p>
            <p><strong>处理PDF数：</strong>{{ .ProcessedPDFs }}</p>
            <p><strong>找到二维码：</strong>{{ .QRCodesFound }}</p>
        </div>

        <table>
            <thead>
                <tr>
                    <th>序号</th>
                    <th>姓名</th>
                    <th>发票号</th>
                    <th>来源</th>
                    <th>日期</th>
                    <th>不含税</th>
                    <th>税额</th>
                    <th>含税总额</th>
                </tr>
            </thead>
            <tbody>
                {{ range $index, $invoice := .Invoices }}
                <tr>
                    <td>{{ add $index 1 }}</td>
                    <td>{{ $invoice.Name }}</td>
                    <td>{{ $invoice.InvoiceNumber }}</td>
                    <td class="source-{{ lower $invoice.Source }}">{{ $invoice.Source }}</td>
                    <td>{{ formatDate $invoice.Date }}</td>
                    <td class="amount-no-tax">{{ formatCurrency $invoice.AmountNoTax }}</td>
                    <td class="tax-amount">{{ formatCurrency $invoice.TaxAmount }}</td>
                    <td class="total-amount">{{ formatCurrency $invoice.TotalAmount }}</td>
                </tr>
                {{ end }}
            </tbody>
        </table>

        <div class="footer">
            <p>由 BatchInvoice PDF 生成</p>
        </div>
    </div>
</body>
</html>`

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatCurrency": func(amount float64) string {
			return FormatCurrency(amount)
		},
		"lower": func(s string) string {
			return s // 简化版
		},
	}

	t, err := template.New("html").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, result); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportToMarkdown 导出为Markdown
func ExportToMarkdown(result *core.ExtractionResult) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("# 发票二维码提取结果\n\n")
	buf.WriteString(fmt.Sprintf("**提取时间：** %s\n\n", result.ExtractedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("- 总邮件数：%d\n", result.TotalEmails))
	buf.WriteString(fmt.Sprintf("- 处理PDF数：%d\n", result.ProcessedPDFs))
	buf.WriteString(fmt.Sprintf("- 找到二维码：%d\n\n", result.QRCodesFound))

	buf.WriteString("## 发票清单\n\n")
	buf.WriteString("| 序号 | 姓名 | 发票号 | 来源 | 日期 | 不含税 | 税额 | 含税总额 |\n")
	buf.WriteString("|------|------|--------|------|------|--------|------|----------|\n")

	for i, invoice := range result.Invoices {
		buf.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s | %s |\n",
			i+1,
			invoice.Name,
			invoice.InvoiceNumber,
			invoice.Source,
			invoice.Date.Format("2006-01-02"),
			FormatCurrency(invoice.AmountNoTax),
			FormatCurrency(invoice.TaxAmount),
			FormatCurrency(invoice.TotalAmount),
		))
	}

	return buf.Bytes(), nil
}
