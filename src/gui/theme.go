package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme 自定义主题
type CustomTheme struct{}

var _ fyne.Theme = (*CustomTheme)(nil)

// Color 返回主题颜色
func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 102, G: 126, B: 234, A: 255} // #667eea
	case theme.ColorNameButton:
		return color.NRGBA{R: 102, G: 126, B: 234, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 118, G: 75, B: 162, A: 255} // #764ba2
	case theme.ColorNameFocus:
		return color.NRGBA{R: 102, G: 126, B: 234, A: 255}
	case theme.ColorNameBackground:
		if variant == theme.VariantLight {
			return color.NRGBA{R: 248, G: 249, B: 250, A: 255}
		}
		return color.NRGBA{R: 33, G: 37, B: 41, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font 返回字体
func (c *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon 返回图标
func (c *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size 返回尺寸
func (c *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
