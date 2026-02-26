package gui

import (
	"batchinvoice-pdf/src/core"
	"batchinvoice-pdf/src/utils"
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gen2brain/go-fitz"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainWindow 主窗口
type MainWindow struct {
	app           fyne.App
	window        fyne.Window
	config        *core.Config
	isConfigFixed bool // 标记配置是否已固定（从登录界面获取）

	// UI 组件
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
	invoiceList *widget.Table
	exportBtn   *widget.Button

	// 数据
	extractionResult *core.ExtractionResult
	qrCodeCache      map[string]fyne.CanvasObject // 缓存生成的二维码图片
}

// HoverTableCell 带鼠标悬停高亮效果的表格单元格
type HoverTableCell struct {
	widget.BaseWidget
	content fyne.CanvasObject
	bg      *canvas.Rectangle
}

func NewHoverTableCell() *HoverTableCell {
	cell := &HoverTableCell{}
	cell.ExtendBaseWidget(cell)
	return cell
}

// SetContent 设置单元格内部内容
func (c *HoverTableCell) SetContent(obj fyne.CanvasObject) {
	c.content = obj
	c.Refresh()
}

// 鼠标移入
func (c *HoverTableCell) MouseIn(_ *desktop.MouseEvent) {
	c.setHover(true)
}

// 鼠标移出
func (c *HoverTableCell) MouseOut() {
	c.setHover(false)
}

// 鼠标移动（这里不需要处理）
func (c *HoverTableCell) MouseMoved(_ *desktop.MouseEvent) {}

func (c *HoverTableCell) setHover(hover bool) {
	if c.bg == nil {
		return
	}

	if hover {
		c.bg.FillColor = theme.PrimaryColor()
		// 仅当内容是文本时才改成白色
		if txt, ok := c.content.(*canvas.Text); ok {
			txt.Color = color.White
		}
	} else {
		c.bg.FillColor = color.Transparent
		if txt, ok := c.content.(*canvas.Text); ok {
			txt.Color = theme.ForegroundColor()
		}
	}
	c.Refresh()
}

// CreateRenderer 创建渲染器
func (c *HoverTableCell) CreateRenderer() fyne.WidgetRenderer {
	c.bg = canvas.NewRectangle(color.Transparent)
	return &hoverTableCellRenderer{cell: c}
}

type hoverTableCellRenderer struct {
	cell *HoverTableCell
}

func (r *hoverTableCellRenderer) Layout(size fyne.Size) {
	if r.cell.bg != nil {
		r.cell.bg.Resize(size)
		r.cell.bg.Move(fyne.NewPos(0, 0))
	}
	if r.cell.content != nil {
		min := r.cell.content.MinSize()
		// 居中显示内容
		r.cell.content.Resize(min)
		r.cell.content.Move(fyne.NewPos(
			(size.Width-min.Width)/2,
			(size.Height-min.Height)/2,
		))
	}
}

func (r *hoverTableCellRenderer) MinSize() fyne.Size {
	if r.cell.content != nil {
		return r.cell.content.MinSize()
	}
	return fyne.NewSize(60, 24)
}

func (r *hoverTableCellRenderer) Refresh() {
	if r.cell.bg != nil {
		r.cell.bg.Refresh()
	}
	if r.cell.content != nil {
		r.cell.content.Refresh()
	}
}

func (r *hoverTableCellRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *hoverTableCellRenderer) Objects() []fyne.CanvasObject {
	var objs []fyne.CanvasObject
	if r.cell.bg != nil {
		objs = append(objs, r.cell.bg)
	}
	if r.cell.content != nil {
		objs = append(objs, r.cell.content)
	}
	return objs
}

func (r *hoverTableCellRenderer) Destroy() {}

// NewMainWindow 创建主窗口（兼容旧版本）
func NewMainWindow(app fyne.App) fyne.Window {
	return NewMainWindowWithConfig(app, nil)
}

// NewMainWindowWithConfig 使用配置创建主窗口
func NewMainWindowWithConfig(app fyne.App, config *core.Config) fyne.Window {
	mw := &MainWindow{
		app:           app,
		config:        config,
		isConfigFixed: config != nil,
	}

	// 如果没有提供配置，创建默认配置
	if mw.config == nil {
		mw.config = core.NewConfig()
	}

	mw.window = app.NewWindow("BatchInvoice PDF - 发票二维码批量提取工具")
	mw.window.Resize(fyne.NewSize(1200, 700))
	mw.window.CenterOnScreen()

	content := mw.buildUI()
	mw.window.SetContent(content)

	return mw.window
}

// buildUI 构建用户界面
func (mw *MainWindow) buildUI() fyne.CanvasObject {
	// 标题
	_ = widget.NewLabelWithStyle(
		"📧 发票二维码批量提取工具",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// 操作按钮
	extractBtn := widget.NewButton("🚀 开始提取", mw.onExtract)
	extractBtn.Importance = widget.HighImportance

	mw.exportBtn = widget.NewButton("📥 导出数据", mw.onExport)
	mw.exportBtn.Disable() // 初始禁用

	// 配置表单（内部会把最大邮件数、税率和按钮排成一行）
	configForm := mw.buildConfigForm(extractBtn, mw.exportBtn)

	// 进度条和状态
	mw.progressBar = widget.NewProgressBar()
	mw.progressBar.Hide()
	mw.statusLabel = widget.NewLabel("就绪")

	progressContainer := container.NewVBox(
		mw.statusLabel,
		mw.progressBar,
	)

	// 初始化二维码缓存
	mw.qrCodeCache = make(map[string]fyne.CanvasObject)

	// 发票表格
	mw.invoiceList = widget.NewTable(
		func() (int, int) {
			if mw.extractionResult == nil {
				return 1, 7 // 1行表头，7列
			}
			return len(mw.extractionResult.Invoices) + 1, 7 // +1是表头
		},
		func() fyne.CanvasObject {
			// 使用支持鼠标悬停高亮的自定义单元格
			return NewHoverTableCell()
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			cell := obj.(*HoverTableCell)

			// 表头
			if id.Row == 0 {
				headers := []string{"序号", "发票号码", "来源", "金额", "发票日期", "收件时间", "二维码"}
				text := canvas.NewText(headers[id.Col], theme.ForegroundColor())
				text.TextStyle = fyne.TextStyle{Bold: true}
				cell.SetContent(text)
				return
			}

			// 数据行
			if mw.extractionResult == nil || id.Row-1 >= len(mw.extractionResult.Invoices) {
				cell.SetContent(canvas.NewText("", theme.ForegroundColor()))
				return
			}

			invoice := mw.extractionResult.Invoices[id.Row-1]

			switch id.Col {
			case 0: // 序号
				cell.SetContent(canvas.NewText(fmt.Sprintf("%d", id.Row), theme.ForegroundColor()))
			case 1: // 发票号码
				cell.SetContent(canvas.NewText(invoice.InvoiceNumber, theme.ForegroundColor()))
			case 2: // 来源
				cell.SetContent(canvas.NewText(invoice.Source, theme.ForegroundColor()))
			case 3: // 金额
				cell.SetContent(canvas.NewText(fmt.Sprintf("¥%.2f", invoice.TotalAmount), theme.ForegroundColor()))
			case 4: // 发票日期
				cell.SetContent(canvas.NewText(invoice.Date.Format("2006-01-02"), theme.ForegroundColor()))
			case 5: // 收件时间
				cell.SetContent(canvas.NewText(invoice.EmailDate, theme.ForegroundColor()))
			case 6: // 二维码
				// 从缓存获取或生成二维码
				qrKey := invoice.InvoiceNumber
				qrImg, exists := mw.qrCodeCache[qrKey]
				if !exists {
					var err error
					qrImg, err = generateQRCodeImage(invoice.QRCodeData, 80, 80)
					if err != nil {
						qrImg = widget.NewLabel("二维码")
					} else {
						mw.qrCodeCache[qrKey] = qrImg
					}
				}
				cell.SetContent(qrImg)
			}
		},
	)

	// 设置列宽
	mw.invoiceList.SetColumnWidth(0, 60)  // 序号
	mw.invoiceList.SetColumnWidth(1, 200) // 发票号码
	mw.invoiceList.SetColumnWidth(2, 120) // 来源
	mw.invoiceList.SetColumnWidth(3, 100) // 金额
	mw.invoiceList.SetColumnWidth(4, 100) // 发票日期
	mw.invoiceList.SetColumnWidth(5, 180) // 收件时间
	mw.invoiceList.SetColumnWidth(6, 100) // 二维码

	// 表格点击事件
	mw.invoiceList.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 && mw.extractionResult != nil && id.Row-1 < len(mw.extractionResult.Invoices) {
			mw.showInvoiceDetail(mw.extractionResult.Invoices[id.Row-1])
		}
	}

	resultContainer := container.NewBorder(
		nil, //widget.NewLabel("📋 发票列表"),
		nil,
		nil,
		nil,
		container.NewScroll(mw.invoiceList),
	)

	// 状态栏（底部显示状态和进度）
	statusBar := progressContainer

	// 组合布局
	mainContent := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			configForm,
			widget.NewSeparator(),
		),
		statusBar,
		nil,
		nil,
		resultContainer,
	)

	return mainContent
}

// buildConfigForm 构建配置表单
func (mw *MainWindow) buildConfigForm(extractBtn, exportBtn *widget.Button) fyne.CanvasObject {
	// 如果配置已固定（从登录界面获取），显示简化表单
	if mw.isConfigFixed {
		return mw.buildSimpleConfigForm(extractBtn, exportBtn)
	}

	// 否则显示完整配置表单
	return mw.buildFullConfigForm(extractBtn, exportBtn)
}

// buildSimpleConfigForm 构建简化配置表单（登录后）
func (mw *MainWindow) buildSimpleConfigForm(extractBtn, exportBtn *widget.Button) fyne.CanvasObject {
	// 显示已登录的邮箱信息
	accountInfo := widget.NewLabelWithStyle(
		fmt.Sprintf("✓ 已登录: %s", mw.config.Username),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// 最大邮件数
	maxEmailsEntry := widget.NewEntry()
	maxEmailsEntry.SetText(fmt.Sprintf("%d", mw.config.MaxEmails))
	maxEmailsEntry.OnChanged = func(s string) {
		if max, err := strconv.Atoi(s); err == nil {
			mw.config.MaxEmails = max
		}
	}

	// 税率
	taxRateEntry := widget.NewEntry()
	taxRateEntry.SetText(fmt.Sprintf("%.2f", mw.config.TaxRate))
	taxRateEntry.OnChanged = func(s string) {
		if rate, err := strconv.ParseFloat(s, 64); err == nil {
			mw.config.TaxRate = rate
		}
	}

	// 调整输入框宽度，使其更长一些
	maxEmailsEntryContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(120, 36)),
		maxEmailsEntry,
	)
	taxRateEntryContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(120, 36)),
		taxRateEntry,
	)

	// 控件行：最大邮件数 + 税率 + 开始提取 + 导出 全部一行
	controlsRow := container.NewHBox(
		widget.NewLabel("最大邮件数"),
		maxEmailsEntryContainer,
		widget.NewLabel("税率(如0.13)"),
		taxRateEntryContainer,
		layout.NewSpacer(),
		extractBtn,
		exportBtn,
	)

	form := container.NewVBox(
		accountInfo,
		widget.NewSeparator(),
		controlsRow,
	)

	return form
}

// buildFullConfigForm 构建完整配置表单（未登录）
func (mw *MainWindow) buildFullConfigForm(extractBtn, exportBtn *widget.Button) fyne.CanvasObject {
	// IMAP服务器
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("例如: imap.gmail.com")
	serverEntry.OnChanged = func(s string) {
		mw.config.IMAPServer = s
	}

	// 端口
	portEntry := widget.NewEntry()
	portEntry.SetText("993")
	portEntry.OnChanged = func(s string) {
		if port, err := strconv.Atoi(s); err == nil {
			mw.config.Port = port
		}
	}

	// 用户名
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("your-email@example.com")
	usernameEntry.OnChanged = func(s string) {
		mw.config.Username = s
	}

	// 密码
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("邮箱密码或应用专用密码")
	passwordEntry.OnChanged = func(s string) {
		mw.config.Password = s
	}

	// 最大邮件数
	maxEmailsEntry := widget.NewEntry()
	maxEmailsEntry.SetText("20")
	maxEmailsEntry.OnChanged = func(s string) {
		if max, err := strconv.Atoi(s); err == nil {
			mw.config.MaxEmails = max
		}
	}

	// 税率
	taxRateEntry := widget.NewEntry()
	taxRateEntry.SetText("0.13")
	taxRateEntry.OnChanged = func(s string) {
		if rate, err := strconv.ParseFloat(s, 64); err == nil {
			mw.config.TaxRate = rate
		}
	}

	// 调整输入框宽度
	maxEmailsEntryContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(120, 0)),
		maxEmailsEntry,
	)
	taxRateEntryContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(120, 0)),
		taxRateEntry,
	)

	// 顶部是基础邮箱配置表单，下面一行放 最大邮件数 + 税率 + 开始提取 + 导出
	baseForm := widget.NewForm(
		widget.NewFormItem("IMAP服务器", serverEntry),
		widget.NewFormItem("端口", portEntry),
		widget.NewFormItem("用户名", usernameEntry),
		widget.NewFormItem("密码", passwordEntry),
	)

	controlsRow := container.NewHBox(
		widget.NewLabel("最大邮件数"),
		maxEmailsEntryContainer,
		widget.NewLabel("税率(如0.13)"),
		taxRateEntryContainer,
		layout.NewSpacer(),
		extractBtn,
		exportBtn,
	)

	form := container.NewVBox(
		baseForm,
		controlsRow,
	)

	return form
}

// onExtract 提取按钮点击事件
func (mw *MainWindow) onExtract() {
	// 验证配置
	if mw.config.IMAPServer == "" || mw.config.Username == "" || mw.config.Password == "" {
		dialog.ShowError(fmt.Errorf("请填写完整的邮箱配置"), mw.window)
		return
	}

	// 重置进度和状态
	mw.progressBar.SetValue(0)
	mw.progressBar.Show()
	mw.statusLabel.SetText("正在连接邮箱...")
	mw.extractionResult = &core.ExtractionResult{}
	mw.invoiceList.Refresh()
	mw.exportBtn.Disable()

	// 在后台执行提取
	go func() {
		defer func() {
			// 在主线程中隐藏进度条
			fyne.Do(func() {
				mw.progressBar.Hide()
			})
		}()

		result, err := mw.extractInvoices()
		if err != nil {
			// 在主线程中显示错误对话框
			fyne.Do(func() {
				mw.statusLabel.SetText(fmt.Sprintf("提取失败: %v", err))
				dialog.ShowError(err, mw.window)
			})
			return
		}

		// 在主线程中显示结果
		fyne.Do(func() {
			// 合并最终统计结果，保持列表为增量渲染的结果
			if mw.extractionResult == nil {
				mw.extractionResult = result
			} else {
				mw.extractionResult.TotalEmails = result.TotalEmails
				mw.extractionResult.ProcessedPDFs = result.ProcessedPDFs
				mw.extractionResult.QRCodesFound = result.QRCodesFound
				mw.extractionResult.ExtractedAt = result.ExtractedAt
				if len(mw.extractionResult.Invoices) == 0 {
					mw.extractionResult.Invoices = result.Invoices
				}
			}

			// 设置所有行的高度以适应二维码
			for i := 0; i <= len(mw.extractionResult.Invoices); i++ {
				mw.invoiceList.SetRowHeight(i, 90)
			}
			mw.invoiceList.Refresh()
			mw.exportBtn.Enable()
			mw.statusLabel.SetText(fmt.Sprintf(
				"✅ 提取完成！共 %d 封邮件，%d 个PDF，找到 %d 个二维码，识别 %d 张发票",
				mw.extractionResult.TotalEmails,
				mw.extractionResult.ProcessedPDFs,
				mw.extractionResult.QRCodesFound,
				len(mw.extractionResult.Invoices),
			))

			if len(mw.extractionResult.Invoices) > 0 {
				dialog.ShowInformation("提取完成",
					fmt.Sprintf("成功提取 %d 张发票！\n点击列表查看详情", len(mw.extractionResult.Invoices)),
					mw.window)
			}
		})
	}()
}

// extractInvoices 执行发票提取
func (mw *MainWindow) extractInvoices() (*core.ExtractionResult, error) {
	// 创建邮件客户端
	emailClient := core.NewEmailClient(mw.config)

	log.Println("========================================")
	log.Printf("开始提取发票 - %s", mw.config.Username)
	log.Println("========================================")

	// 更新状态
	fyne.Do(func() {
		mw.statusLabel.SetText("正在连接邮箱...")
		mw.progressBar.SetValue(0.1)
	})

	log.Printf("正在连接邮箱服务器: %s:%d", mw.config.IMAPServer, mw.config.Port)

	// 连接
	if err := emailClient.Connect(); err != nil {
		log.Printf("❌ 连接失败: %v", err)
		return nil, fmt.Errorf("连接邮箱失败: %w", err)
	}
	defer emailClient.Disconnect()

	log.Println("✅ 邮箱连接成功")

	// 更新状态
	fyne.Do(func() {
		mw.statusLabel.SetText("正在获取邮件列表...")
		mw.progressBar.SetValue(0.2)
	})

	log.Printf("正在获取最新 %d 封邮件（从新到旧）...", mw.config.MaxEmails)

	// 获取邮件
	messages, err := emailClient.FetchRecentEmails(mw.config.MaxEmails)
	if err != nil {
		log.Printf("❌ 获取邮件失败: %v", err)
		return nil, fmt.Errorf("获取邮件失败: %w", err)
	}

	log.Printf("✅ 找到 %d 封邮件（已按时间从新到旧排序）", len(messages))

	// 更新状态
	fyne.Do(func() {
		mw.statusLabel.SetText(fmt.Sprintf("找到 %d 封邮件，开始提取发票...", len(messages)))
		mw.progressBar.SetValue(0.3)
	})

	// 处理每封邮件
	var invoices []core.Invoice
	processedPDFs := 0
	totalQRCodes := 0
	totalPDFs := 0

	// 先统计总PDF数量
	for _, msg := range messages {
		totalPDFs += len(msg.Attachments)
	}

	log.Printf("统计完成：共 %d 个PDF附件需要处理", totalPDFs)
	log.Println("开始处理邮件（从最新到最旧）...")
	log.Println("----------------------------------------")

	currentPDF := 0
	for msgIdx, msg := range messages {
		log.Printf("[邮件 %d/%d] 主题: %s", msgIdx+1, len(messages), msg.Subject)
		for _, attachment := range msg.Attachments {
			currentPDF++
			processedPDFs++

			log.Printf("  [PDF %d/%d] 处理文件: %s (%.2f KB)", currentPDF, totalPDFs, attachment.Filename, float64(len(attachment.Data))/1024)

			// 更新进度
			progress := 0.3 + (0.6 * float64(currentPDF) / float64(totalPDFs))
			fyne.Do(func() {
				mw.progressBar.SetValue(progress)
				mw.statusLabel.SetText(fmt.Sprintf(
					"处理中... 邮件 %d/%d, PDF %d/%d",
					msgIdx+1, len(messages), currentPDF, totalPDFs,
				))
			})

			// 提取二维码
			qrcodes, err := core.ExtractQRCodesFromPDF(attachment.Data)
			if err != nil {
				log.Printf("  ❌ 提取二维码失败: %v", err)
				continue
			}

			totalQRCodes += len(qrcodes)
			if len(qrcodes) > 0 {
				log.Printf("  ✅ 找到 %d 个二维码", len(qrcodes))
			} else {
				log.Printf("  ⚠️  未找到二维码")
			}

			// 处理每个二维码
			for qrIdx, qrData := range qrcodes {
				// 解析二维码数据
				invoiceNum, amount, date, err := core.ParseQRCodeData(qrData)
				if err != nil {
					log.Printf("    [二维码 %d] ❌ 解析失败: %v", qrIdx+1, err)
					continue
				}

				// 计算税额
				amountNoTax, taxAmount := utils.CalculateTax(amount, mw.config.TaxRate)

				// 创建发票记录
				source := core.DetectSource(msg.Subject, attachment.Filename)
				invoice := core.Invoice{
					Name:          extractNameFromSubject(msg.Subject),
					InvoiceNumber: invoiceNum,
					QRCodeData:    qrData,
					Source:        source,
					Date:          date,
					TotalAmount:   amount,
					AmountNoTax:   amountNoTax,
					TaxAmount:     taxAmount,
					EmailSubject:  msg.Subject,
					FileName:      attachment.Filename,
					EmailDate:     msg.Date,
					PDFData:       attachment.Data,
				}

				invoices = append(invoices, invoice)
				currentIndex := len(invoices)
				log.Printf("    [发票 %d] ✅ %s | %s | ¥%.2f | %s",
					currentIndex, invoiceNum, source, amount, date.Format("2006-01-02"))

				// 增量更新界面：每解析一张发票就立刻添加到列表中
				fyne.Do(func() {
					if mw.extractionResult == nil {
						mw.extractionResult = &core.ExtractionResult{}
					}
					mw.extractionResult.Invoices = append(mw.extractionResult.Invoices, invoice)
					// 行 0 为表头，从 1 开始设置行高
					mw.invoiceList.SetRowHeight(currentIndex, 90)
					mw.invoiceList.Refresh()
				})
			}
		}
	}

	// 更新最终进度
	fyne.Do(func() {
		mw.progressBar.SetValue(0.9)
		mw.statusLabel.SetText("正在整理数据...")
	})

	log.Println("----------------------------------------")
	log.Println("📊 提取统计：")
	log.Printf("  总邮件数: %d", len(messages))
	log.Printf("  处理PDF数: %d", processedPDFs)
	log.Printf("  找到二维码: %d", totalQRCodes)
	log.Printf("  识别发票: %d", len(invoices))

	result := &core.ExtractionResult{
		TotalEmails:   len(messages),
		ProcessedPDFs: processedPDFs,
		QRCodesFound:  totalQRCodes,
		Invoices:      invoices,
		ExtractedAt:   time.Now(),
	}

	// 完成
	fyne.Do(func() {
		mw.progressBar.SetValue(1.0)
	})

	if len(invoices) > 0 {
		log.Println("========================================")
		log.Println("✅ 提取完成！发票列表：")
		log.Println("========================================")
		for i, inv := range invoices {
			log.Printf("[%d] %s | %s | ¥%.2f | %s | %s",
				i+1,
				inv.InvoiceNumber,
				inv.Source,
				inv.TotalAmount,
				inv.Date.Format("2006-01-02"),
				inv.FileName,
			)
		}
		log.Println("========================================")
	} else {
		log.Println("⚠️  未找到任何发票")
	}

	return result, nil
}

// showResult 显示结果
func (mw *MainWindow) showResult(result *core.ExtractionResult) {
	msg := fmt.Sprintf(
		"✅ 提取完成!\n\n"+
			"总邮件数: %d\n"+
			"处理PDF数: %d\n"+
			"找到二维码: %d\n"+
			"识别发票: %d",
		result.TotalEmails,
		result.ProcessedPDFs,
		result.QRCodesFound,
		len(result.Invoices),
	)

	dialog.ShowInformation("提取结果", msg, mw.window)
}

// showInvoiceDetail 显示发票详情
func (mw *MainWindow) showInvoiceDetail(invoice core.Invoice) {
	// 创建详情窗口
	detailWindow := mw.app.NewWindow(fmt.Sprintf("发票详情 - %s", invoice.InvoiceNumber))
	detailWindow.Resize(fyne.NewSize(1400, 800))
	detailWindow.CenterOnScreen()

	// 发票基本信息
	infoText := fmt.Sprintf(
		"发票号码：%s\n"+
			"来源：%s\n"+
			"日期：%s\n"+
			"含税金额：¥%.2f\n"+
			"不含税金额：¥%.2f\n"+
			"税额：¥%.2f\n"+
			"邮件收件时间：%s\n"+
			"邮件主题：%s\n"+
			"文件名：%s",
		invoice.InvoiceNumber,
		invoice.Source,
		invoice.Date.Format("2006-01-02"),
		invoice.TotalAmount,
		invoice.AmountNoTax,
		invoice.TaxAmount,
		invoice.EmailDate,
		invoice.EmailSubject,
		invoice.FileName,
	)

	infoLabel := widget.NewLabel(infoText)
	infoLabel.Wrapping = fyne.TextWrapWord

	// 二维码数据
	qrCodeEntry := widget.NewMultiLineEntry()
	qrCodeEntry.SetText(invoice.QRCodeData)
	qrCodeEntry.Wrapping = fyne.TextWrapWord
	qrCodeEntry.Disable()

	// 生成二维码图片
	qrImage, err := generateQRCodeImage(invoice.QRCodeData, 200, 200)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		qrImage = widget.NewLabel("二维码生成失败")
	}

	// 生成PDF预览图片
	pdfPreview := generatePDFPreview(invoice.PDFData)

	// 关闭按钮
	closeBtn := widget.NewButton("关闭", func() {
		detailWindow.Close()
	})

	// 复制二维码数据按钮
	copyBtn := widget.NewButton("复制二维码数据", func() {
		mw.window.Clipboard().SetContent(invoice.QRCodeData)
		dialog.ShowInformation("已复制", "二维码数据已复制到剪贴板", detailWindow)
	})

	buttonBox := container.NewHBox(copyBtn, closeBtn)

	// 布局 - 左右分栏
	leftPanel := container.NewVBox(
		widget.NewLabel("发票信息："),
		infoLabel,
		widget.NewSeparator(),
		widget.NewLabel("二维码图片："),
		qrImage,
		widget.NewSeparator(),
		widget.NewLabel("二维码数据："),
		qrCodeEntry,
	)

	openPDFBtn := widget.NewButton("在系统中打开PDF", func() {
		openPDFWithDefaultViewer(invoice.PDFData)
	})

	rightPanel := container.NewVBox(
		widget.NewLabel("发票预览："),
		pdfPreview,
		openPDFBtn,
	)

	mainContent := container.NewHSplit(
		container.NewScroll(leftPanel),
		container.NewScroll(rightPanel),
	)
	mainContent.SetOffset(0.4)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(
				fmt.Sprintf("📄 发票详情 - %s", invoice.InvoiceNumber),
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			buttonBox,
		),
		nil,
		nil,
		mainContent,
	)

	detailWindow.SetContent(container.NewPadded(content))
	detailWindow.Show()
}

// openPDFWithDefaultViewer 使用系统默认PDF查看器打开PDF
func openPDFWithDefaultViewer(pdfData []byte) {
	if len(pdfData) == 0 {
		log.Println("no PDF data to open")
		return
	}

	tmpFile, err := os.CreateTemp("", "invoice-*.pdf")
	if err != nil {
		log.Printf("failed to create temp pdf file: %v", err)
		return
	}
	if _, err := tmpFile.Write(pdfData); err != nil {
		log.Printf("failed to write temp pdf file: %v", err)
		_ = tmpFile.Close()
		return
	}
	if err := tmpFile.Close(); err != nil {
		log.Printf("failed to close temp pdf file: %v", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", tmpFile.Name())
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", tmpFile.Name())
	default:
		cmd = exec.Command("xdg-open", tmpFile.Name())
	}

	if err := cmd.Start(); err != nil {
		log.Printf("failed to open pdf with default viewer: %v", err)
	}
}

// onExport 导出按钮点击事件
func (mw *MainWindow) onExport() {
	if mw.extractionResult == nil || len(mw.extractionResult.Invoices) == 0 {
		dialog.ShowInformation("导出", "当前没有可导出的发票数据，请先完成一次提取。", mw.window)
		return
	}

	formats := []string{"JSON", "HTML", "Markdown"}
	formatSelector := widget.NewRadioGroup(formats, nil)
	formatSelector.SetSelected("JSON")

	content := container.NewVBox(
		widget.NewLabel("请选择导出格式："),
		formatSelector,
	)

	dialog.NewCustomConfirm("导出数据", "下一步", "取消", content, func(confirmed bool) {
		if !confirmed {
			return
		}

		format := formatSelector.Selected
		if format == "" {
			format = "JSON"
		}

		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			if writer == nil {
				// 用户取消
				return
			}
			defer writer.Close()

			var data []byte
			switch format {
			case "JSON":
				data, err = utils.ExportToJSON(mw.extractionResult)
			case "HTML":
				data, err = utils.ExportToHTML(mw.extractionResult)
			case "Markdown":
				data, err = utils.ExportToMarkdown(mw.extractionResult)
			default:
				err = fmt.Errorf("不支持的导出格式: %s", format)
			}
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			if _, err := writer.Write(data); err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			dialog.ShowInformation("导出成功", fmt.Sprintf("数据已导出到：\n%s", writer.URI().Path()), mw.window)
		}, mw.window)
	}, mw.window).Show()
}

// extractNameFromSubject 从邮件主题提取姓名
func extractNameFromSubject(subject string) string {
	// TODO: 实现姓名提取逻辑
	return "未知"
}

// generateQRCodeImage 生成二维码图片
func generateQRCodeImage(data string, width, height int) (fyne.CanvasObject, error) {
	// 使用 core.GenerateQRCode 生成二维码图片
	imgBytes, err := core.GenerateQRCode(data, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// 创建 Fyne 图片对象
	img := canvas.NewImageFromResource(fyne.NewStaticResource("qrcode.png", imgBytes))
	img.FillMode = canvas.ImageFillOriginal
	img.SetMinSize(fyne.NewSize(float32(width), float32(height)))

	return img, nil
}

// generatePDFPreview 生成PDF预览图片
func generatePDFPreview(pdfData []byte) fyne.CanvasObject {
	if len(pdfData) == 0 {
		return widget.NewLabel("无PDF数据")
	}

	// 将 PDF 写入临时文件，供 go-fitz 打开
	tmpFile, err := os.CreateTemp("", "invoice-preview-*.pdf")
	if err != nil {
		log.Printf("Failed to create temp pdf for preview: %v", err)
		return widget.NewLabel("PDF预览失败（创建临时文件失败）")
	}
	if _, err := tmpFile.Write(pdfData); err != nil {
		log.Printf("Failed to write temp pdf for preview: %v", err)
		_ = tmpFile.Close()
		return widget.NewLabel("PDF预览失败（写入临时文件失败）")
	}
	if err := tmpFile.Close(); err != nil {
		log.Printf("Failed to close temp pdf for preview: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 使用 go-fitz 渲染整页 PDF 为图像
	doc, err := fitz.New(tmpFile.Name())
	if err != nil {
		log.Printf("Failed to open pdf with go-fitz: %v", err)
		return widget.NewLabel("PDF预览失败（无法打开文件）")
	}
	defer doc.Close()

	if doc.NumPage() == 0 {
		return widget.NewLabel("PDF没有任何页面")
	}

	firstPage, err := doc.Image(0)
	if err != nil {
		log.Printf("Failed to render pdf page: %v", err)
		return widget.NewLabel("PDF预览失败（渲染页面出错）")
	}

	// 将整页 image.Image 转换为 PNG 字节
	var buf bytes.Buffer
	if err := png.Encode(&buf, firstPage); err != nil {
		log.Printf("Failed to encode image: %v", err)
		return widget.NewLabel("图像编码失败")
	}

	// 创建Fyne图片对象
	img := canvas.NewImageFromResource(fyne.NewStaticResource("invoice.png", buf.Bytes()))
	img.FillMode = canvas.ImageFillContain

	// 设置合适的显示尺寸
	bounds := firstPage.Bounds()
	aspectRatio := float32(bounds.Dx()) / float32(bounds.Dy())
	displayHeight := float32(600)
	displayWidth := displayHeight * aspectRatio
	img.SetMinSize(fyne.NewSize(displayWidth, displayHeight))

	return img
}
