package main

import (
	"batchinvoice-pdf/src/core"
	"batchinvoice-pdf/src/gui"

	"fyne.io/fyne/v2/app"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

const (
	// 默认邮箱地址（这里可以配置为你的固定邮箱）
	DefaultEmail = "623601391@qq.com"
)

func main() {
	// 创建 Fyne 应用
	a := app.NewWithID("com.batchinvoice.pdf")
	a.Settings().SetTheme(&gui.CustomTheme{})

	// 创建登录窗口
	loginWindow := gui.NewLoginWindow(a, DefaultEmail, func(config *core.Config) {
		// 登录成功后创建并显示主窗口
		mainWindow := gui.NewMainWindowWithConfig(a, config)
		mainWindow.Show()
	})

	loginWindow.ShowAndRun()
}
