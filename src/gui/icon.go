package gui

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"fyne.io/fyne/v2"
)

// GetAppIcon 生成应用程序图标并返回为资源
// 为了方便打包，不依赖外部文件，图标在运行时由代码绘制并自动内嵌进二进制。
func GetAppIcon() fyne.Resource {
	const size = 256

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// 背景渐变色块（近似蓝紫色）
	top := color.RGBA{R: 0x66, G: 0x7e, B: 0xea, A: 0xff}
	bottom := color.RGBA{R: 0x76, G: 0x4b, B: 0xa2, A: 0xff}
	for y := 0; y < size; y++ {
		t := float64(y) / float64(size-1)
		r := uint8(float64(top.R)*(1-t) + float64(bottom.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bottom.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bottom.B)*t)
		rowColor := color.RGBA{R: r, G: g, B: b, A: 0xff}
		for x := 0; x < size; x++ {
			img.Set(x, y, rowColor)
		}
	}

	// 居中的白色圆角矩形，代表文档/PDF
	docMargin := 32
	docRect := image.Rect(docMargin, docMargin, size-docMargin, size-docMargin)
	docImg := image.NewRGBA(docRect)
	draw.Draw(docImg, docImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(img, docRect, docImg, image.Point{}, draw.Over)

	// 右上角小折角效果
	foldSize := 48
	for dy := 0; dy < foldSize; dy++ {
		for dx := 0; dx < foldSize-dy; dx++ {
			x := docRect.Max.X - foldSize + dx
			y := docRect.Min.Y + dy
			img.Set(x, y, color.RGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff})
		}
	}

	// 左下角的“二维码”方块图案
	qrSize := 120
	qrMargin := 32
	qrRect := image.Rect(
		docRect.Min.X+qrMargin,
		docRect.Max.Y-qrMargin-qrSize,
		docRect.Min.X+qrMargin+qrSize,
		docRect.Max.Y-qrMargin,
	)

	// 背景浅灰
	draw.Draw(img, qrRect, &image.Uniform{color.RGBA{0xf5, 0xf5, 0xf5, 0xff}}, image.Point{}, draw.Over)

	// 简化的“二维码”黑白块
	block := qrSize / 6
	black := color.RGBA{0x22, 0x22, 0x22, 0xff}

	setBlock := func(cx, cy int) {
		startX := qrRect.Min.X + cx*block
		startY := qrRect.Min.Y + cy*block
		for y := startY; y < startY+block; y++ {
			for x := startX; x < startX+block; x++ {
				if x >= qrRect.Max.X || y >= qrRect.Max.Y {
					continue
				}
				img.Set(x, y, black)
			}
		}
	}

	// 几个固定块，形成二维码感
	for _, p := range [][2]int{
		{0, 0}, {1, 0}, {0, 1}, {1, 1},       // 左上
		{4, 0}, {5, 0}, {4, 1}, {5, 1},       // 右上
		{0, 4}, {1, 4}, {0, 5}, {1, 5},       // 左下
		{3, 2}, {2, 3}, {3, 3}, {4, 3}, {3, 4}, // 中间散点
	} {
		setBlock(p[0], p[1])
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)

	return fyne.NewStaticResource("batchinvoice-icon.png", buf.Bytes())
}

