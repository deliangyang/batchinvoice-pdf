package core

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// PDFProcessor PDF处理器
type PDFProcessor struct{}

// NewPDFProcessor 创建PDF处理器
func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{}
}

// ExtractImages 从PDF提取图像
func (p *PDFProcessor) ExtractImages(pdfData []byte) ([]image.Image, error) {
	// 使用pdfcpu提取图像
	reader := bytes.NewReader(pdfData)

	// 配置 - 使用宽松验证模式以支持CJK字体
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	// 解析PDF
	ctx, err := api.ReadContext(reader, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	// 验证PDF（使用宽松模式，跳过字体检查）
	if err := api.ValidateContext(ctx); err != nil {
		return nil, fmt.Errorf("invalid PDF: %w", err)
	}

	var images []image.Image

	// 遍历每一页
	for pageNum := 1; pageNum <= ctx.PageCount; pageNum++ {
		// 渲染页面为图像
		pageImage, err := p.renderPageToImage(ctx, pageNum)
		if err != nil {
			log.Printf("Failed to render page %d: %v", pageNum, err)
			continue
		}
		if pageImage != nil {
			images = append(images, pageImage)
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images extracted from PDF")
	}

	log.Printf("Extracted %d images from PDF", len(images))
	return images, nil
}

// renderPageToImage 渲染PDF页面为图像
func (p *PDFProcessor) renderPageToImage(ctx *model.Context, pageNum int) (image.Image, error) {
	// 这是一个简化的实现
	// 实际使用中，可能需要使用外部工具如 pdfium 或 ghostscript
	// 或者使用 CGO 调用 poppler 库

	// 这里提供一个占位实现
	// TODO: 实现真正的PDF页面渲染

	// 创建一个空白图像作为占位
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	return img, nil
}

// ConvertPDFToImages 转换PDF为图像（使用 pdfcpu 纯 Go 实现）
func (p *PDFProcessor) ConvertPDFToImages(pdfData []byte) ([]image.Image, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "pdf-extract-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 保存 PDF 到临时文件
	pdfPath := filepath.Join(tempDir, "input.pdf")
	if err := os.WriteFile(pdfPath, pdfData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write PDF file: %w", err)
	}

	// 使用 pdfcpu 提取 PDF 中的所有图像到临时目录（纯 Go，无外部命令）
	imagesDir := filepath.Join(tempDir, "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create images dir: %w", err)
	}

	// 配置 - 使用宽松验证模式以支持CJK字体
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	// 选中所有页面
	selectedPages := []string{}
	if err := api.ExtractImagesFile(pdfPath, imagesDir, selectedPages, conf); err != nil {
		return nil, fmt.Errorf("failed to extract images from PDF: %w", err)
	}

	// 读取生成的图片文件（可能是 PNG 或 JPEG）
	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read images dir: %w", err)
	}

	var images []image.Image
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		imgPath := filepath.Join(imagesDir, entry.Name())

		// 读取图像文件
		imgFile, err := os.Open(imgPath)
		if err != nil {
			log.Printf("Failed to open image %s: %v", imgPath, err)
			continue
		}

		var img image.Image
		switch ext {
		case ".png":
			img, err = png.Decode(imgFile)
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(imgFile)
		}
		imgFile.Close()

		if err != nil {
			log.Printf("Failed to decode image %s: %v", imgPath, err)
			continue
		}

		images = append(images, img)
		log.Printf("Loaded image from %s: %dx%d", entry.Name(), img.Bounds().Dx(), img.Bounds().Dy())
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images extracted from PDF")
	}

	log.Printf("Successfully converted PDF to %d images", len(images))
	return images, nil
}

// SaveImageAsPNG 保存图像为PNG
func (p *PDFProcessor) SaveImageAsPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}
	return buf.Bytes(), nil
}
