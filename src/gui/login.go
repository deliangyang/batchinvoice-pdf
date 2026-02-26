package gui

import (
	"batchinvoice-pdf/src/core"
	"batchinvoice-pdf/src/utils"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// LoginWindow 登录窗口
type LoginWindow struct {
	app     fyne.App
	window  fyne.Window
	email   string
	onLogin func(*core.Config)
}

// NewLoginWindow 创建登录窗口
func NewLoginWindow(app fyne.App, email string, onLogin func(*core.Config)) fyne.Window {
	lw := &LoginWindow{
		app:     app,
		email:   email,
		onLogin: onLogin,
	}

	lw.window = app.NewWindow("BatchInvoice PDF - 登录")
	lw.window.Resize(fyne.NewSize(450, 300))
	lw.window.CenterOnScreen()

	content := lw.buildLoginUI()
	lw.window.SetContent(content)

	return lw.window
}

// buildLoginUI 构建登录界面
func (lw *LoginWindow) buildLoginUI() fyne.CanvasObject {
	// 标题
	title := widget.NewLabelWithStyle(
		"📧 发票二维码批量提取工具",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitle := widget.NewLabelWithStyle(
		"请输入邮箱密码登录",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	// 邮箱地址显示（不可编辑）
	emailLabel := widget.NewLabelWithStyle(
		lw.email,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// 自动检测IMAP服务器
	imapConfig, found := GetIMAPServer(lw.email)
	var imapInfoText string
	if found {
		imapInfoText = fmt.Sprintf("✓ 已自动识别IMAP服务器: %s:%d", imapConfig.Server, imapConfig.Port)
	} else {
		imapInfoText = fmt.Sprintf("⚠ 未识别邮箱类型，使用默认配置: %s:%d", imapConfig.Server, imapConfig.Port)
	}

	imapInfo := widget.NewLabelWithStyle(
		imapInfoText,
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// 密码输入框
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入邮箱密码或应用专用密码")

	// 记住密码复选框
	rememberPasswordCheck := widget.NewCheck("记住密码", nil)

	// 尝试加载保存的密码
	credManager := utils.NewCredentialManager()
	savedPassword, err := credManager.LoadPassword(lw.email)
	if err != nil {
		log.Printf("Failed to load saved password: %v", err)
	} else if savedPassword != "" {
		passwordEntry.SetText(savedPassword)
		rememberPasswordCheck.SetChecked(true)
		log.Println("✓ 已自动填充保存的密码")
	}

	// 提示信息
	tipLabel := widget.NewLabel("💡 提示：Gmail等邮箱需要使用应用专用密码")
	tipLabel.Wrapping = fyne.TextWrapWord

	// 登录按钮
	loginBtn := widget.NewButton("登录", func() {
		lw.handleLogin(passwordEntry.Text, imapConfig, rememberPasswordCheck.Checked)
	})
	loginBtn.Importance = widget.HighImportance

	// 回车键登录
	passwordEntry.OnSubmitted = func(password string) {
		lw.handleLogin(password, imapConfig, rememberPasswordCheck.Checked)
	}

	// 布局
	form := container.NewVBox(
		widget.NewLabel("邮箱账号"),
		emailLabel,
		widget.NewLabel(""),
		widget.NewLabel("邮箱密码"),
		passwordEntry,
		rememberPasswordCheck,
		widget.NewLabel(""),
		imapInfo,
		widget.NewLabel(""),
		tipLabel,
	)

	content := container.NewBorder(
		container.NewVBox(
			title,
			subtitle,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			loginBtn,
		),
		nil,
		nil,
		container.NewPadded(form),
	)

	return content
}

// handleLogin 处理登录
func (lw *LoginWindow) handleLogin(password string, imapConfig IMAPServerConfig, rememberPassword bool) {
	// 验证密码
	password = strings.TrimSpace(password)
	if password == "" {
		dialog.ShowError(fmt.Errorf("请输入邮箱密码"), lw.window)
		return
	}

	// 创建配置
	config := &core.Config{
		IMAPServer: imapConfig.Server,
		Port:       imapConfig.Port,
		Username:   lw.email,
		Password:   password,
		MaxEmails:  20,
		TaxRate:    0.13,
	}

	// 显示连接进度
	progress := dialog.NewProgressInfinite("连接中", "正在验证邮箱账号...", lw.window)
	progress.Show()

	// 在后台测试连接
	go func() {
		defer func() {
			// 在主线程中隐藏进度对话框
			fyne.Do(func() {
				progress.Hide()
			})
		}()

		// 测试连接
		emailClient := core.NewEmailClient(config)
		err := emailClient.Connect()
		if err != nil {
			// 在主线程中显示错误对话框
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("登录失败: %v", err), lw.window)
			})
			return
		}
		emailClient.Disconnect()

		// 保存或删除密码
		credManager := utils.NewCredentialManager()
		if rememberPassword {
			if err := credManager.SavePassword(lw.email, password); err != nil {
				log.Printf("Failed to save password: %v", err)
			} else {
				log.Println("✓ 密码已加密保存")
			}
		} else {
			// 如果取消勾选记住密码，删除保存的密码
			if err := credManager.DeletePassword(); err != nil {
				log.Printf("Failed to delete saved password: %v", err)
			}
		}

		// 在主线程中关闭登录窗口并调用回调
		fyne.Do(func() {
			// 登录成功，关闭登录窗口
			lw.window.Close()

			// 调用回调函数，传递配置
			if lw.onLogin != nil {
				lw.onLogin(config)
			}
		})
	}()
}
